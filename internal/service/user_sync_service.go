package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"chat-room/config"
	"chat-room/internal/model"
	"chat-room/pkg/global/log"
)

type UserSyncService struct {
	catcalAPIURL string
}

type CatcalUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
}

func NewUserSyncService() *UserSyncService {
	return &UserSyncService{
		catcalAPIURL: config.GetConfig().Catcal.APIURL,
	}
}

// SyncUserFromCatcal 从catcal同步用户信息
func (s *UserSyncService) SyncUserFromCatcal(userID int) (*model.User, error) {
	url := fmt.Sprintf("%s/api/v1/users/%d", s.catcalAPIURL, userID)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Logger.Error("创建请求失败", log.Any("error", err))
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Logger.Error("请求catcal API失败", log.Any("error", err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Logger.Error("catcal API返回错误状态码", log.Int("status", resp.StatusCode))
		return nil, fmt.Errorf("API返回状态码: %d", resp.StatusCode)
	}

	var catcalUser CatcalUser
	if err := json.NewDecoder(resp.Body).Decode(&catcalUser); err != nil {
		log.Logger.Error("解析catcal用户数据失败", log.Any("error", err))
		return nil, err
	}

	// 转换为chat-room的用户模型
	user := &model.User{
		Id:       int32(catcalUser.ID),
		Account:  fmt.Sprintf("%d", catcalUser.ID),
		Username: catcalUser.Username,
		Nickname: catcalUser.Nickname,
		Avatar:   catcalUser.Avatar,
		Email:    catcalUser.Email,
		CreateAt: time.Now(),
	}

	return user, nil
}

// ValidateUserInCatcal 验证用户在catcal中是否存在
func (s *UserSyncService) ValidateUserInCatcal(userID int) (bool, error) {
	url := fmt.Sprintf("%s/api/v1/users/%d", s.catcalAPIURL, userID)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}

// SyncUserToChat 同步用户到聊天系统
func (s *UserSyncService) SyncUserToChat(userID int) error {
	user, err := s.SyncUserFromCatcal(userID)
	if err != nil {
		return err
	}

	// 检查用户是否已存在
	existingUser, err := UserService.GetUserDetails(user.Account)
	if err == nil && existingUser.ID != 0 {
		// 更新用户信息
		existingUser.Nickname = user.Nickname
		// existingUser.Avatar = user.Avatar // PostgresUser没有Avatar字段
		if user.Email != "" {
			existingUser.Email = &user.Email
		}
		return UserService.ModifyUserInfo(&model.User{
			Account:  existingUser.Account,
			Nickname: existingUser.Nickname,
			Email:    user.Email,
		})
	}

	// 创建新用户
	return UserService.Register(user)
}
