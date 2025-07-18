package response

type ResponseMsg struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func SuccessMsg(data interface{}) *ResponseMsg {
	msg := &ResponseMsg{
		Code: 0,
		Msg:  "SUCCESS",
		Data: data,
	}
	return msg
}

func FailMsg(msg string) *ResponseMsg {
	msgObj := &ResponseMsg{
		Code: -1,
		Msg:  msg,
	}
	return msgObj
}

func FailCodeMsg(code int, msg string) *ResponseMsg {
	msgObj := &ResponseMsg{
		Code: code,
		Msg:  msg,
	}
	return msgObj
}

// 日文错误消息
var JapaneseMessages = map[string]string{
	"login_failed":        "ログインに失敗しました",
	"register_failed":     "登録に失敗しました",
	"user_not_found":      "ユーザーが見つかりません",
	"invalid_password":    "パスワードが無効です",
	"user_exists":         "ユーザーは既に存在します",
	"friend_add_failed":   "友達追加に失敗しました",
	"message_send_failed": "メッセージ送信に失敗しました",
	"group_create_failed": "グループ作成に失敗しました",
}

func JapaneseFailMsg(key string) *ResponseMsg {
	msg, exists := JapaneseMessages[key]
	if !exists {
		msg = "エラーが発生しました"
	}
	return &ResponseMsg{
		Code: -1,
		Msg:  msg,
	}
}
