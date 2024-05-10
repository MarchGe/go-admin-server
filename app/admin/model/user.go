package model

const (
	UserStatusNormal  int8 = 0
	UserStatusDisable int8 = 1
)
const (
	UserSexMan   int8 = 0
	UserSexWoman int8 = 1
)

func ExistsUserStatus(status int8) bool {
	if status == UserStatusNormal || status == UserStatusDisable {
		return true
	}
	return false
}

type User struct {
	Base
	Name      string  `json:"name" gorm:"type:varchar(50);index:idx_name;comment:真实姓名"`
	Nickname  string  `json:"nickname" gorm:"type:varchar(50);index:idx_nickname;comment:昵称"`
	Cellphone string  `json:"cellphone" gorm:"type:varchar(50);index:idx_cellphone;comment:手机号"`
	Email     string  `json:"email" gorm:"type:varchar(100);index:idx_email,unique;comment:邮箱"`
	Sex       int8    `json:"sex" gorm:"type:tinyint(4);default:0;comment:性别，0-男，1-女"`
	Birthday  string  `json:"birthday" gorm:"type:varchar(50);comment:生日"`
	Status    int8    `json:"status" gorm:"type:tinyint(4);default:0;comment:账号状态，0-正常，1-禁用"`
	Root      bool    `json:"root" gorm:"<-:create;type:tinyint(4);default:0;comment:是否超级用户，0-否，1-是"`
	DeptId    int64   `json:"deptId" gorm:"type:bigint(20);default:0;comment:部门ID"`
	Dept      *Dept   `json:"dept" gorm:"-:migration;foreignKey:DeptId"`
	RoleList  []*Role `json:"roleList" gorm:"-:migration;many2many:ops_role_user"`
	MenuList  []*Menu `json:"menuList" gorm:"-:migration;many2many:ops_user_menu"`
	JobList   []*Job  `json:"jobList" gorm:"-:migration;many2many:ops_user_job"`
}

func (*User) TableName() string {
	return "ops_user"
}
