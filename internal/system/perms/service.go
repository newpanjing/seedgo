package perms

import (
	"github.com/newpanjing/seedgo/internal/model"
	"github.com/newpanjing/seedgo/internal/shared"
	"github.com/newpanjing/seedgo/pkg/request"
)

type Service struct {
	*shared.BaseLogic[model.Permission]
}

func NewService() *Service {
	logic := &Service{
		shared.NewBaseLogic[model.Permission](),
	}
	logic.DB = logic.BaseLogic.DB
	return logic
}

// GetTree 获取权限树
func (s *Service) GetTree(user *request.UserContext) ([]*model.Permission, error) {

	var perms []*model.Permission
	var err error

	// 1. 获取所有相关权限
	if user.IsSuper {
		// 超级管理员获取所有权限
		err = s.DB.Order("sort").Find(&perms).Error

	} else {

		//如果是普通用户，需要把角色信息填充到用户中
		var userModel model.User
		err = s.DB.Preload("Roles").First(&userModel, user.ID).Error
		if err != nil {
			return nil, err
		}
		user.Roles = userModel.Roles

		// 普通用户根据角色获取权限
		var roleIDs []uint64
		for _, role := range user.Roles {
			roleIDs = append(roleIDs, uint64(role.ID))
		}

		if len(roleIDs) > 0 {
			err = s.DB.Distinct("permission.*").
				Joins("JOIN role_permission ON role_permission.permission_id = permission.id").
				Where("role_permission.role_id IN ?", roleIDs).
				Order("sort").
				Find(&perms).Error
		}
	}

	if err != nil {
		return nil, err
	}

	// 2. 构建树形结构
	return buildTree(perms), nil
}

func buildTree(perms []*model.Permission) []*model.Permission {
	var roots []*model.Permission
	permMap := make(map[model.ID]*model.Permission)

	// 先建立 ID -> Node 的映射
	for _, p := range perms {
		permMap[p.ID] = p
	}

	// 再次遍历，挂载到父节点
	for _, p := range perms {
		if p.ParentID == nil || *p.ParentID == 0 {
			roots = append(roots, p)
		} else {
			if parent, ok := permMap[*p.ParentID]; ok {
				parent.Children = append(parent.Children, p)
			} else {
				// 如果找不到父节点（可能是数据问题或者父节点没权限），也可以作为根节点或者忽略
				// 这里选择作为根节点，避免数据丢失
				// roots = append(roots, p)
			}
		}
	}
	return roots
}
