package database

import (
	"fmt"

	"github.com/krisn2/go-social/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Initialize(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Enable SQL logging
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Auto-migrate models
	err = db.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Like{},
		&models.Comment{},
	)
	if err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	// Create indexes
	if err := createIndexes(db); err != nil {
		return nil, fmt.Errorf("failed to create indexes: %w", err)
	}

	return db, nil
}

func createIndexes(db *gorm.DB) error {
	// Unique composite index for likes
	if err := db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_likes_user_post ON likes (user_id, post_id)").Error; err != nil {
		return err
	}

	// Performance indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts (created_at DESC)").Error; err != nil {
		return err
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_comments_post_created ON comments (post_id, created_at)").Error; err != nil {
		return err
	}

	return nil
}
