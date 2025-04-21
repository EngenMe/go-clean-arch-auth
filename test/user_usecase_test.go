package test

import (
	"testing"

	"github.com/EngenMe/go-clean-arch-auth/internal/data/entity"
	"github.com/EngenMe/go-clean-arch-auth/internal/useCase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id uint) (*entity.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*entity.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) List(page, limit int) ([]entity.User, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]entity.User), args.Error(1)
}

func TestUserUseCase_Register(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userUseCase := useCase.NewUserUseCase(mockRepo)

	user := &entity.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	mockRepo.On("GetByEmail", user.Email).Return(nil, nil).Once()
	mockRepo.On("Create", user).Return(nil).Once()

	err := userUseCase.Register(user)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)

	existingUser := &entity.User{
		ID:       1,
		Name:     "Existing User",
		Email:    "test@example.com",
		Password: "password123",
	}
	mockRepo.On("GetByEmail", user.Email).Return(existingUser, nil).Once()

	err = userUseCase.Register(user)
	assert.Error(t, err)
	assert.Equal(t, "email already exists", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_Login(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userUseCase := useCase.NewUserUseCase(mockRepo)

	testUser := &entity.User{
		ID:       1,
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "$2a$10$0Z7PtES8QhEcJCbx0uZ4c.1EagmkZ7EKU4K7AZJzMdVvI2ZS5KTy2", // bcrypt hash for "password123"
	}

	mockRepo.On("GetByEmail", "test@example.com").Return(testUser, nil).Once()

	_, err := userUseCase.Login("test@example.com", "wrong-password")
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)

	mockRepo.On("GetByEmail", "nonexistent@example.com").Return(
		nil,
		assert.AnError,
	).Once()

	_, err = userUseCase.Login("nonexistent@example.com", "password123")
	assert.Error(t, err)
	assert.Equal(t, "invalid email or password", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_GetUserByID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userUseCase := useCase.NewUserUseCase(mockRepo)

	testUser := &entity.User{
		ID:       1,
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "hashed_password",
	}

	mockRepo.On("GetByID", uint(1)).Return(testUser, nil).Once()

	user, err := userUseCase.GetUserByID(1)
	assert.NoError(t, err)
	assert.Equal(t, testUser, user)

	mockRepo.AssertExpectations(t)

	mockRepo.On("GetByID", uint(999)).Return(nil, assert.AnError).Once()

	_, err = userUseCase.GetUserByID(999)
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_UpdateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userUseCase := useCase.NewUserUseCase(mockRepo)

	existingUser := &entity.User{
		ID:       1,
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "hashed_password",
	}

	updatedUser := &entity.User{
		ID:       1,
		Name:     "Updated User",
		Email:    "test@example.com",
		Password: "",
	}

	mockRepo.On("GetByID", uint(1)).Return(existingUser, nil).Once()
	mockRepo.On(
		"Update",
		mock.AnythingOfType("*entity.User"),
	).Return(nil).Once()

	err := userUseCase.UpdateUser(updatedUser)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)

	mockRepo.On("GetByID", uint(999)).Return(nil, assert.AnError).Once()

	updatedUser.ID = 999
	err = userUseCase.UpdateUser(updatedUser)
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)

	newUser := &entity.User{
		ID:       1,
		Name:     "Test User",
		Email:    "new@example.com",
		Password: "",
	}

	existingUserWithSameEmail := &entity.User{
		ID:       2,
		Name:     "Another User",
		Email:    "new@example.com",
		Password: "hashed_password",
	}

	mockRepo.On("GetByID", uint(1)).Return(existingUser, nil).Once()
	mockRepo.On(
		"GetByEmail",
		"new@example.com",
	).Return(existingUserWithSameEmail, nil).Once()

	err = userUseCase.UpdateUser(newUser)
	assert.Error(t, err)
	assert.Equal(t, "email already exists", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_DeleteUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userUseCase := useCase.NewUserUseCase(mockRepo)

	testUser := &entity.User{
		ID:       1,
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "hashed_password",
	}

	mockRepo.On("GetByID", uint(1)).Return(testUser, nil).Once()
	mockRepo.On("Delete", uint(1)).Return(nil).Once()

	err := userUseCase.DeleteUser(1)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)

	mockRepo.On("GetByID", uint(999)).Return(nil, assert.AnError).Once()

	err = userUseCase.DeleteUser(999)
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_ListUsers(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userUseCase := useCase.NewUserUseCase(mockRepo)

	users := []entity.User{
		{ID: 1, Name: "User 1", Email: "user1@example.com"},
		{ID: 2, Name: "User 2", Email: "user2@example.com"},
		{ID: 3, Name: "User 3", Email: "user3@example.com"},
	}

	// Test successful list
	mockRepo.On("List", 1, 10).Return(users, nil).Once()

	userResponses, err := userUseCase.ListUsers(1, 10)
	assert.NoError(t, err)
	assert.Len(t, userResponses, 3)
	assert.Equal(t, users[0].ID, userResponses[0].ID)
	assert.Equal(t, users[0].Name, userResponses[0].Name)
	assert.Equal(t, users[0].Email, userResponses[0].Email)

	mockRepo.AssertExpectations(t)

	mockRepo.On("List", 1, 10).Return(users, nil).Once()

	userResponses, err = userUseCase.ListUsers(-1, -5)
	assert.NoError(t, err)
	assert.Len(t, userResponses, 3)

	mockRepo.AssertExpectations(t)
}
