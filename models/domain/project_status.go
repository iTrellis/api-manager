package domain

type ProjectStatus struct {
	ID          string `gorm:"column:id;primary_key"`
	Description string `gorm:"column:description"`
}

func (*ProjectStatus) TableName() string {
	return "dim_project_status"
}
