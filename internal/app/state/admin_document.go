package state

import (
	"context"
	"strconv"

	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"

	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
)

type DocumentCabinetState struct {
	postgres *postrgres.Repo
}

func (state DocumentCabinetState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Работа с заявкой":
		return workingRequestDocument, nil, nil
	case "Работа с файлом":
		return workingDocument, nil, nil
	case "Заблокировать пользователя":
		return blockUser, nil, nil
	case "Добавить администратора":
		return documentCabinet, nil, nil
	case "Удалить администратора":
		return documentCabinet, nil, nil
	case "Выйти из кабинета администратора":
		return selectArchive, nil, nil
	default:
		return documentCabinet, nil, nil
	}
}

func (state DocumentCabinetState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Вы в кабинете администратора документоархива, выберите действие")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Работа с заявкой", "", "secondary")
	k.AddTextButton("Работа с файлом", "", "secondary")
	k.AddRow()
	k.AddTextButton("Заблокировать пользователя", "", "secondary")
	k.AddRow()
	k.AddTextButton("Добавить администратора", "", "secondary")
	k.AddTextButton("Удалить администратора", "", "secondary")
	k.AddRow()
	k.AddTextButton("Выйти из кабинета администратора", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state DocumentCabinetState) Name() stateName {
	return documentCabinet
}

///////////

type WorkingRequestDocumentState struct {
}

func (state WorkingRequestDocumentState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Заявки из очереди":
		return workingRequestDocument, nil, nil
	case "Конкретная заявка":
		return workingRequestDocument, nil, nil
	case "Назад":
		return documentCabinet, nil, nil
	default:
		return workingRequestDocument, nil, nil
	}
}

func (state WorkingRequestDocumentState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Вы хотите работать с заявками из очереди или с конкретной заявкой?")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Заявки из очереди", "", "secondary")
	k.AddTextButton("Конкретная заявка", "", "secondary")
	addWBackButton(k)
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state WorkingRequestDocumentState) Name() stateName {
	return workingRequestDocument
}

///////////

type WorkingDocumentState struct {
	postgres *postrgres.Repo
}

func (state WorkingDocumentState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	if messageText == "Назад" {
		return documentCabinet, nil, nil
	}

	documentID, err := strconv.Atoi(messageText)
	if err != nil {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("ID должно быть числом, например, 12")
		return workingDocument, []*params.MessagesSendBuilder{b}, nil
	}

	exists, err := state.postgres.Document.CheckExistence(ctx, documentID)
	if err != nil {
		return workingDocument, nil, err
	}

	if !exists {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("ID c таким документом не найдено")
		return workingDocument, []*params.MessagesSendBuilder{b}, nil
	}

	return workingDocument, nil, nil
}

func (state WorkingDocumentState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите ID файла, над которым хотите поработать. Например: 12")
	k := object.NewMessagesKeyboard(true)
	addWBackButton(k)
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state WorkingDocumentState) Name() stateName {
	return workingDocument
}
