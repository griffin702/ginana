package service

import (
	"context"
	"fmt"
	"ginana/internal/model"
	"ginana/library/database"
	"ginana/library/ecode"
	"sync"
)

func (s *service) GetEFRoles(c context.Context) (roles []*database.CasbinRole, err error) {
	var roleIdList []int64
	s.db.Model(&model.Role{}).Select("id").Pluck("id", &roleIdList)
	var wg sync.WaitGroup
	var ch = make(chan int64, 1)
	wg.Add(len(roleIdList))
	for _, roleId := range roleIdList {
		go func(roleId int64, roles *[]*database.CasbinRole, wg *sync.WaitGroup) {
			r, err := s.GetRole(c, roleId)
			if err != nil {
				return
			}
			role := new(database.CasbinRole)
			role.RoleName = r.RoleName
			for _, p := range r.Policys {
				policy := new(database.CasbinPolicy)
				policy.Router = p.Router
				policy.Method = p.Method
				role.Policys = append(role.Policys, policy)
			}
			ch <- roleId
			*roles = append(*roles, role)
			<-ch
			wg.Done()
		}(roleId, &roles, &wg)
	}
	wg.Wait()
	return
}

func (s *service) GetRole(ctx context.Context, id int64) (role *model.Role, err error) {
	key := fmt.Sprintf("role_%d", id)
	role = new(model.Role)
	err = s.mc.Get(key, role)
	if err != nil {
		role.ID = id
		if err = s.db.Find(role).Related(&role.Policys, "Policys").Error; err != nil {
			err = ecode.Errorf(s.GetError(1001, err.Error()))
			return
		}
		if err = s.mc.Set(key, role); err != nil {
			err = ecode.Errorf(s.GetError(1002, err.Error()))
			return
		}
	}
	return
}
