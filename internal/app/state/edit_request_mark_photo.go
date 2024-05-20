package state

import (
	"context"
	"fmt"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	"strconv"
)

// EditIsPeoplePresentRequestPhotoState пользователь редактирует ответ на вопрос есть ли на фото люди
type EditIsPeoplePresentRequestPhotoState struct {
	postgres *postrgres.Repo
}

func (state EditIsPeoplePresentRequestPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhoto.GetPhotoRequestID(ctx, msg.PeerID)
	if err != nil {
		return editIsPeoplePresentRequestPhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Да":
		return editCountPeopleRequestPhoto, nil, nil
	case "Нет":
		err = state.postgres.RequestPhoto.DeleteMarksOnPhoto(ctx, photoID)
		if err != nil {
			return editIsPeoplePresentRequestPhoto, []*params.MessagesSendBuilder{}, err
		}
		return requestPhotoFromQueue, nil, nil
	case "Назад":
		return editPhoto, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Выберите из предложенных вариантов")
		return editIsPeoplePresentRequestPhoto, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state EditIsPeoplePresentRequestPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("На фото есть люди?")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Да", "", "secondary")
	k.AddTextButton("Нет", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditIsPeoplePresentRequestPhotoState) Name() stateName {
	return editIsPeoplePresentRequestPhoto
}

// EditCountPeopleRequestPhotoState пользователь редактирует количество людей на фото
type EditCountPeopleRequestPhotoState struct {
	postgres *postrgres.Repo
}

func (state EditCountPeopleRequestPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhoto.GetPhotoRequestID(ctx, msg.PeerID)
	if err != nil {
		return editCountPeopleRequestPhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return editIsPeoplePresentRequestPhoto, nil, nil
	default:
		count, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Количество людей должно быть целым числом")
			return editCountPeopleRequestPhoto, []*params.MessagesSendBuilder{b}, nil
		}

		if count < 1 {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Количество людей должно быть больше нуля")
			return editCountPeopleRequestPhoto, []*params.MessagesSendBuilder{b}, nil
		}
		err = state.postgres.RequestPhoto.UpdateCountPeople(ctx, photoID, count)
		if err != nil {
			return editCountPeopleRequestPhoto, []*params.MessagesSendBuilder{}, err
		}
		return editMarkedPeopleRequestPhoto, nil, nil
	}
}

func (state EditCountPeopleRequestPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите количество людей на фото")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditCountPeopleRequestPhotoState) Name() stateName {
	return editCountPeopleRequestPhoto
}

// EditMarkedPeopleRequestPhotoState пользователь редактирует отметки людей на фото
type EditMarkedPeopleRequestPhotoState struct {
	postgres *postrgres.Repo
}

func (state EditMarkedPeopleRequestPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhoto.GetPhotoRequestID(ctx, msg.PeerID)
	if err != nil {
		return editMarkedPeopleRequestPhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Знаю ФИО человека":
		return editIsTeacherRequestPhoto, nil, nil
	case "Не знаю ФИО человека":
		allMarked, err := state.postgres.RequestPhoto.UpdateMarkedPeople(ctx, photoID, "???")
		if err != nil {
			return editMarkedPeopleRequestPhoto, []*params.MessagesSendBuilder{}, err
		}
		if allMarked {
			return requestPhotoFromQueue, nil, nil
		}
		return editMarkedPeopleRequestPhoto, nil, nil
	case "Закончить отмечать людей":
		return requestPhotoFromQueue, nil, nil
	case "Назад":
		return editCountPeopleRequestPhoto, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Выберите из предложенных вариантов")
		return editMarkedPeopleRequestPhoto, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state EditMarkedPeopleRequestPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	photoID, err := state.postgres.RequestPhoto.GetPhotoRequestID(ctx, vkID)
	if err != nil {
		return []*params.MessagesSendBuilder{}, err
	}

	markedPerson, err := state.postgres.RequestPhoto.GetMarkedPerson(ctx, photoID)
	if err != nil {
		return []*params.MessagesSendBuilder{}, err
	}

	msg := fmt.Sprintf("Отметьте %v-го человека слева", markedPerson+1)
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message(msg)
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Знаю ФИО человека", "", "secondary")
	k.AddRow()
	k.AddTextButton("Не знаю ФИО человека", "", "secondary")
	k.AddRow()
	k.AddTextButton("Закончить отмечать людей", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditMarkedPeopleRequestPhotoState) Name() stateName {
	return editMarkedPeopleRequestPhoto
}

// EditIsTeacherRequestPhotoState пользователь редактирует ответ на вопрос учитель ли это
type EditIsTeacherRequestPhotoState struct {
	postgres *postrgres.Repo
}

func (state EditIsTeacherRequestPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Да":
		return editTeacherNameRequestPhoto, nil, nil
	case "Нет":
		return editStudentNameRequestPhoto, nil, nil
	case "Назад":
		return editMarkedPeopleRequestPhoto, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Выберите из предложенных вариантов")
		return editIsTeacherRequestPhoto, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state EditIsTeacherRequestPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Это преподаватель?")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Да", "", "secondary")
	k.AddRow()
	k.AddTextButton("Нет", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditIsTeacherRequestPhotoState) Name() stateName {
	return editIsTeacherRequestPhoto
}

// EditTeacherNameRequestPhotoState пользователь выбирает преподавателя из списка
type EditTeacherNameRequestPhotoState struct {
	postgres *postrgres.Repo
}

func (state EditTeacherNameRequestPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhoto.GetPhotoRequestID(ctx, msg.PeerID)
	if err != nil {
		return editTeacherNameRequestPhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Ввести ФИО вручную":
		return editUserTeacherNameRequestPhoto, nil, nil
	case "Назад":
		err = state.postgres.RequestPhoto.DeletePointer(msg.PeerID)
		if err != nil {
			return editTeacherNameRequestPhoto, nil, err
		}
		return editIsTeacherRequestPhoto, nil, nil
	case "⬅️":
		err = state.postgres.RequestPhoto.ChangePointer(msg.PeerID, false)
		if err != nil {
			return editTeacherNameRequestPhoto, nil, err
		}
		return editTeacherNameRequestPhoto, nil, nil
	case "➡️":
		err = state.postgres.RequestPhoto.ChangePointer(msg.PeerID, true)
		if err != nil {
			return editTeacherNameRequestPhoto, nil, err
		}
		return editTeacherNameRequestPhoto, nil, nil
	default:
		teacherID, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введите номер преподавателя числом")
			return editTeacherNameRequestPhoto, []*params.MessagesSendBuilder{b}, nil
		}

		maxID, err := state.postgres.RequestPhoto.GetTeacherMaxID()
		if err != nil {
			return editTeacherNameRequestPhoto, []*params.MessagesSendBuilder{}, err
		}

		if !(teacherID >= 1 && teacherID <= maxID) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Такого преподавателя нет в списке")
			return editTeacherNameRequestPhoto, []*params.MessagesSendBuilder{b}, nil
		}

		teacherName, err := state.postgres.RequestPhoto.GetTeacherName(ctx, teacherID)
		if err != nil {
			return editTeacherNameRequestPhoto, []*params.MessagesSendBuilder{}, err
		}

		err = state.postgres.RequestPhoto.UpdateTeachers(ctx, photoID, teacherName)
		if err != nil {
			return editTeacherNameRequestPhoto, []*params.MessagesSendBuilder{}, err
		}

		allMarked, err := state.postgres.RequestPhoto.UpdateMarkedPeople(ctx, photoID, teacherName)
		if err != nil {
			return editTeacherNameRequestPhoto, []*params.MessagesSendBuilder{}, err
		}

		err = state.postgres.RequestPhoto.DeletePointer(msg.PeerID)
		if err != nil {
			return editTeacherNameRequestPhoto, nil, err
		}

		if allMarked {
			return requestPhotoFromQueue, nil, nil
		}

		return editMarkedPeopleRequestPhoto, nil, nil
	}
}

func (state EditTeacherNameRequestPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	teacherNames, pointer, count, err := state.postgres.RequestPhoto.GetTeacherNames(vkID)
	if err != nil {
		return []*params.MessagesSendBuilder{}, err
	}

	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите номер преподавателя из списка ниже:\n" + teacherNames)
	k := object.NewMessagesKeyboard(true)
	if count > 10 {
		k.AddRow()
		if pointer != 0 {
			k.AddTextButton("⬅️", "", "secondary")
		}
		if count-pointer > 10 {
			k.AddTextButton("➡️", "", "secondary")
		}
	}
	k.AddRow()
	k.AddTextButton("Ввести ФИО вручную", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditTeacherNameRequestPhotoState) Name() stateName {
	return editTeacherNameRequestPhoto
}

// EditUserTeacherNameRequestPhotoState пользователь вводит ФИО преподавателя
type EditUserTeacherNameRequestPhotoState struct {
	postgres *postrgres.Repo
}

func (state EditUserTeacherNameRequestPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return editUserTeacherNameRequestPhoto, nil, nil
	}

	photoID, err := state.postgres.RequestPhoto.GetPhotoRequestID(ctx, msg.PeerID)
	if err != nil {
		return editUserTeacherNameRequestPhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return editTeacherNameRequestPhoto, nil, nil
	default:
		err = state.postgres.RequestPhoto.UpdateTeachers(ctx, photoID, messageText)
		if err != nil {
			return editUserTeacherNameRequestPhoto, []*params.MessagesSendBuilder{}, err
		}

		allMarked, err := state.postgres.RequestPhoto.UpdateMarkedPeople(ctx, photoID, messageText)
		if err != nil {
			return editUserTeacherNameRequestPhoto, []*params.MessagesSendBuilder{}, err
		}

		if allMarked {
			return requestPhotoFromQueue, nil, nil
		}

		return editMarkedPeopleRequestPhoto, nil, nil
	}
}

func (state EditUserTeacherNameRequestPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите ФИО преподавателя")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditUserTeacherNameRequestPhotoState) Name() stateName {
	return editUserTeacherNameRequestPhoto
}

// EditStudentNameRequestPhotoState пользователь вводит ФИО студента
type EditStudentNameRequestPhotoState struct {
	postgres *postrgres.Repo
}

func (state EditStudentNameRequestPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return editStudentNameRequestPhoto, nil, nil
	}

	photoID, err := state.postgres.RequestPhoto.GetPhotoRequestID(ctx, msg.PeerID)
	if err != nil {
		return editStudentNameRequestPhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return editTeacherNameRequestPhoto, nil, nil
	default:
		allMarked, err := state.postgres.RequestPhoto.UpdateMarkedPeople(ctx, photoID, messageText)
		if err != nil {
			return editStudentNameRequestPhoto, []*params.MessagesSendBuilder{}, err
		}

		if allMarked {
			return requestPhotoFromQueue, nil, nil
		}

		return editMarkedPeopleRequestPhoto, nil, nil
	}
}

func (state EditStudentNameRequestPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите фамилию и имя студента")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditStudentNameRequestPhotoState) Name() stateName {
	return editStudentNameRequestPhoto
}
