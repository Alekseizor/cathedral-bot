package state

import (
	"context"
	"github.com/Alekseizor/cathedral-bot/internal/app/ds"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var validExtensionArchive = map[string]struct{}{
	"rar": struct{}{},
	"zip": struct{}{},
}

// LoadArchiveState пользователь загружает архив
type LoadArchiveState struct {
	postgres *postrgres.Repo
	vk       *api.VK
}

func (state LoadArchiveState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "Назад" {
		return documentStart, nil, nil
	}
	attachment := msg.Attachments

	if len(attachment) == 0 {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Загрузите ваш архив, прикрепив его к сообщению")
		return loadArchive, []*params.MessagesSendBuilder{b}, nil
	}

	if len(attachment) > 1 {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Можно загрузить лишь один архив")
		return loadArchive, []*params.MessagesSendBuilder{b}, nil
	}

	if _, ok := validExtensionArchive[attachment[0].Doc.Ext]; !ok {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Поддерживаются только архивы форматов rar и zip.")
		return loadArchive, []*params.MessagesSendBuilder{b}, nil
	}

	err := state.postgres.Document.UploadDocument(ctx, state.vk, attachment[0].Doc, msg.PeerID)
	if err != nil {
		return loadArchive, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	default:
		return nameArchive, nil, nil
	}
}

func (state LoadArchiveState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Загрузите ваш архив. Убедитесь, что все документы в архиве могут быть описаны одинаковыми параметрами.\nПоддерживаются форматы rar и zip.")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state LoadArchiveState) Name() stateName {
	return loadArchive
}

// NameArchiveState пользователь указывает название документа
type NameArchiveState struct {
	postgres *postrgres.Repo
}

func (state NameArchiveState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Назад":
		err := state.postgres.Document.DeleteDocumentRequest(ctx, msg.PeerID)
		if err != nil {
			return nameArchive, []*params.MessagesSendBuilder{}, err
		}
		return loadArchive, nil, nil
	case "Пропустить":
		return authorArchive, nil, nil
	default:
		err := state.postgres.Document.UpdateName(ctx, msg.PeerID, msg.Text)
		if err != nil {
			return nameArchive, []*params.MessagesSendBuilder{}, err
		}
		return authorArchive, nil, nil
	}
}

func (state NameArchiveState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите название загружаемого архива(пропустить - будет использовано название архива)")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	k.AddRow()
	k.AddTextButton("Пропустить", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state NameArchiveState) Name() stateName {
	return nameArchive
}

// AuthorArchiveState пользователь указывает ФИО автора документа
type AuthorArchiveState struct {
	postgres *postrgres.Repo
}

func (state AuthorArchiveState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Назад":
		return nameArchive, nil, nil
	case "Пропустить":
		return yearArchive, nil, nil
	default:
		if len(messageText) > 60 {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("ФИО автора слишком длинное, повторите ввод")
			return authorArchive, []*params.MessagesSendBuilder{b}, nil
		}
		russianRegex := regexp.MustCompile("^[а-яА-Я\\s]+$")
		if !russianRegex.MatchString(messageText) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("ФИО автора должно состоять из русских букв, повторите ввод")
			return authorArchive, []*params.MessagesSendBuilder{b}, nil
		}
		err := state.postgres.Document.UpdateAuthor(ctx, msg.PeerID, msg.Text)
		if err != nil {
			return authorArchive, []*params.MessagesSendBuilder{}, err
		}
		return yearArchive, nil, nil
	}
}

func (state AuthorArchiveState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите ФИО автора архива. ФИО может быть неполным")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	k.AddRow()
	k.AddTextButton("Пропустить", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state AuthorArchiveState) Name() stateName {
	return authorArchive
}

// YearArchiveState пользователь указывает год создания архива
type YearArchiveState struct {
	postgres *postrgres.Repo
}

func (state YearArchiveState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Назад":
		return authorArchive, nil, nil
	case "Пропустить":
		return categoryArchive, nil, nil
	default:
		year, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введите год числом в формате YYYY")
			return yearArchive, []*params.MessagesSendBuilder{b}, nil
		}
		currentYear := time.Now().Year()
		if !(year >= 1800 && year <= currentYear) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введите существующий год в формате YYYY")
			return yearArchive, []*params.MessagesSendBuilder{b}, nil
		}
		err = state.postgres.Document.UpdateYear(ctx, msg.PeerID, year)
		if err != nil {
			return yearArchive, []*params.MessagesSendBuilder{}, err
		}
		return categoryArchive, nil, nil
	}
}

