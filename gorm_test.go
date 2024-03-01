package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGetStudentCountByDepartmentID(t *testing.T) {
	// Setup the in-memory SQLite database for testing
	dsn := "root:possible04@tcp(localhost:3306)/golang"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("error")
	}
	// Migrate the schema
	db.AutoMigrate(&Student{})

	// Test the function
	count := GetStudentCountByDepartmentID(db, 1)
	assert.Equal(t, int64(2), count, "should be equal")
}
