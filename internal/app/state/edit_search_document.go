package state

import (
	"context"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	"regexp"
)

// EditSearchDocumentState пользователь выбирает параметр поиска документа для редактирования
type EditSearchDocumentState struct {
	postgres *postrgres.Repo
}

func (state EditSearchDocumentState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Назад":
		return checkSearchDocument, nil, nil
	case "Изменить название":
		return editNameSearchDocument, nil, nil
	case "Изменить ФИО автора":
		return editAuthorSearchDocument, nil, nil
	case "Изменить год/временной интервал":
		return editYearSearchDocument, nil, nil
	case "Изменить список категорий":
		return editCategoriesSearchDocument, nil, nil
	case "Изменить список хештегов":
		return editHashtagSearchDocument, nil, nil
	default:
		return editSearchDocument, nil, nil
	}
}

func (state EditSearchDocumentState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Выберите параметр для редактирования")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	k.AddRow()
	k.AddTextButton("Изменить название", "", "secondary")
	k.AddTextButton("Изменить ФИО автора", "", "secondary")
	k.AddRow()
	k.AddTextButton("Изменить год/временной интервал", "", "secondary")
	k.AddTextButton("Изменить список категорий", "", "secondary")
	k.AddTextButton("Изменить список хештегов", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditSearchDocumentState) Name() stateName {
	return editSearchDocument
}

// EditNameSearchDocumentState пользователь указывает другое название документа для поиска
type EditNameSearchDocumentState struct {
	postgres *postrgres.Repo
}

func (state EditNameSearchDocumentState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Назад":
		return checkSearchDocument, nil, nil
	default:
		err := state.postgres.SearchDocument.UpdateNameSearch(ctx, messageText, msg.PeerID)
		if err != nil {
			return editNameSearchDocument, []*params.MessagesSendBuilder{}, err
		}
		return checkSearchDocument, nil, nil
	}
}

func (state EditNameSearchDocumentState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Укажите другое название документа")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditNameSearchDocumentState) Name() stateName {
	return editNameSearchDocument
}

// EditAuthorSearchDocumentState пользователь указывает другое ФИО автора документа для поиска
type EditAuthorSearchDocumentState struct {
	postgres *postrgres.Repo
}

func (state EditAuthorSearchDocumentState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Назад":
		return checkSearchDocument, nil, nil
	default:
		if len(messageText) > 60 {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("ФИО автора слишком длинное, повторите ввод")
			return editAuthorSearchDocument, []*params.MessagesSendBuilder{b}, nil
		}
		russianRegex := regexp.MustCompile("^[а-яА-Я\\s]+$")
		if !russianRegex.MatchString(messageText) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("ФИО автора должно состоять из русских букв, повторите ввод")
			return editAuthorSearchDocument, []*params.MessagesSendBuilder{b}, nil
		}
		err := state.postgres.SearchDocument.UpdateAuthorSearch(ctx, messageText, msg.PeerID)
		if err != nil {
			return editAuthorSearchDocument, []*params.MessagesSendBuilder{}, err
		}
		return checkSearchDocument, nil, nil
	}
}

func (state EditAuthorSearchDocumentState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Укажите другое ФИО автора")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditAuthorSearchDocumentState) Name() stateName {
	return editAuthorSearchDocument
}
