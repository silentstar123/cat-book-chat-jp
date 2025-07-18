package v1

import (
	"net/http"

	"chat-room/internal/service"
	"chat-room/pkg/common/request"
	"chat-room/pkg/common/response"
	"chat-room/pkg/global/log"

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
