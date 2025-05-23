package services

import (
	"log"
	"slack-review-notify/models"
	"time"

	"gorm.io/gorm"
)

// アーカイブされたチャンネルの設定を非アクティブにする
func CleanupArchivedChannels(db *gorm.DB) {
	var configs []models.ChannelConfig
	db.Where("is_active = ?", true).Find(&configs)
	
	for _, config := range configs {
		isArchived, err := IsChannelArchived(config.SlackChannelID)
		if err != nil {
			log.Printf("channel status check error (channel: %s): %v", config.SlackChannelID, err)
			continue
		}
		
		if isArchived {
			log.Printf("channel %s is archived", config.SlackChannelID)
			
			// 非アクティブに更新
			config.IsActive = false
			config.UpdatedAt = time.Now()
			if err := db.Save(&config).Error; err != nil {
				log.Printf("channel config update error: %v", err)
			} else {
				log.Printf("channel %s config is deactivated", config.SlackChannelID)
			}
		}
	}
} 
