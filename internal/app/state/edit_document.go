package state

import (
	"context"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// EditDocumentState пользователь выбирает параметр для редактирования
type EditDocumentState struct {
	postgres *postrgres.Repo
}

func (state EditDocumentState) Handler(msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Назад":
		return checkDocument, nil, nil
	case "Изменить название":
		return editNameDocument, nil, nil
	case "Изменить ФИО автора":
		return editAuthorDocument, nil, nil
	case "Изменить год":
		return editYearDocument, nil, nil
	case "Изменить категорию":
		return editCategoryDocument, nil, nil
	case "Изменить описание":
		return editDescriptionDocument, nil, nil
	case "Изменить хэштеги":
		return editHashtagDocument, nil, nil
	default:
		return editDocument, nil, nil
	}
}

func (state EditDocumentState) Show(vkID int) ([]*params.MessagesSendBuilder, error) {
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
	k.AddTextButton("Изменить год", "", "secondary")
	k.AddTextButton("Изменить категорию", "", "secondary")
	k.AddRow()
	k.AddTextButton("Изменить описание", "", "secondary")
	k.AddTextButton("Изменить хэштеги", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditDocumentState) Name() stateName {
	return editDocument
}

// EditNameDocumentState пользователь редактирует название документа
type EditNameDocumentState struct {
	postgres *postrgres.Repo
}

func (state EditNameDocumentState) Handler(msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	reqID, err := state.postgres.Document.GetDocumentLastID(context.Background(), msg.PeerID)
	if err != nil {
		return editNameDocument, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return editDocument, nil, nil
	default:
		err = state.postgres.Document.EditName(context.Background(), messageText, reqID)
		if err != nil {
			return editNameDocument, []*params.MessagesSendBuilder{}, err
		}
		return checkDocument, nil, nil
	}
}

func (state EditNameDocumentState) Show(vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите новое название загружаемого документа")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditNameDocumentState) Name() stateName {
	return editNameDocument
}

// EditAuthorDocumentState пользователь редактирует ФИО автора документа
type EditAuthorDocumentState struct {
	postgres *postrgres.Repo
}

func (state EditAuthorDocumentState) Handler(msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	reqID, err := state.postgres.Document.GetDocumentLastID(context.Background(), msg.PeerID)
	if err != nil {
		return editAuthorDocument, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return editDocument, nil, nil
	default:
		if len(messageText) > 60 {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("ФИО автора слишком длинное, повторите ввод")
			return editAuthorDocument, []*params.MessagesSendBuilder{b}, nil
		}
		russianRegex := regexp.MustCompile("^[а-яА-Я\\s]+$")
		if !russianRegex.MatchString(messageText) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("ФИО автора должно состоять из русских букв, повторите ввод")
			return editAuthorDocument, []*params.MessagesSendBuilder{b}, nil
		}
		err = state.postgres.Document.EditAuthor(context.Background(), messageText, reqID)
		if err != nil {
			return editAuthorDocument, []*params.MessagesSendBuilder{}, err
		}
		return checkDocument, nil, nil
	}
}

func (state EditAuthorDocumentState) Show(vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите новое ФИО автора загружаемого документа")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditAuthorDocumentState) Name() stateName {
	return editAuthorDocument
}

// EditYearDocumentState пользователь редактирует год создания документа
type EditYearDocumentState struct {
	postgres *postrgres.Repo
}

func (state EditYearDocumentState) Handler(msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	reqID, err := state.postgres.Document.GetDocumentLastID(context.Background(), msg.PeerID)
	if err != nil {
		return editYearDocument, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return editDocument, nil, nil
	default:
		year, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введите год числом в формате YYYY")
			return editYearDocument, []*params.MessagesSendBuilder{b}, nil
		}
		currentYear := time.Now().Year()
		if !(year >= 1800 && year <= currentYear) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введите существующий год в формате YYYY")
			return editYearDocument, []*params.MessagesSendBuilder{b}, nil
		}
		err = state.postgres.Document.EditYear(context.Background(), year, reqID)
		if err != nil {
			return editYearDocument, []*params.MessagesSendBuilder{}, err
		}
		return checkDocument, nil, nil
	}
}

func (state EditYearDocumentState) Show(vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите новый год создания документа")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditYearDocumentState) Name() stateName {
	return editYearDocument
}

