-- PostgreSQL版本的聊天表结构

-- 消息表
CREATE TABLE IF NOT EXISTS messages (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    from_account TEXT NOT NULL,
    to_account TEXT NOT NULL,
    content TEXT,
    url TEXT,
    pic TEXT,
    message_type SMALLINT DEFAULT 1, -- 1单聊，2群聊
    content_type SMALLINT DEFAULT 1, -- 1文字，2语音，3视频
    is_read SMALLINT DEFAULT 0, -- 0未读，1已读
    file BYTEA, -- 文件二进制数据
    file_suffix TEXT -- 文件后缀
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_messages_from_account ON messages(from_account);
CREATE INDEX IF NOT EXISTS idx_messages_to_account ON messages(to_account);
CREATE INDEX IF NOT EXISTS idx_messages_deleted_at ON messages(deleted_at);
CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at);

-- 好友关系表
CREATE TABLE IF NOT EXISTS user_friends (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    user_account TEXT NOT NULL,
    friend_account TEXT NOT NULL
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_user_friends_user_account ON user_friends(user_account);
CREATE INDEX IF NOT EXISTS idx_user_friends_friend_account ON user_friends(friend_account);

-- 群组表
CREATE TABLE IF NOT EXISTS groups (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    uuid TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    notice TEXT,
    owner_account TEXT NOT NULL
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_groups_uuid ON groups(uuid);
CREATE INDEX IF NOT EXISTS idx_groups_owner_account ON groups(owner_account);

-- 群组成员表
CREATE TABLE IF NOT EXISTS group_members (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    group_uuid TEXT NOT NULL,
    user_account TEXT NOT NULL,
    nickname TEXT,
    mute SMALLINT DEFAULT 0 -- 0不禁言，1禁言
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_group_members_group_uuid ON group_members(group_uuid);
CREATE INDEX IF NOT EXISTS idx_group_members_user_account ON group_members(user_account);

-- 群组消息读取状态表
CREATE TABLE IF NOT EXISTS group_message_read_status (
    group_uuid TEXT NOT NULL,
    user_account TEXT NOT NULL,
    last_read_message_id BIGINT NOT NULL,
    last_read_time TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (group_uuid, user_account)
); 