package authen_and_post_svc

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/maxuanquang/social-network/internal/pkg/types"
	pb_aap "github.com/maxuanquang/social-network/pkg/types/proto/pb/authen_and_post"
	"github.com/segmentio/kafka-go"
)

func (a *AuthenticateAndPostService) CreatePost(ctx context.Context, postInfo *pb_aap.PostDetailInfo) (*pb_aap.ActionResult, error) {
	// Check if the user exists
	existed, _ := a.checkUserId(postInfo.GetUserId())
	if !existed {
		return nil, errors.New("user does not exist")
	}

	// Execute the add post command
	newPost := types.Post{
		UserID:           postInfo.GetUserId(),
		ContentText:      postInfo.GetContentText(),
		ContentImagePath: strings.Join(postInfo.GetContentImagePath(), " "),
		Visible:          postInfo.GetVisible(),
	}
	err := a.db.Create(&newPost).Error
	if err != nil {
		return nil, err
	}

	// Send post to message queue to announce to followers
	post, err := a.GetPost(ctx, &pb_aap.PostInfo{Id: int64(newPost.ID)})
	if err != nil {
		return nil, err
	}
	jsonPost, err := json.Marshal(post)
	if err != nil {
		return nil, err
	}
	err = a.kafkaWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte("post"),
		Value: jsonPost,
		Headers: []kafka.Header{
			{Key: "Content-Type", Value: []byte("application/json")},
		},
	})
	if err != nil {
		return nil, err
	}

	return &pb_aap.ActionResult{Status: pb_aap.ActionStatus_SUCCEEDED}, nil
}

func (a *AuthenticateAndPostService) EditPost(ctx context.Context, postInfo *pb_aap.PostDetailInfo) (*pb_aap.ActionResult, error) {
	// Check if the user exists
	existed, _ := a.checkUserId(postInfo.GetUserId())
	if !existed {
		return nil, errors.New("user does not exist")
	}

	// Check if the post exists
	existed, _ = a.checkPostId(postInfo.GetId())
	if !existed {
		return nil, errors.New("user does not exist")
	}

	// Execute the edit post command
	if postInfo.GetContentText() != "" {
		err := a.db.Exec("update post set content_text = ? where id = ? and user_id = ?",
			postInfo.GetContentText(),
			postInfo.GetId(),
			postInfo.GetUserId(),
		).Error
		if err != nil {
			return nil, err
		}
	}

	if len(postInfo.GetContentImagePath()) > 0 {
		err := a.db.Exec("update post set content_image_path = ? where id = ? and user_id = ?",
			strings.Join(postInfo.GetContentImagePath(), " "),
			postInfo.GetId(),
			postInfo.GetUserId(),
		).Error
		if err != nil {
			return nil, err
		}
	}

	return &pb_aap.ActionResult{Status: pb_aap.ActionStatus_SUCCEEDED}, nil
}

func (a *AuthenticateAndPostService) DeletePost(ctx context.Context, postInfo *pb_aap.PostInfo) (*pb_aap.ActionResult, error) {
	// Check if the user exists
	existed, _ := a.checkUserId(postInfo.GetUserId())
	if !existed {
		return nil, errors.New("user does not exist")
	}

	// Check if the post exists
	existed, _ = a.checkPostId(postInfo.GetId())
	if !existed {
		return nil, errors.New("post does not exist")
	}

	// Execute the delete post command
	rowsAffected := a.db.Exec("update post set `visible` = false where id = ? and user_id = ?",
		postInfo.GetId(),
		postInfo.GetUserId(),
	).RowsAffected
	if rowsAffected == 0 {
		return nil, errors.New("can not delete post")
	}

	return &pb_aap.ActionResult{Status: pb_aap.ActionStatus_SUCCEEDED}, nil
}

func (a *AuthenticateAndPostService) GetPost(ctx context.Context, postInfo *pb_aap.PostInfo) (*pb_aap.PostDetailInfo, error) {
	// Check if the post exists
	postModel := types.Post{}
	err := a.db.Raw("select * from post where id = ?", postInfo.GetId()).Scan(&postModel).Error
	if err != nil {
		return nil, err
	}
	if postModel.ID == 0 {
		return nil, errors.New("can not find post")
	}

	// Get comments
	var commentModels []types.Comment
	err = a.db.Raw("select * from comment where post_id = ?", postInfo.GetId()).Scan(&commentModels).Error
	if err != nil {
		return nil, err
	}

	var comments []*pb_aap.CommentInfo
	for _, comment := range commentModels {
		comments = append(comments, &pb_aap.CommentInfo{
			Id:      int64(comment.ID),
			PostId:  int64(comment.PostID),
			UserId:  int64(comment.UserID),
			Content: comment.Content,
		})
	}

	// Get likes
	var likeModels []types.Like
	err = a.db.Raw("select * from `like` where post_id = ?", postInfo.GetId()).Scan(&likeModels).Error
	if err != nil {
		return nil, err
	}

	var likes []*pb_aap.LikeInfo
	for _, like := range likeModels {
		likes = append(likes, &pb_aap.LikeInfo{
			PostId: int64(like.PostID),
			UserId: int64(like.UserID),
		})
	}

	// If the post exists, return the post
	postDetailInfo := pb_aap.PostDetailInfo{
		Id:               int64(postModel.ID),
		UserId:           int64(postModel.UserID),
		ContentText:      postModel.ContentText,
		ContentImagePath: strings.Split(postModel.ContentImagePath, " "),
		Visible:          postModel.Visible,
		CreatedAt:        postModel.CreatedAt.Unix(),
		Comments:         comments,
		Likes:            likes,
	}
	return &postDetailInfo, nil
}

func (a *AuthenticateAndPostService) CommentPost(ctx context.Context, commentInfo *pb_aap.CommentInfo) (*pb_aap.ActionResult, error) {
	// Check if user exists
	existed, _ := a.checkUserId(commentInfo.GetUserId())
	if !existed {
		return nil, errors.New("user does not exist")
	}

	// Check if post exists
	existed, _ = a.checkPostId(commentInfo.GetPostId())
	if !existed {
		return nil, errors.New("post does not exist")
	}

	// Execute command
	rowsAffected := a.db.Exec("insert into comment (post_id, user_id, content) values (?, ?, ?)", commentInfo.GetPostId(), commentInfo.GetUserId(), commentInfo.GetContent()).RowsAffected
	if rowsAffected == 0 {
		return nil, errors.New("can not comment post")
	}

	return a.NewActionResult(pb_aap.ActionStatus_SUCCEEDED), nil
}

func (a *AuthenticateAndPostService) LikePost(ctx context.Context, likeInfo *pb_aap.LikeInfo) (*pb_aap.ActionResult, error) {
	// Check if user exists
	existed, _ := a.checkUserId(likeInfo.GetUserId())
	if !existed {
		return nil, errors.New("user does not exist")
	}

	// Check if post exists
	existed, _ = a.checkPostId(likeInfo.GetPostId())
	if !existed {
		return nil, errors.New("post does not exist")
	}

	// Execute command
	rowsAffected := a.db.Exec("insert into `like` (post_id, user_id) values (?, ?)", likeInfo.GetPostId(), likeInfo.GetUserId()).RowsAffected
	if rowsAffected == 0 {
		return nil, errors.New("can not like post")
	}

	return a.NewActionResult(pb_aap.ActionStatus_SUCCEEDED), nil
}
