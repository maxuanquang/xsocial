package types

type RedisUser struct {
	ID             int64  `json:"id"`
	HashedPassword string `json:"hashed_password"`
	Salt           []byte `json:"salt"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	DOB            int64  `json:"dob"`
	Email          string `json:"email"`
	UserName       string `json:"user_name"`
}

type RedisFollowing struct {
	UserID     int64 `json:"user_id"`
	FollowerID int64 `json:"follower_id"`
}

type RedisPost struct {
	ID               int64  `json:"id"`
	CreatedAt        int64  `json:"created_at"`
	UpdatedAt        int64  `json:"updated_at"`
	DeletedAt        int64  `json:"deleted_at"`
	UserID           int64  `json:"user_id"`
	ContentText      string `json:"content_text"`
	ContentImagePath string `json:"content_image_path"`
	Visible          bool   `json:"visible"`
	CommentsIds      string `json:"comments_ids"`
	LikedUsersIds    string `json:"liked_users_ids"`
}

type RedisComment struct {
	ID          int64  `json:"id"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
	DeletedAt   int64  `json:"deleted_at"`
	PostID      int64  `json:"post_id"`
	UserID      int64  `json:"user_id"`
	ContentText string `json:"content_text"`
}

type RedisLike struct {
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
	DeletedAt int64 `json:"deleted_at"`
	PostID    int64 `json:"post_id"`
	UserID    int64 `json:"user_id"`
}
