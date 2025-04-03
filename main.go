package main

import (
	"encoding/json"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

type Country struct {
	ID             uint   `gorm:"primaryKey"`
	CommonName     string `gorm:"index"`
	OfficialName   string
	FlagPNG        string
	FlagSVG        string
	FlagAlt        string
	NativeCommon   string
	NativeOfficial string
}

type CountryRaw struct {
	Flags struct {
		Png string `json:"png"`
		Svg string `json:"svg"`
		Alt string `json:"alt"`
	} `json:"flags"`
	Name struct {
		Common     string `json:"common"`
		Official   string `json:"official"`
		NativeName map[string]struct {
			Common   string `json:"common"`
			Official string `json:"official"`
		} `json:"nativeName"`
	} `json:"name"`
}

func main() {
	db, _ := gorm.Open(sqlite.Open("countries.db"), &gorm.Config{})
	db.AutoMigrate(&Country{})

	file, _ := os.ReadFile("countries.json")
	var raw []CountryRaw
	json.Unmarshal(file, &raw)

	for _, c := range raw {
		native := c.Name.NativeName["eng"] // можно добавить fallback

		country := Country{
			CommonName:     c.Name.Common,
			OfficialName:   c.Name.Official,
			FlagPNG:        c.Flags.Png,
			FlagSVG:        c.Flags.Svg,
			FlagAlt:        c.Flags.Alt,
			NativeCommon:   native.Common,
			NativeOfficial: native.Official,
		}
		db.Create(&country)
	}
}
