package domain

type APIParamsType struct {
	ID          int64  `gorm:"column:id;primary_key"`
	Description string `gorm:"column:description"`
}

func (*APIParamsType) TableName() string {
	return "dim_api_params_type"
}
