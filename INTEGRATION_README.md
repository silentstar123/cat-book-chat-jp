# チャットサーバーとメインサーバーの統合ガイド

## 📋 概要

このドキュメントでは、チャットサーバー（catchat-main）とメインサーバー（catcal）の統合方法について説明します。

## 🏗️ アーキテクチャ

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   フロントエンド   │    │   メインサーバー   │    │   チャットサーバー   │
│  (商家端/個人端)  │◄──►│   (catcal)     │◄──►│  (catchat-main) │
│  ポート: 3000    │    │  ポート: 8082    │    │  ポート: 8888    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🔧 統合のポイント

### 1. ユーザー管理の統合
- **メインサーバー**: ユーザー登録・ログイン・認証を担当
- **チャットサーバー**: メインサーバーからユーザー情報を取得して使用

### 2. データベースの統合
- **メインサーバー**: PostgreSQL（catcalデータベース）
- **チャットサーバー**: PostgreSQL（catcalデータベースを共有）

### 3. API連携
- チャットサーバーはメインサーバーのAPIを呼び出してユーザー情報を取得
- フロントエンドは両方のサーバーと通信

## 🚀 起動手順

### 1. メインサーバーの起動
```bash
cd catcal
go run main.go
# または
./catcal
```

### 2. チャットサーバーの起動
```bash
cd catchat-main
export SKIP_MYSQL=true
./chat-server
```

### 3. フロントエンドの起動
```bash
# 商家端
cd cat-book-merchant-shell-master-5
npm run dev

# 個人端
cd cat-book-client-shell-master-9
npm run dev
```

## 📡 API設定

### フロントエンド設定

#### 商家端 (`cat-book-merchant-shell-master-5/src/apis/request.js`)
```javascript
export const BASEURL = "http://localhost:8082";        // メインサーバー
export const SOCKETURL = "ws://localhost:8888";        // WebSocket
export const CHAT_SERVER_URL = "http://localhost:8888"; // チャットサーバー
```

#### 個人端 (`cat-book-client-shell-master-9/src/apis/request.js`)
```javascript
export const BASEURL = "http://localhost:8082";        // メインサーバー
export const SOCKETURL = "ws://localhost:8888";        // WebSocket
export const CHAT_SERVER_URL = "http://localhost:8888"; // チャットサーバー
```

### チャットAPI

#### セッションリスト取得
```javascript
chatApi.getConversations({
    account: userAccount,
    page: 1,
    pageSize: 20
})
```

#### メッセージ取得
```javascript
chatApi.getMessages(sessionId, {
    Account: userAccount,
    ToAccount: sessionId,
    MessageType: 1
})
```

#### メッセージ送信
```javascript
chatApi.sendMessage({
    from: userAccount,
    to: targetAccount,
    content: messageContent,
    contentType: 1,
    messageType: 1
})
```

## 🔄 データフロー

### 1. ユーザーログイン
```
フロントエンド → メインサーバー → 認証成功 → チャットサーバーでユーザー情報取得
```

### 2. メッセージ送信
```
フロントエンド → チャットサーバー → データベース保存 → WebSocket配信
```

### 3. メッセージ受信
```
WebSocket → フロントエンド → UI更新
```

## 🧪 テスト

### 統合テストの実行
```bash
cd catchat-main
./test_integration.sh
```

### 手動テスト
1. メインサーバーでユーザー登録
2. チャットサーバーでメッセージ送信
3. フロントエンドでチャット機能確認

## 🐛 トラブルシューティング

### よくある問題

#### 1. チャットサーバーがメインサーバーに接続できない
- メインサーバーが起動しているか確認
- ポート8082が開いているか確認
- ファイアウォール設定を確認

#### 2. フロントエンドがチャットサーバーに接続できない
- チャットサーバーが起動しているか確認
- ポート8888が開いているか確認
- CORS設定を確認

#### 3. ユーザー情報が取得できない
- メインサーバーのAPIが正常に動作しているか確認
- データベース接続を確認
- ログを確認

### ログ確認
```bash
# メインサーバーログ
tail -f catcal/logs/app.log

# チャットサーバーログ
tail -f catchat-main/logs/app.log
```

## 📝 設定ファイル

### チャットサーバー設定 (`catchat-main/config.toml`)
```toml
[postgres]
host = "localhost"
port = 5432
user = "catcal"
password = "catcal123456"
dbname = "catcal"

[server]
port = 8888
```

### メインサーバー設定 (`catcal/config.dev.yaml`)
```yaml
database:
  host: localhost
  port: 5432
  user: catcal
  password: catcal123456
  name: catcal
```

## 🔒 セキュリティ

### 認証
- JWTトークンを使用
- メインサーバーでトークン生成
- チャットサーバーでトークン検証

### データ保護
- HTTPS通信の使用（本番環境）
- データベース接続の暗号化
- 入力値の検証

## 📈 パフォーマンス

### 最適化ポイント
- WebSocket接続の再利用
- メッセージのバッチ処理
- データベースクエリの最適化

### 監視
- 接続数の監視
- メッセージ処理時間の監視
- エラー率の監視

## 🚀 デプロイ

### 本番環境設定
1. 環境変数の設定
2. データベースの移行
3. SSL証明書の設定
4. ロードバランサーの設定

### Docker化
```dockerfile
# チャットサーバー
FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build -o chat-server cmd/main.go
EXPOSE 8888
CMD ["./chat-server"]
```

## 📞 サポート

問題が発生した場合は、以下の手順で調査してください：

1. ログファイルの確認
2. ネットワーク接続の確認
3. データベース接続の確認
4. APIエンドポイントの確認

## 🔄 更新履歴

- 2025-01-25: 初期統合実装
- 2025-01-25: フロントエンド統合
- 2025-01-25: テストスクリプト追加 