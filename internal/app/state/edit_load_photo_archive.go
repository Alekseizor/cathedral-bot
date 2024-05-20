package state

import (
	"context"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	"strconv"
	"time"
)

// EditPhotoArchiveState пользователь выбирает параметр для редактирования
type EditPhotoArchiveState struct {
	postgres *postrgres.Repo
}

func (state EditPhotoArchiveState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Год события":
		return editEventYearPhotoArchive, nil, nil
	case "Программа обучения":
		return editStudyProgramPhotoArchive, nil, nil
	case "Название события":
		return editEventNamePhotoArchive, nil, nil
	case "Описание":
		return editDescriptionPhotoArchive, nil, nil
	case "Назад":
		return checkPhotoArchive, nil, nil
	default:
		return editPhotoArchive, nil, nil
	}
}

func (state EditPhotoArchiveState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
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
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditPhotoArchiveState) Name() stateName {
	return editPhotoArchive
}

// EditEventYearPhotoArchiveState пользователь редактирует год события для фотографии
type EditEventYearPhotoArchiveState struct {
	postgres *postrgres.Repo
}

func (state EditEventYearPhotoArchiveState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhotoArchive.GetPhotoLastID(ctx, msg.PeerID)
	if err != nil {
		return editEventYearPhotoArchive, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return editPhotoArchive, nil, nil
	default:
		year, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введён год недопустимого формата")
			return editEventYearPhotoArchive, []*params.MessagesSendBuilder{b}, nil
		}

		currentYear := time.Now().Year()
		if !(year >= 1800 && year <= currentYear) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введён несуществующий год")
			return editEventYearPhotoArchive, []*params.MessagesSendBuilder{b}, nil
		}

		err = state.postgres.RequestPhotoArchive.UpdateYear(ctx, photoID, year)
		if err != nil {
			return editEventYearPhotoArchive, []*params.MessagesSendBuilder{}, err
		}
		return checkPhotoArchive, nil, nil
	}
}

func (state EditEventYearPhotoArchiveState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите год события в формате YYYY")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditEventYearPhotoArchiveState) Name() stateName {
	return editEventYearPhotoArchive
}

// EditStudyProgramPhotoArchiveState пользователь редактирует программу обучения
type EditStudyProgramPhotoArchiveState struct {
	postgres *postrgres.Repo
}

func (state EditStudyProgramPhotoArchiveState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhotoArchive.GetPhotoLastID(ctx, msg.PeerID)
	if err != nil {
		return editStudyProgramPhotoArchive, []*params.MessagesSendBuilder{}, err
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
		return editPhotoArchive, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Такой программы обучения нет в предложенных вариантах")
		return editStudyProgramPhotoArchive, []*params.MessagesSendBuilder{b}, nil
	}

	err = state.postgres.RequestPhotoArchive.UpdateStudyProgram(ctx, photoID, educationProgram)
	if err != nil {
		return editStudyProgramPhotoArchive, []*params.MessagesSendBuilder{}, err
	}
	return checkPhotoArchive, nil, nil
}

func (state EditStudyProgramPhotoArchiveState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
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

func (state EditStudyProgramPhotoArchiveState) Name() stateName {
	return editStudyProgramPhotoArchive
}

// EditEventNamePhotoArchiveState пользователь редактирует название события, выбирая из существующих
type EditEventNamePhotoArchiveState struct {
	postgres *postrgres.Repo
}

func (state EditEventNamePhotoArchiveState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhotoArchive.GetPhotoLastID(ctx, msg.PeerID)
	if err != nil {
		return editEventNamePhotoArchive, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Своё событие":
		return editUserEventNamePhotoArchive, nil, nil
	case "Назад":
		return editPhotoArchive, nil, nil
	default:
		maxID, err := state.postgres.RequestPhotoArchive.GetEventMaxID()
		if err != nil {
			return editEventNamePhotoArchive, []*params.MessagesSendBuilder{}, err
		}

		eventNumber, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введите номер события числом")
			return editEventNamePhotoArchive, []*params.MessagesSendBuilder{b}, nil
		}

		if !(eventNumber >= 1 && eventNumber <= maxID) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Такого события нет в списке")
			return editEventNamePhotoArchive, []*params.MessagesSendBuilder{b}, nil
		}

		err = state.postgres.RequestPhotoArchive.UpdateEvent(ctx, photoID, eventNumber-1)
		if err != nil {
			return editEventNamePhotoArchive, []*params.MessagesSendBuilder{}, err
		}
		return checkPhotoArchive, nil, nil
	}
}

func (state EditEventNamePhotoArchiveState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	categories, err := state.postgres.RequestPhotoArchive.GetEventNames()
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

func (state EditEventNamePhotoArchiveState) Name() stateName {
	return editEventNamePhotoArchive
}

// EditUserEventNamePhotoArchiveState пользователь редактирует название события, предлагая своё название
type EditUserEventNamePhotoArchiveState struct {
	postgres *postrgres.Repo
}

func (state EditUserEventNamePhotoArchiveState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return editUserEventNamePhotoArchive, nil, nil
	}

	photoID, err := state.postgres.RequestPhotoArchive.GetPhotoLastID(ctx, msg.PeerID)
	if err != nil {
		return editUserEventNamePhotoArchive, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return editEventNamePhotoArchive, nil, nil
	default:
		err := state.postgres.RequestPhotoArchive.UpdateUserEvent(ctx, photoID, messageText)
		if err != nil {
			return editUserEventNamePhotoArchive, []*params.MessagesSendBuilder{}, err
		}
		return checkPhotoArchive, nil, nil
	}
}

func (state EditUserEventNamePhotoArchiveState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите название своего события")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditUserEventNamePhotoArchiveState) Name() stateName {
	return editUserEventNamePhotoArchive
}

// EditDescriptionPhotoArchiveState пользователь редактирует описание фотографии
type EditDescriptionPhotoArchiveState struct {
	postgres *postrgres.Repo
}

func (state EditDescriptionPhotoArchiveState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return editDescriptionPhotoArchive, nil, nil
	}

	photoID, err := state.postgres.RequestPhotoArchive.GetPhotoLastID(ctx, msg.PeerID)
	if err != nil {
		return editDescriptionPhotoArchive, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return editPhotoArchive, nil, nil
	default:
		err := state.postgres.RequestPhotoArchive.UpdateDescription(ctx, photoID, messageText)
		if err != nil {
			return editDescriptionPhotoArchive, []*params.MessagesSendBuilder{}, err
		}
		return checkPhotoArchive, nil, nil
	}
}

func (state EditDescriptionPhotoArchiveState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите общее описание фотографий в архиве")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditDescriptionPhotoArchiveState) Name() stateName {
	return editDescriptionPhotoArchive
}
