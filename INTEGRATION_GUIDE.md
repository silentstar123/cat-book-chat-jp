# 猫书项目整合指南

## 项目架构

### 后端服务架构
```
┌─────────────────┐    ┌─────────────────┐
│   catcal        │    │  catchat-main   │
│   (主业务服务)   │    │   (聊天服务)    │
│   端口: 8082    │    │   端口: 8888    │
└─────────────────┘    └─────────────────┘
         │                       │
         └───────────┬───────────┘
                     │
         ┌─────────────────┐
         │   PostgreSQL    │
         │   数据库        │
         └─────────────────┘
```

## 服务功能分工

### catcal (主业务服务) - 端口 8082
- **用户管理**: 注册、登录、个人信息
- **宠物管理**: 宠物信息、品种管理
- **社交功能**: 帖子、评论、点赞、收藏
- **订单管理**: 服务订单、支付处理
- **商家功能**: 店铺管理、商品管理
- **支付系统**: 多种支付方式集成
- **文件上传**: 图片、视频上传

### catchat-main (聊天服务) - 端口 8888
- **实时聊天**: WebSocket长连接
- **消息管理**: 文本、图片、文件、语音、视频消息
- **群组管理**: 群聊、群成员管理
- **好友系统**: 好友添加、好友列表
- **在线状态**: 用户在线状态管理
- **消息同步**: 与catcal用户数据同步

## 数据同步机制

### 用户数据同步
```go
// catchat-main 从 catcal 同步用户信息
func (s *UserSyncService) SyncUserFromCatcal(userID int) (*model.User, error) {
    url := fmt.Sprintf("%s/api/v1/users/%d", s.catcalAPIURL, userID)
    // 获取用户信息并转换为聊天系统格式
}
```

### 数据库共享
- 两个服务使用相同的PostgreSQL数据库
- catcal: 业务数据表
- catchat-main: 聊天相关表

## 启动方式

### 1. 分别启动
```bash
# 启动catcal主业务服务
cd catcal
./catcal

# 启动catchat-main聊天服务
cd catchat-main
./dev.sh
```

### 2. 整合启动 (推荐)
```bash
cd catchat-main
./start_integrated.sh
```

## 配置说明

### catcal配置 (config.yaml)
```yaml
database:
  host: localhost
  port: 5432
  user: catcal
  password: catcal123456
  db: catcal

server:
  port: 8082
```

### catchat-main配置 (config.toml)
```toml
[postgres]
host = "localhost"
port = 5432
user = "catcal"
password = "catcal123456"
dbname = "catcal"

[catcal]
api_url = "http://localhost:8082"
api_version = "v1"

[server]
port = 8888
host = "0.0.0.0"
```

## API接口

### catcal API (端口 8082)
- `GET /api/v1/users/{id}` - 获取用户信息
- `POST /api/v1/users` - 创建用户
- `PUT /api/v1/users/{id}` - 更新用户信息
- `GET /api/v1/pets` - 获取宠物列表
- `POST /api/v1/orders` - 创建订单
- `GET /api/v1/shops` - 获取店铺列表

### catchat-main API (端口 8888)
- `GET /api/v1/user/login` - 用户登录
- `POST /api/v1/user/register` - 用户注册
- `GET /api/v1/user/{uuid}` - 获取用户详情
- `POST /api/v1/message` - 发送消息
- `GET /api/v1/message` - 获取消息列表
- `POST /api/v1/group` - 创建群组

## WebSocket连接

### 聊天WebSocket
```
ws://localhost:8888/ws
```

### 消息格式 (Protocol Buffer)
```protobuf
message Message {
    string avatar = 1;       // 头像
    string fromUsername = 2; // 发送者用户名
    string from = 3;         // 发送者UUID
    string to = 4;           // 接收者UUID
    string content = 5;      // 消息内容
    int32 contentType = 6;   // 消息类型
    string type = 7;         // 消息类型标识
    int32 messageType = 8;   // 1.单聊 2.群聊
    string url = 9;          // 文件URL
    string fileSuffix = 10;  // 文件后缀
    bytes file = 11;         // 文件二进制数据
}
```

## 开发环境

### 依赖要求
- Go 1.16+
- PostgreSQL 12+
- Air (热更新工具)

### 安装依赖
```bash
# 安装Air
curl -sSfL https://raw.githubusercontent.com/air-verse/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# 安装Go依赖
cd catcal && go mod tidy
cd ../catchat-main && go mod tidy
```

### 数据库初始化
```sql
-- 创建数据库
CREATE DATABASE catcal;

-- 运行catcal迁移
-- catcal会自动创建业务表

-- 运行catchat-main迁移
-- 执行chat.sql中的聊天相关表
```

## 部署说明

### Docker部署
```bash
# 启动PostgreSQL
docker-compose up -d mysql8-chat

# 构建并启动服务
docker-compose up -d
```

### 生产环境
- 使用Nginx反向代理
- 配置SSL证书
- 设置防火墙规则
- 配置日志轮转

## 监控和日志

### 日志位置
- catcal: `catcal/logs/`
- catchat-main: `catchat-main/logs/`

### 健康检查
- catcal: `http://localhost:8082/health`
- catchat-main: `http://localhost:8888/health`

## 故障排除

### 常见问题
1. **端口冲突**: 确保8082和8888端口未被占用
2. **数据库连接失败**: 检查PostgreSQL服务状态
3. **用户同步失败**: 检查catcal服务是否正常运行
4. **WebSocket连接失败**: 检查防火墙设置

### 调试命令
```bash
# 检查服务状态
ps aux | grep catcal
ps aux | grep air

# 查看日志
tail -f catchat-main/logs/catcal.log
tail -f catchat-main/logs/chat.log

# 检查端口占用
lsof -i :8082
lsof -i :8888
```

## 更新日志

### v1.0.0 (2025-01-29)
- 初始版本发布
- 支持用户数据同步
- 集成聊天功能
- 添加日文错误消息支持 