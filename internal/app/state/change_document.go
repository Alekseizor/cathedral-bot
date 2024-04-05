package state

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
)

///////////

type ChangeTitleDocumentState struct {
	postgres *postrgres.Repo
}

func (state ChangeTitleDocumentState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "Назад" {
		return workingDocument, nil, nil
	}

	fileID, err := state.postgres.ObjectAdmin.Get(ctx, msg.PeerID)
	if err != nil {
		return "", nil, fmt.Errorf("[object_admin.Get]: %w", err)
	}

	err = state.postgres.Documents.EditTitle(ctx, messageText, fileID)
	if err != nil {
		return "", nil, fmt.Errorf("[documents.EditTitle]: %w", err)
	}

	resp := fmt.Sprintf("Название для документа №%d изменено на - %s", fileID, messageText)

	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message(resp)

	return documentCabinet, []*params.MessagesSendBuilder{b}, nil
}

func (state ChangeTitleDocumentState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите новое название для документа")
	k := object.NewMessagesKeyboard(true)
	addWBackButton(k)
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state ChangeTitleDocumentState) Name() stateName {
	return changeTitleDocument
}

///////////

type ChangeDescriptionDocumentState struct {
	postgres *postrgres.Repo
}

func (state ChangeDescriptionDocumentState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "Назад" {
		return workingDocument, nil, nil
	}

	fileID, err := state.postgres.ObjectAdmin.Get(ctx, msg.PeerID)
	if err != nil {
		return "", nil, fmt.Errorf("[object_admin.Get]: %w", err)
	}

	err = state.postgres.Documents.EditDescription(ctx, messageText, fileID)
	if err != nil {
		return "", nil, fmt.Errorf("[documents.EditTitle]: %w", err)
	}

	resp := fmt.Sprintf("Описание для документа №%d изменено на - %s", fileID, messageText)

	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message(resp)

	return documentCabinet, []*params.MessagesSendBuilder{b}, nil
}

func (state ChangeDescriptionDocumentState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите новое описание для документа")
	k := object.NewMessagesKeyboard(true)
	addWBackButton(k)
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state ChangeDescriptionDocumentState) Name() stateName {
	return changeDescriptionDocument
}

///////////

type ChangeAuthorDocumentState struct {
	postgres *postrgres.Repo
}

func (state ChangeAuthorDocumentState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "Назад" {
		return workingDocument, nil, nil
	}

	fileID, err := state.postgres.ObjectAdmin.Get(ctx, msg.PeerID)
	if err != nil {
		return "", nil, fmt.Errorf("[object_admin.Get]: %w", err)
	}

	err = state.postgres.Documents.EditAuthor(ctx, messageText, fileID)
	if err != nil {
		return "", nil, fmt.Errorf("[documents.EditTitle]: %w", err)
	}

	resp := fmt.Sprintf("Автор для документа №%d изменен на - %s", fileID, messageText)

	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message(resp)

	return documentCabinet, []*params.MessagesSendBuilder{b}, nil
}

func (state ChangeAuthorDocumentState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите автора документа")
	k := object.NewMessagesKeyboard(true)
	addWBackButton(k)
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state ChangeAuthorDocumentState) Name() stateName {
	return changeAuthorDocument
}

///////////

type ChangeYearDocumentState struct {
	postgres *postrgres.Repo
}

func (state ChangeYearDocumentState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "Назад" {
		return workingDocument, nil, nil
	}

	year, err := strconv.Atoi(messageText)
	if err != nil {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Год должен быть числом, например, 2020")
		return changeYearDocument, []*params.MessagesSendBuilder{b}, nil
	}

	currentYear := time.Now().Year()
	if year > currentYear {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Введите существующий год в формате YYYY")
		return changeYearDocument, []*params.MessagesSendBuilder{b}, nil
	}

	fileID, err := state.postgres.ObjectAdmin.Get(ctx, msg.PeerID)
	if err != nil {
		return "", nil, fmt.Errorf("[object_admin.Get]: %w", err)
	}

	err = state.postgres.Documents.EditYear(ctx, year, fileID)
	if err != nil {
		return "", nil, fmt.Errorf("[documents.EditTitle]: %w", err)
	}

	resp := fmt.Sprintf("Год издания для документа №%d изменен на - %s", fileID, messageText)

	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message(resp)

	return documentCabinet, []*params.MessagesSendBuilder{b}, nil
}

func (state ChangeYearDocumentState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите год публикации документа")
	k := object.NewMessagesKeyboard(true)
	addWBackButton(k)
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state ChangeYearDocumentState) Name() stateName {
	return changeYearDocument
}

///////////

type ChangeCategoryDocumentState struct {
	postgres *postrgres.Repo
}

func (state ChangeCategoryDocumentState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "Назад" {
		return workingDocument, nil, nil
	}

	year, err := strconv.Atoi(messageText)
	if err != nil {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Год должен быть числом, например, 2020")
		return workingDocument, []*params.MessagesSendBuilder{b}, nil
	}

	fileID, err := state.postgres.ObjectAdmin.Get(ctx, msg.PeerID)
	if err != nil {
		return "", nil, fmt.Errorf("[object_admin.Get]: %w", err)
	}

	err = state.postgres.Documents.EditYear(ctx, year, fileID)
	if err != nil {
		return "", nil, fmt.Errorf("[documents.EditTitle]: %w", err)
	}

	resp := fmt.Sprintf("Категория для документа №%d изменена на - %s", fileID, messageText)

	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message(resp)

	return documentCabinet, []*params.MessagesSendBuilder{b}, nil
}

func (state ChangeCategoryDocumentState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	categories, err := state.postgres.RequestsDocuments.GetCategoryNames()
	if err != nil {
		return nil, err
	}
	b.Message("Вот существующие категории:\n" + categories)

	b1 := params.NewMessagesSendBuilder()
	b1.RandomID(0)
	b1.Message("Если хотите выбрать из существующих, отправьте ее номер, иначе напишите название для новой категории")
	k := object.NewMessagesKeyboard(true)
	addWBackButton(k)
	b1.Keyboard(k)

	return []*params.MessagesSendBuilder{b, b1}, nil
}

func (state ChangeCategoryDocumentState) Name() stateName {
	return changeCategoryDocument
}
