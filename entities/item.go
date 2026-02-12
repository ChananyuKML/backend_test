package entities

type Item struct {
	ID          uint `gorm:"primaryKey"`
	OwnerID     uint `gorm:"index"`
	Description string
}

type HoleInfo struct {
	ProductID uint `gorm:"primaryKey"`
	AngleID   uint `gorm:"index"`
	HolePath  string
	SegPath   string
	ImgPath   string
}
