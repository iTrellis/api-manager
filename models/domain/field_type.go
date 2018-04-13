package domain

type FieldType struct {
	ID          int    `gorm:"column:id;primary_key"`
	Type        string `gorm:"column:type"`
	Description string `gorm:"column:description"`
}

func (*FieldType) TableName() string {
	return "dim_field_type"
}
