package authz

import (
	"fmt"
	"github.com/MarchGe/go-admin-server/app/common/R"
	"github.com/MarchGe/go-admin-server/app/common/constant"
	"github.com/MarchGe/go-admin-server/app/common/database"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"sync"
)

var _enforcer *casbin.Enforcer
var mux = sync.Mutex{}

func RequiresPermissions(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetBool(constant.IsRootUser) {
			c.Next()
			return
		}
		userId := c.GetInt64(constant.SessionUserId)
		enforcer := GetEnforcer()
		for _, permission := range permissions {
			obj, act, err := parsePermission(permission)
			if err != nil {
				panic(fmt.Errorf("parse permission string error: %v", err))
			}
			if ok, e := enforcer.Enforce(UserSub(userId), obj, act); ok && e == nil {
				c.Next()
				return
			}
		}
		R.Fail(c, fmt.Sprintf("没有权限执行该操作，缺少权限: %s", strings.Join(permissions, " | ")), http.StatusForbidden)
		return
	}
}

func parsePermission(permission string) (string, string, error) {
	if !strings.Contains(permission, ":") {
		return "", "", fmt.Errorf("permission string '%s' must contains ':'", permission)
	}
	parts := strings.Split(permission, ":")
	if len(parts) > 2 {
		return "", "", fmt.Errorf("permission string '%s' contains at most one ':'", permission)
	}
	return parts[0], parts[1], nil
}

func GetEnforcer() *casbin.Enforcer {
	if _enforcer != nil {
		return _enforcer
	}
	mux.Lock()
	defer mux.Unlock()
	if _enforcer != nil {
		return _enforcer
	}
	m, err := model.NewModelFromString(ModelString)
	if err != nil {
		panic(err)
	}
	adapter, err := gormadapter.NewAdapterByDB(database.GetMysql())
	if err != nil {
		panic(err)
	}
	_enforcer, err = casbin.NewEnforcer(m, adapter)
	if err != nil {
		panic(err)
	}
	return _enforcer
}

func UpdateSubPolicies(sub string, permissionStrings []string) error {
	enforcer := GetEnforcer()
	adapter := enforcer.GetAdapter().(*gormadapter.Adapter)
	err := adapter.Transaction(enforcer, func(enforcer casbin.IEnforcer) error {
		// 删除sub关联的权限策略
		if _, err := enforcer.DeletePermissionsForUser(sub); err != nil {
			return fmt.Errorf("delete permissons for user error when update sub policies, %w", err)
		}

		// 添加sub的权限策略
		permissions := make([][]string, len(permissionStrings))
		for i, symbol := range permissionStrings {
			permission := strings.Split(symbol, ":")
			permissions[i] = permission
		}
		if len(permissions) > 0 {
			if _, err := enforcer.AddPermissionsForUser(sub, permissions...); err != nil {
				return fmt.Errorf("add permissions for user error when update sub policies, %w", err)
			}
		}
		return nil
	})
	return err
}

func DeleteUser(sub string) error {
	if _, err := GetEnforcer().DeleteUser(sub); err != nil {
		return fmt.Errorf("casbin delete user error, %w", err)
	}
	return nil
}

func DeleteRole(sub string) error {
	if _, err := GetEnforcer().DeleteRole(sub); err != nil {
		return fmt.Errorf("casbin delete role error, %w", err)
	}
	return nil
}

func DeletePermission(permissionString string) error {
	permission := strings.Split(permissionString, ":")
	if _, err := GetEnforcer().DeletePermission(permission...); err != nil {
		return fmt.Errorf("casbin delete psermission error, %w", err)
	}
	return nil
}

func UpdateGroupingPolicies(userSub string, roleSubs []string) error {
	enforcer := GetEnforcer()
	adapter := enforcer.GetAdapter().(*gormadapter.Adapter)
	err := adapter.Transaction(enforcer, func(enforcer casbin.IEnforcer) error {
		// 删除用户已绑定的角色
		if _, err := enforcer.DeleteRolesForUser(userSub); err != nil {
			return fmt.Errorf("delete roles for user error when update grouping policies, %w", err)
		}
		// 用户绑定新角色
		rules := make([][]string, len(roleSubs))
		for i, roleSub := range roleSubs {
			rules[i] = []string{userSub, roleSub}
		}
		if len(rules) > 0 {
			if _, err := enforcer.AddGroupingPolicies(rules); err != nil {
				return fmt.Errorf("add grouping policies error when update grouping policies, %w", err)
			}
		}
		return nil
	})
	return err
}

func UserSub(userId int64) string {
	return fmt.Sprintf("user:%d", userId)
}

func RoleSub(roleId int64) string {
	return fmt.Sprintf("role:%d", roleId)
}
