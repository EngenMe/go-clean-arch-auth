package mediator

import (
	"errors"
	"github.com/EngenMe/go-clean-arch-auth/internal/core/features/user/command/handler"
	"github.com/EngenMe/go-clean-arch-auth/internal/core/features/user/command/requests"
	"github.com/EngenMe/go-clean-arch-auth/internal/useCase"
	"reflect"
)

type Handler interface {
	Handle(request interface{}) (interface{}, error)
}

type Mediator struct {
	handlers map[reflect.Type]Handler
}

func NewMediator() *Mediator {
	return &Mediator{handlers: make(map[reflect.Type]Handler)}
}

func (m *Mediator) Register(request interface{}, handler Handler) {
	requestType := reflect.TypeOf(request)
	m.handlers[requestType] = handler
}

func (m *Mediator) Send(request interface{}) (interface{}, error) {
	requestType := reflect.TypeOf(request)
	reqHandler, exists := m.handlers[requestType]
	if !exists {
		return nil, errors.New("no handler registered for request")
	}
	return reqHandler.Handle(request)
}

// TODO: bug
func (m *Mediator) RegisterUserHandlers(userUseCase useCase.UserUseCase) {
	m.Register(
		&requests.RegisterUserRequest{},
		handler.NewRegisterUserHandler(userUseCase),
	)
}
