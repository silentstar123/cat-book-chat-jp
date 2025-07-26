package response

import "time"

type MessageResponse struct {
	ID          int32     `json:"id" gorm:"primarykey"`
	FromAccount string    `json:"FromAccount" gorm:"index"`
	ToAccount   string    `json:"ToAccount" gorm:"index"`
	Content     string    `json:"content" gorm:"type:varchar(2500)"`
	ContentType int16     `json:"contentType" gorm:"comment:'消息内容类型：1文字，2语音，3视频'"`
	CreatedAt   time.Time `json:"createAt"`
	// FromUsername string    `json:"fromUsername"`
	// ToUsername   string    `json:"toUsername"`
	// Avatar       string    `json:"avatar"`
	// Url          string    `json:"url"`
}

type ChatSummary struct {
	Account       string    `json:"account"`
	Nickname      string    `json:"nickname"`
	Avatar        string    `json:"avatar"`
	LatestTime    time.Time `json:"latestTime"`
	LatestContent string    `json:"latestContent"`
	UnreadCount   int64     `json:"unreadCount"`
}
