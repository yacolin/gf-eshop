// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package do

import(
"github.com/gogf/gf/v2/frame/g"
"github.com/gogf/gf/v2/os/gtime"
)

// NotificationReads is the golang structure of table base_notification_reads for DAO operations like Where/Data.
type NotificationReads struct {
g.Meta `orm:"table:base_notification_reads, do:true"`
    Id interface{} // 主键        
    NotificationId interface{} // 关联通知ID  
    UserId interface{} // 用户ID      
    ReadAt         *gtime.Time // 阅读时间    
}