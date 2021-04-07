package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ismtabo/poc-notification-system/pkg/model"
	"github.com/ismtabo/poc-notification-system/pkg/service"
)

type GroupRepository interface {
	Create(string) (*model.Group, error)
	Read(model.ID) (*model.Group, error)
	Update(*model.Group) (*model.Group, error)
	Delete(model.ID) error
}

var groups map[model.ID]*model.Group = make(map[model.ID]*model.Group)

type groupRepository struct {
	notifications service.NotificationService
}

func NewGroupRepository(notificationService service.NotificationService) GroupRepository {
	return &groupRepository{notifications: notificationService}
}

func (r *groupRepository) Create(name string) (*model.Group, error) {
	var id model.ID = model.ID(uuid.NewString())
	now := time.Now()
	group := model.Group{ID: id, Name: name, CreatedAt: &now}
	notification := &model.Notification{}
	if err := notification.MarshalCurrent(group); err != nil {
		return nil, err
	}
	if err := r.notifications.Publish("group", notification); err != nil {
		return nil, err
	}
	groups[id] = &group
	return &group, nil
}

func (r *groupRepository) Read(id model.ID) (*model.Group, error) {
	group, found := groups[id]
	if !found {
		return nil, errors.New("group not found")
	}
	return group, nil
}

func (r *groupRepository) Update(group *model.Group) (*model.Group, error) {
	now := time.Now()
	oldGroup, err := r.Read(group.ID)
	if err != nil {
		return nil, err
	}
	group.UpdatedAt = &now
	notification := &model.Notification{}
	if err := notification.MarshalOld(oldGroup); err != nil {
		return nil, err
	}
	if err := notification.MarshalCurrent(group); err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	if err := r.notifications.Publish("group", notification); err != nil {
		return nil, err
	}
	groups[group.ID] = group
	return group, nil
}

func (r *groupRepository) Delete(id model.ID) error {
	oldGroup, err := r.Read(id)
	if err != nil {
		return err
	}
	notification := &model.Notification{}
	if err := notification.MarshalOld(oldGroup); err != nil {
		return err
	}
	if err != nil {
		return err
	}
	if err := r.notifications.Publish("group", notification); err != nil {
		return err
	}
	return nil
}
