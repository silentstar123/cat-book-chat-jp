package service

import (
	"bytes"
	"chat-room/pkg/errors"
	"chat-room/pkg/global/log"
	"encoding/json"
	"fmt"
	"net/http"

	"chat-room/internal/model"
)

type userService struct {
}

var UserService = new(userService)

// 从主服务器获取用户信息
func (u *userService) GetUserFromMainServer(account string) (*model.PostgresUser, error) {
	// 主服务器API地址
	mainServerURL := "http://localhost:8082"

	// 调用主服务器的用户API
	resp, err := http.Get(fmt.Sprintf("%s/api/v1/users/%s", mainServerURL, account))
	if err != nil {
		return nil, fmt.Errorf("主サーバーへの接続に失敗しました: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ユーザー情報の取得に失敗しました: %d", resp.StatusCode)
	}

	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Account  string `json:"account"`
			NickName string `json:"nickname"`
			Email    string `json:"email"`
			Phone    string `json:"phone"`
			Role     string `json:"role"`
			Status   uint   `json:"status"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("レスポンスの解析に失敗しました: %w", err)
	}

	if result.Code != 200 {
		return nil, fmt.Errorf("ユーザーが見つかりません: %s", result.Msg)
	}

	// 转换为PostgresUser格式
	user := &model.PostgresUser{
		Account:    result.Data.Account,
		Nickname:   result.Data.NickName,
		Role:       result.Data.Role,
		Status:     int64(result.Data.Status),
		EmailValid: true,
		PhoneValid: true,
	}

	if result.Data.Email != "" {
		user.Email = &result.Data.Email
	}
	if result.Data.Phone != "" {
		user.Phone = &result.Data.Phone
	}

	return user, nil
}

// 验证用户登录（从主服务器验证）
func (u *userService) ValidateUserLogin(account, password string) (bool, error) {
	// 主服务器API地址
	mainServerURL := "http://localhost:8082"

	// 调用主服务器的登录API
	loginData := map[string]interface{}{
		"loginType": "account",
		"account":   account,
		"password":  password,
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		return false, fmt.Errorf("ログインデータの作成に失敗しました: %w", err)
	}

	resp, err := http.Post(fmt.Sprintf("%s/api/v1/login", mainServerURL), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return false, fmt.Errorf("ログインAPIへの接続に失敗しました: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, nil
}

// 获取用户详情（从主服务器获取）
func (u *userService) GetUserDetails(account string) (*model.PostgresUser, error) {
	return u.GetUserFromMainServer(account)
}

// 修改后的登录方法
func (u *userService) Login(user *model.User) bool {
	log.Logger.Debug("user", log.Any("user in service", user))

	// 从主服务器验证用户登录
	isValid, err := u.ValidateUserLogin(user.Username, user.Password)
	if err != nil {
		log.Logger.Error("login error", log.Any("error", err))
		return false
	}

	if !isValid {
		return false
	}

	// 获取用户详细信息
	mainUser, err := u.GetUserFromMainServer(user.Username)
	if err != nil {
		log.Logger.Error("get user details error", log.Any("error", err))
		return false
	}

	// 更新用户信息
	user.Account = mainUser.Account
	user.Nickname = mainUser.Nickname
	if mainUser.Email != nil {
		user.Email = *mainUser.Email
	}

	return true
}

// 修改后的注册方法（实际应该调用主服务器的注册API）
func (u *userService) Register(user *model.User) error {
	// 聊天服务器不直接处理用户注册，应该通过主服务器
	return errors.New("ユーザー登録はメインサーバーで行ってください")
}

// 获取用户列表（从主服务器获取）
func (u *userService) GetUserList() ([]*model.PostgresUser, error) {
	// 这里可以实现从主服务器获取用户列表的逻辑
	// 暂时返回空列表
	return []*model.PostgresUser{}, nil
}

// 修改用户信息（调用主服务器API）
func (u *userService) ModifyUserInfo(user *model.User) error {
	// 调用主服务器的用户更新API
	mainServerURL := "http://localhost:8082"

	updateData := map[string]interface{}{
		"nickname": user.Nickname,
		"email":    user.Email,
	}

	jsonData, err := json.Marshal(updateData)
	if err != nil {
		return fmt.Errorf("更新データの作成に失敗しました: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/users/update", mainServerURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("リクエストの作成に失敗しました: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	// 这里需要添加Authorization header，从context中获取token

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("更新APIへの接続に失敗しました: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ユーザー情報の更新に失敗しました: %d", resp.StatusCode)
	}

	return nil
}

// 添加好友（调用主服务器API）
func (u *userService) AddFriend(user *model.User) error {
	// 调用主服务器的关注API
	mainServerURL := "http://localhost:8082"

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/follows/%s", mainServerURL, user.Account), nil)
	if err != nil {
		return fmt.Errorf("リクエストの作成に失敗しました: %w", err)
	}

	// 这里需要添加Authorization header，从context中获取token

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("フォローAPIへの接続に失敗しました: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("フォローに失敗しました: %d", resp.StatusCode)
	}

	return nil
}

// 根据用户名或群名搜索
func (u *userService) GetUserOrGroupByName(name string) ([]*model.PostgresUser, error) {
	// 这里可以实现搜索逻辑
	// 暂时返回空列表
	return []*model.PostgresUser{}, nil
}

// 修改用户头像
func (u *userService) ModifyUserAvatar(user *model.User) error {
	// 调用主服务器的头像更新API
	mainServerURL := "http://localhost:8082"

	updateData := map[string]interface{}{
		"avatar": user.Avatar,
	}

	jsonData, err := json.Marshal(updateData)
	if err != nil {
		return fmt.Errorf("アバターデータの作成に失敗しました: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/users/update", mainServerURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("リクエストの作成に失敗しました: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	// 这里需要添加Authorization header，从context中获取token

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("アバター更新APIへの接続に失敗しました: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("アバターの更新に失敗しました: %d", resp.StatusCode)
	}

	return nil
}
