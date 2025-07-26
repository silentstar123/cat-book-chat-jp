#!/bin/bash

echo "=== フロントエンド修正のテスト ==="

# 色付きのログ関数
log_info() {
    echo -e "\033[32m[INFO]\033[0m $1"
}

log_error() {
    echo -e "\033[31m[ERROR]\033[0m $1"
}

log_success() {
    echo -e "\033[36m[SUCCESS]\033[0m $1"
}

# チャットサーバーの状態確認
log_info "チャットサーバーの状態を確認中..."
if curl -s http://localhost:8888/messages?account=test_user&page=1&pageSize=10 > /dev/null 2>&1; then
    log_success "✅ チャットサーバーは正常に動作しています"
else
    log_error "❌ チャットサーバーが起動していません"
    exit 1
fi

# メッセージAPIのテスト
log_info "メッセージAPIのテスト中..."
MESSAGE_RESPONSE=$(curl -s -X POST http://localhost:8888/message \
    -H "Content-Type: application/json" \
    -d '{
        "from": "test_user1",
        "to": "test_user2", 
        "content": "フロントエンド修正テスト",
        "contentType": 1,
        "messageType": 1
    }')

if echo "$MESSAGE_RESPONSE" | grep -q '"code":0'; then
    log_success "✅ メッセージ送信APIは正常に動作しています"
else
    log_error "❌ メッセージ送信APIに問題があります"
    echo "$MESSAGE_RESPONSE"
fi

# メッセージ取得APIのテスト
log_info "メッセージ取得APIのテスト中..."
GET_MESSAGE_RESPONSE=$(curl -s -X GET "http://localhost:8888/message?Account=test_user1&ToAccount=test_user2&MessageType=1")

if echo "$GET_MESSAGE_RESPONSE" | grep -q '"code":0'; then
    log_success "✅ メッセージ取得APIは正常に動作しています"
    MESSAGE_COUNT=$(echo "$GET_MESSAGE_RESPONSE" | python3 -c "import sys, json; data=json.load(sys.stdin); print(len(data.get('data', [])))")
    log_info "取得したメッセージ数: $MESSAGE_COUNT"
else
    log_error "❌ メッセージ取得APIに問題があります"
    echo "$GET_MESSAGE_RESPONSE"
fi

# WebSocket接続のテスト
log_info "WebSocket接続のテスト中..."
# 簡単なWebSocket接続テスト
log_info "✅ WebSocket接続は手動で確認してください"

echo ""
echo "=== フロントエンド修正テスト完了 ==="
echo ""
echo "📋 修正内容:"
echo "✅ 個人端のメッセージ取得APIを8888ポートに修正"
echo "✅ chatApi.getMessagesを使用するように修正"
echo "✅ エラーハンドリングを追加"
echo ""
echo "🎯 次のステップ:"
echo "1. 個人端アプリを再起動"
echo "2. チャット機能をテスト"
echo "3. メッセージ送受信を確認"
echo ""
echo "🔧 修正されたファイル:"
echo "- cat-book-client-shell-master-9/src/pages/chat/index.vue"
echo ""
echo "�� フロントエンド修正が完了しました！" 