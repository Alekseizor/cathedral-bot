package state

import (
	"context"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	"strconv"
	"time"
)

// EditParamsPhotoState пользователь выбирает параметр для редактирования
type EditParamsPhotoState struct {
	postgres *postrgres.Repo
}

func (state EditParamsPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Год события":
		return editEventYearParamsPhoto, nil, nil
	case "Программа обучения":
		return editStudyProgramParamsPhoto, nil, nil
	case "Название события":
		return editEventNameParamsPhoto, nil, nil
	case "Описание":
		return editDescriptionParamsPhoto, nil, nil
	case "Отмеченные люди":
		return editIsPeoplePresentParamsPhoto, nil, nil
	case "Назад":
		return personalAccountPhoto, nil, nil
	default:
		return editParamsPhoto, nil, nil
	}
}

func (state EditParamsPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
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

func (state EditParamsPhotoState) Name() stateName {
	return editParamsPhoto
}

// EditEventYearParamsPhotoState пользователь редактирует год события для фотографии
type EditEventYearParamsPhotoState struct {
	postgres *postrgres.Repo
}

func (state EditEventYearParamsPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhoto.GetPhotoClientID(ctx, msg.PeerID)
	if err != nil {
		return editEventYearRequestPhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return editParamsPhoto, nil, nil
	default:
		year, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введён год недопустимого формата")
			return editEventYearParamsPhoto, []*params.MessagesSendBuilder{b}, nil
		}

		currentYear := time.Now().Year()
		if !(year >= 1800 && year <= currentYear) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введён несуществующий год")
			return editEventYearParamsPhoto, []*params.MessagesSendBuilder{b}, nil
		}

		err = state.postgres.RequestPhoto.UpdateYear(ctx, photoID, year)
		if err != nil {
			return editEventYearParamsPhoto, []*params.MessagesSendBuilder{}, err
		}
		return personalAccountPhoto, nil, nil
	}
}

func (state EditEventYearParamsPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите год события в формате YYYY")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditEventYearParamsPhotoState) Name() stateName {
	return editEventYearParamsPhoto
}

// EditStudyProgramParamsPhotoState пользователь редактирует программу обучения
type EditStudyProgramParamsPhotoState struct {
	postgres *postrgres.Repo
}

func (state EditStudyProgramParamsPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhoto.GetPhotoClientID(ctx, msg.PeerID)
	if err != nil {
		return editStudyProgramParamsPhoto, []*params.MessagesSendBuilder{}, err
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
		return editParamsPhoto, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Такой программы обучения нет в предложенных вариантах")
		return editStudyProgramParamsPhoto, []*params.MessagesSendBuilder{b}, nil
	}

	err = state.postgres.RequestPhoto.UpdateStudyProgram(ctx, photoID, educationProgram)
	if err != nil {
		return editStudyProgramParamsPhoto, []*params.MessagesSendBuilder{}, err
	}
	return personalAccountPhoto, nil, nil
}

func (state EditStudyProgramParamsPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
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

func (state EditStudyProgramParamsPhotoState) Name() stateName {
	return editStudyProgramParamsPhoto
}

// EditEventNameParamsPhotoState пользователь редактирует название события, выбирая из существующих
type EditEventNameParamsPhotoState struct {
	postgres *postrgres.Repo
}

func (state EditEventNameParamsPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhoto.GetPhotoClientID(ctx, msg.PeerID)
	if err != nil {
		return editEventNameParamsPhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Своё событие":
		return editUserEventNameParamsPhoto, nil, nil
	case "Назад":
		return editParamsPhoto, nil, nil
	default:
		maxID, err := state.postgres.RequestPhoto.GetEventMaxID()
		if err != nil {
			return editEventNameParamsPhoto, []*params.MessagesSendBuilder{}, err
		}

		eventNumber, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введите номер события числом")
			return editEventNameParamsPhoto, []*params.MessagesSendBuilder{b}, nil
		}

		if !(eventNumber >= 1 && eventNumber <= maxID) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Такого события нет в списке")
			return editEventNameParamsPhoto, []*params.MessagesSendBuilder{b}, nil
		}

		err = state.postgres.RequestPhoto.UpdateEvent(ctx, photoID, eventNumber-1)
		if err != nil {
			return editEventNameParamsPhoto, []*params.MessagesSendBuilder{}, err
		}
		return personalAccountPhoto, nil, nil
	}
}

func (state EditEventNameParamsPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
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

func (state EditEventNameParamsPhotoState) Name() stateName {
	return editEventNameParamsPhoto
}

// EditUserEventNameParamsPhotoState пользователь редактирует название события, предлагая своё название
type EditUserEventNameParamsPhotoState struct {
	postgres *postrgres.Repo
}

func (state EditUserEventNameParamsPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return editUserEventNameParamsPhoto, nil, nil
	}

	photoID, err := state.postgres.RequestPhoto.GetPhotoClientID(ctx, msg.PeerID)
	if err != nil {
		return editUserEventNameParamsPhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return editEventNameParamsPhoto, nil, nil
	default:
		err := state.postgres.RequestPhoto.UpdateUserEvent(ctx, photoID, messageText)
		if err != nil {
			return editUserEventNameParamsPhoto, []*params.MessagesSendBuilder{}, err
		}
		return personalAccountPhoto, nil, nil
	}
}

func (state EditUserEventNameParamsPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите название своего события")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditUserEventNameParamsPhotoState) Name() stateName {
	return editUserEventNameParamsPhoto
}

// EditDescriptionParamsPhotoState пользователь редактирует описание фотографии
type EditDescriptionParamsPhotoState struct {
	postgres *postrgres.Repo
}

func (state EditDescriptionParamsPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return editDescriptionParamsPhoto, nil, nil
	}

	photoID, err := state.postgres.RequestPhoto.GetPhotoClientID(ctx, msg.PeerID)
	if err != nil {
		return editDescriptionParamsPhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return editParamsPhoto, nil, nil
	default:
		err := state.postgres.RequestPhoto.UpdateDescription(ctx, photoID, messageText)
		if err != nil {
			return editDescriptionParamsPhoto, []*params.MessagesSendBuilder{}, err
		}
		return personalAccountPhoto, nil, nil
	}
}

func (state EditDescriptionParamsPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите описание фотографии")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditDescriptionParamsPhotoState) Name() stateName {
	return editDescriptionParamsPhoto
}
