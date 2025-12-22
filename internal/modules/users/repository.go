package users

import "gorm.io/gorm"

type Repository interface {
	Create(user *User) error
	FindById(id string) (*User, error)
	SoftDelete(id string) error
	Restore(id string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *repository) FindById(id string) (*User, error) {
	var user User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) SoftDelete(id string) error {
	return r.db.Delete(&User{}, "id = ?", id).Error
}

func (r *repository) Restore(id string) error {
	return r.db.Unscoped().
		Model(&User{}).
		Where("id = ?", id).
		Update("deleted_at", nil).
		Error
}
