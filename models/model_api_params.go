package models

import (
	"github.com/go-trellis/api-manager/models/domain"

	"github.com/go-trellis/connector/tgorm"
	"github.com/jinzhu/gorm"
)

// RepoAPI 用户函数库
type RepoAPIParams interface {
	Add(*domain.APIParameters) error
	Get(params map[string]interface{}) (*domain.APIParameters, error)
	Delete(id int64) error
	GetList(params map[string]interface{}) ([]domain.APIParameters, error)
}

var defaultMAPIParams *MAPIParams

// MAPI 操作的类
type MAPIParams struct {
	tgorm.TGorm
	apiParams domain.APIParameters
}

// NewMAPIParams 获取操作对象
func NewMAPIParams(dbs map[string]*gorm.DB) RepoAPIParams {
	if defaultMAPIParams == nil {
		defaultMAPIParams = &MAPIParams{}
	}
	defaultMAPIParams.SetDBs(dbs)
	return defaultMAPIParams
}

// AddProject 增加
func (p *MAPIParams) Add(api *domain.APIParameters) error {
	db := p.Session().Create(api)
	return db.Error
}

// GetProject 通过id获取信息
func (p *MAPIParams) Get(params map[string]interface{}) (*domain.APIParameters, error) {
	api := &domain.APIParameters{}
	db := p.Session().Where(params).First(api)
	if db.Error != nil {
		return nil, db.Error
	}
	return api, nil
}

// GetProject 通过id获取信息
func (p *MAPIParams) Delete(id int64) error {
	db := p.Session().Model(p.apiParams).Delete(id)
	return db.Error
}

func (p *MAPIParams) GetList(params map[string]interface{}) ([]domain.APIParameters, error) {
	var apiParams []domain.APIParameters
	db := p.Session().Where(params).Find(&apiParams)
	if db.Error != nil {
		return nil, db.Error
	}
	return apiParams, nil
}
