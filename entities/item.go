package entities

type Item struct {
	ProductID   uint   `gorm:"primaryKey;column:product_id" json:"productId"`
	ProductName string `gorm:"index" json:"productName"`
	ProductDesc string `json:"productDesc"`
}

type HoleInfo struct {
	ID       uint `gorm:"primaryKey"`
	AngleID  uint `gorm:"index"`
	HolePath string
	SegPath  string
	ImgPath  string
}
