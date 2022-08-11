// @Description
// @Author weitao.yin@shopee.com
// @Since 2022/6/13

package entity

// 登录请求
type LoginReq struct {
	UserName string `form:"userName"  binding:"required,min=4,max=32"`
	Password string `form:"password"  binding:"required,min=4,max=32"`
}

// 登录响应
type LoginResp struct {
	SessionId string `json:"sessionId"`
}

type InfoReq struct {
	SessionId string `form:"sessionId"`
}

// 用户信息
type UserInfoResp struct {
	UserName  string `json:"userName"`
	NickName  string `json:"nickName"`
	AvatarUrl string `json:"avatarUrl"`
	UserId    string `json:"userId"`
}

type EditReq struct {
	SessionId string `form:"sessionId"`
	NickName  string `form:"nickName" binding:"min=2,max=32"`
	AvatarUrl string `form:"avatarUrl" binding:"max=2000"`
}

type EditResp struct {
}

type UploadReq struct {
	UserName string
	FileName string
	FileExt  string
	Content  []byte
}

type UploadResp struct {
	AvatarUrl string `json:"avatarUrl"`
}

type RegisterReq struct {
	UserName string `form:"userName"  binding:"required,min=4,max=32"`
	Password string `form:"password"  binding:"required,min=4,max=32"`
	NickName string `form:"nickName" binding:"min=2,max=20"`
}

type RegisterResp struct {
}

type LogoutReq struct {
	SessionId string `form:"sessionId"`
}

type LogoutResp struct {
}

type PingResp struct {
	Text string `json:"text"`
}
