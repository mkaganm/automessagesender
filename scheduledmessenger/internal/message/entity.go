package message

import "time"

type Message struct {
	ID          uint      `gorm:"primaryKey"`
	MessageID   string    `gorm:"uniqueIndex;not null"`
	Number      string    `gorm:"not null"`
	Context     string    `gorm:"type:LONGTEXT"`
	IsSend      bool      `gorm:"default:false"`
	IsProcessed bool      `gorm:"default:false"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
