package types

import (
	"time"

	"gorm.io/gorm"
)

// CREATE TABLE IF NOT EXISTS `user` (
//     id BIGINT AUTO_INCREMENT,
//     hashed_password VARCHAR(1000) NOT NULL,
//     salt VARBINARY(1000) NOT NULL,
//     first_name VARCHAR(50) NOT NULL,
//     last_name VARCHAR(50) NOT NULL,
//     dob DATE NOT NULL,
//     email VARCHAR(100) NOT NULL,
//     user_name VARCHAR(50) UNIQUE NOT NULL,
//     PRIMARY KEY (id),
//     INDEX idx_username (user_name)
// );

type User struct {
	ID             int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement"`
	HashedPassword string    `gorm:"column:hashed_password;type:varchar(1000);not null"`
	Salt           []byte    `gorm:"column:salt;type:varbinary(1000);not null"`
	FirstName      string    `gorm:"column:first_name;type:varchar(50);not null"`
	LastName       string    `gorm:"column:last_name;type:varchar(50);not null"`
	DOB            time.Time `gorm:"column:dob;type:date;not null"`
	Email          string    `gorm:"column:email;type:varchar(100);unique;not null"`
	UserName       string    `gorm:"column:user_name;type:varchar(50);unique;not null"`

	Following []*User `gorm:"many2many:following;foreignKey:id;joinForeignKey:user_id;References:id;joinReferences:follower_id"`
	Follower  []*User `gorm:"many2many:following;foreignKey:id;joinForeignKey:follower_id;References:id;joinReferences:user_id"`
}

func (User) TableName() string {
	return "user"
}

// CREATE TABLE IF NOT EXISTS `following` (
//     user_id BIGINT NOT NULL,
//     follower_id BIGINT NOT NULL,
//     PRIMARY KEY (user_id, follower_id),
//     FOREIGN KEY (user_id) REFERENCES `user`(id),
//     FOREIGN KEY (follower_id) REFERENCES `user`(id)
// );

type Following struct {
	UserID     int64 `gorm:"column:user_id;type:bigint;primaryKey"`
	FollowerID int64 `gorm:"column:follower_id;type:bigint;primaryKey"`

	User     User `gorm:"foreignKey:user_id;references:id"`
	Follower User `gorm:"foreignKey:follower_id;references:id"`
}

func (Following) TableName() string {
	return "following"
}

// CREATE TABLE IF NOT EXISTS `post` (
//     id BIGINT AUTO_INCREMENT,
//     created_at TIMESTAMP NOT NULL,
//     updated_at TIMESTAMP NOT NULL,
//     deleted_at TIMESTAMP NOT NULL,
//     user_id BIGINT NOT NULL,
//     content_text TEXT(100000) NOT NULL,
//     content_image_path VARCHAR(1000),
//     `visible` BOOLEAN NOT NULL,
//     PRIMARY KEY (id),
//     FOREIGN KEY (user_id) REFERENCES `user`(id)
// );

type Post struct {
	gorm.Model
	UserID           int64  `gorm:"column:user_id;type:bigint;not null"`
	ContentText      string `gorm:"column:content_text;type:text(100000);not null"`
	ContentImagePath string `gorm:"column:content_image_path;type:text(1000)"`
	Visible          bool   `gorm:"column:visible;type:boolean;not null"`

	User User `gorm:"foreignKey:user_id;references:id"`
}

func (Post) TableName() string {
	return "post"
}

// CREATE TABLE IF NOT EXISTS `comment` (
//     id BIGINT AUTO_INCREMENT,
//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
//     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
//     deleted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
//     post_id BIGINT NOT NULL,
//     user_id BIGINT NOT NULL,
//     content TEXT(100000) NOT NULL,
//     PRIMARY KEY (id),
//     FOREIGN KEY (post_id) REFERENCES `post`(id),
//     FOREIGN KEY (user_id) REFERENCES `user`(id)
// );

type Comment struct {
	gorm.Model
	PostID  int64  `gorm:"column:post_id;type:bigint;not null"`
	UserID  int64  `gorm:"column:user_id;type:bigint;not null"`
	Content string `gorm:"column:content;type:text(100000);not null"`

	Post Post `gorm:"foreignKey:post_id;references:id"`
	User User `gorm:"foreignKey:user_id;references:id"`
}

func (Comment) TableName() string {
	return "comment"
}

// CREATE TABLE IF NOT EXISTS `like` (
//     post_id BIGINT NOT NULL,
//     user_id BIGINT NOT NULL,
//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
//     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
//     deleted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
//     PRIMARY KEY (post_id, user_id),
//     FOREIGN KEY (post_id) REFERENCES `post`(id),
//     FOREIGN KEY (user_id) REFERENCES `user`(id)
// );

type Like struct {
	PostID    int64     `gorm:"column:post_id;type:bigint;primaryKey"`
	UserID    int64     `gorm:"column:user_id;type:bigint;primaryKey"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;not null;default:current_timestamp"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;not null;default:current_timestamp"`
	DeletedAt time.Time `gorm:"column:deleted_at;type:timestamp;not null;default:current_timestamp"`

	Post Post `gorm:"foreignKey:post_id;references:id"`
	User User `gorm:"foreignKey:user_id;references:id"`
}

func (Like) TableName() string {
	return "like"
}
