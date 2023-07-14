package authen_and_post_svc

import (
	"context"
	"errors"

	"github.com/maxuanquang/social-network/internal/pkg/types"
	pb_aap "github.com/maxuanquang/social-network/pkg/types/proto/pb/authen_and_post"
)

func (a *AuthenticateAndPostService) GetUserFollower(ctx context.Context, userInfo *pb_aap.UserInfo) (*pb_aap.UserFollowerInfo, error) {
	// Check if the username which is changing information exists in the database
	existed, userModel := a.checkUserId(userInfo.GetId())
	if !existed {
		return nil, errors.New("user does not exist")
	}

	// If the user exists, return the followers
	var followers []types.User
	err := a.db.Raw("select u.* from following f join user u on f.follower_id = u.id where f.user_id = ?", userModel.ID).Scan(&followers).Error
	if err != nil {
		return nil, err
	}

	returnUserFolower := pb_aap.UserFollowerInfo{}
	for _, follower := range followers {
		followerInfo := pb_aap.UserInfo{Id: follower.ID, UserName: follower.UserName}
		returnUserFolower.Followers = append(returnUserFolower.Followers, &followerInfo)
	}

	return &returnUserFolower, nil
}

func (a *AuthenticateAndPostService) FollowUser(ctx context.Context, info *pb_aap.UserAndFollowerInfo) (*pb_aap.ActionResult, error) {
	// Check if the user exists
	existed, _ := a.checkUserId(info.GetUser().GetId())
	if !existed {
		return nil, errors.New("user does not exist")
	}

	// Check if the follower exists
	existed, _ = a.checkUserId(info.GetFollower().GetId())
	if !existed {
		return nil, errors.New("follower does not exist")
	}

	// Execute the follow command
	err := a.db.Exec("insert into following (user_id, follower_id) values (?, ?)",
		info.GetUser().GetId(),
		info.GetFollower().GetId(),
	).Error
	if err != nil {
		return nil, err
	}

	return &pb_aap.ActionResult{Status: pb_aap.ActionStatus_SUCCEEDED}, nil
}

func (a *AuthenticateAndPostService) UnfollowUser(ctx context.Context, info *pb_aap.UserAndFollowerInfo) (*pb_aap.ActionResult, error) {
	// Check if the user exists
	existed, _ := a.checkUserId(info.GetUser().GetId())
	if !existed {
		return nil, errors.New("user does not exist")
	}

	// Check if the follower exists
	existed, _ = a.checkUserId(info.GetFollower().GetId())
	if !existed {
		return nil, errors.New("follower does not exist")
	}

	// Execute the unfollow command
	err := a.db.Exec("delete from following where user_id = ? and follower_id = ?",
		info.GetUser().GetId(),
		info.GetFollower().GetId(),
	).Error
	if err != nil {
		return nil, err
	}

	return &pb_aap.ActionResult{Status: pb_aap.ActionStatus_SUCCEEDED}, nil
}

func (a *AuthenticateAndPostService) GetUserPost(ctx context.Context, userInfo *pb_aap.UserInfo) (*pb_aap.UserPostDetailInfo, error) {
	// Check if the user exists
	existed, _ := a.checkUserId(userInfo.GetId())
	if !existed {
		return nil, errors.New("user does not exist")
	}

	// Execute command
	var postModels []types.Post
	err := a.db.Raw("select * from post where user_id = ?", userInfo.GetId()).Scan(&postModels).Error
	if err != nil {
		return nil, err
	}

	// Return
	var posts []*pb_aap.PostDetailInfo
	for _, model := range postModels {
		postId := int64(model.ID)
		post, err := a.GetPost(ctx, &pb_aap.PostInfo{Id: postId})
		if err == nil {
			posts = append(posts, post)
		}
	}

	return &pb_aap.UserPostDetailInfo{Posts: posts}, nil
}
