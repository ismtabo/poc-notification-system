package repository

import (
	"errors"
	"time"

	"github.com/golang-collections/collections/set"
	"github.com/ismtabo/poc-notification-system/pkg/model"
	"github.com/ismtabo/poc-notification-system/pkg/service"
)

type BelongingRepository interface {
	Create(*model.User, *model.Group) (*model.Belonging, error)
	Read(*model.User, *model.Group) (*model.Belonging, error)
	Delete(*model.User, *model.Group) error
}

var belongings *set.Set = set.New()

type belongingRepository struct {
	notifications service.NotificationService
}

func NewBelongingRepository(notificationService service.NotificationService) BelongingRepository {
	return &belongingRepository{notifications: notificationService}
}

func (r *belongingRepository) Create(user *model.User, group *model.Group) (*model.Belonging, error) {
	now := time.Now()
	belonging := &model.Belonging{UserID: user.ID, GroupID: group.ID, CreatedAt: &now}
	notification := &model.Notification{}
	if err := notification.MarshalCurrent(belonging); err != nil {
		return nil, err
	}
	if err := r.notifications.Publish("belonging", notification); err != nil {
		return nil, err
	}
	belongings.Insert(belonging)
	return belonging, nil
}

func (r *belongingRepository) Read(user *model.User, group *model.Group) (*model.Belonging, error) {
	belonging := findBelonging(user.ID, group.ID)
	if belonging == nil {
		return nil, errors.New("belonging not found")
	}
	return belonging, nil
}

func (r *belongingRepository) Delete(user *model.User, group *model.Group) error {
	oldBelonging, err := r.Read(user, group)
	if err != nil {
		return err
	}
	notification := &model.Notification{}
	if err := notification.MarshalOld(oldBelonging); err != nil {
		return err
	}
	if err := r.notifications.Publish("belonging", notification); err != nil {
		return err
	}
	return nil
}

func findBelonging(userID model.ID, groupID model.ID) *model.Belonging {
	var belonging *model.Belonging
	belongings.Do(func(b interface{}) {
		other := b.(*model.Belonging)
		if other.UserID == userID && other.GroupID == groupID {
			belonging = other
		}

	})
	return belonging
}
