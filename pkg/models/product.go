package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name                  string `json:"name"`
	Price                 int    `json:"price"`
	ModelNo               int`json:""`
	Image1                string`json:""`
	Image2                string`json:""`
	Image3                string`json:""`
	Stock                 int `json:"stock"`
	CategoryID            int`json:""`
	SubCategoryID         int`json:""`
	Description           string `json:"description"`
	Brand                 string `json:"brand"`
	Chipset_brand         string `json:"chipset_brand"`
	Model_gpu             string `json:"model_gpu"`
	Series                string `json:"series"`
	Generation            string `json:"generation"`
	Memmory_type          string `json:"memmory_type"`
	Thermal_design_power  string `json:"thermal_design_power"`
	Released              string `json:"released"`
	Architecture          string `json:"architecture"`
	Memmory_size          int    `json:"memmory_size"`
	Recomented_resolution string `json:"recomented_resolution"`
	DirectX               string `json:"directx"`
	Memmory_bus_width     string `json:"memmory_bus_width"`
	Production_status     string `json:"production_status"`
	Text_mapping_unit     string `json:"text_mapping unit"`
	Slots                 string `json:"slots"`
	Rops                  string `json:"rops"`
	Power_Connecters      string `json:"powerconnecters"`
	Discount              int    `json:"discount"`
	Discount_Price        int    `json:"discountprice"`
}
