package model

type LoginLog struct {
	Base
	UserId    int64  `json:"userId" gorm:"type:bigint(20);comment:登录用户的ID"`
	Nickname  string `json:"nickname" gorm:"type:varchar(50);not null;index:idx_uname;comment:用户昵称"`
	RealName  string `json:"realName" gorm:"type:varchar(50);comment:用户真名"`
	DeptName  string `json:"deptName" gorm:"type:varchar(50);comment:部门名称"`
	Ip        string `json:"ip" gorm:"type:varchar(50);comment:登录IP"`
	Address   string `json:"address" gorm:"type:varchar(255);comment:登录地点"`
	UserAgent string `json:"userAgent" gorm:"type:varchar(255);comment:浏览器的userAgent"`
}

func (*LoginLog) TableName() string {
	return "ops_login_log"
}
