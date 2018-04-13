package models

import (
	"github.com/go-trellis/api-manager/models/domain"
	"github.com/go-trellis/connector/tgorm"
	"github.com/jinzhu/gorm"
)

// RepoFieldType 用户函数库
type RepoFieldType interface {
	GetList(params map[string]interface{}) ([]domain.FieldType, error)
}

var defaultMFieldType *MFieldType

// MFieldType 操作的类
type MFieldType struct {
	tgorm.TGorm
	FieldType domain.FieldType
}

// NewMFieldType 获取操作对象
func NewMFieldType(dbs map[string]*gorm.DB) RepoFieldType {
	if defaultMFieldType == nil {
		defaultMFieldType = &MFieldType{}
	}
	defaultMFieldType.SetDBs(dbs)
	return defaultMFieldType
}

// GetProject 通过id获取信息
func (p *MFieldType) GetList(params map[string]interface{}) ([]domain.FieldType, error) {
	var pts []domain.FieldType
	db := p.Session().Where(params).Find(&pts)
	if db.Error != nil {
		return nil, db.Error
	}
	return pts, nil
}
