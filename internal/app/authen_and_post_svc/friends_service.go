package authen_and_post_svc

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/maxuanquang/social-network/internal/pkg/types"
	pb_aap "github.com/maxuanquang/social-network/pkg/types/proto/pb/authen_and_post"
)

func (a *AuthenticateAndPostService) GetUserFollower(ctx context.Context, info *pb_aap.GetUserFollowerRequest) (*pb_aap.GetUserFollowerResponse, error) {
	followersKey := fmt.Sprintf("followers:%d", info.GetUserId())
	keyExist := (a.redisClient.Exists(context.Background(), followersKey).Val() == 1)
	if keyExist {
		a.redisClient.Expire(context.Background(), followersKey, 15*time.Minute)

		var followersIds []int64
		for _, id := range a.redisClient.LRange(context.Background(), followersKey, 0, -1).Val() {
			intId, err := strconv.Atoi(id)
			if err != nil {
				a.logger.Debug(err.Error())
				continue
			}
			followersIds = append(followersIds, int64(intId))
		}
		return &pb_aap.GetUserFollowerResponse{
			Status:       pb_aap.GetUserFollowerResponse_OK,
			FollowersIds: followersIds,
		}, nil
	}

	exist, _ := a.findUserById(info.GetUserId())
	if !exist {
		return &pb_aap.GetUserFollowerResponse{
			Status: pb_aap.GetUserFollowerResponse_USER_NOT_FOUND,
		}, nil
	}

	var user types.User
	result := a.db.Preload("Followers").First(&user, info.GetUserId())
	if result.Error != nil {
		return nil, result.Error
	}

	// Caching
	var cacheFollowersIds []interface{}
	for _, follower := range user.Followers {
		cacheFollowersIds = append(cacheFollowersIds, follower.ID)
	}
	a.redisClient.RPush(context.Background(), followersKey, cacheFollowersIds...)
	a.redisClient.Expire(context.Background(), followersKey, 15*time.Minute)

	var followersIds []int64
	for _, follower := range user.Followers {
		followersIds = append(followersIds, int64(follower.ID))
	}
	return &pb_aap.GetUserFollowerResponse{
		Status:       pb_aap.GetUserFollowerResponse_OK,
		FollowersIds: followersIds,
	}, nil
}

func (a *AuthenticateAndPostService) GetUserFollowing(ctx context.Context, info *pb_aap.GetUserFollowingRequest) (*pb_aap.GetUserFollowingResponse, error) {
	followingsKey := fmt.Sprintf("followings:%d", info.GetUserId())
	keyExist := (a.redisClient.Exists(context.Background(), followingsKey).Val() == 1)
	if keyExist {
		a.redisClient.Expire(context.Background(), followingsKey, 15*time.Minute)

		var followingsIds []int64
		for _, id := range a.redisClient.LRange(context.Background(), followingsKey, 0, -1).Val() {
			intId, err := strconv.Atoi(id)
			if err != nil {
				a.logger.Debug(err.Error())
				continue
			}
			followingsIds = append(followingsIds, int64(intId))
		}
		return &pb_aap.GetUserFollowingResponse{
			Status:        pb_aap.GetUserFollowingResponse_OK,
			FollowingsIds: followingsIds,
		}, nil
	}

	exist, _ := a.findUserById(info.GetUserId())
	if !exist {
		return &pb_aap.GetUserFollowingResponse{
			Status: pb_aap.GetUserFollowingResponse_USER_NOT_FOUND,
		}, nil
	}

	var user types.User
	result := a.db.Preload("Followings").First(&user, info.GetUserId())
	if result.Error != nil {
		return nil, result.Error
	}

	// Caching
	var cacheFollowingsIds []interface{}
	for _, following := range user.Followings {
		cacheFollowingsIds = append(cacheFollowingsIds, following.ID)
	}
	a.redisClient.RPush(context.Background(), followingsKey, cacheFollowingsIds...)
	a.redisClient.Expire(context.Background(), followingsKey, 15*time.Minute)

	var followingsIds []int64
	for _, following := range user.Followings {
		followingsIds = append(followingsIds, int64(following.ID))
	}
	return &pb_aap.GetUserFollowingResponse{
		Status:        pb_aap.GetUserFollowingResponse_OK,
		FollowingsIds: followingsIds,
	}, nil
}

