package main

import (
	"log"
	"time"

	"github.com/ismtabo/poc-notification-system/pkg/model"
	"github.com/ismtabo/poc-notification-system/pkg/repository"
	"github.com/ismtabo/poc-notification-system/pkg/service"
	"github.com/jinzhu/copier"
)

func main() {
	notificationService := service.NewChannelNotificationService()
	userRepository := repository.NewUserRepository(notificationService)
	groupRepository := repository.NewGroupRepository(notificationService)
	belongingRepository := repository.NewBelongingRepository(notificationService)

	_, err := notificationService.Subscribe(func(topic string, notification *model.Notification) {
		if topic == "user" && (notification.Type() == model.InsertNotification || notification.Type() == model.UpdateNotification) {
			var user model.User
			notification.UnmarshalCurrent(&user)
			log.Printf("Changes with user '%s' may be an insert or an update: %s", user.ID, user)
		}

		if topic == "belonging" && notification.Type() == model.DeleteNotification {
			var belonging model.Belonging
			notification.UnmarshalOld(&belonging)
			log.Printf("Belonging delete with userId '%s' and groupId '%s' created at %s", belonging.UserID, belonging.GroupID, belonging.CreatedAt)
		}
	})
	if err != nil {
		log.Fatal(err)
	}

	// Users
	user, err := userRepository.Create("José")
	if err != nil {
		log.Fatal(err)
	}

	userCpy := model.User{}
	if err := copier.Copy(&userCpy, user); err != nil {
		log.Fatal(err)
	}
	userCpy.Name = "Pepe"
	updatedUser, err := userRepository.Update(&userCpy)
	if err != nil {
		log.Fatal(err)
	}

	if err := userRepository.Delete(updatedUser.ID); err != nil {
		log.Fatal(err)
	}

	// Groups
	group, err := groupRepository.Create("DevOps")
	if err != nil {
		log.Fatal(err)
	}

	groupCpy := model.Group{}
	if err := copier.Copy(&groupCpy, group); err != nil {
		log.Fatal(err)
	}
	groupCpy.Name = "SRE"
	updatedGroup, err := groupRepository.Update(&groupCpy)
	if err != nil {
		log.Fatal(err)
	}

	if err := groupRepository.Delete(updatedGroup.ID); err != nil {
		log.Fatal(err)
	}

	// Belongings
	if _, err := belongingRepository.Create(user, group); err != nil {
		log.Fatal(err)
	}

	if err := belongingRepository.Delete(user, group); err != nil {
		log.Fatal(err)
	}

	time.Sleep(1 * time.Second)
}
