package models

import (
	"gorm.io/gorm"
)

type LevelName string

const (
	BRONZE LevelName = "BRONZE"
	SILVER LevelName = "SILVER"
	GOLD   LevelName = "GOLD"
)

type Level struct {
	gorm.Model
	LevelName LevelName `json:"level_name" gorm:"not null;unique"`
}

func (l *Level) GetLevels(db *gorm.DB) ([]string, error) {
	var levels []Level
	err := db.Find(&levels).Error
	if err != nil {
		return nil, err
	}

	levelNames := make([]string, len(levels))
	for i, level := range levels {
		levelNames[i] = string(level.LevelName)
	}

	return levelNames, nil
}
