package authen_and_post_svc

import (
	"context"
	"errors"
	"time"

	// "encoding/json"
	// "errors"
	"strings"

	"github.com/maxuanquang/social-network/internal/pkg/types"
	pb_aap "github.com/maxuanquang/social-network/pkg/types/proto/pb/authen_and_post"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	// "github.com/segmentio/kafka-go"
)

func (a *AuthenticateAndPostService) CreatePost(ctx context.Context, info *pb_aap.CreatePostRequest) (*pb_aap.CreatePostResponse, error) {
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

	// // Send post to message queue to announce to followers
	// post, err := a.GetPost(ctx, &pb_aap.PostInfo{Id: int64(newPost.ID)})
	// if err != nil {
	// 	return nil, err
	// }
	// jsonPost, err := json.Marshal(post)
	// if err != nil {
	// 	return nil, err
	// }
	// err = a.kafkaWriter.WriteMessages(ctx, kafka.Message{
	// 	Key:   []byte("post"),
	// 	Value: jsonPost,
	// 	Headers: []kafka.Header{
	// 		{Key: "Content-Type", Value: []byte("application/json")},
	// 	},
	// })
	// if err != nil {
	// 	return nil, err
	// }

	return &pb_aap.CreatePostResponse{
		Status: pb_aap.CreatePostResponse_OK,
		PostId: int64(newPost.ID),
	}, nil
}

func (a *AuthenticateAndPostService) EditPost(ctx context.Context, info *pb_aap.EditPostRequest) (*pb_aap.EditPostResponse, error) {
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
	return &pb_aap.EditPostResponse{
		Status: pb_aap.EditPostResponse_OK,
	}, nil
}

func (a *AuthenticateAndPostService) DeletePost(ctx context.Context, info *pb_aap.DeletePostRequest) (*pb_aap.DeletePostResponse, error) {
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

	err := a.db.Delete(&post).Error
	if err != nil {
		return nil, err
	}
	return &pb_aap.DeletePostResponse{
		Status: pb_aap.DeletePostResponse_OK,
	}, nil
}

func (a *AuthenticateAndPostService) GetPostDetailInfo(ctx context.Context, info *pb_aap.GetPostDetailInfoRequest) (*pb_aap.GetPostDetailInfoResponse, error) {
	exist, _ := a.findPostById(info.GetPostId())
	if !exist {
		return &pb_aap.GetPostDetailInfoResponse{Status: pb_aap.GetPostDetailInfoResponse_POST_NOT_FOUND}, nil
	}

	var post types.Post
	result := a.db.Preload("Comments").Preload("LikedUsers").First(&post, info.GetPostId())
	if result.Error != nil {
		return nil, result.Error
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

	return &pb_aap.CommentPostResponse{
		Status:    pb_aap.CommentPostResponse_OK,
		CommentId: int64(newComment.ID),
	}, nil
}

func (a *AuthenticateAndPostService) LikePost(ctx context.Context, info *pb_aap.LikePostRequest) (*pb_aap.LikePostResponse, error) {
	exist, user := a.findUserById(info.GetUserId())
	if !exist {
		return &pb_aap.LikePostResponse{Status: pb_aap.LikePostResponse_USER_NOT_FOUND}, nil
	}
	exist, _ = a.findPostById(info.GetPostId())
	if !exist {
		return &pb_aap.LikePostResponse{Status: pb_aap.LikePostResponse_POST_NOT_FOUND}, nil
	}

	var post types.Post
	err := a.db.Preload("LikedUsers").First(&post, info.GetPostId()).Error
	if err != nil {
		return nil, err
	}
	err = a.db.Model(&post).Association("LikedUsers").Append(&user)
	if err != nil {
		return nil, err
	}

	return &pb_aap.LikePostResponse{
		Status: pb_aap.LikePostResponse_OK,
	}, nil
}
