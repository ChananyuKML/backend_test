package entities

type Item struct {
	ID          uint `gorm:"primaryKey"`
	OwnerID     uint `gorm:"index"`
	Description string
}
