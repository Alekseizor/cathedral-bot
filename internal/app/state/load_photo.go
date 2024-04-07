package state

import (
	"context"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	"strconv"
	"time"
)

var validExtension = map[string]struct{}{
	"jpg":  struct{}{},
	"jpeg": struct{}{},
	"png":  struct{}{},
	"tif":  struct{}{},
	"tiff": struct{}{},
}

// LoadPhotoState пользователь загружает фотографию
type LoadPhotoState struct {
	postgres *postrgres.Repo
	vk       *api.VK
}

func (state LoadPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "Назад" {
		return photoStart, nil, nil
	}
	attachment := msg.Attachments

	if len(attachment) == 0 {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Загрузите фотографию, прикрепив её к сообщению")
		return loadPhoto, []*params.MessagesSendBuilder{b}, nil
	}

	if len(attachment) > 1 {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Можно загрузить лишь одну фотографию, для загрузки множества фотографий воспользуйтесь загрузкой архива")
		return loadPhoto, []*params.MessagesSendBuilder{b}, nil
	}

	if attachment[0].Type == "photo" {
		err := state.postgres.RequestPhoto.UploadPhoto(ctx, state.vk, attachment[0].Photo, msg.PeerID)
		if err != nil {
			return loadPhoto, []*params.MessagesSendBuilder{}, err
		}

		return isPeoplePresentPhoto, nil, nil
	}

	if attachment[0].Type == "doc" {
		if _, ok := validExtension[attachment[0].Doc.Ext]; !ok {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Данная фотография недопустимого формата")
			return loadPhoto, []*params.MessagesSendBuilder{b}, nil
		}

		err := state.postgres.RequestPhoto.UploadPhotoAsFile(ctx, state.vk, attachment[0].Doc, msg.PeerID)
		if err != nil {
			return loadPhoto, []*params.MessagesSendBuilder{}, err
		}

		return isPeoplePresentPhoto, nil, nil
	}

	return loadPhoto, nil, nil
}

func (state LoadPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Загрузите фото. Допустимые  форматы фото: jpg, jpeg, png, tif, tiff")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state LoadPhotoState) Name() stateName {
	return loadPhoto
}

// IsPeoplePresentPhotoState пользователь отвечает да или нет на вопрос есть ли на фото люди
type IsPeoplePresentPhotoState struct {
	postgres *postrgres.Repo
}

func (state IsPeoplePresentPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Назад":
		return loadPhoto, nil, nil
	case "Да":
		return isPeoplePresentPhoto, nil, nil
	case "Нет":
		return eventYearPhoto, nil, nil
	default:
		return isPeoplePresentPhoto, nil, nil
	}
}

func (state IsPeoplePresentPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("На фото есть люди?")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Да", "", "secondary")
	k.AddRow()
	k.AddTextButton("Нет", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state IsPeoplePresentPhotoState) Name() stateName {
	return isPeoplePresentPhoto
}

// EventYearPhotoState пользователь указывает год создания документа
type EventYearPhotoState struct {
	postgres *postrgres.Repo
}

func (state EventYearPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Назад":
		return isPeoplePresentPhoto, nil, nil
	case "Пропустить":
		return studyProgramPhoto, nil, nil
	default:
		year, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введён год недопустимого формата")
			return eventYearPhoto, []*params.MessagesSendBuilder{b}, nil
		}

		currentYear := time.Now().Year()
		if !(year >= 1900 && year <= currentYear) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введён несуществующий год")
			return eventYearPhoto, []*params.MessagesSendBuilder{b}, nil
		}
		err = state.postgres.RequestPhoto.UpdateYear(ctx, msg.PeerID, year)
		if err != nil {
			return eventYearPhoto, []*params.MessagesSendBuilder{}, err
		}
		return studyProgramPhoto, nil, nil
	}
}

func (state EventYearPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите год события в формате YYYY")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Пропустить", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EventYearPhotoState) Name() stateName {
	return eventYearPhoto
}

// StudyProgramPhotoState пользователь указывает программу обучения
type StudyProgramPhotoState struct {
	postgres *postrgres.Repo
}

