package entities

type Item struct {
	ProductID   uint   `gorm:"primaryKey;column:product_id"`
	ProductName string `gorm:"index"`
	ProductDesc string
}

type HoleInfo struct {
	ID       uint `gorm:"primaryKey"`
	AngleID  uint `gorm:"index"`
	HolePath string
	SegPath  string
	ImgPath  string
}
