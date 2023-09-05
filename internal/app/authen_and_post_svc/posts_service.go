package authen_and_post_svc

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/maxuanquang/social-network/internal/pkg/types"
	pb_aap "github.com/maxuanquang/social-network/pkg/types/proto/pb/authen_and_post"
	pb_nfp "github.com/maxuanquang/social-network/pkg/types/proto/pb/newsfeed_publishing"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func (a *AuthenticateAndPostService) CreatePost(ctx context.Context, info *pb_aap.CreatePostRequest) (*pb_aap.CreatePostResponse, error) {
	a.logger.Debug("start creating post")
	defer a.logger.Debug("end creating post")

	exist, _ := a.findUserById(info.GetUserId())
	if !exist {
		return &pb_aap.CreatePostResponse{Status: pb_aap.CreatePostResponse_USER_NOT_FOUND}, nil
	}

	newPost := types.Post{
		UserID:           info.GetUserId(),
		ContentText:      info.GetContentText(),
		ContentImagePath: strings.Join(info.GetContentImagePath(), " "),
	}
	if !info.GetVisible() {
		newPost.DeletedAt.Valid = true
		newPost.DeletedAt.Time = time.Now()
	}
	result := a.db.Create(&newPost)
	if result.Error != nil {
		return nil, result.Error
	}

	// Send user_id and post_id to NewsfeedPublishingClient to announce to followers
	a.nfPubClient.PublishPost(ctx, &pb_nfp.PublishPostRequest{
		UserId:    newPost.UserID,
		PostId:    int64(newPost.ID),
		CreatedAt: timestamppb.New(newPost.CreatedAt),
	})

	return &pb_aap.CreatePostResponse{
		Status: pb_aap.CreatePostResponse_OK,
		PostId: int64(newPost.ID),
	}, nil
}

func (a *AuthenticateAndPostService) EditPost(ctx context.Context, info *pb_aap.EditPostRequest) (*pb_aap.EditPostResponse, error) {
	a.logger.Debug("start editing post")
	defer a.logger.Debug("end editing post")

	exist, user := a.findUserById(info.GetUserId())
	if !exist {
		return &pb_aap.EditPostResponse{Status: pb_aap.EditPostResponse_USER_NOT_FOUND}, nil
	}
	var post types.Post
	result := a.db.Unscoped().First(&post, info.GetPostId())
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &pb_aap.EditPostResponse{Status: pb_aap.EditPostResponse_POST_NOT_FOUND}, nil
	}
	if user.ID != uint(post.UserID) {
		return &pb_aap.EditPostResponse{Status: pb_aap.EditPostResponse_NOT_ALLOWED}, nil
	}

	if info.ContentText != nil {
		post.ContentText = info.GetContentText()
	}
	if info.ContentImagePath != nil {
		post.ContentImagePath = info.GetContentImagePath()
	}
	if info.Visible != nil {
		if info.GetVisible() {
			post.DeletedAt.Valid = false
		} else {
			post.DeletedAt.Valid = true
			post.DeletedAt.Time = time.Now()
		}
	}

	err := a.db.Save(&post).Error
	if err != nil {
		return nil, err
	}
	a.redisClient.Del(ctx, fmt.Sprintf("posts:%d", post.ID), fmt.Sprintf("comments_ids:%d", post.ID), fmt.Sprintf("liked_users_ids:%d", post.ID))

	return &pb_aap.EditPostResponse{
		Status: pb_aap.EditPostResponse_OK,
	}, nil
}

func (a *AuthenticateAndPostService) DeletePost(ctx context.Context, info *pb_aap.DeletePostRequest) (*pb_aap.DeletePostResponse, error) {
	a.logger.Debug("start deleting post")
	defer a.logger.Debug("end deleting post")

	exist, user := a.findUserById(info.GetUserId())
	if !exist {
		return &pb_aap.DeletePostResponse{Status: pb_aap.DeletePostResponse_USER_NOT_FOUND}, nil
	}
	exist, post := a.findPostById(info.GetPostId())
	if !exist {
		return &pb_aap.DeletePostResponse{Status: pb_aap.DeletePostResponse_POST_NOT_FOUND}, nil
	}
	if user.ID != uint(post.UserID) {
		return &pb_aap.DeletePostResponse{Status: pb_aap.DeletePostResponse_NOT_ALLOWED}, nil
	}

	// Delete image on S3
	imgPaths := strings.Split(post.ContentImagePath, " ")
	for _, path := range imgPaths {
		components := strings.Split(path, "/")
		key := components[len(components)-1]

		_, err := a.s3Client.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(os.Getenv("AWS_S3_BUCKET")),
			Key:    aws.String(key),
		})
		if err != nil {
			a.logger.Error(err.Error())
		}
	}

	// Delete post in db
	err := a.db.Delete(&post).Error
	if err != nil {
		return nil, err
	}
	a.redisClient.Del(ctx, fmt.Sprintf("posts:%d", post.ID), fmt.Sprintf("comments_ids:%d", post.ID), fmt.Sprintf("liked_users_ids:%d", post.ID))

	return &pb_aap.DeletePostResponse{
		Status: pb_aap.DeletePostResponse_OK,
	}, nil
}

