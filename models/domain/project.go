package domain

import "time"

type Project struct {
	ID               int       `gorm:"column:id;primary_key"`
	Name             string    `gorm:"column:name"`
	Address          string    `gorm:"column:address"`
	Host             string    `gorm:"column:host"`
	ContactName      string    `gorm:"column:contact_name"`
	ContactCellphone string    `gorm:"column:contact_cellphone"`
	Brokerage        int64     `gorm:"column:brokerage"`
	Deposit          int64     `gorm:"column:deposit"`
	Refund           int64     `gorm:"column:refund"`
	Status           string    `gorm:"column:status"`
	CreatedAt        time.Time `gorm:"column:create_time"`
	UpdatedAt        time.Time `gorm:"column:update_time"`
}

func (*Project) TableName() string {
	return "project"
}
