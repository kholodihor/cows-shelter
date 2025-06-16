//go:build ignore

package main

import (
	"log"

	"gorm.io/gorm"
)

func CreateSupportTables(db *gorm.DB) error {
	// Create SupportCard table
	if err := db.AutoMigrate(&SupportCard{}); err != nil {
		return err
	}

	// Create SupportStep table
	if err := db.AutoMigrate(&SupportStep{}); err != nil {
		return err
	}

	// Add any initial data if needed
	supportCards := []SupportCard{
		{
			Title:    "card_title_1",
			Subtitle: "card_subtitle_1",
			Banner:   "card_banner_1",
			Image:    "support/support_1.png",
		},
		{
			Title:    "card_title_2",
			Subtitle: "card_subtitle_2",
			Banner:   "card_banner_2",
			Image:    "support/support_2.png",
		},
	}

	supportSteps := []SupportStep{
		{StepText: "step_1", Order: 1},
		{StepText: "step_2", Order: 2},
		{StepText: "step_3", Order: 3},
	}

	// Insert support cards
	for _, card := range supportCards {
		if err := db.FirstOrCreate(&card, "title = ?", card.Title).Error; err != nil {
			return err
		}
	}

	// Insert support steps
	for _, step := range supportSteps {
		if err := db.FirstOrCreate(&step, "step_text = ?", step.StepText).Error; err != nil {
			return err
		}
	}

	log.Println("Support tables created and seeded successfully")
	return nil
}

type SupportCard struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey"`
	Title     string `json:"title"`
	Subtitle  string `json:"subtitle"`
	Banner    string `json:"banner"`
	Image     string `json:"image"`
}

type SupportStep struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey"`
	StepText string `json:"step_text"`
	Order    int    `json:"order"`
}
