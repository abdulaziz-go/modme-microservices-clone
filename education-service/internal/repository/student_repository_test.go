package repository

import (
	"education-service/config"
	"testing"
)

func TestHappyBirthdayAlert(t *testing.T) {
	// This test only verifies that the function runs without panic.
	// Further mocking would be required to test SMS delivery or DB interaction.

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("HappyBirthdayAlert panicked: %v", r)
		}
	}()
	db, err := NewPostgresDB(&config.DatabaseConfig{
		Host:     "158.220.111.34",
		Port:     9015,
		User:     "postgres",
		Password: "password",
		DBName:   "sphere_education_db",
		SSLMode:  "disable",
	})
	if err != nil {
		panic("erorr while connecting db")
		return
	}
	repo := &StudentRepository{
		db: db,
	}
	repo.HappyBirthdayAlert()
}
