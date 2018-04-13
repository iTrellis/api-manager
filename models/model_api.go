package models

import (
	"github.com/go-trellis/api-manager/models/domain"

	"github.com/go-trellis/connector/tgorm"
	"github.com/jinzhu/gorm"
)

// RepoAPI 用户函数库
type RepoAPI interface {
	Add(*domain.API) error
	Get(params map[string]interface{}) (*domain.API, error)
	Update(api *domain.API, params map[string]interface{}) error
	GetList(params map[string]interface{}, page, number int) ([]domain.API, error)
}

var defaultMAPI *MAPI

// MAPI 操作的类
type MAPI struct {
	tgorm.TGorm
	api domain.API
}

// NewMAPI 获取操作对象
func NewMAPI(dbs map[string]*gorm.DB) RepoAPI {
	if defaultMAPI == nil {
		defaultMAPI = &MAPI{}
	}
	defaultMAPI.SetDBs(dbs)
	return defaultMAPI
}

// AddProject 增加
func (p *MAPI) Add(api *domain.API) error {
	db := p.Session().Create(api)
	return db.Error
}

// GetProject 通过id获取信息
func (p *MAPI) Get(params map[string]interface{}) (*domain.API, error) {
	api := &domain.API{}
	db := p.Session().Where(params).First(api)
	if db.Error != nil {
		return nil, db.Error
	}
	return api, nil
}

func (p *MAPI) GetList(params map[string]interface{}, offset, limit int) ([]domain.API, error) {
	var apis []domain.API
	db := p.Session().Where(params).Offset(offset).Limit(limit).Find(&apis)
	if db.Error != nil {
		return nil, db.Error
	}
	return apis, nil
}

func (p *MAPI) Update(api *domain.API, params map[string]interface{}) error {
	db := p.Session().Model(p.api).Where(api.ID).Update(params)
	return db.Error
}
