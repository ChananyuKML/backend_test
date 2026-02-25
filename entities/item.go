package entities

import "io"

type Item struct {
	ProductID       uint   `gorm:"primaryKey;column:product_id" json:"productId"`
	ProductName     string `gorm:"index" json:"productName"`
	ProductDesc     string `json:"productDesc"`
	ProductImageKey string `json:"productImageKey"`
}

type HoleInfo struct {
	ID       uint `gorm:"primaryKey"`
	AngleID  uint `gorm:"index"`
	HolePath string
	SegPath  string
	ImgPath  string
}

type FileStream struct {
	Reader      io.Reader
	ContentType string
	Size        int64
}
