package state

import (
	"context"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	"strconv"
	"time"
)

// EditPhotoState пользователь выбирает параметр для редактирования
type EditPhotoState struct {
	postgres *postrgres.Repo
}

func (state EditPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Год события":
		return editEventYearPhoto, nil, nil
	case "Программа обучения":
		return editStudyProgramPhoto, nil, nil
	case "Название события":
		return editEventNamePhoto, nil, nil
	case "Описание":
		return editDescriptionPhoto, nil, nil
	case "Отмеченные люди":
		return editIsPeoplePresentPhoto, nil, nil
	case "Назад":
		return checkPhoto, nil, nil
	default:
		return editPhoto, nil, nil
	}
}

func (state EditPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Выберите параметр для редактирования")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Год события", "", "secondary")
	k.AddTextButton("Программа обучения", "", "secondary")
	k.AddRow()
	k.AddTextButton("Название события", "", "secondary")
	k.AddTextButton("Описание", "", "secondary")
	k.AddRow()
	k.AddTextButton("Отмеченные люди", "", "secondary")
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditPhotoState) Name() stateName {
	return editPhoto
}

// EditEventYearPhotoState пользователь редактирует год события для фотографии
type EditEventYearPhotoState struct {
	postgres *postrgres.Repo
}

func (state EditEventYearPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhoto.GetPhotoLastID(ctx, msg.PeerID)
	if err != nil {
		return editEventYearPhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return editPhoto, nil, nil
	default:
		year, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введён год недопустимого формата")
			return editEventYearPhoto, []*params.MessagesSendBuilder{b}, nil
		}

		currentYear := time.Now().Year()
		if !(year >= 1800 && year <= currentYear) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введён несуществующий год")
			return editEventYearPhoto, []*params.MessagesSendBuilder{b}, nil
		}

		err = state.postgres.RequestPhoto.UpdateYear(ctx, photoID, year)
		if err != nil {
			return editEventYearPhoto, []*params.MessagesSendBuilder{}, err
		}
		return checkPhoto, nil, nil
	}
}

func (state EditEventYearPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите год события в формате YYYY")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditEventYearPhotoState) Name() stateName {
	return editEventYearPhoto
}

// EditStudyProgramPhotoState пользователь редактирует программу обучения
type EditStudyProgramPhotoState struct {
	postgres *postrgres.Repo
}

func (state EditStudyProgramPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhoto.GetPhotoLastID(ctx, msg.PeerID)
	if err != nil {
		return editStudyProgramPhoto, []*params.MessagesSendBuilder{}, err
	}

	var educationProgram string

	switch messageText {
	case "Бакалавриат":
		educationProgram = "Бакалавриат"
	case "Магистратура":
		educationProgram = "Магистратура"
	case "Специалитет":
		educationProgram = "Специалитет"
	case "Аспирантура":
		educationProgram = "Аспирантура"
	case "Назад":
		return editPhoto, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Такой программы обучения нет в предложенных вариантах")
		return editStudyProgramPhoto, []*params.MessagesSendBuilder{b}, nil
	}

	err = state.postgres.RequestPhoto.UpdateStudyProgram(ctx, photoID, educationProgram)
	if err != nil {
		return editStudyProgramPhoto, []*params.MessagesSendBuilder{}, err
	}
	return checkPhoto, nil, nil
}

func (state EditStudyProgramPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Выберите программу обучения")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Бакалавриат", "", "secondary")
	k.AddTextButton("Магистратура", "", "secondary")
	k.AddRow()
	k.AddTextButton("Специалитет", "", "secondary")
	k.AddTextButton("Аспирантура", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditStudyProgramPhotoState) Name() stateName {
	return editStudyProgramPhoto
}

// EditEventNamePhotoState пользователь редактирует название события, выбирая из существующих
type EditEventNamePhotoState struct {
	postgres *postrgres.Repo
}

func (state EditEventNamePhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhoto.GetPhotoLastID(ctx, msg.PeerID)
	if err != nil {
		return editEventNamePhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Своё событие":
		return editUserEventNamePhoto, nil, nil
	case "Назад":
		return editPhoto, nil, nil
	default:
		maxID, err := state.postgres.RequestPhoto.GetEventMaxID()
		if err != nil {
			return editEventNamePhoto, []*params.MessagesSendBuilder{}, err
		}

		eventNumber, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введите номер события числом")
			return editEventNamePhoto, []*params.MessagesSendBuilder{b}, nil
		}

		if !(eventNumber >= 1 && eventNumber <= maxID) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Такого события нет в списке")
			return editEventNamePhoto, []*params.MessagesSendBuilder{b}, nil
		}

		err = state.postgres.RequestPhoto.UpdateEvent(ctx, photoID, eventNumber-1)
		if err != nil {
			return editEventNamePhoto, []*params.MessagesSendBuilder{}, err
		}
		return checkPhoto, nil, nil
	}
}

func (state EditEventNamePhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
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
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditEventNamePhotoState) Name() stateName {
	return editEventNamePhoto
}

// EditUserEventNamePhotoState пользователь редактирует название события, предлагая своё название
type EditUserEventNamePhotoState struct {
	postgres *postrgres.Repo
}

func (state EditUserEventNamePhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return editUserEventNamePhoto, nil, nil
	}

	photoID, err := state.postgres.RequestPhoto.GetPhotoLastID(ctx, msg.PeerID)
	if err != nil {
		return editUserEventNamePhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return editEventNamePhoto, nil, nil
	default:
		err := state.postgres.RequestPhoto.UpdateUserEvent(ctx, photoID, messageText)
		if err != nil {
			return editUserEventNamePhoto, []*params.MessagesSendBuilder{}, err
		}
		return checkPhoto, nil, nil
	}
}

func (state EditUserEventNamePhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите название своего события")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditUserEventNamePhotoState) Name() stateName {
	return editUserEventNamePhoto
}

// EditDescriptionPhotoState пользователь редактирует описание фотографии
type EditDescriptionPhotoState struct {
	postgres *postrgres.Repo
}

func (state EditDescriptionPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return editDescriptionPhoto, nil, nil
	}

	photoID, err := state.postgres.RequestPhoto.GetPhotoLastID(ctx, msg.PeerID)
	if err != nil {
		return editDescriptionPhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return editPhoto, nil, nil
	default:
		err := state.postgres.RequestPhoto.UpdateDescription(ctx, photoID, messageText)
		if err != nil {
			return editDescriptionPhoto, []*params.MessagesSendBuilder{}, err
		}
		return checkPhoto, nil, nil
	}
}

func (state EditDescriptionPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите описание фотографии")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditDescriptionPhotoState) Name() stateName {
	return editDescriptionPhoto
}
