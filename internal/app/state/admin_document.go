package state

import (
	"context"
	"fmt"
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

	exists, err := state.postgres.Documents.CheckExistence(ctx, documentID)
	if err != nil {
		return workingDocument, nil, err
	}
	if !exists {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("ID c таким документом не найдено")
		return workingDocument, []*params.MessagesSendBuilder{b}, nil
	}

	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	output, _, err := state.postgres.Documents.GetOutput(ctx, documentID)
	//TODO: добавить attachment
	//b.Attachment(attachment)
	b.Message(output)
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Изменить", documentID, "secondary")
	k.AddTextButton("Удалить", documentID, "secondary")
	addWBackButton(k)
	b.Keyboard(k)

	return actionOnDocument, []*params.MessagesSendBuilder{b}, nil
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

///////////

type ActionOnDocumentState struct {
	postgres *postrgres.Repo
}

func (state ActionOnDocumentState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "Назад" {
		return workingDocument, nil, nil
	}

	payload := msg.Payload

	documentID, err := strconv.Atoi(payload)
	if err != nil {
		return "", nil, err
	}

	switch messageText {
	case "Удалить":
		err = state.postgres.Documents.Delete(ctx, documentID)
		if err != nil {
			return "", nil, err
		}
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Документ успешно удален")
		return workingDocument, []*params.MessagesSendBuilder{b}, nil
	case "Изменить":
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Что Вы хотите изменить в документе?")
		k := object.NewMessagesKeyboard(true)
		k.AddRow()
		k.AddTextButton("Название", documentID, "secondary")
		k.AddTextButton("Описание", documentID, "secondary")
		k.AddRow()
		k.AddTextButton("Автор", documentID, "secondary")
		k.AddTextButton("Год", documentID, "secondary")
		k.AddRow()
		k.AddTextButton("Категория", documentID, "secondary")
		k.AddTextButton("Хештеги", documentID, "secondary")
		addWBackButton(k)
		b.Keyboard(k)
		return changeDocument, []*params.MessagesSendBuilder{b}, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Выберете действие")
		return workingDocument, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state ActionOnDocumentState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	return nil, nil
}

func (state ActionOnDocumentState) Name() stateName {
	return actionOnDocument
}

///////////

type ChangeDocumentState struct {
	postgres *postrgres.Repo
}

func (state ChangeDocumentState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "Назад" {
		return workingDocument, nil, nil
	}

	payload := msg.Payload

	documentID, err := strconv.Atoi(payload)
	if err != nil {
		return "", nil, err
	}

	err = state.postgres.ObjectAdmin.Update(ctx, documentID, msg.PeerID)
	if err != nil {
		return "", nil, fmt.Errorf("[object_admin.Update]: %w", err)
	}

	switch messageText {
	case "Название":
		return changeTitleDocument, nil, nil
	case "Описание":
		return changeDescriptionDocument, nil, nil
	case "Автор":
		return changeAuthorDocument, nil, nil
	case "Год":
		return changeYearDocument, nil, nil
	case "Категория":
		return changeCategoryDocument, nil, nil
	case "Хештеги":
		return workingDocument, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Выберете действие")
		return workingDocument, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state ChangeDocumentState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	return nil, nil
}

func (state ChangeDocumentState) Name() stateName {
	return changeDocument
}