func (a *AuthenticateAndPostService) GetPostDetailInfo(ctx context.Context, info *pb_aap.GetPostDetailInfoRequest) (*pb_aap.GetPostDetailInfoResponse, error) {
	a.logger.Debug("start getting post")
	defer a.logger.Debug("end getting post")

	var post types.Post
	post, cached := a.findPostInCache(info.GetPostId())
	if !cached {
		a.logger.Debug("post is not cached, getting post from db")
		exist, _ := a.findPostById(info.GetPostId())
		if !exist {
			return &pb_aap.GetPostDetailInfoResponse{Status: pb_aap.GetPostDetailInfoResponse_POST_NOT_FOUND}, nil
		}

		result := a.db.Preload("Comments").Preload("LikedUsers").First(&post, info.GetPostId())
		if result.Error != nil {
			return nil, result.Error
		}
		a.logger.Debug("done getting post from db")

		a.cachePost(&post)
	}

	var comments []*pb_aap.Comment
	for i := range post.Comments {
		comments = append(comments, &pb_aap.Comment{
			CommentId:   int64(post.Comments[i].ID),
			UserId:      post.Comments[i].UserID,
			ContentText: post.Comments[i].ContentText,
			PostId:      int64(post.ID),
		})
	}

	var likedUsers []*pb_aap.Like
	for i := range post.LikedUsers {
		likedUsers = append(likedUsers, &pb_aap.Like{
			UserId: int64(post.LikedUsers[i].ID),
			PostId: int64(post.ID),
		})
	}

	return &pb_aap.GetPostDetailInfoResponse{
		Status: pb_aap.GetPostDetailInfoResponse_OK,
		Post: &pb_aap.PostDetailInfo{
			PostId:           int64(post.ID),
			UserId:           post.UserID,
			ContentText:      post.ContentText,
			ContentImagePath: strings.Split(post.ContentImagePath, " "),
			CreatedAt:        timestamppb.New(post.CreatedAt),
			Comments:         comments,
			LikedUsers:       likedUsers,
		},
	}, nil
}

func (a *AuthenticateAndPostService) CommentPost(ctx context.Context, info *pb_aap.CommentPostRequest) (*pb_aap.CommentPostResponse, error) {
	a.logger.Debug("start commenting post")
	defer a.logger.Debug("end commenting post")

	exist, _ := a.findUserById(info.GetUserId())
	if !exist {
		return &pb_aap.CommentPostResponse{Status: pb_aap.CommentPostResponse_USER_NOT_FOUND}, nil
	}
	exist, _ = a.findPostById(info.GetPostId())
	if !exist {
		return &pb_aap.CommentPostResponse{Status: pb_aap.CommentPostResponse_POST_NOT_FOUND}, nil
	}

	var newComment = types.Comment{
		PostID:      info.GetPostId(),
		UserID:      info.GetUserId(),
		ContentText: info.GetContentText(),
	}
	err := a.db.Create(&newComment).Error
	if err != nil {
		return nil, err
	}

	a.cacheObject(&newComment, fmt.Sprintf("comment:%d", newComment.ID), 15*time.Minute)
	postKey := fmt.Sprintf("post:%d", newComment.PostID)
	commentsIdsKey := fmt.Sprintf("comments_ids:%d", newComment.PostID)
	likedUsersIdsKey := fmt.Sprintf("liked_users_ids:%d", newComment.PostID)
	keyExist := a.redisClient.Exists(context.Background(), postKey, commentsIdsKey, likedUsersIdsKey).Val()
	if keyExist == 3 {
		a.redisClient.RPush(context.Background(), commentsIdsKey, newComment.ID)
		a.redisClient.Expire(context.Background(), postKey, 15*time.Minute)
		a.redisClient.Expire(context.Background(), commentsIdsKey, 15*time.Minute)
		a.redisClient.Expire(context.Background(), likedUsersIdsKey, 15*time.Minute)
	}

	return &pb_aap.CommentPostResponse{
		Status:    pb_aap.CommentPostResponse_OK,
		CommentId: int64(newComment.ID),
	}, nil
}

