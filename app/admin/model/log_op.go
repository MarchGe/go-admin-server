package model

type OpLog struct {
	Base
	UserId    int64  `json:"userId" gorm:"type:bigint(20);comment:登录用户的ID"`
	Nickname  string `json:"nickname" gorm:"type:varchar(50);not null;index:idx_uname;comment:用户昵称"`
	RealName  string `json:"realName" gorm:"type:varchar(50);comment:用户真名"`
	DeptName  string `json:"deptName" gorm:"type:varchar(50);comment:部门名称"`
	Ip        string `json:"ip" gorm:"type:varchar(50);comment:登录IP"`
	Address   string `json:"address" gorm:"type:varchar(255);comment:登录地点"`
	UserAgent string `json:"userAgent" gorm:"type:varchar(255);comment:浏览器的userAgent"`
	Action    string `json:"action" gorm:"type:varchar(50);comment:操作名称"`
	Target    string `json:"target" gorm:"type:varchar(50);comment:操作对象"`
	Path      string `json:"path" gorm:"type:varchar(255);comment:请求url"`
	Query     string `json:"query" gorm:"type:varchar(255);comment:请求url"`
	Body      string `json:"body" gorm:"type:varchar(500);comment:body参数信息"`
}

func (*OpLog) TableName() string {
	return "ops_op_log"
}
