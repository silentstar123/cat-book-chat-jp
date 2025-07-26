package v1

import (
	"net/http"

	"chat-room/internal/service"
	"chat-room/pkg/common/request"
	"chat-room/pkg/common/response"
	"chat-room/pkg/global/log"
	"chat-room/pkg/protocol"

	"github.com/gin-gonic/gin"
)

// 获取消息列表
func GetMessage(c *gin.Context) {
	log.Logger.Info(c.Query("account"))
	var messageRequest request.MessageRequest
	err := c.BindQuery(&messageRequest)
	println("Account:", messageRequest.Account)
	println("ToAccount:", messageRequest.ToAccount)
	println("MessageType:", messageRequest.MessageType)
	if nil != err {
		log.Logger.Error("bindQueryError", log.Any("bindQueryError", err))
	}
	log.Logger.Info("messageRequest params: ", log.Any("messageRequest", messageRequest))

	messages, err := service.MessageService.GetMessages(messageRequest)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessMsg(messages))
}

func GetMessageList(c *gin.Context) {

	account := c.Query("account")

	// 调用服务层获取聊天会话
	summaries, err := service.MessageService.GetMessageList(account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": summaries})
}

// 发送消息
func SendMessage(c *gin.Context) {
	var messageRequest struct {
		From        string `json:"from" binding:"required"`
		To          string `json:"to" binding:"required"`
		Content     string `json:"content" binding:"required"`
		ContentType int16  `json:"contentType" default:"1"`
		MessageType int16  `json:"messageType" default:"1"`
	}

	if err := c.ShouldBindJSON(&messageRequest); err != nil {
		c.JSON(http.StatusOK, response.FailMsg("参数错误: "+err.Error()))
		return
	}

	// 创建消息对象
	message := protocol.Message{
		From:        messageRequest.From,
		To:          messageRequest.To,
		Content:     messageRequest.Content,
		ContentType: int32(messageRequest.ContentType),
		MessageType: int32(messageRequest.MessageType),
	}

	// 保存消息到数据库
	err := service.MessageService.SaveMessage(message)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg("发送消息失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessMsg("消息发送成功"))
}
