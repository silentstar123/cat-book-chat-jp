package v1

import (
	"net/http"

	"chat-room/internal/model"
	"chat-room/internal/service"
	"chat-room/pkg/common/request"
	"chat-room/pkg/common/response"
	"chat-room/pkg/global/log"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var user model.User
	// c.BindJSON(&user)
	c.ShouldBindJSON(&user)
	log.Logger.Debug("user", log.Any("user", user))

	if service.UserService.Login(&user) {
		c.JSON(http.StatusOK, response.SuccessMsg(user))
		return
	}

	c.JSON(http.StatusOK, response.FailMsg("Login failed"))
}

func Register(c *gin.Context) {
	var user model.User
	c.ShouldBindJSON(&user)
	err := service.UserService.Register(&user)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessMsg(user))
}

func ModifyUserInfo(c *gin.Context) {
	var user model.User
	c.ShouldBindJSON(&user)
	log.Logger.Debug("user", log.Any("user", user))
	if err := service.UserService.ModifyUserInfo(&user); err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessMsg(nil))
}

func GetUserDetails(c *gin.Context) {
	uuid := c.Param("uuid")

	user, err := service.UserService.GetUserDetails(uuid)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessMsg(user))
}

// 通过用户名获取用户信息
func GetUserOrGroupByName(c *gin.Context) {
	name := c.Query("name")

	users, err := service.UserService.GetUserOrGroupByName(name)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessMsg(users))
}

func GetUserList(c *gin.Context) {
	users, err := service.UserService.GetUserList()
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessMsg(users))
}

func AddFriend(c *gin.Context) {
	var userFriendRequest request.FriendRequest
	c.ShouldBindJSON(&userFriendRequest)

	// 创建用户对象用于添加好友
	user := &model.User{
		Account: userFriendRequest.FriendAccount,
	}

	err := service.UserService.AddFriend(user)
	if nil != err {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessMsg(nil))
}
