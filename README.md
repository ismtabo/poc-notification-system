# POC - Heterogeneous notification message

This repository consists of a proof of concept of a system of notifications capable of notifying changes in entities of multiple types.

These notifications have the following JSON Schema:
```json
{
    "old": null | any,
    "current": null | any
}
```

The aim of this notifications is to be able to notify three types of notification types:
- Insert: when new `current` data is created and there is no `old` data
- Update: when there are both `old` and `current` data
- Delete: when there is `old` data but no `current`

To resume, the following type shows the different types of notifications depending on its properties:

|  Type  | `old`  | `current` |
| :----: | :----: | :-------: |
| Insert | `null` |    any    |
| Update |  any   |    any    |
| Delete |  any   |  `null`   |

Also the purpose of this repository is to be able to notify these changes independently of the type of entity involved in the change. For these POC, there have been included three types:
- User: representing a possible user of a system with an identifier, name, and creation and update dates
- Group: representing a possible group of users of a system with an identifier, name, and creation and update dates
- Belonging: representing a relation between users and groups with a creation date.

## Prerequisites
To run this repository you should have installed Golang >1.16. Use `go mod tidy` to update the local dependencies.

## Usage
To run this project run the following command:
```
# Inside cmd/poc-notification-system folder
go run main.go
```

## Folders

```
.
├── cmd/poc-notification-system    Main package
├── pkg
│   ├── model                      Models of the repository
│   ├── repository                 Data repositories
│   └── service                    Notification service
├── go.mod
├── go.sum
└── README.md
```

### cmd/poc-notification-system
This folder contains the main package with the main function to run. This function will follow the structure bellow:
1. Initialize notification service and all the repositories
2. Subscribe to changes in the notification service
   1. The subscription consists of a callback function which logs notifications for two cases:
      1. Insert and Update notifications in Users
      2. Delete notifications of Belongings
3. Do CRUD operations over User entities
4. Do CRUD operations over Group entities
5. Do CRUD operations over Belonging entities


### pkg/model
This folder contains all the structs of the application. This structs are:

#### Notification
This struct represent a change notification. As said in the introduction this struct has two properties `Old` and `Current`. Both of them are generic types using `json.RawJson` type from `encoding/json` package, so later both properties can be unmarshal to proper types.

To add content to the notification, the methods `MarshalOld` and `MarshalCurrent` should be used, according to the table showed before. In the other hand, to unmarshal any of the properties you may use `UnmarshalOld` and `UnmasrhalCurrent` instead.

```go
var entity Type
notification := &model.Notification{}
err := notification.MarshalCurrent(&entity)

err = notification.UnmarshalCurrent(&entity)
```

Also this struct has a method `Type` to return its type depending on the state of its properties as commented before.

#### Entity types
As commented before the `model` package contains example entity types such as User, Group and Belonging.

### repository
This folder contains the data repositories for each type of entity. Each repository has the basic CRUD operations, unless the Belonging repository that has not update method. All the repositories have a dependency with the notification services, as each of them will publish a message into the proper topic (`user`, `group` and `belonging`) when an insert, update or delete operation is done.

### service
This folder contains the notification service. This service internally use a Message struct to store the information about the topic and the notification that is needed to be publish. This struct is place in the folder `pkg/service/dto`. Also, this service consists of two functions:
- Publish: given a topic and a notification, a message is sent to the internal channel of the service.
- Subscribe: given a callback function with a topic and a notification as parameters, the service run a goroutine which runs forever and execute the callback function for each message that arrives to then internal channel. This method returns an input channel to stop/cancel the infinite loop.

## License

Distributed under the MIT License. See `LICENSE` for more information.

## Contact

Ismael Taboada - [ismtabo](https://github.com/ismtabo)

Project Link: [https://github.com/your_username/repo_name](https://github.com/your_username/repo_name)