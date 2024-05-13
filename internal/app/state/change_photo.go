package state

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
)

type ChangeEventYearPhotoState struct {
	postgres *postrgres.Repo
}

func (state ChangeEventYearPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.ObjectAdmin.Get(ctx, msg.PeerID)
	if err != nil {
		return changeEventYearPhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return workingAlbums, nil, nil
	default:
		year, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введён год недопустимого формата")
			return changeEventYearPhoto, []*params.MessagesSendBuilder{b}, nil
		}

		currentYear := time.Now().Year()
		if !(year >= 1800 && year <= currentYear) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введён несуществующий год")
			return changeEventYearPhoto, []*params.MessagesSendBuilder{b}, nil
		}

		err = state.postgres.StudentAlbums.UpdateYear(ctx, photoID, year)
		if err != nil {
			return changeEventYearPhoto, []*params.MessagesSendBuilder{}, err
		}

		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message(fmt.Sprintf("Для альбома №%d изменен год на - %d", photoID, year))
		return workingAlbums, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state ChangeEventYearPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите год события в формате YYYY")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state ChangeEventYearPhotoState) Name() stateName {
	return changeEventYearPhoto
}

type ChangeStudyProgramPhotoState struct {
	postgres *postrgres.Repo
}

func (state ChangeStudyProgramPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.ObjectAdmin.Get(ctx, msg.PeerID)
	if err != nil {
		return changeStudyProgramPhoto, []*params.MessagesSendBuilder{}, err
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
		return workingAlbums, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Такой программы обучения нет в предложенных вариантах")
		return changeStudyProgramPhoto, []*params.MessagesSendBuilder{b}, nil
	}

	err = state.postgres.StudentAlbums.UpdateStudyProgram(ctx, photoID, educationProgram)
	if err != nil {
		return changeStudyProgramPhoto, []*params.MessagesSendBuilder{}, err
	}

	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message(fmt.Sprintf("Для альбома №%d изменена программа обучения на - %s", photoID, educationProgram))
	return workingAlbums, []*params.MessagesSendBuilder{b}, nil
}

func (state ChangeStudyProgramPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
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
	addBackButton(k)
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state ChangeStudyProgramPhotoState) Name() stateName {
	return changeStudyProgramPhoto
}

type ChangeEventNamePhotoState struct {
	postgres *postrgres.Repo
}

func (state ChangeEventNamePhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.ObjectAdmin.Get(ctx, msg.PeerID)
	if err != nil {
		return changeEventNamePhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Своё событие":
		//TODO: реализовать создание нового события
		// новое событие должно быть присвоено альбому
		return changeUserEventNamePhoto, nil, nil
	case "Назад":
		return workingAlbums, nil, nil
	default:
		maxID, err := state.postgres.RequestPhoto.GetEventMaxID()
		if err != nil {
			return changeEventNamePhoto, []*params.MessagesSendBuilder{}, err
		}

		eventNumber, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введите номер события числом")
			return changeEventNamePhoto, []*params.MessagesSendBuilder{b}, nil
		}

		if !(eventNumber >= 1 && eventNumber <= maxID) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Такого события нет в списке")
			return changeEventNamePhoto, []*params.MessagesSendBuilder{b}, nil
		}

		err = state.postgres.StudentAlbums.UpdateEvent(ctx, photoID, eventNumber)
		if err != nil {
			return changeEventNamePhoto, []*params.MessagesSendBuilder{}, err
		}

		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message(fmt.Sprintf("Для альбома №%d изменен номер события на - %s", photoID, eventNumber))
		return workingAlbums, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state ChangeEventNamePhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
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
	addBackButton(k)
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state ChangeEventNamePhotoState) Name() stateName {
	return changeEventNamePhoto
}

type ChangeDescriptionPhotoState struct {
	postgres *postrgres.Repo
}

func (state ChangeDescriptionPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return changeDescriptionPhoto, nil, nil
	}

	photoID, err := state.postgres.ObjectAdmin.Get(ctx, msg.PeerID)
	if err != nil {
		return changeDescriptionPhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return workingAlbums, nil, nil
	default:
		err := state.postgres.StudentAlbums.UpdateDescription(ctx, photoID, messageText)
		if err != nil {
			return changeDescriptionPhoto, []*params.MessagesSendBuilder{}, err
		}
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message(fmt.Sprintf("Для альбома №%d изменено описание на - %s", photoID, messageText))
		return workingAlbums, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state ChangeDescriptionPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите описание альбома")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state ChangeDescriptionPhotoState) Name() stateName {
	return changeDescriptionPhoto
}

type ChangeUserEventNamePhotoState struct {
	postgres *postrgres.Repo
}

func (state ChangeUserEventNamePhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return changeUserEventNamePhoto, nil, nil
	}

	photoID, err := state.postgres.RequestPhoto.GetPhotoLastID(ctx, msg.PeerID)
	if err != nil {
		return changeUserEventNamePhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return changeEventNamePhoto, nil, nil
	default:
		err := state.postgres.RequestPhoto.UpdateUserEvent(ctx, photoID, messageText)
		if err != nil {
			return changeUserEventNamePhoto, []*params.MessagesSendBuilder{}, err
		}
		return checkPhoto, nil, nil
	}
}

func (state ChangeUserEventNamePhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите название своего события")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state ChangeUserEventNamePhotoState) Name() stateName {
	return changeUserEventNamePhoto
}
