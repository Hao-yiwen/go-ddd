package aggregate

import (
	"yiwen/go-ddd/internal/domain/entity"
	"yiwen/go-ddd/internal/domain/event"
	"yiwen/go-ddd/internal/domain/valueobject"
)

// UserAggregate 用户聚合根
// 聚合根是DDD中的核心概念
// 1. 聚合是一组相关对象的集合
// 2. 聚合根是聚合的入口点，外部只能通过聚合根访问聚合内的对象
// 3. 聚合根负责维护聚合内的一致性
// 4. 聚合根可以发布领域事件
type UserAggregate struct {
	User   *entity.User
	Events []event.Event
}

func NewUserAggregate(user *entity.User) *UserAggregate {
	return &UserAggregate{
		User:   user,
		Events: []event.Event{},
	}
}

// Register 注册用户
func Register(uuid, username string, email valueobject.Email, password valueobject.Password) *UserAggregate {
	user := entity.NewUser(uuid, username, email, password)
	agg := NewUserAggregate(user)
	agg.addEvent(event.NewUserRegisteredEvent(uuid, username, email.String()))

	return agg
}

// UpdateProfile 更新用户资料
func (a *UserAggregate) UpdateProfile(nickname, avatar string) {
	oldNickname := a.User.Nickname
	a.User.UpdateProfile(nickname, avatar)

	a.addEvent(event.NewUserProfileUpdatedEvent(a.User.UUID, oldNickname, nickname))
}

// ChangePassword 修改密码
func (a *UserAggregate) ChangePassword(newPassword valueobject.Password) {
	a.User.ChangePassword(newPassword)

	a.addEvent(event.NewUserPasswordChangedEvent(a.User.UUID))
}

// Activate 激活用户
func (a *UserAggregate) Activate() {
	if a.User.Status == entity.UserStatusActive {
		return
	}

	a.User.Activate()
	a.addEvent(event.NewUserActivatedEvent(a.User.UUID))
}

// Deactivate 禁用用户
func (a *UserAggregate) Deactivate() {
	if a.User.Status == entity.UserStatusInactive {
		return
	}

	a.User.Deactivate()
	a.addEvent(event.NewUserDeactivatedEvent(a.User.UUID))
}

func (a *UserAggregate) Ban(reason string) {
	if a.User.Status == entity.UserStatusBanned {
		return
	}

	a.User.Ban()
	a.addEvent(event.NewUserBannedEvent(a.User.UUID, reason))
}

func (a *UserAggregate) PromoteToAdmin() {
	if a.User.IsAdmin() {
		return
	}

	a.User.PromoteToAdmin()
	a.addEvent(event.NewUserPromotedEvent(a.User.UUID))
}

func (a *UserAggregate) addEvent(e event.Event) {
	a.Events = append(a.Events, e)
}

func (a *UserAggregate) ClearEvents() {
	a.Events = make([]event.Event, 0)
}

func (a *UserAggregate) GetEvents() []event.Event {
	return a.Events
}
