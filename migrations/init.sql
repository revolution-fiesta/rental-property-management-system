CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- 创建 rooms 表
CREATE TABLE rooms (
    id SERIAL PRIMARY KEY,                         -- 主键，自动递增
    type VARCHAR(20) NOT NULL,                     -- 房间类型
    quantity INTEGER NOT NULL,                     -- 房间数量
    price DECIMAL(10, 2) NOT NULL,                 -- 房间价格
    is_deleted BOOLEAN DEFAULT false,              -- 是否已租，默认为 false
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 创建时间，默认为当前时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 更新时间，默认为当前时间
    tags VARCHAR(255),                             -- 房间标签，如方向、是否精修、是否近地铁等
    area DECIMAL(10, 2) NOT NULL                   -- 房间占地面积(单位平方米)
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