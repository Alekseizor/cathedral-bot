package state

import (
	"context"
	"fmt"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	"strconv"
)

// EditIsPeoplePresentPhotoState пользователь редактирует ответ на вопрос есть ли на фото люди
type EditIsPeoplePresentPhotoState struct {
	postgres *postrgres.Repo
}

func (state EditIsPeoplePresentPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhoto.GetPhotoLastID(ctx, msg.PeerID)
	if err != nil {
		return editIsPeoplePresentPhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Да":
		return editCountPeoplePhoto, nil, nil
	case "Нет":
		err = state.postgres.RequestPhoto.DeleteMarksOnPhoto(ctx, photoID)
		if err != nil {
			return editIsPeoplePresentPhoto, []*params.MessagesSendBuilder{}, err
		}
		return checkPhoto, nil, nil
	case "Назад":
		return editPhoto, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Выберите из предложенных вариантов")
		return editIsPeoplePresentPhoto, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state EditIsPeoplePresentPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
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

func (state EditIsPeoplePresentPhotoState) Name() stateName {
	return editIsPeoplePresentPhoto
}

// EditCountPeoplePhotoState пользователь редактирует количество людей на фото
type EditCountPeoplePhotoState struct {
	postgres *postrgres.Repo
}

func (state EditCountPeoplePhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhoto.GetPhotoLastID(ctx, msg.PeerID)
	if err != nil {
		return editCountPeoplePhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return editIsPeoplePresentPhoto, nil, nil
	default:
		count, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Количество людей должно быть целым числом")
			return editCountPeoplePhoto, []*params.MessagesSendBuilder{b}, nil
		}

		if count < 1 {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Количество людей должно быть больше нуля")
			return editCountPeoplePhoto, []*params.MessagesSendBuilder{b}, nil
		}
		err = state.postgres.RequestPhoto.UpdateCountPeople(ctx, photoID, count)
		if err != nil {
			return editCountPeoplePhoto, []*params.MessagesSendBuilder{}, err
		}
		return editMarkedPeoplePhoto, nil, nil
	}
}

func (state EditCountPeoplePhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите количество людей на фото")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditCountPeoplePhotoState) Name() stateName {
	return editCountPeoplePhoto
}

// EditMarkedPeoplePhotoState пользователь редактирует отметки людей на фото
type EditMarkedPeoplePhotoState struct {
	postgres *postrgres.Repo
}

func (state EditMarkedPeoplePhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhoto.GetPhotoLastID(ctx, msg.PeerID)
	if err != nil {
		return editMarkedPeoplePhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Знаю ФИО человека":
		return editIsTeacherPhoto, nil, nil
	case "Не знаю ФИО человека":
		allMarked, err := state.postgres.RequestPhoto.UpdateMarkedPeople(ctx, photoID, "???")
		if err != nil {
			return editMarkedPeoplePhoto, []*params.MessagesSendBuilder{}, err
		}
		if allMarked {
			return checkPhoto, nil, nil
		}
		return editMarkedPeoplePhoto, nil, nil
	case "Закончить отмечать людей":
		return checkPhoto, nil, nil
	case "Назад":
		return editCountPeoplePhoto, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Выберите из предложенных вариантов")
		return editMarkedPeoplePhoto, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state EditMarkedPeoplePhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	photoID, err := state.postgres.RequestPhoto.GetPhotoLastID(ctx, vkID)
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

func (state EditMarkedPeoplePhotoState) Name() stateName {
	return editMarkedPeoplePhoto
}

// EditIsTeacherPhotoState пользователь редактирует ответ на вопрос учитель ли это
type EditIsTeacherPhotoState struct {
	postgres *postrgres.Repo
}

func (state EditIsTeacherPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Да":
		return editTeacherNamePhoto, nil, nil
	case "Нет":
		return editStudentNamePhoto, nil, nil
	case "Назад":
		return editMarkedPeoplePhoto, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Выберите из предложенных вариантов")
		return editIsTeacherPhoto, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state EditIsTeacherPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
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

func (state EditIsTeacherPhotoState) Name() stateName {
	return editIsTeacherPhoto
}

// EditTeacherNamePhotoState пользователь выбирает преподавателя из списка
type EditTeacherNamePhotoState struct {
	postgres *postrgres.Repo
}

func (state EditTeacherNamePhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhoto.GetPhotoLastID(ctx, msg.PeerID)
	if err != nil {
		return editTeacherNamePhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Ввести ФИО вручную":
		return editUserTeacherNamePhoto, nil, nil
	case "Назад":
		err = state.postgres.RequestPhoto.DeletePointer(msg.PeerID)
		if err != nil {
			return editTeacherNamePhoto, nil, err
		}
		return editIsTeacherPhoto, nil, nil
	case "⬅️":
		err = state.postgres.RequestPhoto.ChangePointer(msg.PeerID, false)
		if err != nil {
			return editTeacherNamePhoto, nil, err
		}
		return editTeacherNamePhoto, nil, nil
	case "➡️":
		err = state.postgres.RequestPhoto.ChangePointer(msg.PeerID, true)
		if err != nil {
			return editTeacherNamePhoto, nil, err
		}
		return editTeacherNamePhoto, nil, nil
	default:
		teacherID, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введите номер преподавателя числом")
			return editTeacherNamePhoto, []*params.MessagesSendBuilder{b}, nil
		}

		maxID, err := state.postgres.RequestPhoto.GetTeacherMaxID()
		if err != nil {
			return editTeacherNamePhoto, []*params.MessagesSendBuilder{}, err
		}

		if !(teacherID >= 1 && teacherID <= maxID) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Такого преподавателя нет в списке")
			return editTeacherNamePhoto, []*params.MessagesSendBuilder{b}, nil
		}

		teacherName, err := state.postgres.RequestPhoto.GetTeacherName(ctx, teacherID)
		if err != nil {
			return editTeacherNamePhoto, []*params.MessagesSendBuilder{}, err
		}

		err = state.postgres.RequestPhoto.UpdateTeachers(ctx, photoID, teacherName)
		if err != nil {
			return editTeacherNamePhoto, []*params.MessagesSendBuilder{}, err
		}

		allMarked, err := state.postgres.RequestPhoto.UpdateMarkedPeople(ctx, photoID, teacherName)
		if err != nil {
			return editTeacherNamePhoto, []*params.MessagesSendBuilder{}, err
		}

		err = state.postgres.RequestPhoto.DeletePointer(msg.PeerID)
		if err != nil {
			return editTeacherNamePhoto, nil, err
		}

		if allMarked {
			return checkPhoto, nil, nil
		}

		return editMarkedPeoplePhoto, nil, nil
	}
}

func (state EditTeacherNamePhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
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

func (state EditTeacherNamePhotoState) Name() stateName {
	return editTeacherNamePhoto
}

// EditUserTeacherNamePhotoState пользователь вводит ФИО преподавателя
type EditUserTeacherNamePhotoState struct {
	postgres *postrgres.Repo
}

func (state EditUserTeacherNamePhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return editUserTeacherNamePhoto, nil, nil
	}

	photoID, err := state.postgres.RequestPhoto.GetPhotoLastID(ctx, msg.PeerID)
	if err != nil {
		return editUserTeacherNamePhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return editTeacherNamePhoto, nil, nil
	default:
		err = state.postgres.RequestPhoto.UpdateTeachers(ctx, photoID, messageText)
		if err != nil {
			return editUserTeacherNamePhoto, []*params.MessagesSendBuilder{}, err
		}

		allMarked, err := state.postgres.RequestPhoto.UpdateMarkedPeople(ctx, photoID, messageText)
		if err != nil {
			return editUserTeacherNamePhoto, []*params.MessagesSendBuilder{}, err
		}

		if allMarked {
			return checkPhoto, nil, nil
		}

		return editMarkedPeoplePhoto, nil, nil
	}
}

func (state EditUserTeacherNamePhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите ФИО преподавателя")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditUserTeacherNamePhotoState) Name() stateName {
	return editUserTeacherNamePhoto
}

// EditStudentNamePhotoState пользователь вводит ФИО студента
type EditStudentNamePhotoState struct {
	postgres *postrgres.Repo
}

func (state EditStudentNamePhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return editStudentNamePhoto, nil, nil
	}

	photoID, err := state.postgres.RequestPhoto.GetPhotoLastID(ctx, msg.PeerID)
	if err != nil {
		return editStudentNamePhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return editTeacherNamePhoto, nil, nil
	default:
		allMarked, err := state.postgres.RequestPhoto.UpdateMarkedPeople(ctx, photoID, messageText)
		if err != nil {
			return editStudentNamePhoto, []*params.MessagesSendBuilder{}, err
		}

		if allMarked {
			return checkPhoto, nil, nil
		}

		return editMarkedPeoplePhoto, nil, nil
	}
}

func (state EditStudentNamePhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите фамилию и имя студента")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditStudentNamePhotoState) Name() stateName {
	return editStudentNamePhoto
}
