package state

import (
	"context"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var validExtension = map[string]struct{}{
	"doc":  struct{}{},
	"docx": struct{}{},
	"ppt":  struct{}{},
	"pptx": struct{}{},
	"txt":  struct{}{},
	"pdf":  struct{}{},
}

// LoadDocumentState пользователь загружает документ
type LoadDocumentState struct {
	postgres *postrgres.Repo
	vk       *api.VK
}

func (state LoadDocumentState) Handler(msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "Назад" {
		return documentStart, nil, nil
	}
	attachment := msg.Attachments

	if len(attachment) == 0 {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Загрузите ваш документ, прикрепив его к сообщению")
		return loadDocument, []*params.MessagesSendBuilder{b}, nil
	}

	if len(attachment) > 1 {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Можно загрузить лишь один документ, для загрузки множества документов воспользуйтесь загрузкой архива")
		return loadDocument, []*params.MessagesSendBuilder{b}, nil
	}

	if _, ok := validExtension[attachment[0].Doc.Ext]; !ok {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Поддерживаются только документы формата doc/docx, pdf, txt, ppt/pptx.")
		return loadDocument, []*params.MessagesSendBuilder{b}, nil
	}

	err := state.postgres.Document.UploadDocument(context.Background(), state.vk, attachment[0].Doc, msg.PeerID)
	if err != nil {
		return loadDocument, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	default:
		return nameDocument, nil, nil
	}
}

func (state LoadDocumentState) Show(vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Загрузите ваш документ")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state LoadDocumentState) Name() stateName {
	return loadDocument
}

// NameDocumentState пользователь указывает название документа
type NameDocumentState struct {
	postgres *postrgres.Repo
}

func (state NameDocumentState) Handler(msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Назад":
		err := state.postgres.Document.DeleteDocumentRequest(context.Background(), msg.PeerID)
		if err != nil {
			return nameDocument, []*params.MessagesSendBuilder{}, err
		}
		return loadDocument, nil, nil
	case "Пропустить":
		return authorDocument, nil, nil
	default:
		err := state.postgres.Document.UpdateName(context.Background(), msg.PeerID, msg.Text)
		if err != nil {
			return nameDocument, []*params.MessagesSendBuilder{}, err
		}
		return authorDocument, nil, nil
	}
}

func (state NameDocumentState) Show(vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите название загружаемого документа(пропустить - будет использовано название документа)")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	k.AddRow()
	k.AddTextButton("Пропустить", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state NameDocumentState) Name() stateName {
	return nameDocument
}

// AuthorDocumentState пользователь указывает ФИО автора документа
type AuthorDocumentState struct {
	postgres *postrgres.Repo
}

func (state AuthorDocumentState) Handler(msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Назад":
		return nameDocument, nil, nil
	case "Пропустить":
		return yearDocument, nil, nil
	default:
		if len(messageText) > 60 {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("ФИО автора слишком длинное, повторите ввод")
			return authorDocument, []*params.MessagesSendBuilder{b}, nil
		}
		russianRegex := regexp.MustCompile("^[а-яА-Я\\s]+$")
		if !russianRegex.MatchString(messageText) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("ФИО автора должно состоять из русских букв, повторите ввод")
			return authorDocument, []*params.MessagesSendBuilder{b}, nil
		}
		err := state.postgres.Document.UpdateAuthor(context.Background(), msg.PeerID, msg.Text)
		if err != nil {
			return authorDocument, []*params.MessagesSendBuilder{}, err
		}
		return yearDocument, nil, nil
	}
}

func (state AuthorDocumentState) Show(vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите ФИО автора. ФИО может быть неполным")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	k.AddRow()
	k.AddTextButton("Пропустить", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state AuthorDocumentState) Name() stateName {
	return authorDocument
}

// YearDocumentState пользователь указывает год создания документа
type YearDocumentState struct {
	postgres *postrgres.Repo
}

func (state YearDocumentState) Handler(msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Назад":
		return authorDocument, nil, nil
	case "Пропустить":
		return categoryDocument, nil, nil
	default:
		year, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введите год числом в формате YYYY")
			return yearDocument, []*params.MessagesSendBuilder{b}, nil
		}
		currentYear := time.Now().Year()
		if !(year >= 1800 && year <= currentYear) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введите существующий год в формате YYYY")
			return yearDocument, []*params.MessagesSendBuilder{b}, nil
		}
		err = state.postgres.Document.UpdateYear(context.Background(), msg.PeerID, year)
		if err != nil {
			return yearDocument, []*params.MessagesSendBuilder{}, err
		}
		return categoryDocument, nil, nil
	}
}

