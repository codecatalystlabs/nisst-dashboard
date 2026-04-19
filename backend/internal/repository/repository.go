package repository

import "gorm.io/gorm"

type Registry struct { DB *gorm.DB }

func NewRegistry(db *gorm.DB) *Registry { return &Registry{DB: db} }
