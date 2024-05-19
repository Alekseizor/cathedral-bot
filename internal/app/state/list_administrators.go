package state

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
)

type AddDocumentAdministratorState struct {
	postgres *postrgres.Repo
}

func (state AddDocumentAdministratorState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	if messageText == "Назад" {
		return documentCabinet, nil, nil
	}

	vkID, err := strconv.Atoi(messageText)
	if err != nil || vkID < 100000000 || vkID >= 1000000000 {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("VK ID должно быть числом, например, 221486551")
		return addDocumentAdministrator, []*params.MessagesSendBuilder{b}, nil
	}

	err = state.postgres.Admin.AddDocumentAdmin(ctx, vkID)
	if err != nil {
		return "", nil, err
	}

	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	resp := fmt.Sprintf("Новый администратор документоархива добавлен, его vkID - %d", vkID)
	b.Message(resp)
	return addDocumentAdministrator, []*params.MessagesSendBuilder{b}, nil
}

func (state AddDocumentAdministratorState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите vkID нового администратора")
	k := object.NewMessagesKeyboard(true)
	addBackButton(k)
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state AddDocumentAdministratorState) Name() stateName {
	return addDocumentAdministrator
}

type RemoveDocumentAdministratorState struct {
	postgres *postrgres.Repo
}

func (state RemoveDocumentAdministratorState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	if messageText == "Назад" {
		return documentCabinet, nil, nil
	}

	vkID, err := strconv.Atoi(messageText)
	if err != nil || vkID < 100000000 || vkID >= 1000000000 {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("VK ID должно быть числом, например, 221486551")
		return removeDocumentAdministrator, []*params.MessagesSendBuilder{b}, nil
	}

	exists, err := state.postgres.Admin.CheckExistence(ctx, vkID)
	if err != nil {
		return "", nil, err
	}

	if !exists {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message(fmt.Sprintf("Администратора с VK ID - %d - не найдено", vkID))
		return removeDocumentAdministrator, []*params.MessagesSendBuilder{b}, nil
	}

	err = state.postgres.Admin.DeleteDocumentAdmin(ctx, vkID)
	if err != nil {
		return "", nil, err
	}

	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	resp := fmt.Sprintf("Сняты права администрирования документоархива с vkID - %d", vkID)
	b.Message(resp)
	return removeDocumentAdministrator, []*params.MessagesSendBuilder{b}, nil
}

func (state RemoveDocumentAdministratorState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите vkID администратора, который будет удален из списка администраторов архива документов")
	k := object.NewMessagesKeyboard(true)
	addBackButton(k)
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state RemoveDocumentAdministratorState) Name() stateName {
	return removeDocumentAdministrator
}

type AddPhotoAdministratorState struct {
	postgres *postrgres.Repo
}

func (state AddPhotoAdministratorState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	if messageText == "Назад" {
		return albumsCabinet, nil, nil
	}

	vkID, err := strconv.Atoi(messageText)
	if err != nil || vkID < 100000000 || vkID >= 1000000000 {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("VK ID должно быть числом, например, 221486551")
		return addPhotoAdministrator, []*params.MessagesSendBuilder{b}, nil
	}

	err = state.postgres.Admin.AddAlbumsAdmin(ctx, vkID)
	if err != nil {
		return "", nil, err
	}

	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	resp := fmt.Sprintf("Новый администратор фотоархива добавлен, его vkID - %d", vkID)
	b.Message(resp)
	return addPhotoAdministrator, []*params.MessagesSendBuilder{b}, nil
}

func (state AddPhotoAdministratorState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите vkID нового администратора")
	k := object.NewMessagesKeyboard(true)
	addBackButton(k)
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state AddPhotoAdministratorState) Name() stateName {
	return addPhotoAdministrator
}

type RemovePhotoAdministratorState struct {
	postgres *postrgres.Repo
}

func (state RemovePhotoAdministratorState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	if messageText == "Назад" {
		return albumsCabinet, nil, nil
	}

	vkID, err := strconv.Atoi(messageText)
	if err != nil || vkID < 100000000 || vkID >= 1000000000 {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("VK ID должно быть числом, например, 221486551")
		return removePhotoAdministrator, []*params.MessagesSendBuilder{b}, nil
	}

	exists, err := state.postgres.Admin.CheckExistence(ctx, vkID)
	if err != nil {
		return "", nil, err
	}

	if !exists {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message(fmt.Sprintf("Администратора с VK ID - %d - не найдено", vkID))
		return removePhotoAdministrator, []*params.MessagesSendBuilder{b}, nil
	}

	err = state.postgres.Admin.DeleteAlbumsAdmin(ctx, vkID)
	if err != nil {
		return "", nil, err
	}

	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	resp := fmt.Sprintf("Сняты права администрирования фотоархива с vkID - %d", vkID)
	b.Message(resp)
	return removePhotoAdministrator, []*params.MessagesSendBuilder{b}, nil
}

func (state RemovePhotoAdministratorState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите vkID администратора, который будет удален из списка администраторов фотоархива")
	k := object.NewMessagesKeyboard(true)
	addBackButton(k)
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state RemovePhotoAdministratorState) Name() stateName {
	return removePhotoAdministrator
}
