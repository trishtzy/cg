package main

import "gorm.io/gorm"

type Country struct {
	gorm.Model
	Label string `json:"label"`
	Value string `json:"value"`
}
