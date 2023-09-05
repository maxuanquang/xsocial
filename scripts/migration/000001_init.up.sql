-- Use the database
USE engineerpro;

-- Create the user table
CREATE TABLE IF NOT EXISTS `user` (
    id BIGINT AUTO_INCREMENT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    hashed_password VARCHAR(1000) NOT NULL,
    salt VARBINARY(1000) NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    date_of_birth TIMESTAMP,
    email VARCHAR(100) NOT NULL,
    user_name VARCHAR(50) UNIQUE NOT NULL,
    profile_picture VARCHAR(500),
    cover_picture VARCHAR(500),
    PRIMARY KEY (id),
    INDEX idx_username (user_name)
);

-- Create following table
CREATE TABLE IF NOT EXISTS `following` (
    user_id BIGINT NOT NULL,
    follower_id BIGINT NOT NULL,
    PRIMARY KEY (user_id, follower_id),
    FOREIGN KEY (user_id) REFERENCES `user`(id),
    FOREIGN KEY (follower_id) REFERENCES `user`(id)
);

-- Create the post table
CREATE TABLE IF NOT EXISTS `post` (
    id BIGINT AUTO_INCREMENT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    user_id BIGINT NOT NULL,
    content_text TEXT(100000) NOT NULL,
    content_image_path VARCHAR(1000),
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES `user`(id)
);

-- Create the comment table
CREATE TABLE IF NOT EXISTS `comment` (
    id BIGINT AUTO_INCREMENT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    post_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    content_text TEXT(100000) NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (post_id) REFERENCES `post`(id),
    FOREIGN KEY (user_id) REFERENCES `user`(id)
);

-- Create the like table
CREATE TABLE IF NOT EXISTS `like` (
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    post_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    PRIMARY KEY (post_id, user_id),
    FOREIGN KEY (post_id) REFERENCES `post`(id),
    FOREIGN KEY (user_id) REFERENCES `user`(id)
);