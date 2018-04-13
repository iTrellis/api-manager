package domain

type HTTPMethod struct {
	ID            int64  `gorm:"column:id;primary_key"`
	Method        string `gorm:"column:method"`
	DisplayMethod string `gorm:"column:display_method"`
}

func (*HTTPMethod) TableName() string {
	return "dim_http_method"
}