// EditCategoryDocumentState пользователь редактирует категорию документа
type EditCategoryDocumentState struct {
	postgres *postrgres.Repo
}

func (state EditCategoryDocumentState) Handler(msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	reqID, err := state.postgres.Document.GetDocumentLastID(context.Background(), msg.PeerID)
	if err != nil {
		return editCategoryDocument, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return editDocument, nil, nil
	case "Своя категория":
		return editUserCategoryDocument, nil, nil
	default:
		maxID, err := state.postgres.Document.GetCategoryMaxID()
		if err != nil {
			return editCategoryDocument, []*params.MessagesSendBuilder{}, err
		}
		categoryNumber, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введите номер категории числом, повторите ввод")
			return editCategoryDocument, []*params.MessagesSendBuilder{b}, nil
		}
		if !(categoryNumber >= 1 && categoryNumber <= maxID) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Категории с таким номером нет в списке, повторите ввод")
			return editCategoryDocument, []*params.MessagesSendBuilder{b}, nil
		}
		err = state.postgres.Document.EditCategory(context.Background(), categoryNumber, reqID)
		if err != nil {
			return editCategoryDocument, []*params.MessagesSendBuilder{}, err
		}
		return checkDocument, nil, nil
	}
}

func (state EditCategoryDocumentState) Show(vkID int) ([]*params.MessagesSendBuilder, error) {
	categories, err := state.postgres.Document.GetCategoryNames()
	if err != nil {
		return []*params.MessagesSendBuilder{}, err
	}
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите новый номер категории документа из списка ниже:\n" + categories)
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	k.AddRow()
	k.AddTextButton("Своя категория", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditCategoryDocumentState) Name() stateName {
	return editCategoryDocument
}

// EditUserCategoryDocumentState пользователь редактирует и добавляет категорию документа
type EditUserCategoryDocumentState struct {
	postgres *postrgres.Repo
}

func (state EditUserCategoryDocumentState) Handler(msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	reqID, err := state.postgres.Document.GetDocumentLastID(context.Background(), msg.PeerID)
	if err != nil {
		return editUserCategoryDocument, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return editCategoryDocument, nil, nil
	default:
		err = state.postgres.Document.EditUserCategory(context.Background(), messageText, reqID)
		if err != nil {
			return editUserCategoryDocument, []*params.MessagesSendBuilder{}, err
		}
		return checkDocument, nil, nil
	}
}

func (state EditUserCategoryDocumentState) Show(vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите новое название своей категории. Оно будет рассмотрено администратором")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditUserCategoryDocumentState) Name() stateName {
	return editUserCategoryDocument
}

// EditDescriptionDocumentState пользователь редактирует описание документа
type EditDescriptionDocumentState struct {
	postgres *postrgres.Repo
}

func (state EditDescriptionDocumentState) Handler(msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	reqID, err := state.postgres.Document.GetDocumentLastID(context.Background(), msg.PeerID)
	if err != nil {
		return editDescriptionDocument, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return editDocument, nil, nil
	default:
		err = state.postgres.Document.EditDescription(context.Background(), messageText, reqID)
		if err != nil {
			return editDescriptionDocument, []*params.MessagesSendBuilder{}, err
		}
		return checkDocument, nil, nil
	}
}

func (state EditDescriptionDocumentState) Show(vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите новое описание документа")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditDescriptionDocumentState) Name() stateName {
	return editDescriptionDocument
}

// EditHashtagDocumentState пользователь редактирует хештеги документа
type EditHashtagDocumentState struct {
	postgres *postrgres.Repo
}

func (state EditHashtagDocumentState) Handler(msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	reqID, err := state.postgres.Document.GetDocumentLastID(context.Background(), msg.PeerID)
	if err != nil {
		return editHashtagDocument, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return editDocument, nil, nil
	default:
		hashtags := strings.Split(messageText, " ")
		err = state.postgres.Document.EditHashtags(context.Background(), hashtags, reqID)
		if err != nil {
			return editHashtagDocument, []*params.MessagesSendBuilder{}, err
		}
		return checkDocument, nil, nil
	}
}

func (state EditHashtagDocumentState) Show(vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите новые названия хештегов через пробел (например, фамилия преподавателя или название предмета)")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditHashtagDocumentState) Name() stateName {
	return editHashtagDocument
}
