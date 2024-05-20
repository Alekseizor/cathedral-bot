package state

import (
	"context"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	"strconv"
	"time"
)

// EditRequestPhotoState пользователь выбирает параметр для редактирования
type EditRequestPhotoState struct {
	postgres *postrgres.Repo
}

func (state EditRequestPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Год события":
		return editEventYearRequestPhoto, nil, nil
	case "Программа обучения":
		return editStudyProgramRequestPhoto, nil, nil
	case "Название события":
		return editEventNameRequestPhoto, nil, nil
	case "Описание":
		return editDescriptionRequestPhoto, nil, nil
	case "Отмеченные люди":
		return editIsPeoplePresentRequestPhoto, nil, nil
	case "Назад":
		return viewRequestsPhoto, nil, nil
	default:
		return editRequestPhoto, nil, nil
	}
}

func (state EditRequestPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
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

func (state EditRequestPhotoState) Name() stateName {
	return editRequestPhoto
}

// EditEventYearRequestPhotoState пользователь редактирует год события для фотографии
type EditEventYearRequestPhotoState struct {
	postgres *postrgres.Repo
}

func (state EditEventYearRequestPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhoto.GetPhotoRequestID(ctx, msg.PeerID)
	if err != nil {
		return editEventYearRequestPhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return editRequestPhoto, nil, nil
	default:
		year, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введён год недопустимого формата")
			return editEventYearRequestPhoto, []*params.MessagesSendBuilder{b}, nil
		}

		currentYear := time.Now().Year()
		if !(year >= 1800 && year <= currentYear) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введён несуществующий год")
			return editEventYearRequestPhoto, []*params.MessagesSendBuilder{b}, nil
		}

		err = state.postgres.RequestPhoto.UpdateYear(ctx, photoID, year)
		if err != nil {
			return editEventYearRequestPhoto, []*params.MessagesSendBuilder{}, err
		}
		return viewRequestsPhoto, nil, nil
	}
}

func (state EditEventYearRequestPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите год события в формате YYYY")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditEventYearRequestPhotoState) Name() stateName {
	return editEventYearRequestPhoto
}

// EditStudyProgramRequestPhotoState пользователь редактирует программу обучения
type EditStudyProgramRequestPhotoState struct {
	postgres *postrgres.Repo
}

func (state EditStudyProgramRequestPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhoto.GetPhotoRequestID(ctx, msg.PeerID)
	if err != nil {
		return editStudyProgramRequestPhoto, []*params.MessagesSendBuilder{}, err
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
		return editRequestPhoto, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Такой программы обучения нет в предложенных вариантах")
		return editStudyProgramRequestPhoto, []*params.MessagesSendBuilder{b}, nil
	}

	err = state.postgres.RequestPhoto.UpdateStudyProgram(ctx, photoID, educationProgram)
	if err != nil {
		return editStudyProgramRequestPhoto, []*params.MessagesSendBuilder{}, err
	}
	return viewRequestsPhoto, nil, nil
}

func (state EditStudyProgramRequestPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
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

func (state EditStudyProgramRequestPhotoState) Name() stateName {
	return editStudyProgramRequestPhoto
}

// EditEventNameRequestPhotoState пользователь редактирует название события, выбирая из существующих
type EditEventNameRequestPhotoState struct {
	postgres *postrgres.Repo
}

func (state EditEventNameRequestPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhoto.GetPhotoRequestID(ctx, msg.PeerID)
	if err != nil {
		return editEventNameRequestPhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Своё событие":
		return editUserEventNameRequestPhoto, nil, nil
	case "Назад":
		return editRequestPhoto, nil, nil
	default:
		maxID, err := state.postgres.RequestPhoto.GetEventMaxID()
		if err != nil {
			return editEventNameRequestPhoto, []*params.MessagesSendBuilder{}, err
		}

		eventNumber, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введите номер события числом")
			return editEventNameRequestPhoto, []*params.MessagesSendBuilder{b}, nil
		}

		if !(eventNumber >= 1 && eventNumber <= maxID) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Такого события нет в списке")
			return editEventNameRequestPhoto, []*params.MessagesSendBuilder{b}, nil
		}

		err = state.postgres.RequestPhoto.UpdateEvent(ctx, photoID, eventNumber-1)
		if err != nil {
			return editEventNameRequestPhoto, []*params.MessagesSendBuilder{}, err
		}
		return viewRequestsPhoto, nil, nil
	}
}

func (state EditEventNameRequestPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
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

func (state EditEventNameRequestPhotoState) Name() stateName {
	return editEventNameRequestPhoto
}

// EditUserEventNameRequestPhotoState пользователь редактирует название события, предлагая своё название
type EditUserEventNameRequestPhotoState struct {
	postgres *postrgres.Repo
}

func (state EditUserEventNameRequestPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return editUserEventNameRequestPhoto, nil, nil
	}

	photoID, err := state.postgres.RequestPhoto.GetPhotoRequestID(ctx, msg.PeerID)
	if err != nil {
		return editUserEventNameRequestPhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return editEventNameRequestPhoto, nil, nil
	default:
		err := state.postgres.RequestPhoto.UpdateUserEvent(ctx, photoID, messageText)
		if err != nil {
			return editUserEventNameRequestPhoto, []*params.MessagesSendBuilder{}, err
		}
		return viewRequestsPhoto, nil, nil
	}
}

func (state EditUserEventNameRequestPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите название своего события")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditUserEventNameRequestPhotoState) Name() stateName {
	return editUserEventNameRequestPhoto
}

// EditDescriptionRequestPhotoState пользователь редактирует описание фотографии
type EditDescriptionRequestPhotoState struct {
	postgres *postrgres.Repo
}

func (state EditDescriptionRequestPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return editDescriptionRequestPhoto, nil, nil
	}

	photoID, err := state.postgres.RequestPhoto.GetPhotoRequestID(ctx, msg.PeerID)
	if err != nil {
		return editDescriptionRequestPhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return editRequestPhoto, nil, nil
	default:
		err := state.postgres.RequestPhoto.UpdateDescription(ctx, photoID, messageText)
		if err != nil {
			return editDescriptionRequestPhoto, []*params.MessagesSendBuilder{}, err
		}
		return viewRequestsPhoto, nil, nil
	}
}

func (state EditDescriptionRequestPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите описание фотографии")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditDescriptionRequestPhotoState) Name() stateName {
	return editDescriptionRequestPhoto
}
