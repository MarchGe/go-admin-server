package model

type ExceptionLog struct {
	Base
	UserId    int64  `json:"userId" gorm:"type:bigint(20);comment:登录用户的ID"`
	Nickname  string `json:"nickname" gorm:"type:varchar(50);not null;index:idx_uname;comment:用户昵称"`
	Ip        string `json:"ip" gorm:"type:varchar(50);comment:登录IP"`
	UserAgent string `json:"userAgent" gorm:"type:varchar(255);comment:浏览器的userAgent"`
	Path      string `json:"path" gorm:"type:varchar(255);comment:请求url"`
	Query     string `json:"query" gorm:"type:varchar(255);comment:请求url"`
	Body      string `json:"body" gorm:"type:varchar(500);comment:body参数信息"`
	Error     string `json:"error" gorm:"type:varchar(5000);comment:错误内容"`
}

func (*ExceptionLog) TableName() string {
	return "ops_exception_log"
}
