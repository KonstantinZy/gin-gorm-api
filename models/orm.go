package models

import (
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	modelMu   sync.RWMutex
	ormModels = map[string]interface{}{}
	DB        *gorm.DB
)

// open/create database for storing data
func StartDB(dbFilename string) {
	dbLocal, err := gorm.Open(sqlite.Open(dbFilename), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	DB = dbLocal
	RegisterModels()
}

// register all models
func RegisterModels() {
	AddModelToORM("task", Task{})
	AddModelToORM("task.subtask", SubTask{})
}

// adding model for migration to DB
func AddModelToORM(name string, model interface{}) {
	modelMu.Lock()
	defer modelMu.Unlock()

	if _, exists := ormModels[name]; exists {
		panic("orm: Adding model twice: " + name)
	}

	ormModels[name] = model
}

// migrate all registered models
func MigrateOrm() {
	for _, model := range ormModels {
		DB.AutoMigrate(model)
	}
}
