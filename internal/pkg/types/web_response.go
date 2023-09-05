package types

// import "time"

// MessageResponse return message in response.
type MessageResponse struct {
	Message string `json:"message"`
}

// PostDetailInfoResponse return post detail in response.
type PostDetailInfoResponse struct {
	PostID           int64             `json:"post_id"`
	UserID           int64             `json:"user_id"`
	ContentText      string            `json:"content_text"`
	ContentImagePath []string          `json:"content_image_path"`
	CreatedAt        string            `json:"created_at"`
	Comments         []CommentResponse `json:"comments"`
	UsersLiked       []int64           `json:"users_liked"`
}

type CommentResponse struct {
	CommentId   int64  `json:"comment_id"`
	UserId      int64  `json:"user_id"`
	PostId      int64  `json:"post_id"`
	ContentText string `json:"content_text"`
}

type UserFollowerResponse struct {
	FollowersIds []int64 `json:"followers_ids"`
}

type UserFollowingResponse struct {
	FollowingsIds []int64 `json:"followings_ids"`
}

type UserPostsResponse struct {
	PostsIds []int64 `json:"posts_ids"`
}

type NewsfeedResponse struct {
	PostsIds []int64 `json:"posts_ids"`
}

type LoginResponse struct {
	Message string         `json:"message"`
	User    UserDetailInfo `json:"user"`
}

type UserDetailInfo struct {
	UserID         int64  `json:"user_id"`
	UserName       string `json:"user_name"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	DateOfBirth    string `json:"date_of_birth"`
	Email          string `json:"email"`
	ProfilePicture string `json:"profile_picture"`
	CoverPicture   string `json:"cover_picture"`
}
