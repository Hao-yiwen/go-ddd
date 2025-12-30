package event

import "time"

// Event 领域事件
// 领域事件是DDD中的重要概念
// 1. 表示领域中发生的有意义的事情
// 2. 事件是不可变的 发生过的事情不能改变
// 3. 可以用于解耦不同的领域/服务
// 4. 支持事件溯源 event sourcing
type Event interface {
	EventName() string
	OccurredAt() time.Time
	AggregateID() string
}

type BaseEvent struct {
	Name        string    `json:"name"`
	OccurredOn  time.Time `json:"occurred_at"`
	AggregateId string    `json:"aggregate_id"`
}

func (e BaseEvent) EventName() string {
	return e.Name
}

func (e BaseEvent) OccurredAt() time.Time {
	return e.OccurredOn
}

func (e BaseEvent) AggregateID() string {
	return e.AggregateId
}

// UserRegisteredEvent 用户注册事件
type UserRegisteredEvent struct {
	BaseEvent
	UserName string `json:"user_name"`
	Email    string `json:"email"`
}

func NewUserRegisteredEvent(uuid, username, email string) *UserRegisteredEvent {
	return &UserRegisteredEvent{
		BaseEvent: BaseEvent{
			Name:        "UserRegistered",
			OccurredOn:  time.Now(),
			AggregateId: uuid,
		},
		UserName: username,
		Email:    email,
	}
}

type UserProfileUpdatedEvent struct {
	BaseEvent
	OldNickname string `json:"old_nickname"`
	NewNickname string `json:"new_nickname"`
}

func NewUserProfileUpdatedEvent(uuid, oldNickname, newNickname string) *UserProfileUpdatedEvent {
	return &UserProfileUpdatedEvent{
		BaseEvent: BaseEvent{
			Name:        "UserProfileUpdated",
			OccurredOn:  time.Now(),
			AggregateId: uuid,
		},
		OldNickname: oldNickname,
		NewNickname: newNickname,
	}
}

type UserPasswordChangedEvent struct {
	BaseEvent
}

func NewUserPasswordChangedEvent(uuid string) *UserPasswordChangedEvent {
	return &UserPasswordChangedEvent{
		BaseEvent: BaseEvent{
			Name:        "UserPasswordChanged",
			OccurredOn:  time.Now(),
			AggregateId: uuid,
		},
	}
}

type UserActivatedEvent struct {
	BaseEvent
}

func NewUserActivatedEvent(uuid string) *UserActivatedEvent {
	return &UserActivatedEvent{
		BaseEvent: BaseEvent{
			Name:        "UserActivated",
			OccurredOn:  time.Now(),
			AggregateId: uuid,
		},
	}
}

type UserDeactivatedEvent struct {
	BaseEvent
}

func NewUserDeactivatedEvent(uuid string) *UserDeactivatedEvent {
	return &UserDeactivatedEvent{
		BaseEvent: BaseEvent{
			Name:        "UserDeactivated",
			OccurredOn:  time.Now(),
			AggregateId: uuid,
		},
	}
}

// UserBannedEvent 用户禁用事件
type UserBannedEvent struct {
	BaseEvent
	Reason string `json:"reason"`
}

func NewUserBannedEvent(uuid, reason string) *UserBannedEvent {
	return &UserBannedEvent{
		BaseEvent: BaseEvent{
			Name:        "user.banned",
			OccurredOn:  time.Now(),
			AggregateId: uuid,
		},
		Reason: reason,
	}
}

// UserPromotedEvent 用户提升为管理员事件
type UserPromotedEvent struct {
	BaseEvent
}

func NewUserPromotedEvent(uuid string) *UserPromotedEvent {
	return &UserPromotedEvent{
		BaseEvent: BaseEvent{
			Name:        "user.promoted",
			OccurredOn:  time.Now(),
			AggregateId: uuid,
		},
	}
}

type EventHandler interface {
	Handle(event Event) error
}

type EventPublisher interface {
	Publish(event Event) error
}
