package model

import "encoding/json"

type NotificationType string

const (
	InsertNotification  NotificationType = NotificationType("insert")
	UpdateNotification  NotificationType = NotificationType("update")
	DeleteNotification  NotificationType = NotificationType("delete")
	UnknownNotification NotificationType = NotificationType("unknown")
)

type Notification struct {
	Old     json.RawMessage
	Current json.RawMessage
}

func (n *Notification) Type() NotificationType {
	if string(n.Old) == "null" && string(n.Current) != "null" {
		return InsertNotification
	} else if string(n.Old) != "null" && string(n.Current) == "null" {
		return DeleteNotification
	} else if string(n.Old) != "null" && string(n.Current) != "null" {
		return UpdateNotification
	}
	return UnknownNotification
}

func (n *Notification) MarshalOld(old interface{}) error {
	data, err := json.Marshal(old)
	if err != nil {
		return err
	}
	n.Old.UnmarshalJSON(data)
	return nil
}

func (n *Notification) UnmarshalOld(v interface{}) error {
	return json.Unmarshal(n.Old, v)
}

func (n *Notification) MarshalCurrent(current interface{}) error {
	data, err := json.Marshal(current)
	if err != nil {
		return err
	}
	n.Current.UnmarshalJSON(data)
	return nil
}

func (n *Notification) UnmarshalCurrent(v interface{}) error {
	return json.Unmarshal(n.Current, v)
}
