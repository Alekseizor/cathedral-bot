package state

import (
	"context"
	"fmt"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	"regexp"
	"strconv"
	"strings"
	"time"
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
	k.AddRow()
	k.AddTextButton("Изменить список категорий", "", "secondary")
	k.AddRow()
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
	case "Очистить параметр":
		err := state.postgres.SearchDocument.NullNameSearch(ctx, msg.PeerID)
		if err != nil {
			return editNameSearchDocument, []*params.MessagesSendBuilder{}, err
		}
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
	k.AddRow()
	k.AddTextButton("Очистить параметр", "", "secondary")
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
	case "Очистить параметр":
		err := state.postgres.SearchDocument.NullAuthorSearch(ctx, msg.PeerID)
		if err != nil {
			return editAuthorSearchDocument, []*params.MessagesSendBuilder{}, err
		}
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
	k.AddRow()
	k.AddTextButton("Очистить параметр", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditAuthorSearchDocumentState) Name() stateName {
	return editAuthorSearchDocument
}

// EditYearSearchDocumentState пользователь указывает другой год создания документа для поиска
type EditYearSearchDocumentState struct {
	postgres *postrgres.Repo
}

func (state EditYearSearchDocumentState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Назад":
		return checkSearchDocument, nil, nil
	case "Очистить параметр":
		err := state.postgres.SearchDocument.NullYearSearch(ctx, msg.PeerID)
		if err != nil {
			return editYearSearchDocument, []*params.MessagesSendBuilder{}, err
		}
		return checkSearchDocument, nil, nil
	default:
		if len(messageText) == 4 {
			year, err := strconv.Atoi(messageText)
			if err != nil {
				b := params.NewMessagesSendBuilder()
				b.RandomID(0)
				b.Message("Укажите год числом в формате YYYY")
				return editYearSearchDocument, []*params.MessagesSendBuilder{b}, nil
			}
			currentYear := time.Now().Year()
			if !(year >= 1800 && year <= currentYear) {
				b := params.NewMessagesSendBuilder()
				b.RandomID(0)
				b.Message("Укажите существующий год в формате YYYY")
				return editYearSearchDocument, []*params.MessagesSendBuilder{b}, nil
			}
			err = state.postgres.SearchDocument.UpdateYearSearch(ctx, year, msg.PeerID)
			if err != nil {
				return editYearSearchDocument, []*params.MessagesSendBuilder{}, err
			}
		} else if years := strings.Split(messageText, "-"); len(years) == 2 && years[0] != years[1] {
			startYear, err := strconv.Atoi(years[0])
			if err != nil {
				b := params.NewMessagesSendBuilder()
				b.RandomID(0)
				b.Message("Укажите временной интервал в формате YYYY-YYYY")
				return editYearSearchDocument, []*params.MessagesSendBuilder{b}, nil
			}
			currentYear := time.Now().Year()
			if !(startYear >= 1800 && startYear <= currentYear) {
				b := params.NewMessagesSendBuilder()
				b.RandomID(0)
				b.Message("Укажите существующий временной интервал в формате YYYY-YYYY")
				return editYearSearchDocument, []*params.MessagesSendBuilder{b}, nil
			}
			endYear, err := strconv.Atoi(years[1])
			if err != nil {
				b := params.NewMessagesSendBuilder()
				b.RandomID(0)
				b.Message("Укажите временной интервал в формате YYYY-YYYY")
				return editYearSearchDocument, []*params.MessagesSendBuilder{b}, nil
			}
			if !(endYear >= 1800 && endYear <= currentYear) {
				b := params.NewMessagesSendBuilder()
				b.RandomID(0)
				b.Message("Укажите существующий временной интервал в формате YYYY-YYYY")
				return editYearSearchDocument, []*params.MessagesSendBuilder{b}, nil
			}
			err = state.postgres.SearchDocument.UpdateYearRangeSearch(ctx, startYear, endYear, msg.PeerID)
			if err != nil {
				return editYearSearchDocument, []*params.MessagesSendBuilder{}, err
			}
		} else {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Параметры не действительны")
			return editYearSearchDocument, []*params.MessagesSendBuilder{b}, nil
		}
		return checkSearchDocument, nil, nil
	}
}

