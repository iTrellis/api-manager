package domain

type APIParameters struct {
	ID            int64  `gorm:"column:id;primary_key"`
	ParentID      int64  `gorm:"column:parent_id"`
	APIID         int64  `gorm:"column:api_id"`
	APIParamsType int    `gorm:"column:api_params_type"`
	Key           string `gorm:"column:key"`
	FieldType     string `gorm:"column:field_type"`
	IsList        bool   `gorm:"column:is_list"`
	Required      bool   `gorm:"column:required"`
	Description   string `gorm:"column:description"`
	Sample        string `gorm:"column:sample"`
}

func (*APIParameters) TableName() string {
	return "api_parameters"
}
