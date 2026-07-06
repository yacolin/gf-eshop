package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/model/entity"
)

type UserProfileReq struct {
	g.Meta `path:"/user" tags:"User" method:"get" summary:"获取当前用户资料"`
}
type UserProfileRes struct {
	*entity.Users
	Info *entity.Infos `json:"info,omitempty"`
}

type UserUpdateInfoReq struct {
	g.Meta   `path:"/user" tags:"User" method:"put" summary:"更新个人信息"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Gender   int    `json:"gender"`
	Birthday string `json:"birthday"`
	Bio      string `json:"bio"`
	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
	ZipCode  string `json:"zip_code"`
	Language string `json:"language"`
	Timezone string `json:"timezone"`
}
type UserUpdateInfoRes struct{}