func (state StudyProgramPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	var studyProg string

	switch messageText {
	case "Бакалавриат":
		studyProg = "Бакалавриат"
	case "Магистратура":
		studyProg = "Магистратура"
	case "Специалитет":
		studyProg = "Специалитет"
	case "Аспирантура":
		studyProg = "Аспирантура"
	case "Пропустить":
		return eventNamePhoto, nil, nil
	case "Назад":
		return eventYearPhoto, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Такой программы обучения нет в предложенных вариантах")
		return eventYearPhoto, []*params.MessagesSendBuilder{b}, nil
	}

	err := state.postgres.RequestPhoto.UpdateStudyProgram(ctx, msg.PeerID, studyProg)
	if err != nil {
		return studyProgramPhoto, []*params.MessagesSendBuilder{}, err
	}
	return eventNamePhoto, nil, nil
}

func (state StudyProgramPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите год события в формате YYYY")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Бакалавриат", "", "secondary")
	k.AddRow()
	k.AddTextButton("Магистратура", "", "secondary")
	k.AddRow()
	k.AddTextButton("Специалитет", "", "secondary")
	k.AddRow()
	k.AddTextButton("Аспирантура", "", "secondary")
	k.AddRow()
	k.AddTextButton("Пропустить", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state StudyProgramPhotoState) Name() stateName {
	return studyProgramPhoto
}

// EventNamePhotoState пользователь указывает существующее название события
type EventNamePhotoState struct {
	postgres *postrgres.Repo
}

func (state EventNamePhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Своё событие":
		return userEventNamePhoto, nil, nil
	case "Пропустить":
		return descriptionPhoto, nil, nil
	case "Назад":
		return studyProgramPhoto, nil, nil
	default:
		maxID, err := state.postgres.RequestPhoto.GetEventMaxID()
		if err != nil {
			return eventNamePhoto, []*params.MessagesSendBuilder{}, err
		}

		eventNumber, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введите номер события числом")
			return eventNamePhoto, []*params.MessagesSendBuilder{b}, nil
		}

		if !(eventNumber >= 1 && eventNumber <= maxID) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Такого события нет в списке")
			return eventNamePhoto, []*params.MessagesSendBuilder{b}, nil
		}

		err = state.postgres.RequestPhoto.UpdateEvent(ctx, msg.PeerID, eventNumber)
		if err != nil {
			return eventNamePhoto, []*params.MessagesSendBuilder{}, err
		}
		return descriptionPhoto, nil, nil
	}
}

func (state EventNamePhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	categories, err := state.postgres.RequestPhoto.GetEventNames()
	if err != nil {
		return []*params.MessagesSendBuilder{}, err
	}
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите номер события из списка ниже:\n" + categories)
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Своё событие", "", "secondary")
	k.AddRow()
	k.AddTextButton("Пропустить", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EventNamePhotoState) Name() stateName {
	return eventNamePhoto
}

// UserEventNamePhotoState пользователь указывает своё название события
type UserEventNamePhotoState struct {
	postgres *postrgres.Repo
}

func (state UserEventNamePhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Пропустить":
		return descriptionPhoto, nil, nil
	case "Назад":
		return eventNamePhoto, nil, nil
	default:
		err := state.postgres.RequestPhoto.UpdateUserEvent(ctx, msg.PeerID, messageText)
		if err != nil {
			return userEventNamePhoto, []*params.MessagesSendBuilder{}, err
		}
		return descriptionPhoto, nil, nil
	}
}

func (state UserEventNamePhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите название своего события. Оно будет рассмотрено администратором")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Пропустить", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state UserEventNamePhotoState) Name() stateName {
	return userEventNamePhoto
}

// DescriptionPhotoState пользователь вводит описание фотографии
type DescriptionPhotoState struct {
	postgres *postrgres.Repo
}

func (state DescriptionPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Пропустить":
		return descriptionPhoto, nil, nil
	case "Назад":
		return eventNamePhoto, nil, nil
	default:
		err := state.postgres.RequestPhoto.UpdateDescription(ctx, msg.PeerID, messageText)
		if err != nil {
			return descriptionPhoto, []*params.MessagesSendBuilder{}, err
		}
		return descriptionPhoto, nil, nil
	}
}

func (state DescriptionPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите описание фотографии")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Пропустить", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state DescriptionPhotoState) Name() stateName {
	return descriptionPhoto
}