func (a *AuthenticateAndPostService) LikePost(ctx context.Context, info *pb_aap.LikePostRequest) (*pb_aap.LikePostResponse, error) {
	a.logger.Debug("start liking post")
	defer a.logger.Debug("end liking post")

	exist, _ := a.findUserById(info.GetUserId())
	if !exist {
		return &pb_aap.LikePostResponse{Status: pb_aap.LikePostResponse_USER_NOT_FOUND}, nil
	}
	exist, _ = a.findPostById(info.GetPostId())
	if !exist {
		return &pb_aap.LikePostResponse{Status: pb_aap.LikePostResponse_POST_NOT_FOUND}, nil
	}

	var like types.Like
	a.db.Raw("select * from `like` where user_id = ? and post_id = ?", info.GetUserId(), info.GetPostId()).Scan(&like)
	fmt.Println(like)
	if like.UserId != 0 && like.PostId != 0 {
		return &pb_aap.LikePostResponse{
			Status: pb_aap.LikePostResponse_ALREADY_LIKED,
		}, nil
	} else {
		rowsAffected := a.db.Exec("insert into `like` (user_id, post_id) values (?, ?)",
			info.GetUserId(),
			info.GetPostId(),
		).RowsAffected
		if rowsAffected == 0 {
			return nil, fmt.Errorf("can not insert into `like` table")
		}
	}

	postKey := fmt.Sprintf("post:%d", info.GetPostId())
	commentsIdsKey := fmt.Sprintf("comments_ids:%d", info.GetPostId())
	likedUsersIdsKey := fmt.Sprintf("liked_users_ids:%d", info.GetPostId())
	keyExist := a.redisClient.Exists(context.Background(), postKey, commentsIdsKey, likedUsersIdsKey).Val()
	if keyExist == 3 {
		a.redisClient.RPush(context.Background(), likedUsersIdsKey, info.GetUserId())
		a.redisClient.Expire(context.Background(), postKey, 15*time.Minute)
		a.redisClient.Expire(context.Background(), commentsIdsKey, 15*time.Minute)
		a.redisClient.Expire(context.Background(), likedUsersIdsKey, 15*time.Minute)
	}

	return &pb_aap.LikePostResponse{
		Status: pb_aap.LikePostResponse_OK,
	}, nil
}

func (a *AuthenticateAndPostService) GetS3PresignedUrl(ctx context.Context, info *pb_aap.GetS3PresignedUrlRequest) (*pb_aap.GetS3PresignedUrlResponse, error) {
	a.logger.Debug("start getting s3 presigned url")
	defer a.logger.Debug("end getting s3 presigned url")

	// userExist, _ := a.findUserById(info.GetUserId())
	// if !userExist {
	// 	return &pb_aap.GetS3PresignedUrlResponse{Status: pb_aap.GetS3PresignedUrlResponse_USER_NOT_FOUND}, nil
	// }

	req, _ := a.s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET")),
		Key:    aws.String(a.getRandomS3Key(64)),
	})
	url, err := req.Presign(time.Minute)
	expirationTime := time.Now().Add(time.Minute)
	if err != nil {
		a.logger.Error(err.Error())
		return nil, err
	}

	return &pb_aap.GetS3PresignedUrlResponse{
		Status:         pb_aap.GetS3PresignedUrlResponse_OK,
		Url:            url,
		ExpirationTime: timestamppb.New(expirationTime),
	}, nil
}

// findPostById checks if a post with provided postId exists in database
func (a *AuthenticateAndPostService) findPostById(postId int64) (exist bool, post types.Post) {
	result := a.db.First(&post, postId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, types.Post{}
	}
	return true, post
}

// cachPost caches post and relevant information in cache
func (a *AuthenticateAndPostService) cachePost(post *types.Post) error {
	a.logger.Debug("start caching post")
	defer a.logger.Debug("end caching post")

	// Only cache post when len(comments_ids) > 0 and len(liked_users_ids) > 0
	if len(post.Comments) == 0 || len(post.LikedUsers) == 0 {
		return fmt.Errorf("do not cache because length of comments or liked_users equal 0")
	}

	// Cache comments ids
	var commentsIds []interface{}
	for _, comment := range post.Comments {
		commentsIds = append(commentsIds, comment.ID)
	}
	commentsIdsKey := fmt.Sprintf("comments_ids:%d", post.ID)
	a.redisClient.RPush(context.Background(), commentsIdsKey, commentsIds...)
	a.redisClient.Expire(context.Background(), commentsIdsKey, 15*time.Minute)

	// Cache liked users ids
	var likeUsersIds []interface{}
	for _, user := range post.LikedUsers {
		likeUsersIds = append(likeUsersIds, user.ID)
	}
	likedUsersKey := fmt.Sprintf("liked_users_ids:%d", post.ID)
	a.redisClient.RPush(context.Background(), likedUsersKey, likeUsersIds...)
	a.redisClient.Expire(context.Background(), likedUsersKey, 15*time.Minute)

	// Cache each comment
	for _, comment := range post.Comments {
		commentKey := fmt.Sprintf("comment:%d", comment.ID)
		a.cacheObject(comment, commentKey, 15*time.Minute)
	}

	// Cache each user
	for _, user := range post.LikedUsers {
		userKey := fmt.Sprintf("user:%d", user.ID)
		a.cacheObject(user, userKey, 15*time.Minute)
	}

	// Cache metadata of post (comments and liked users are not included - they have their own redis key above)
	postKey := fmt.Sprintf("post:%d", post.ID)
	a.cacheObject(&types.Post{
		Model:            post.Model,
		ContentText:      post.ContentText,
		ContentImagePath: post.ContentImagePath,
		UserID:           post.UserID,
	}, postKey, 15*time.Minute)

	return nil
}

