package models

import (
	"time"
)

// CREATE TABLE IF NOT EXISTS `user`  (
//     id BIGINT AUTO_INCREMENT,
//     hashed_password VARCHAR(1000) NOT NULL,
//     salt VARCHAR(1000) NOT NULL,
//     first_name VARCHAR(50) NOT NULL,
//     last_name VARCHAR(50) NOT NULL,
//     dob DATE NOT NULL,
//     email VARCHAR(100) NOT NULL,
//     user_name VARCHAR(50) UNIQUE NOT NULL,
//     PRIMARY KEY (id)
// );

type User struct {
	ID             int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement"`
	HashedPassword string    `gorm:"column:hashed_password;type:varchar(1000);not null"`
	Salt           string    `gorm:"column:salt;type:varchar(1000);not null"`
	FirstName      string    `gorm:"column:first_name;type:varchar(50);not null"`
	LastName       string    `gorm:"column:last_name;type:varchar(50);not null"`
	DOB            time.Time `gorm:"column:dob;type:date;not null"`
	Email          string    `gorm:"column:email;type:varchar(100);not null"`
	UserName       string    `gorm:"column:user_name;type:varchar(50);unique;not null"`
}

func (User) TableName() string {
	return "user"
}

// CREATE TABLE IF NOT EXISTS `post` (
//     id BIGINT AUTO_INCREMENT,
//     user_id BIGINT NOT NULL,
//     content_text TEXT(100000) NOT NULL,
//     content_image_path VARCHAR(1000),
//     created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
//     `visible` BOOLEAN NOT NULL,
//     PRIMARY KEY (id),
//     FOREIGN KEY (user_id) REFERENCES `user`(id)
// );

type Post struct {
	ID               int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement"`
	UserID           int64     `gorm:"column:user_id;type:bigint;not null"`
	ContentText      string    `gorm:"column:content_text;type:text(100000);not null"`
	ContentImagePath string    `gorm:"column:content_image_path;type:varchar(1000)"`
	CreatedAt        time.Time `gorm:"column:created_at;type:datetime;default:current_timestamp;not null"`
	Visible          bool      `gorm:"column:visible;type:boolean;not null"`

	User User `gorm:"foreign_key:user_id;references:id"`
}

func (Post) TableName() string {
	return "post"
}

// CREATE TABLE IF NOT EXISTS `comment` (
//     id BIGINT AUTO_INCREMENT,
//     post_id BIGINT NOT NULL,
//     user_id BIGINT NOT NULL,
//     content TEXT(100000) NOT NULL,
//     created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
//     PRIMARY KEY (id),
//     FOREIGN KEY (post_id) REFERENCES `post`(id),
//     FOREIGN KEY (user_id) REFERENCES `user`(id)
// );

type Comment struct {
	ID        int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement"`
	PostID    int64     `gorm:"column:post_id;type:bigint;not null"`
	UserID    int64     `gorm:"column:user_id;type:bigint;not null"`
	Content   string    `gorm:"column:content;type:text(100000);not null"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:current_timestamp"`

	Post Post `gorm:"foreign_key:post_id;references:id"`
	User User `gorm:"foreign_key:user_id;references:id"`
}

func (Comment) TableName() string {
	return "comment"
}

// CREATE TABLE IF NOT EXISTS `like` (
//     post_id BIGINT NOT NULL,
//     user_id BIGINT NOT NULL,
//     created_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
//     FOREIGN KEY (post_id) REFERENCES `post`(id),
//     FOREIGN KEY (user_id) REFERENCES `user`(id),
//     CONSTRAINT unique_post_id_user_id UNIQUE (post_id, user_id)
// );

type Like struct {
	PostID    int64     `gorm:"column:post_id;type:bigint;not null;index:unique_post_id_user_id,unique"`
	UserID    int64     `gorm:"column:user_id;type:bigint;not null;index:unique_post_id_user_id,unique"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:current_timestamp"`

	Post Post `gorm:"constraint:foreign_key:post_id;references:id"`
	User User `gorm:"foreign_key:user_id;references:id"`
}

func (Like) TableName() string {
	return "like"
}

// CREATE TABLE IF NOT EXISTS `user_user` (
//     user_id BIGINT NOT NULL,
//     follower_id BIGINT NOT NULL,
//     FOREIGN KEY (user_id) REFERENCES `user`(id),
//     FOREIGN KEY (follower_id) REFERENCES `user`(id),
//     CONSTRAINT unique_user_id_follower_id UNIQUE (user_id, follower_id)
// );

type UserUser struct {
	UserID     int64 `gorm:"column:user_id;type:bigint;not null;index:unique_user_id_follower_id,unique"`
	FollowerID int64 `gorm:"column:follower_id;type:bigint;not null;index:unique_user_id_follower_id,unique"`

	User     User `gorm:"foreign_key:user_id;references:id"`
	Follower User `gorm:"foreign_key:follower_id;references:id"`
}

func (UserUser) TableName() string {
	return "user_user"
}
