package models

import (
	"github.com/go-trellis/api-manager/models/domain"
	"github.com/go-trellis/connector/tgorm"
	"github.com/jinzhu/gorm"
)

// RepoProjectStatus 用户函数库
type RepoProjectStatus interface {
	GetList(params map[string]interface{}) ([]domain.ProjectStatus, error)
}

var defaultMProjectStatus *MProjectStatus

// MProjectStatus 操作的类
type MProjectStatus struct {
	tgorm.TGorm
	projectStatus domain.ProjectStatus
}

// NewMProjectStatus 获取操作对象
func NewMProjectStatus(dbs map[string]*gorm.DB) RepoProjectStatus {
	if defaultMProjectStatus == nil {
		defaultMProjectStatus = &MProjectStatus{}
	}
	defaultMProjectStatus.SetDBs(dbs)
	return defaultMProjectStatus
}

// GetProject 通过id获取信息
func (p *MProjectStatus) GetList(params map[string]interface{}) ([]domain.ProjectStatus, error) {
	var pts []domain.ProjectStatus
	db := p.Session().Where(params).Find(&pts)
	if db.Error != nil {
		return nil, db.Error
	}
	return pts, nil
}
