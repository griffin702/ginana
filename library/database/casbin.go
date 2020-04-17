package database

import (
	"context"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"path"
	"time"
)

type casbinService interface {
	GetEFRoles(ctx context.Context) (roles []*EFRolePolicy, err error)
	GetEFUsers(ctx context.Context) (users []*EFUseRole, err error)
}

type EFRolePolicy struct {
	RoleName string
	Router   string
	Method   string
}

type EFUseRole struct {
	UserID   int64
	RoleName string
}

// Config casbin config.
type CasbinConfig struct {
	Model            string
	Enable           bool
	AutoLoad         bool
	AutoLoadInternal int
}

// NewCasbinConn with CasbinConfig and Custom Adapter
// Start Goroutine to Watching CasbinModel and CasbinConfig
func NewCasbinConn(svc casbinService, dir string, c *CasbinConfig) (e *casbin.SyncedEnforcer, err error) {
	if !c.Enable {
		return
	}
	adapter := NewCasbinAdapter(svc)
	//adapter := gormadapter.NewAdapterByDB(db)
	m := path.Join(dir, c.Model)
	e, err = casbin.NewSyncedEnforcer(m, adapter)
	if err != nil {
		return
	}
	e.EnableAutoSave(false)
	e.EnableAutoBuildRoleLinks(true)
	if c.AutoLoad {
		_ = e.InitWithModelAndAdapter(e.GetModel(), adapter)
		e.StartAutoLoadPolicy(time.Duration(c.AutoLoadInternal) * time.Second)
	} else {
		err = adapter.LoadPolicy(e.GetModel())
		if err != nil {
			return
		}
	}
	err = e.BuildRoleLinks()
	return
}

// NewCasbinAdapter 创建casbin适配器
func NewCasbinAdapter(svc casbinService) *CasbinAdapter {
	return &CasbinAdapter{
		svc: svc,
	}
}

// CasbinAdapter casbin适配器
type CasbinAdapter struct {
	svc casbinService
}

// LoadPolicy loads all policy rules from the storage.
func (a *CasbinAdapter) LoadPolicy(model model.Model) error {
	ctx := context.Background()
	err := a.loadRolePolicy(ctx, model)
	if err != nil {
		return err
	}
	err = a.loadUserPolicy(ctx, model)
	if err != nil {
		return err
	}
	return nil
}

// loadRolePolicy loads all policy rules of role.
func (a *CasbinAdapter) loadRolePolicy(ctx context.Context, model model.Model) error {
	roles, err := a.svc.GetEFRoles(ctx)
	if err != nil {
		return err
	}
	for _, role := range roles {
		if role.Router == "" || role.Method == "" {
			continue
		}
		line := fmt.Sprintf("p,%s,%s,%s", role.RoleName, role.Router, role.Method)
		persist.LoadPolicyLine(line, model)
	}
	return nil
}

// loadRolePolicy loads all policy rules of user.
func (a *CasbinAdapter) loadUserPolicy(ctx context.Context, model model.Model) error {
	users, err := a.svc.GetEFUsers(ctx)
	if err != nil {
		return err
	}
	for _, user := range users {
		line := fmt.Sprintf("g,%d,%s", user.UserID, user.RoleName)
		persist.LoadPolicyLine(line, model)
	}
	return nil
}

// SavePolicy saves all policy rules to the storage.
func (a *CasbinAdapter) SavePolicy(model model.Model) error {
	return nil
}

// AddPolicy adds a policy rule to the storage.
// This is part of the Auto-Save feature.
func (a *CasbinAdapter) AddPolicy(sec string, ptype string, rule []string) error {
	return nil
}

// RemovePolicy removes a policy rule from the storage.
// This is part of the Auto-Save feature.
func (a *CasbinAdapter) RemovePolicy(sec string, ptype string, rule []string) error {
	return nil
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
// This is part of the Auto-Save feature.
func (a *CasbinAdapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return nil
}
