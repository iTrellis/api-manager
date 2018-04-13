package domain

type API struct {
	ID           int64  `gorm:"column:id;primary_key"`
	ProjectID    int    `gorm:"column:project_id"`
	ProjectName  string `gorm:"column:project_name"`
	Name         string `gorm:"column:name"`
	Path         string `gorm:"column:path"`
	HttpMethodID int    `gorm:"column:http_method_id"`
	HttpMethod   string `gorm:"column:http_method"`
	Description  string `gorm:"column:description"`
	Status       string `gorm:"column:status"`
}

func (*API) TableName() string {
	return "api"
}
