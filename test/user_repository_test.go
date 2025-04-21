package test

import (
	"github.com/EngenMe/go-clean-arch-auth/internal/Infrastructure/repository/respositories"
	"testing"
	"time"

	"github.com/EngenMe/go-clean-arch-auth/internal/data/entity"
	"github.com/EngenMe/go-clean-arch-auth/pkg/database"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestDB() (*gorm.DB, error) {
	return database.NewPostgresDB()
}

func TestUserRepository_Create(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	tx := db.Begin()
	defer tx.Rollback()

	repo := respositories.NewUserRepository(tx)

	testUser := &entity.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	err = repo.Create(testUser)
	assert.NoError(t, err)
	assert.NotZero(t, testUser.ID)
	assert.False(t, testUser.CreatedAt.IsZero())
}

func TestUserRepository_GetByID(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	tx := db.Begin()
	defer tx.Rollback()

	repo := respositories.NewUserRepository(tx)

	testUser := &entity.User{
		Name:      "Test User",
		Email:     "test@example.com",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = repo.Create(testUser)
	assert.NoError(t, err)

	user, err := repo.GetByID(testUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, testUser.ID, user.ID)
	assert.Equal(t, testUser.Name, user.Name)
	assert.Equal(t, testUser.Email, user.Email)
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	tx := db.Begin()
	defer tx.Rollback()

	repo := respositories.NewUserRepository(tx)

	testUser := &entity.User{
		Name:      "Test User",
		Email:     "test@example.com",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = repo.Create(testUser)
	assert.NoError(t, err)

	user, err := repo.GetByEmail(testUser.Email)
	assert.NoError(t, err)
	assert.Equal(t, testUser.ID, user.ID)
	assert.Equal(t, testUser.Email, user.Email)
}

func TestUserRepository_Update(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	tx := db.Begin()
	defer tx.Rollback()

	repo := respositories.NewUserRepository(tx)

	testUser := &entity.User{
		Name:      "Test User",
		Email:     "test@example.com",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = repo.Create(testUser)
	assert.NoError(t, err)

	testUser.Name = "Updated Name"
	err = repo.Update(testUser)
	assert.NoError(t, err)

	updatedUser, err := repo.GetByID(testUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", updatedUser.Name)
}

func TestUserRepository_Delete(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	tx := db.Begin()
	defer tx.Rollback()

	repo := respositories.NewUserRepository(tx)

	testUser := &entity.User{
		Name:      "Test User",
		Email:     "test@example.com",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = repo.Create(testUser)
	assert.NoError(t, err)

	err = repo.Delete(testUser.ID)
	assert.NoError(t, err)

	_, err = repo.GetByID(testUser.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

func TestUserRepository_List(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	tx := db.Begin()
	defer tx.Rollback()

	repo := respositories.NewUserRepository(tx)

	for i := 0; i < 5; i++ {
		testUser := &entity.User{
			Name:      "Test User " + string(rune('A'+i)),
			Email:     "test" + string(rune('a'+i)) + "@example.com",
			Password:  "password123",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err = repo.Create(testUser)
		assert.NoError(t, err)
	}

	users, err := repo.List(1, 3)
	assert.NoError(t, err)
	assert.Len(t, users, 3)

	users, err = repo.List(2, 3)
	assert.NoError(t, err)
	assert.Len(t, users, 2)
}