// findPostInCache checks if a post with provided postId exists in cache
func (a *AuthenticateAndPostService) findPostInCache(postId int64) (types.Post, bool) {
	// Check all relevant cache keys
	keyExist := a.redisClient.Exists(context.Background(),
		fmt.Sprintf("post:%d", int(postId)),
		fmt.Sprintf("comments_ids:%d", int(postId)),
		fmt.Sprintf("liked_users_ids:%d", int(postId)),
	).Val()
	if keyExist != 3 {
		return types.Post{}, false
	}

	// Get post from redis
	var post types.Post
	postKey := fmt.Sprintf("post:%d", int(postId))
	err := a.getObjectFromCache(postKey, &post, 15*time.Minute)
	if err != nil {
		return types.Post{}, false
	}

	// Get comments from redis
	commentsIdsKey := fmt.Sprintf("comments_ids:%d", int(postId))
	commentsIds := a.redisClient.LRange(context.Background(), commentsIdsKey, 0, -1).Val()
	a.redisClient.Expire(context.Background(), commentsIdsKey, 15*time.Minute)

	var comments []*types.Comment
	for _, comment_id := range commentsIds {
		var comment types.Comment
		commentKey := "comment:" + comment_id
		err := a.getObjectFromCache(commentKey, &comment, 15*time.Minute)
		if err != nil {
			return types.Post{}, false
		}
		comments = append(comments, &comment)
	}

	// Get post's likes from redis
	likedUsersIdsKey := fmt.Sprintf("liked_users_ids:%d", int(postId))
	likedUsersIds, err := a.redisClient.LRange(context.Background(), likedUsersIdsKey, 0, -1).Result()
	if err != nil {
		return types.Post{}, false
	}
	a.redisClient.Expire(context.Background(), likedUsersIdsKey, 15*time.Minute)
	var likedUsers []*types.User
	for _, user_id := range likedUsersIds {
		var user types.User
		userKey := "user:" + user_id
		err := a.getObjectFromCache(userKey, &user, 15*time.Minute)
		if err != nil {
			return types.Post{}, false
		}
		likedUsers = append(likedUsers, &user)
	}

	// return
	post.Comments = comments
	post.LikedUsers = likedUsers
	return post, true
}

func (a *AuthenticateAndPostService) cacheObject(objectPointer interface{}, cacheKey string, expireTime time.Duration) error {
	objValue := reflect.ValueOf(objectPointer)
	if objValue.Kind() != reflect.Ptr || objValue.IsNil() {
		return fmt.Errorf("obj must be a non-nil pointer to a struct")
	}

	keyExist := a.redisClient.Exists(context.Background(), cacheKey).Val()
	if keyExist == 1 {
		a.redisClient.Expire(context.Background(), cacheKey, expireTime)
		return nil
	} else {
		byteObj, err := json.Marshal(objectPointer)
		if err != nil {
			a.logger.Debug(err.Error())
			return err
		} else {
			a.redisClient.Set(context.Background(), cacheKey, byteObj, expireTime)
			return nil
		}
	}
}

func (a *AuthenticateAndPostService) getObjectFromCache(cacheKey string, objectPointer interface{}, newExpireTime time.Duration) error {
	objValue := reflect.ValueOf(objectPointer)
	if objValue.Kind() != reflect.Ptr || objValue.IsNil() {
		return fmt.Errorf("obj must be a non-nil pointer to a struct")
	}

	keyExist := a.redisClient.Exists(context.Background(), cacheKey).Val()
	if keyExist == 0 {
		a.logger.Debug(fmt.Sprintf("%s does not exist in redis", cacheKey))
		return fmt.Errorf("%s does not exist in redis", cacheKey)
	}

	byteObj, err := a.redisClient.Get(context.Background(), cacheKey).Bytes()
	if err != nil {
		a.logger.Debug(err.Error())
		return err
	}
	a.redisClient.Expire(context.Background(), cacheKey, newExpireTime)
	err = json.Unmarshal(byteObj, &objectPointer)
	if err != nil {
		a.logger.Debug(err.Error())
		return err
	}

	return nil
}

func (a *AuthenticateAndPostService) getRandomS3Key(length int) string {
	randomBytes := make([]byte, 64)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(randomBytes)[:length]
}