func (state EditYearSearchDocumentState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Укажите другой год создания документа в формате YYYY или временной интервал в формате YYYY-YYYY")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	k.AddRow()
	k.AddTextButton("Очистить параметр", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditYearSearchDocumentState) Name() stateName {
	return editYearSearchDocument
}

// EditCategoriesSearchDocumentState пользователь указывает другую категорию документа для поиска
type EditCategoriesSearchDocumentState struct {
	postgres *postrgres.Repo
}

func (state EditCategoriesSearchDocumentState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Назад":
		return checkSearchDocument, nil, nil
	case "Очистить параметр":
		err := state.postgres.SearchDocument.NullCategoriesSearch(ctx, msg.PeerID)
		if err != nil {
			return editCategoriesSearchDocument, []*params.MessagesSendBuilder{}, err
		}
		return checkSearchDocument, nil, nil
	default:
		maxID, err := state.postgres.RequestsDocuments.GetCategoryMaxID()
		if err != nil {
			return editCategoriesSearchDocument, []*params.MessagesSendBuilder{}, err
		}
		var categoriesNumbers []int
		categories := strings.Split(messageText, " ")
		for _, c := range categories {
			categoryNumber, err := strconv.Atoi(c)
			if err != nil {
				b := params.NewMessagesSendBuilder()
				b.RandomID(0)
				b.Message("Укажите номера категории из списка через пробел")
				return editCategoriesSearchDocument, []*params.MessagesSendBuilder{b}, nil
			}
			if !(categoryNumber >= 1 && categoryNumber <= maxID) {
				b := params.NewMessagesSendBuilder()
				b.RandomID(0)
				b.Message(fmt.Sprintf("Категории с номером %v нет в списке, повторите ввод", categoryNumber))
				return editCategoriesSearchDocument, []*params.MessagesSendBuilder{b}, nil
			}
			categoriesNumbers = append(categoriesNumbers, categoryNumber)
		}
		err = state.postgres.SearchDocument.UpdateCategoriesSearch(ctx, categoriesNumbers, msg.PeerID)
		if err != nil {
			return editCategoriesSearchDocument, []*params.MessagesSendBuilder{}, err
		}
		return checkSearchDocument, nil, nil
	}
}

func (state EditCategoriesSearchDocumentState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	categories, err := state.postgres.RequestsDocuments.GetCategoryNames()
	if err != nil {
		return []*params.MessagesSendBuilder{}, err
	}
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Укажите другие номера категорий документов из списка ниже через пробел:\n" + categories)
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	k.AddRow()
	k.AddTextButton("Очистить параметр", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditCategoriesSearchDocumentState) Name() stateName {
	return editCategoriesSearchDocument
}

// EditHashtagSearchDocumentState пользователь указывает другие хештеги документа для поиска
type EditHashtagSearchDocumentState struct {
	postgres *postrgres.Repo
}

func (state EditHashtagSearchDocumentState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Назад":
		return checkSearchDocument, nil, nil
	case "Очистить параметр":
		err := state.postgres.SearchDocument.NullHashtagsSearch(ctx, msg.PeerID)
		if err != nil {
			return editHashtagSearchDocument, []*params.MessagesSendBuilder{}, err
		}
		return checkSearchDocument, nil, nil
	default:
		hashtags := strings.Split(messageText, " ")
		err := state.postgres.SearchDocument.UpdateHashtagsSearch(ctx, hashtags, msg.PeerID)
		if err != nil {
			return editHashtagSearchDocument, []*params.MessagesSendBuilder{}, err
		}
		return checkSearchDocument, nil, nil
	}
}

func (state EditHashtagSearchDocumentState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Укажите другие названия хештегов через пробел (например, фамилия преподавателя)")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	k.AddRow()
	k.AddTextButton("Очистить параметр", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditHashtagSearchDocumentState) Name() stateName {
	return editHashtagSearchDocument
}
