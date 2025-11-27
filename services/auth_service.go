package services

import (
	"errors"
	"ocr-api/config"
	"ocr-api/models"
	"ocr-api/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Register(username, email, password, fullName string) (*models.User, error) {
	db := config.GetDB()

	var existingUser models.User
	if err := db.Where("username = ? OR email = ?", username, email).First(&existingUser).Error; err == nil {
		return nil, errors.New("username or email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		FullName: fullName,
	}

	if err := db.Create(user).Error; err != nil {
		return nil, errors.New("failed to create user")
	}

	return user, nil
}

func (s *AuthService) Login(username, password string) (*models.User, string, error) {
	db := config.GetDB()

	var user models.User
	if err := db.Where("username = ? OR email = ?", username, username).First(&user).Error; err != nil {
		return nil, "", errors.New("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("invalid username or password")
	}

	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, "", errors.New("failed to generate token")
	}

	return &user, token, nil
}

func (s *AuthService) GetUserByID(userID uint) (*models.User, error) {
	db := config.GetDB()

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}
