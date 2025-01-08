-- Создание таблицы users
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password VARCHAR(100) NOT NULL,
    avatar VARCHAR(40) DEFAULT '1.jpg',
    bio VARCHAR(1024),
    privilege_level INT DEFAULT 0,
    blocked BOOLEAN
);

CREATE TABLE IF NOT EXISTS user_ips (
    user_id INT,
    ip_address VARCHAR(45) NOT NULL,
    last_login TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Создание таблицы subscription_plans
CREATE TABLE IF NOT EXISTS subscription_plans (
    id SERIAL PRIMARY KEY,
    user_id INT,
    name VARCHAR(50) NOT NULL,
    price NUMERIC(10,2) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Создание таблицы user_subscriptions
CREATE TABLE IF NOT EXISTS user_subscriptions (
    user_id INT,
    subscription_plan_id INT,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (subscription_plan_id) REFERENCES subscription_plans(id)
);

-- Создание таблицы user_followers
CREATE TABLE IF NOT EXISTS user_followers (
    follower_id INT,
    followee_id INT,
    FOREIGN KEY (follower_id) REFERENCES users(id),
    FOREIGN KEY (followee_id) REFERENCES users(id)
);

-- Создание таблицы posts
CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    description VARCHAR(2048),
    is_private BOOLEAN DEFAULT FALSE,
    is_nsfw BOOLEAN DEFAULT FALSE,
    reg_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- Создание таблицы user_posts
CREATE TABLE IF NOT EXISTS user_posts (
    user_id INT,
    post_id INT,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (post_id) REFERENCES posts(id)
);

-- Создание таблицы plan_posts
CREATE TABLE IF NOT EXISTS plan_posts (
    subscription_plan_id INT,
    post_id INT,
    FOREIGN KEY (subscription_plan_id) REFERENCES subscription_plans(id),
    FOREIGN KEY (post_id) REFERENCES posts(id)
);

-- Создание таблицы post_images
CREATE TABLE IF NOT EXISTS post_images (
    id SERIAL PRIMARY KEY,
    post_id INT,
    image_url VARCHAR(255) NOT NULL,
    FOREIGN KEY (post_id) REFERENCES posts(id)
);

-- Создание таблицы tags
CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

-- Создание таблицы post_tags
CREATE TABLE IF NOT EXISTS post_tags (
    post_id INT,
    tag_id INT,
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (tag_id) REFERENCES tags(id)
);

-- Создание таблицы user_tags
CREATE TABLE IF NOT EXISTS user_tags (
    user_id INT,
    tag_id INT,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (tag_id) REFERENCES tags(id)
);

-- Создание таблицы post_likes
CREATE TABLE IF NOT EXISTS post_likes (
    post_id INT,
    user_id INT,
    PRIMARY KEY (post_id, user_id),
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Создание таблицы post_stats
CREATE TABLE IF NOT EXISTS post_stats (
    post_id INT,
    likes INT DEFAULT 0,
    views INT DEFAULT 0,
    FOREIGN KEY (post_id) REFERENCES posts(id)
);

-- Создание таблицы ratios
CREATE TABLE IF NOT EXISTS ratios (
    post_id INT,
    ratio INT DEFAULT 10000,
    FOREIGN KEY (post_id) REFERENCES posts(id)
);

-- Создание таблицы comments
CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    content VARCHAR(2048) NOT NULL,
    reg_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы post_comments
CREATE TABLE IF NOT EXISTS post_comments (
    post_id INT,
    comment_id INT,
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (comment_id) REFERENCES comments(id)
);

-- Создание таблицы user_comments
CREATE TABLE IF NOT EXISTS user_comments (
    user_id INT,
    comment_id INT,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (comment_id) REFERENCES comments(id)
);
