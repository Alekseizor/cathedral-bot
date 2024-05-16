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

// NameSearchDocumentState пользователь указывает название документа для поиска
type NameSearchDocumentState struct {
	postgres *postrgres.Repo
}

func (state NameSearchDocumentState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Назад":
		err := state.postgres.SearchDocument.DeleteSearch(ctx, msg.PeerID)
		if err != nil {
			return nameSearchDocument, []*params.MessagesSendBuilder{}, err
		}
		return documentStart, nil, nil
	case "Пропустить":
		return authorSearchDocument, nil, nil
	default:
		err := state.postgres.SearchDocument.UpdateNameSearch(ctx, messageText, msg.PeerID)
		if err != nil {
			return nameSearchDocument, []*params.MessagesSendBuilder{}, err
		}
		return authorSearchDocument, nil, nil
	}
}

func (state NameSearchDocumentState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Укажите название документа")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	k.AddRow()
	k.AddTextButton("Пропустить", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state NameSearchDocumentState) Name() stateName {
	return nameSearchDocument
}

// AuthorSearchDocumentState пользователь указывает ФИО автора документа для поиска
type AuthorSearchDocumentState struct {
	postgres *postrgres.Repo
}

func (state AuthorSearchDocumentState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Назад":
		return nameSearchDocument, nil, nil
	case "Пропустить":
		return yearSearchDocument, nil, nil
	case "Перейти к поиску":
		return doSearchDocument, nil, nil
	default:
		if len(messageText) > 60 {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("ФИО автора слишком длинное, повторите ввод")
			return authorSearchDocument, []*params.MessagesSendBuilder{b}, nil
		}
		russianRegex := regexp.MustCompile("^[а-яА-Я\\s]+$")
		if !russianRegex.MatchString(messageText) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("ФИО автора должно состоять из русских букв, повторите ввод")
			return authorSearchDocument, []*params.MessagesSendBuilder{b}, nil
		}
		err := state.postgres.SearchDocument.UpdateAuthorSearch(ctx, messageText, msg.PeerID)
		if err != nil {
			return authorSearchDocument, []*params.MessagesSendBuilder{}, err
		}
		return yearSearchDocument, nil, nil
	}
}

func (state AuthorSearchDocumentState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Укажите ФИО автора")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	k.AddRow()
	k.AddTextButton("Пропустить", "", "secondary")
	k.AddRow()
	k.AddTextButton("Перейти к поиску", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state AuthorSearchDocumentState) Name() stateName {
	return authorSearchDocument
}

// YearSearchDocumentState пользователь указывает год создания(интервал по годам) документа для поиска
type YearSearchDocumentState struct {
	postgres *postrgres.Repo
}

func (state YearSearchDocumentState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Назад":
		return authorSearchDocument, nil, nil
	case "Пропустить":
		return categoriesSearchDocument, nil, nil
	case "Перейти к поиску":
		return doSearchDocument, nil, nil
	default:
		if len(messageText) == 4 {
			year, err := strconv.Atoi(messageText)
			if err != nil {
				b := params.NewMessagesSendBuilder()
				b.RandomID(0)
				b.Message("Укажите год числом в формате YYYY")
				return yearSearchDocument, []*params.MessagesSendBuilder{b}, nil
			}
			currentYear := time.Now().Year()
			if !(year >= 1800 && year <= currentYear) {
				b := params.NewMessagesSendBuilder()
				b.RandomID(0)
				b.Message("Укажите существующий год в формате YYYY")
				return yearSearchDocument, []*params.MessagesSendBuilder{b}, nil
			}
			err = state.postgres.SearchDocument.UpdateYearSearch(ctx, year, msg.PeerID)
			if err != nil {
				return yearSearchDocument, []*params.MessagesSendBuilder{}, err
			}
		} else if years := strings.Split(messageText, "-"); len(years) == 2 && years[0] != years[1] {
			startYear, err := strconv.Atoi(years[0])
			if err != nil {
				b := params.NewMessagesSendBuilder()
				b.RandomID(0)
				b.Message("Укажите временной интервал в формате YYYY-YYYY")
				return yearSearchDocument, []*params.MessagesSendBuilder{b}, nil
			}
			currentYear := time.Now().Year()
			if !(startYear >= 1800 && startYear <= currentYear) {
				b := params.NewMessagesSendBuilder()
				b.RandomID(0)
				b.Message("Укажите существующий временной интервал в формате YYYY-YYYY")
				return yearSearchDocument, []*params.MessagesSendBuilder{b}, nil
			}
			endYear, err := strconv.Atoi(years[1])
			if err != nil {
				b := params.NewMessagesSendBuilder()
				b.RandomID(0)
				b.Message("Укажите временной интервал в формате YYYY-YYYY")
				return yearSearchDocument, []*params.MessagesSendBuilder{b}, nil
			}
			if !(endYear >= 1800 && endYear <= currentYear) {
				b := params.NewMessagesSendBuilder()
				b.RandomID(0)
				b.Message("Укажите существующий временной интервал в формате YYYY-YYYY")
				return yearSearchDocument, []*params.MessagesSendBuilder{b}, nil
			}
			err = state.postgres.SearchDocument.UpdateYearRangeSearch(ctx, startYear, endYear, msg.PeerID)
			if err != nil {
				return yearSearchDocument, []*params.MessagesSendBuilder{}, err
			}
		} else {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Параметры не действительны")
			return yearSearchDocument, []*params.MessagesSendBuilder{b}, nil
		}
		return categoriesSearchDocument, nil, nil
	}
}

