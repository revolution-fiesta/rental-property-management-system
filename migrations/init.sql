CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE rooms (
    id SERIAL PRIMARY KEY,
    type VARCHAR(20) NOT NULL,
    quantity INTEGER NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    is_deleted BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,                -- 对应 User.ID，使用 SERIAL 类型自增
    username VARCHAR(100) NOT NULL,        -- 对应 User.Username，假设用户名最大长度为 100
    password_hash VARCHAR(255) NOT NULL,   -- 对应 User.PasswordHash，密码哈希值通常较长，这里假设最大长度为 255
    email VARCHAR(100) UNIQUE NOT NULL,    -- 对应 User.Email，邮箱地址最大长度为 100，并且设置唯一约束
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP  -- 假设添加一个创建时间字段，自动设置当前时间戳
);

CREATE TABLE passwords (
    id SERIAL PRIMARY KEY,
    room_id INTEGER REFERENCES rooms(id) ON DELETE CASCADE,
    password VARCHAR(20) NOT NULL,
    is_temp BOOLEAN DEFAULT true,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);