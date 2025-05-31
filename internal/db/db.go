package db

import (
	"fmt"
	"path/filepath"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

// DB handles database operations
type DB struct {
	conn *gorm.DB
}

// NewDB creates a new database connection
func NewDB(path string) (*DB, error) {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	// Open database connection
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run migrations
	if err := db.AutoMigrate(&User{}, &Subscription{}); err != nil {
		return nil, fmt.Errorf("failed to run database migrations: %w", err)
	}

	return &DB{conn: db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	sqlDB, err := db.conn.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	return sqlDB.Close()
}

// User operations

// AddUser adds a new user to the database
func (db *DB) AddUser(user *User) error {
	result := db.conn.Create(user)
	return result.Error
}

// GetUser gets a user by ID
func (db *DB) GetUser(id int) (*User, error) {
	var user User
	result := db.conn.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetUserByEmail gets a user by email
func (db *DB) GetUserByEmail(email string) (*User, error) {
	var user User
	result := db.conn.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// UpdateUser updates a user in the database
func (db *DB) UpdateUser(user *User) error {
	user.UpdatedAt = time.Now()
	result := db.conn.Save(user)
	return result.Error
}

// Subscription operations

// AddSubscription adds a new subscription to the database
func (db *DB) AddSubscription(subscription *Subscription) error {
	result := db.conn.Create(subscription)
	return result.Error
}

// GetSubscription gets a subscription by ID
func (db *DB) GetSubscription(id int) (*Subscription, error) {
	var subscription Subscription
	result := db.conn.First(&subscription, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &subscription, nil
}

// UpdateSubscription updates a subscription in the database
func (db *DB) UpdateSubscription(subscription *Subscription) error {
	subscription.UpdatedAt = time.Now()
	result := db.conn.Save(subscription)
	return result.Error
}

// DeleteSubscription removes a subscription from the database
func (db *DB) DeleteSubscription(id int) error {
	result := db.conn.Delete(&Subscription{}, id)
	return result.Error
}

// ListSubscriptions gets all subscriptions
func (db *DB) ListSubscriptions() ([]Subscription, error) {
	var subscriptions []Subscription
	result := db.conn.Find(&subscriptions)
	return subscriptions, result.Error
}

// ListActiveSubscriptions gets all active subscriptions
func (db *DB) ListActiveSubscriptions() ([]Subscription, error) {
	var subscriptions []Subscription
	result := db.conn.Where("active = ?", true).Find(&subscriptions)
	return subscriptions, result.Error
}

// GetSubscriptionsByUserID gets all subscriptions for a user
func (db *DB) GetSubscriptionsByUserID(userID int) ([]Subscription, error) {
	var subscriptions []Subscription
	result := db.conn.Where("user_id = ?", userID).Find(&subscriptions)
	return subscriptions, result.Error
}