func (state YearArchiveState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите год создания архива в формате YYYY")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	k.AddRow()
	k.AddTextButton("Пропустить", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state YearArchiveState) Name() stateName {
	return yearArchive
}

// CategoryArchiveState пользователь указывает существующую категорию для документов в архиве
type CategoryArchiveState struct {
	postgres *postrgres.Repo
}

func (state CategoryArchiveState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Назад":
		return yearArchive, nil, nil
	case "Своя категория":
		return userCategoryArchive, nil, nil
	default:
		maxID, err := state.postgres.Document.GetCategoryMaxID()
		if err != nil {
			return categoryArchive, []*params.MessagesSendBuilder{}, err
		}
		categoryNumber, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введите номер категории числом, повторите ввод")
			return categoryArchive, []*params.MessagesSendBuilder{b}, nil
		}
		if !(categoryNumber >= 1 && categoryNumber <= maxID) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Категории с таким номером нет в списке, повторите ввод")
			return categoryArchive, []*params.MessagesSendBuilder{b}, nil
		}

		err = state.postgres.Document.UpdateCategory(ctx, msg.PeerID, categoryNumber)
		if err != nil {
			return categoryArchive, []*params.MessagesSendBuilder{}, err
		}
		return descriptionArchive, nil, nil
	}
}

func (state CategoryArchiveState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	categories, err := state.postgres.Document.GetCategoryNames()
	if err != nil {
		return []*params.MessagesSendBuilder{}, err
	}
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите номер категории документов в архиве из списка ниже:\n" + categories)
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	k.AddRow()
	k.AddTextButton("Своя категория", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state CategoryArchiveState) Name() stateName {
	return categoryArchive
}

// UserCategoryArchiveState пользователь указывает свою категорию документов в архиве
type UserCategoryArchiveState struct {
	postgres *postrgres.Repo
}

func (state UserCategoryArchiveState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Назад":
		return categoryArchive, nil, nil
	default:
		err := state.postgres.Document.UpdateUserCategory(ctx, msg.PeerID, messageText)
		if err != nil {
			return userCategoryArchive, []*params.MessagesSendBuilder{}, err
		}
		return descriptionArchive, nil, nil
	}
}

func (state UserCategoryArchiveState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите название своей категории. Оно будет рассмотрено администратором")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state UserCategoryArchiveState) Name() stateName {
	return userCategoryArchive
}

// DescriptionArchiveState пользователь указывает хештег документа
type DescriptionArchiveState struct {
	postgres *postrgres.Repo
}

func (state DescriptionArchiveState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Назад":
		return categoryArchive, nil, nil
	case "Пропустить":
		return hashtagArchive, nil, nil
	default:
		err := state.postgres.Document.UpdateDescription(ctx, msg.PeerID, messageText)
		if err != nil {
			return descriptionArchive, []*params.MessagesSendBuilder{}, err
		}
		return hashtagArchive, nil, nil
	}
}

func (state DescriptionArchiveState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите описание архива")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	k.AddRow()
	k.AddTextButton("Пропустить", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state DescriptionArchiveState) Name() stateName {
	return descriptionArchive
}

// HashtagArchiveState пользователь указывает хештег документа
type HashtagArchiveState struct {
	postgres *postrgres.Repo
}

func (state HashtagArchiveState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Назад":
		return descriptionArchive, nil, nil
	case "Пропустить":
		return checkArchive, nil, nil
	default:
		hashtags := strings.Split(messageText, " ")

		err := state.postgres.Document.UpdateHashtags(ctx, msg.PeerID, hashtags)
		if err != nil {
			return hashtagArchive, []*params.MessagesSendBuilder{}, err
		}
		return checkArchive, nil, nil
	}
}

func (state HashtagArchiveState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
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

func (state HashtagArchiveState) Name() stateName {
	return hashtagArchive
}

// CheckArchiveState пользователь проверяет заявку на загрузку архива
type CheckArchiveState struct {
	postgres *postrgres.Repo
}

func (state CheckArchiveState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Назад":
		return hashtagArchive, nil, nil
	case "Отправить":
		reqID, err := state.postgres.Document.GetDocumentLastID(ctx, msg.PeerID)
		if err != nil {
			return checkArchive, []*params.MessagesSendBuilder{}, err
		}
		err = state.postgres.Document.UpdateStatus(ctx, ds.StatusUserConfirmed, reqID)
		if err != nil {
			return checkArchive, []*params.MessagesSendBuilder{}, err
		}
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Ваша заявка на загрузку архива отправлена на одобрение администратору. Вы можете отслеживать статус своей заявки в личном кабинете")
		return documentStart, []*params.MessagesSendBuilder{b}, nil
	case "Редактировать заявку":
		return editDocument, nil, nil
	default:
		return checkArchive, nil, nil
	}
}

func (state CheckArchiveState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	output, attachment, err := state.postgres.Document.CheckParams(ctx, vkID)
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

func (state CheckArchiveState) Name() stateName {
	return checkArchive
}
