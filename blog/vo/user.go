package vo

import (
	"time"
	"wee-server/blog/entity"
)

type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserVO struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	CreateTime string `json:"createTime"`
	Avatar     string `json:"avatar"`
	MaxTime    int    `json:"maxTime"`
}

func NewUser(u *entity.User) *UserVO {
	t := time.Unix(int64(u.CreateTime), 0).Format("2006-01-02 15:04:05")
	return &UserVO{
		Id:         u.Id,
		Name:       u.Name,
		CreateTime: t,
		Avatar:     u.Avatar.String,
		MaxTime:    3600,
	}
}
