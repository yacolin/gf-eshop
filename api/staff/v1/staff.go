package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type StaffLoginReq struct {
	g.Meta   `path:"/staff/login" tags:"Staff" method:"post" summary:"系统用户登录"`
	Username string `json:"username" v:"required" description:"用户名"`
	Password string `json:"password" v:"required" description:"密码"`
}
type StaffLoginRes struct {
	AccessToken  string `json:"access_token"  description:"访问令牌"`
	ExpireIn     int64  `json:"expire_in"     description:"access_token 过期时间（秒）"`
	RefreshToken string `json:"refresh_token" description:"刷新令牌"`
	RefreshIn    int64  `json:"refresh_in"    description:"refresh_token 过期时间（秒）"`
	StaffId      int64  `json:"staff_id"      description:"用户ID"`
	Username     string `json:"username"      description:"用户名"`
	RealName     string `json:"real_name"     description:"真实姓名"`
}

type StaffRefreshTokenReq struct {
	g.Meta       `path:"/staff/refresh" tags:"Staff" method:"post" summary:"刷新令牌"`
	RefreshToken string `json:"refresh_token" v:"required" description:"刷新令牌"`
}
type StaffRefreshTokenRes struct {
	AccessToken  string `json:"access_token"  description:"新的访问令牌"`
	ExpireIn     int64  `json:"expire_in"     description:"access_token 过期时间（秒）"`
	RefreshToken string `json:"refresh_token" description:"新的刷新令牌"`
	RefreshIn    int64  `json:"refresh_in"    description:"refresh_token 过期时间（秒）"`
}

type StaffLogoutReq struct {
	g.Meta `path:"/staff/logout" tags:"Staff" method:"post" summary:"系统用户退出登录"`
}
type StaffLogoutRes struct{}

type StaffProfileReq struct {
	g.Meta `path:"/staff/profile" tags:"Staff" method:"get" summary:"获取当前用户信息"`
}
type StaffProfileRes struct {
	Id              int64       `json:"id"               description:"用户ID"`
	Username        string      `json:"username"         description:"用户名"`
	RealName        string      `json:"real_name"        description:"真实姓名"`
	Email           string      `json:"email"            description:"邮箱"`
	Phone           string      `json:"phone"            description:"手机号"`
	Avatar          string      `json:"avatar"           description:"头像URL"`
	Status          int         `json:"status"           description:"状态"`
	LastLoginIp     string      `json:"last_login_ip"    description:"最后登录IP"`
	DepartmentIds   []int64     `json:"department_ids"   description:"所属部门ID列表"`
	DepartmentNames []string    `json:"department_names" description:"所属部门名称列表"`
}

type StaffPermissionsReq struct {
	g.Meta `path:"/staff/permissions" tags:"Staff" method:"get" summary:"获取当前用户权限和角色"`
}
type StaffPermissionsRes struct {
	Roles       []string `json:"roles"       description:"角色名称列表"`
	Permissions []string `json:"permissions" description:"权限标识列表"`
}

type StaffListReq struct {
	g.Meta   `path:"/staff" tags:"Staff" method:"get" summary:"员工列表"`
	Page     int    `json:"page"      description:"页码"`
	PageSize int    `json:"page_size" description:"每页条数"`
	Keyword  string `json:"keyword"   description:"搜索关键词（用户名/姓名）"`
	Status   *int   `json:"status"    description:"状态：1-正常 0-禁用"`
}

type StaffListItem struct {
	Id              int64       `json:"id"               description:"主键"`
	Username        string      `json:"username"         description:"登录用户名"`
	RealName        string      `json:"real_name"        description:"真实姓名"`
	Email           string      `json:"email"            description:"邮箱"`
	Phone           string      `json:"phone"            description:"手机号"`
	Avatar          string      `json:"avatar"           description:"头像URL"`
	Status          int         `json:"status"           description:"1-正常 0-禁用"`
	LastLoginIp     string      `json:"last_login_ip"    description:"最后登录IP"`
	LastLoginAt     *gtime.Time `json:"last_login_at"    description:"最后登录时间"`
	CreatedAt       *gtime.Time `json:"created_at"       description:"创建时间"`
	RoleIds         []int64     `json:"role_ids"         description:"角色ID列表"`
	RoleNames       []string    `json:"role_names"       description:"角色名称列表"`
	DepartmentIds   []int64     `json:"department_ids"   description:"所属部门ID列表"`
	DepartmentNames []string    `json:"department_names" description:"所属部门名称列表"`
}
type StaffListRes struct {
	List  []*StaffListItem `json:"list"`
	Total int              `json:"total"`
}

type StaffAssignRolesReq struct {
	g.Meta  `path:"/staff/{id}/roles" tags:"Staff" method:"put" summary:"分配角色"`
	Id      int64   `json:"id"  v:"required" description:"员工ID"`
	RoleIds []int64 `json:"role_ids" v:"required" description:"角色ID列表"`
}
type StaffAssignRolesRes struct{}