func (state YearSearchDocumentState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Укажите год создания документа в формате YYYY или временной интервал в формате YYYY-YYYY")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	k.AddRow()
	k.AddTextButton("Пропустить", "", "secondary")
	k.AddRow()
	k.AddTextButton("Перейти к поиску", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state YearSearchDocumentState) Name() stateName {
	return yearSearchDocument
}

// CategoriesSearchDocumentState пользователь указывает список категорий для поиска документа
type CategoriesSearchDocumentState struct {
	postgres *postrgres.Repo
}

func (state CategoriesSearchDocumentState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Назад":
		return yearSearchDocument, nil, nil
	case "Пропустить":
		return hashtagSearchDocument, nil, nil
	case "Перейти к поиску":
		return doSearchDocument, nil, nil
	default:
		maxID, err := state.postgres.RequestsDocuments.GetCategoryMaxID()
		if err != nil {
			return categoriesSearchDocument, []*params.MessagesSendBuilder{}, err
		}
		var categoriesNumbers []int
		categories := strings.Split(messageText, " ")
		for _, c := range categories {
			categoryNumber, err := strconv.Atoi(c)
			if err != nil {
				b := params.NewMessagesSendBuilder()
				b.RandomID(0)
				b.Message("Укажите номера категории из списка через пробел")
				return categoriesSearchDocument, []*params.MessagesSendBuilder{b}, nil
			}
			if !(categoryNumber >= 1 && categoryNumber <= maxID) {
				b := params.NewMessagesSendBuilder()
				b.RandomID(0)
				b.Message(fmt.Sprintf("Категории с номером %v нет в списке, повторите ввод", categoryNumber))
				return categoriesSearchDocument, []*params.MessagesSendBuilder{b}, nil
			}
			categoriesNumbers = append(categoriesNumbers, categoryNumber)
		}
		err = state.postgres.SearchDocument.UpdateCategoriesSearch(ctx, categoriesNumbers, msg.PeerID)
		if err != nil {
			return categoriesSearchDocument, []*params.MessagesSendBuilder{}, err
		}
		return hashtagSearchDocument, nil, nil
	}
}

func (state CategoriesSearchDocumentState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	categories, err := state.postgres.RequestsDocuments.GetCategoryNames()
	if err != nil {
		return []*params.MessagesSendBuilder{}, err
	}
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Укажите номера категорий документов из списка ниже через пробел:\n" + categories)
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	k.AddRow()
	k.AddTextButton("Пропустить", "", "secondary")
	k.AddRow()
	k.AddTextButton("Перейти к поиску", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state CategoriesSearchDocumentState) Name() stateName {
	return categoriesSearchDocument
}

// HashtagSearchDocumentState пользователь указывает список хештегов для поиска документа
type HashtagSearchDocumentState struct {
	postgres *postrgres.Repo
}

func (state HashtagSearchDocumentState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Назад":
		return categoriesSearchDocument, nil, nil
	case "Пропустить":
		return checkSearchDocument, nil, nil
	case "Перейти к поиску":
		return doSearchDocument, nil, nil
	default:
		hashtags := strings.Split(messageText, " ")
		err := state.postgres.SearchDocument.UpdateHashtagsSearch(ctx, hashtags, msg.PeerID)
		if err != nil {
			return hashtagSearchDocument, []*params.MessagesSendBuilder{}, err
		}
		return checkSearchDocument, nil, nil
	}
}

func (state HashtagSearchDocumentState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Укажите названия хештегов через пробел (например, фамилия преподавателя)")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	k.AddRow()
	k.AddTextButton("Пропустить", "", "secondary")
	k.AddRow()
	k.AddTextButton("Перейти к поиску", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state HashtagSearchDocumentState) Name() stateName {
	return hashtagSearchDocument
}

// CheckSearchDocumentState пользователь проверяет параметры для поиска документа
type CheckSearchDocumentState struct {
	postgres *postrgres.Repo
}

func (state CheckSearchDocumentState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Назад":
		return hashtagSearchDocument, nil, nil
	case "Найти":
		isNull, err := state.postgres.SearchDocument.SearchParamsIsNULL(ctx, msg.PeerID)
		if err != nil {
			return checkSearchDocument, []*params.MessagesSendBuilder{}, err
		}
		if isNull {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Необходимо указать хотя бы один параметр поиска")
			return checkSearchDocument, []*params.MessagesSendBuilder{b}, nil
		} else {
			return doSearchDocument, nil, nil
		}
	case "Редактировать параметры":
		return editSearchDocument, nil, nil
	default:
		return checkSearchDocument, nil, nil
	}
}

func (state CheckSearchDocumentState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	output, err := state.postgres.SearchDocument.CheckSearchParams(ctx, vkID)
	if err != nil {
		return []*params.MessagesSendBuilder{}, err
	}
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Проверьте правильность введенных параметров для поиска:\n" + output)
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	k.AddRow()
	k.AddTextButton("Найти", "", "secondary")
	k.AddRow()
	k.AddTextButton("Редактировать параметры", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state CheckSearchDocumentState) Name() stateName {
	return checkSearchDocument
}
