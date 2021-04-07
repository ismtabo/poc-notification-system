package service

import (
	"encoding/json"

	"github.com/ismtabo/poc-notification-system/pkg/model"
	"github.com/ismtabo/poc-notification-system/pkg/service/dto"
)

type NotificationService interface {
	Publish(string, *model.Notification) error
	Subscribe(func(string, *model.Notification)) (chan<- bool, error)
}

type channelNotificationService struct {
	ch chan dto.Message
}

func NewChannelNotificationService() NotificationService {
	return &channelNotificationService{ch: make(chan dto.Message)}
}

func (n *channelNotificationService) Publish(topic string, notification *model.Notification) error {
	data, err := json.Marshal(notification)
	if err != nil {
		return err
	}
	message := dto.Message{Topic: topic, Data: data}
	n.ch <- message
	return nil
}

func (n *channelNotificationService) Subscribe(cb func(string, *model.Notification)) (chan<- bool, error) {
	cancel := make(chan bool)
	go func() {
	loop:
		for {
			select {
			case message := <-n.ch:
				var notification model.Notification
				json.Unmarshal(message.Data, &notification)
				cb(message.Topic, &notification)
			case <-cancel:
				break loop
			}
		}
	}()
	return cancel, nil
}
