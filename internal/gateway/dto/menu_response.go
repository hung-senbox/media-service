package dto

import "time"

type MenuResponse struct {
	Language LanguageSetting     `json:"language"`
	Menus    []ComponentResponse `json:"menus"`
}

type ComponentResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Key        string `json:"key"`
	Value      string `json:"value"`
	Order      int    `json:"order"`
	IsShow     bool   `json:"is_show"`
	LanguageID uint   `json:"language_id"`
}

type LanguageSetting struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	LangKey     string    `gorm:"unique;not null" json:"lang_key"`
	RegionKey   string    `gorm:"type:varchar(255);not null" json:"region_key"`
	IsPublished bool      `json:"is_published" default:"true"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
}
