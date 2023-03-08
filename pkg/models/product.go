package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name                  string
	Price                 uint
	ModelNo               uint
	Image1                string
	Image2                string
	Image3                string
	Stock                 uint
	CategoryID            uint
	Description           string
	Brand                 string
	Chipset_brand         string
	Model_gpu             string
	Series                string
	Generation            string
	Memmory_type          string
	Thermal_design_power  string
	Released              string
	Architecture          string
	Memmory_size          uint
	Recomented_resolution string
	DirectX               string
	Memmory_bus_width     string
	Production_status     string
	Text_mapping_unit     string
	Slots                 string
	Rops                  string
	Power_Connecters      string
	Discount              uint
}