func (a *AuthenticateAndPostService) FollowUser(ctx context.Context, info *pb_aap.FollowUserRequest) (*pb_aap.FollowUserResponse, error) {
	// Check if the user exists
	exist, _ := a.findUserById(info.GetUserId())
	if !exist {
		return &pb_aap.FollowUserResponse{Status: pb_aap.FollowUserResponse_USER_NOT_FOUND}, nil
	}
	exist, friend := a.findUserById(info.GetFollowingId())
	if !exist {
		return &pb_aap.FollowUserResponse{Status: pb_aap.FollowUserResponse_USER_NOT_FOUND}, nil
	}

	var user types.User
	a.db.Preload("Followings").First(&user, info.GetUserId())
	for _, following := range user.Followings {
		if following.ID == uint(info.GetFollowingId()) {
			return &pb_aap.FollowUserResponse{Status: pb_aap.FollowUserResponse_ALREADY_FOLLOWED}, nil
		}
	}

	err := a.db.Model(&user).Association("Followings").Append(&friend)
	if err != nil {
		return nil, err
	}

	// Update cache
	followingsKey := fmt.Sprintf("followings:%d", info.GetUserId())
	keyExist := (a.redisClient.Exists(context.Background(), followingsKey).Val() == 1)
	if keyExist {
		a.redisClient.RPush(context.Background(), followingsKey, info.GetFollowingId())
	}
	followersKey := fmt.Sprintf("followers:%d", info.GetFollowingId())
	keyExist = (a.redisClient.Exists(context.Background(), followersKey).Val() == 1)
	if keyExist {
		a.redisClient.RPush(context.Background(), followersKey, info.GetUserId())
	}

	return &pb_aap.FollowUserResponse{
		Status: pb_aap.FollowUserResponse_OK,
	}, nil
}

func (a *AuthenticateAndPostService) UnfollowUser(ctx context.Context, info *pb_aap.UnfollowUserRequest) (*pb_aap.UnfollowUserResponse, error) {
	exist, _ := a.findUserById(info.GetUserId())
	if !exist {
		return &pb_aap.UnfollowUserResponse{Status: pb_aap.UnfollowUserResponse_USER_NOT_FOUND}, nil
	}
	exist, friend := a.findUserById(info.GetFollowingId())
	if !exist {
		return &pb_aap.UnfollowUserResponse{Status: pb_aap.UnfollowUserResponse_USER_NOT_FOUND}, nil
	}

	var user types.User
	a.db.Preload("Followings").First(&user, info.GetUserId())
	currentlyFollowing := false
	for _, following := range user.Followings {
		if following.ID == uint(info.GetFollowingId()) {
			currentlyFollowing = true
			break
		}
	}
	if !currentlyFollowing {
		return &pb_aap.UnfollowUserResponse{Status: pb_aap.UnfollowUserResponse_NOT_FOLLOWED}, nil
	}

	err := a.db.Model(&user).Association("Followings").Delete(&friend)
	if err != nil {
		return nil, err
	}

	// Update cache
	followingsKey := fmt.Sprintf("followings:%d", info.GetUserId())
	keyExist := (a.redisClient.Exists(context.Background(), followingsKey).Val() == 1)
	if keyExist {
		a.redisClient.LRem(context.Background(), followingsKey, 0, info.GetFollowingId())
	}
	followersKey := fmt.Sprintf("followers:%d", info.GetFollowingId())
	keyExist = (a.redisClient.Exists(context.Background(), followersKey).Val() == 1)
	if keyExist {
		a.redisClient.LRem(context.Background(), followersKey, 0, info.GetUserId())
	}

	return &pb_aap.UnfollowUserResponse{
		Status: pb_aap.UnfollowUserResponse_OK,
	}, nil
}

func (a *AuthenticateAndPostService) GetUserPosts(ctx context.Context, info *pb_aap.GetUserPostsRequest) (*pb_aap.GetUserPostsResponse, error) {
	exist, _ := a.findUserById(info.GetUserId())
	if !exist {
		return &pb_aap.GetUserPostsResponse{Status: pb_aap.GetUserPostsResponse_USER_NOT_FOUND}, nil
	}

	var user types.User
	a.db.Preload("Posts").First(&user, info.GetUserId())

	// Return
	var posts_ids []int64
	for _, post := range user.Posts {
		posts_ids = append(posts_ids, int64(post.ID))
	}

	return &pb_aap.GetUserPostsResponse{
		Status:   pb_aap.GetUserPostsResponse_OK,
		PostsIds: posts_ids,
	}, nil
}
