package types

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	HashedPassword string    `gorm:"size:1000;not null" json:"hashed_password"`
	Salt           []byte    `gorm:"size:1000;not null" json:"salt"`
	FirstName      string    `gorm:"size:50;not null" json:"first_name"`
	LastName       string    `gorm:"size:50;not null" json:"last_name"`
	DateOfBirth    time.Time `gorm:"not null" json:"dob"`
	Email          string    `gorm:"size:50;not null" json:"email"`
	UserName       string    `gorm:"size:50;not null" json:"user_name"`
	Posts          []*Post
	Followers      []*User `gorm:"many2many:following"`
	Followings     []*User `gorm:"many2many:following;joinForeignKey:follower_id;joinReferences:user_id"`
}

func (User) TableName() string {
	return "user"
}

type Post struct {
	gorm.Model
	ContentText      string `gorm:"size:100000" json:"content_text"`
	ContentImagePath string `gorm:"size:1000" json:"content_image_path"`
	UserID           int64  `gorm:"not null" json:"user_id"`
	Comments         []*Comment
	LikedUsers       []*User `gorm:"many2many:like"`
}

func (Post) TableName() string {
	return "post"
}

type Comment struct {
	gorm.Model
	ContentText string `gorm:"size:100000;not null" json:"content_text"`
	PostID      int64  `gorm:"not null" json:"post_id"`
	UserID      int64  `gorm:"not null" json:"user_id"`
}

func (Comment) TableName() string {
	return "comment"
}

type Like struct {
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	DeletedAt time.Time
	PostId    int64 `gorm:"not null"`
	UserId    int64 `gorm:"not null"`
}
