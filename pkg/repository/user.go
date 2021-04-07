package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ismtabo/poc-notification-system/pkg/model"
	"github.com/ismtabo/poc-notification-system/pkg/service"
)

type UserRepository interface {
	Create(string) (*model.User, error)
	Read(model.ID) (*model.User, error)
	Update(*model.User) (*model.User, error)
	Delete(model.ID) error
}

var users map[model.ID]*model.User = make(map[model.ID]*model.User)

type userRepository struct {
	notifications service.NotificationService
}

func NewUserRepository(notificationService service.NotificationService) UserRepository {
	return &userRepository{notifications: notificationService}
}

func (r *userRepository) Create(name string) (*model.User, error) {
	var id model.ID = model.ID(uuid.NewString())
	now := time.Now()
	user := model.User{ID: id, Name: name, CreatedAt: &now}
	notification := &model.Notification{}
	if err := notification.MarshalCurrent(user); err != nil {
		return nil, err
	}
	if err := r.notifications.Publish("user", notification); err != nil {
		return nil, err
	}
	users[id] = &user
	return &user, nil
}

func (r *userRepository) Read(id model.ID) (*model.User, error) {
	user, found := users[id]
	if !found {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *userRepository) Update(user *model.User) (*model.User, error) {
	now := time.Now()
	oldUser, err := r.Read(user.ID)
	if err != nil {
		return nil, err
	}
	user.UpdatedAt = &now
	notification := &model.Notification{}
	if err := notification.MarshalOld(oldUser); err != nil {
		return nil, err
	}
	if err := notification.MarshalCurrent(user); err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	if err := r.notifications.Publish("user", notification); err != nil {
		return nil, err
	}
	users[user.ID] = user
	return user, nil
}

func (r *userRepository) Delete(id model.ID) error {
	oldUser, err := r.Read(id)
	if err != nil {
		return err
	}
	notification := &model.Notification{}
	if err := notification.MarshalOld(oldUser); err != nil {
		return err
	}
	if err != nil {
		return err
	}
	if err := r.notifications.Publish("user", notification); err != nil {
		return err
	}
	return nil
}
