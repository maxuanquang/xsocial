package authen_and_post_svc

import (
	"context"
	"errors"
	"strings"

	"github.com/maxuanquang/social-network/internal/pkg/types"
	pb_aap "github.com/maxuanquang/social-network/pkg/types/proto/pb/authen_and_post"
)

func (a *AuthenticateAndPostService) CreatePost(ctx context.Context, info *pb_aap.PostDetailInfo) (*pb_aap.ActionResult, error) {
	// Check if the user exists
	existed, _ := a.checkUserId(info.GetUserId())
	if !existed {
		return nil, errors.New("user does not exist")
	}

	// Execute the add post command
	err := a.db.Exec("insert into post (user_id, content_text, content_image_path, visible) values (?, ?, ?, ?)",
		info.GetUserId(),
		info.GetContentText(),
		strings.Join(info.GetContentImagePath(), " "),
		info.GetVisible(),
	).Error
	if err != nil {
		return nil, err
	}

	return &pb_aap.ActionResult{Status: pb_aap.ActionStatus_SUCCEEDED}, nil
}

func (a *AuthenticateAndPostService) EditPost(ctx context.Context, info *pb_aap.PostDetailInfo) (*pb_aap.ActionResult, error) {
	// Check if the user exists
	existed, _ := a.checkUserId(info.GetUserId())
	if !existed {
		return nil, errors.New("user does not exist")
	}

	// Check if the post exists
	existed, _ = a.checkPostId(info.GetPostId())
	if !existed {
		return nil, errors.New("user does not exist")
	}

	// Execute the edit post command
	if info.GetContentText() != "" {
		err := a.db.Exec("update post set content_text = ? where id = ? and user_id = ?",
			info.GetContentText(),
			info.GetPostId(),
			info.GetUserId(),
		).Error
		if err != nil {
			return nil, err
		}
	}

	if len(info.GetContentImagePath()) > 0 {
		err := a.db.Exec("update post set content_image_path = ? where id = ? and user_id = ?",
			strings.Join(info.GetContentImagePath(), " "),
			info.GetPostId(),
			info.GetUserId(),
		).Error
		if err != nil {
			return nil, err
		}
	}

	return &pb_aap.ActionResult{Status: pb_aap.ActionStatus_SUCCEEDED}, nil
}

func (a *AuthenticateAndPostService) DeletePost(ctx context.Context, info *pb_aap.PostInfo) (*pb_aap.ActionResult, error) {
	// Check if the user exists
	existed, _ := a.checkUserId(info.GetUserId())
	if !existed {
		return nil, errors.New("user does not exist")
	}

	// Check if the post exists
	existed, _ = a.checkPostId(info.GetPostId())
	if !existed {
		return nil, errors.New("post does not exist")
	}

	// Execute the delete post command
	rowsAffected := a.db.Exec("update post set `visible` = false where id = ? and user_id = ?",
		info.GetPostId(),
		info.GetUserId(),
	).RowsAffected
	if rowsAffected == 0 {
		return nil, errors.New("can not delete post")
	}

	return &pb_aap.ActionResult{Status: pb_aap.ActionStatus_SUCCEEDED}, nil
}

func (a *AuthenticateAndPostService) GetPost(ctx context.Context, info *pb_aap.PostInfo) (*pb_aap.PostDetailInfo, error) {
	// Check if the post exists
	postModel := types.Post{}
	err := a.db.Raw("select * from post where id = ?", info.GetPostId()).Scan(&postModel).Error
	if err != nil {
		return nil, err
	}
	if postModel.ID == 0 {
		return nil, errors.New("can not find post")
	}

	// Get comments
	var commentModels []types.Comment
	err = a.db.Raw("select * from comment where post_id = ?", info.GetPostId()).Scan(&commentModels).Error
	if err != nil {
		return nil, err
	}

	var comments []*pb_aap.CommentInfo
	for _, comment := range commentModels {
		comments = append(comments, &pb_aap.CommentInfo{
			CommentId: int64(comment.ID),
			PostId:    int64(comment.PostID),
			UserId:    int64(comment.UserID),
			Content:   comment.Content,
		})
	}

	// Get likes
	var likeModels []types.Like
	err = a.db.Raw("select * from `like` where post_id = ?", info.GetPostId()).Scan(&likeModels).Error
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
		PostId:           int64(postModel.ID),
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

func (a *AuthenticateAndPostService) CommentPost(ctx context.Context, info *pb_aap.CommentInfo) (*pb_aap.ActionResult, error) {
	// Check if user exists
	existed, _ := a.checkUserId(info.GetUserId())
	if !existed {
		return nil, errors.New("user does not exist")
	}

	// Check if post exists
	existed, _ = a.checkPostId(info.GetPostId())
	if !existed {
		return nil, errors.New("post does not exist")
	}

	// Execute command
	rowsAffected := a.db.Exec("insert into comment (post_id, user_id, content) values (?, ?, ?)", info.GetPostId(), info.GetUserId(), info.GetContent()).RowsAffected
	if rowsAffected == 0 {
		return nil, errors.New("can not comment post")
	}

	return a.NewActionResult(pb_aap.ActionStatus_SUCCEEDED), nil
}

func (a *AuthenticateAndPostService) LikePost(ctx context.Context, info *pb_aap.LikeInfo) (*pb_aap.ActionResult, error) {
	// Check if user exists
	existed, _ := a.checkUserId(info.GetUserId())
	if !existed {
		return nil, errors.New("user does not exist")
	}

	// Check if post exists
	existed, _ = a.checkPostId(info.GetPostId())
	if !existed {
		return nil, errors.New("post does not exist")
	}

	// Execute command
	rowsAffected := a.db.Exec("insert into `like` (post_id, user_id) values (?, ?)", info.GetPostId(), info.GetUserId(), ).RowsAffected
	if rowsAffected == 0 {
		return nil, errors.New("can not like post")
	}

	return a.NewActionResult(pb_aap.ActionStatus_SUCCEEDED), nil
}