package models

import (
	"github.com/go-trellis/api-manager/models/domain"

	"github.com/go-trellis/connector/tgorm"
	"github.com/jinzhu/gorm"
)

// RepoProject 用户函数库
type RepoProject interface {
	Add(*domain.Project) error
	Get(params map[string]interface{}) (*domain.Project, error)
	GetList(params map[string]interface{}, offset, limit int) ([]domain.Project, error)
	Update(id int, params map[string]interface{}) error
}

var defaultMProject *MProject

// MProject 操作的类
type MProject struct {
	tgorm.TGorm
	project domain.Project
}

// NewMProject 获取操作对象
func NewMProject(dbs map[string]*gorm.DB) RepoProject {
	if defaultMProject == nil {
		defaultMProject = &MProject{}
	}
	defaultMProject.SetDBs(dbs)
	return defaultMProject
}

// GetProject 通过id获取信息
func (p *MProject) GetList(params map[string]interface{}, offset, limit int) ([]domain.Project, error) {
	var pts []domain.Project
	db := p.Session().Where(params).Offset(offset).Limit(limit).Find(&pts)
	if db.Error != nil {
		return nil, db.Error
	}
	return pts, nil
}

// GetProject 通过id获取信息
func (p *MProject) Get(params map[string]interface{}) (*domain.Project, error) {
	pt := &domain.Project{}
	db := p.Session().Where(params).First(pt)
	if db.Error != nil {
		return nil, db.Error
	}
	return pt, nil
}

// AddProject 增加
func (p *MProject) Add(pt *domain.Project) error {
	db := p.Session().Create(pt)
	return db.Error
}

func (p *MProject) Update(id int, params map[string]interface{}) error {
	db := p.Session().Model(p.project).Where(id).Update(params)
	return db.Error
}
