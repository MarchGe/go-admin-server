package req

type UserUpsertReq struct {
	Sex                int8    `json:"sex" binding:"omitempty,oneof=0 1" label:"性别"`
	Birthday           string  `json:"birthday" binding:"omitempty,regex=^[0-9]{4}-[0-9]{2}-[0-9]{2}$" label:"出生日期"`
	Cellphone          string  `json:"cellphone" binding:"omitempty,regex=^1[0-9]{10}$" label:"手机号"`
	Email              string  `json:"email" binding:"required,email,max=50" label:"邮箱"`
	Name               string  `json:"name" binding:"omitempty,max=20" label:"姓名"`
	Nickname           string  `json:"nickname" binding:"required,max=20" label:"昵称"`
	DeptId             int64   `json:"deptId" binding:"required,min=1" label:"部门"`
	RoleIds            []int64 `json:"roleIds" binding:"required,min=1,max=1000" label:"角色列表"`
	JobIds             []int64 `json:"jobIds" binding:"required,min=1,max=1000" label:"岗位列表"`
	ShouldSendPassword bool    `json:"shouldSendPassword" binding:"omitempty" label:"邮件发送标志"` // 仅新增的时候有效
}
