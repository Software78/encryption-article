package repository

import (
	db "github.com/Software78/encryption-test/src/db"
	"github.com/Software78/encryption-test/src/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(user *models.User) error
	Login(login *models.Login) (*models.User, error)
	Register(register *models.Register) (*models.User, error)
	GetUserByID(id uuid.UUID) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
}

// Concrete implementation
type userRepository struct {
	db db.Database
}

func NewUserRepository(db db.Database) UserRepository {
	return &userRepository{
		db: db,
	}
}



func (r *userRepository) Create(user *models.User) error {
	user.ID = uuid.New()
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 15)
	user.Password = string(hash)
	if err := r.db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) Login(login *models.Login) (*models.User, error) {
	user := &models.User{}
	if err := r.db.Where("email = ?", login.Email).First(&user).Error; err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Register(register *models.Register) (*models.User, error) {
	user := &models.User{
		FirstName: register.FirstName,
		LastName:  register.LastName,
		Email:     register.Email,
		Password:  register.Password,
		ID: 	  uuid.New(),
		
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 15)
	user.Password = string(hash)
	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}