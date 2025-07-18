package request

type MessageRequest struct {
	MessageType int32  `form:"MessageType"`
	Account     string `form:"Account"`
	ToAccount   string `form:"ToAccount"`
}