func (state YearDocumentState) Show(vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите год создания документа в формате YYYY")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	k.AddRow()
	k.AddTextButton("Пропустить", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state YearDocumentState) Name() stateName {
	return yearDocument
}

// CategoryDocumentState пользователь указывает существующую категорию документа
type CategoryDocumentState struct {
	postgres *postrgres.Repo
}

func (state CategoryDocumentState) Handler(msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Назад":
		return yearDocument, nil, nil
	case "Своя категория":
		return userCategoryDocument, nil, nil
	default:
		maxID, err := state.postgres.Document.GetCategoryMaxID()
		if err != nil {
			return categoryDocument, []*params.MessagesSendBuilder{}, err
		}
		categoryNumber, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введите номер категории числом, повторите ввод")
			return categoryDocument, []*params.MessagesSendBuilder{b}, nil
		}
		if !(categoryNumber >= 1 && categoryNumber <= maxID) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Категории с таким номером нет в списке, повторите ввод")
			return categoryDocument, []*params.MessagesSendBuilder{b}, nil
		}

		err = state.postgres.Document.UpdateCategory(context.Background(), msg.PeerID, categoryNumber)
		if err != nil {
			return categoryDocument, []*params.MessagesSendBuilder{}, err
		}
		return descriptionDocument, nil, nil
	}
}

func (state CategoryDocumentState) Show(vkID int) ([]*params.MessagesSendBuilder, error) {
	categories, err := state.postgres.Document.GetCategoryNames()
	if err != nil {
		return []*params.MessagesSendBuilder{}, err
	}
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите номер категории документа из списка ниже:\n" + categories)
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	k.AddRow()
	k.AddTextButton("Своя категория", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state CategoryDocumentState) Name() stateName {
	return categoryDocument
}

// UserCategoryDocumentState пользователь указывает свою категорию документа
type UserCategoryDocumentState struct {
	postgres *postrgres.Repo
}

func (state UserCategoryDocumentState) Handler(msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Назад":
		return categoryDocument, nil, nil
	default:
		err := state.postgres.Document.UpdateUserCategory(context.Background(), msg.PeerID, messageText)
		if err != nil {
			return userCategoryDocument, []*params.MessagesSendBuilder{}, err
		}
		return descriptionDocument, nil, nil
	}
}

func (state UserCategoryDocumentState) Show(vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите название своей категории. Оно будет рассмотрено администратором")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state UserCategoryDocumentState) Name() stateName {
	return userCategoryDocument
}

// DescriptionDocumentState пользователь указывает хештег документа
type DescriptionDocumentState struct {
	postgres *postrgres.Repo
}

func (state DescriptionDocumentState) Handler(msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Назад":
		return categoryDocument, nil, nil
	case "Пропустить":
		return hashtagDocument, nil, nil
	default:
		err := state.postgres.Document.UpdateDescription(context.Background(), msg.PeerID, messageText)
		if err != nil {
			return descriptionDocument, []*params.MessagesSendBuilder{}, err
		}
		return hashtagDocument, nil, nil
	}
}

func (state DescriptionDocumentState) Show(vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите описание документа")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	k.AddRow()
	k.AddTextButton("Пропустить", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state DescriptionDocumentState) Name() stateName {
	return descriptionDocument
}

// HashtagDocumentState пользователь указывает хештег документа
type HashtagDocumentState struct {
	postgres *postrgres.Repo
}

func (state HashtagDocumentState) Handler(msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Назад":
		return descriptionDocument, nil, nil
	case "Пропустить":
		return checkDocument, nil, nil
	default:
		hashtags := strings.Split(messageText, " ")

		err := state.postgres.Document.UpdateHashtags(context.Background(), msg.PeerID, hashtags)
		if err != nil {
			return hashtagDocument, []*params.MessagesSendBuilder{}, err
		}
		return checkDocument, nil, nil
	}
}

func (state HashtagDocumentState) Show(vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите названия хештегов через пробел (например, фамилия преподавателя или название предмета)")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	k.AddRow()
	k.AddTextButton("Пропустить", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state HashtagDocumentState) Name() stateName {
	return hashtagDocument
}

// CheckDocumentState пользователь
type CheckDocumentState struct {
	postgres *postrgres.Repo
}

func (state CheckDocumentState) Handler(msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Назад":
		return hashtagDocument, nil, nil
	case "Отправить":
		return checkDocument, nil, nil
	case "Редактировать заявку":
		return editDocument, nil, nil
	default:
		return checkDocument, nil, nil
	}
}

func (state CheckDocumentState) Show(vkID int) ([]*params.MessagesSendBuilder, error) {
	output, attachment, err := state.postgres.Document.CheckParams(context.Background(), vkID)
	if err != nil {
		return []*params.MessagesSendBuilder{}, err
	}
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Проверьте правильность введенных параметров:\n" + output)
	b.Attachment(attachment)
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	k.AddRow()
	k.AddTextButton("Отправить", "", "secondary")
	k.AddRow()
	k.AddTextButton("Редактировать заявку", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state CheckDocumentState) Name() stateName {
	return checkDocument
}
