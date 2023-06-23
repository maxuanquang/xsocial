CREATE TABLE IF NOT EXISTS `user`  (
    id BIGINT AUTO_INCREMENT,
    hashed_password VARCHAR(1000) NOT NULL,
    salt VARBINARY(1000) NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    dob DATE NOT NULL,
    email VARCHAR(100) NOT NULL,
    user_name VARCHAR(50) UNIQUE NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS `post` (
    id BIGINT AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    content_text TEXT(100000) NOT NULL,
    content_image_path VARCHAR(1000),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    `visible` BOOLEAN NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES `user`(id)
);

CREATE TABLE IF NOT EXISTS `comment` (
    id BIGINT AUTO_INCREMENT,
    post_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    content TEXT(100000) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (post_id) REFERENCES `post`(id),
    FOREIGN KEY (user_id) REFERENCES `user`(id)
);

CREATE TABLE IF NOT EXISTS `like` (
    post_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES `post`(id),
    FOREIGN KEY (user_id) REFERENCES `user`(id),
    CONSTRAINT unique_post_id_user_id UNIQUE (post_id, user_id)
);

CREATE TABLE IF NOT EXISTS `following` (
    user_id BIGINT NOT NULL,
    follower_id BIGINT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES `user`(id),
    FOREIGN KEY (follower_id) REFERENCES `user`(id),
    CONSTRAINT unique_user_id_follower_id UNIQUE (user_id, follower_id)
);