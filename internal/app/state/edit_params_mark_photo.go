package state

import (
	"context"
	"fmt"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	"strconv"
)

// EditIsPeoplePresentParamsPhotoState пользователь редактирует ответ на вопрос есть ли на фото люди
type EditIsPeoplePresentParamsPhotoState struct {
	postgres *postrgres.Repo
}

func (state EditIsPeoplePresentParamsPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhoto.GetPhotoRequestID(ctx, msg.PeerID)
	if err != nil {
		return editIsPeoplePresentParamsPhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Да":
		return editCountPeopleParamsPhoto, nil, nil
	case "Нет":
		err = state.postgres.RequestPhoto.DeleteMarksOnPhoto(ctx, photoID)
		if err != nil {
			return editIsPeoplePresentParamsPhoto, []*params.MessagesSendBuilder{}, err
		}
		return personalAccountPhoto, nil, nil
	case "Назад":
		return editPhoto, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Выберите из предложенных вариантов")
		return editIsPeoplePresentParamsPhoto, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state EditIsPeoplePresentParamsPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
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

func (state EditIsPeoplePresentParamsPhotoState) Name() stateName {
	return editIsPeoplePresentParamsPhoto
}

// EditCountPeopleParamsPhotoState пользователь редактирует количество людей на фото
type EditCountPeopleParamsPhotoState struct {
	postgres *postrgres.Repo
}

func (state EditCountPeopleParamsPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhoto.GetPhotoRequestID(ctx, msg.PeerID)
	if err != nil {
		return editCountPeopleParamsPhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return editIsPeoplePresentParamsPhoto, nil, nil
	default:
		count, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Количество людей должно быть целым числом")
			return editCountPeopleParamsPhoto, []*params.MessagesSendBuilder{b}, nil
		}

		if count < 1 {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Количество людей должно быть больше нуля")
			return editCountPeopleParamsPhoto, []*params.MessagesSendBuilder{b}, nil
		}
		err = state.postgres.RequestPhoto.UpdateCountPeople(ctx, photoID, count)
		if err != nil {
			return editCountPeopleParamsPhoto, []*params.MessagesSendBuilder{}, err
		}
		return editMarkedPeopleParamsPhoto, nil, nil
	}
}

func (state EditCountPeopleParamsPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите количество людей на фото")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditCountPeopleParamsPhotoState) Name() stateName {
	return editCountPeopleParamsPhoto
}

// EditMarkedPeopleParamsPhotoState пользователь редактирует отметки людей на фото
type EditMarkedPeopleParamsPhotoState struct {
	postgres *postrgres.Repo
}

func (state EditMarkedPeopleParamsPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhoto.GetPhotoRequestID(ctx, msg.PeerID)
	if err != nil {
		return editMarkedPeopleParamsPhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Знаю ФИО человека":
		return editIsTeacherParamsPhoto, nil, nil
	case "Не знаю ФИО человека":
		allMarked, err := state.postgres.RequestPhoto.UpdateMarkedPeople(ctx, photoID, "???")
		if err != nil {
			return editMarkedPeopleParamsPhoto, []*params.MessagesSendBuilder{}, err
		}
		if allMarked {
			return personalAccountPhoto, nil, nil
		}
		return editMarkedPeopleParamsPhoto, nil, nil
	case "Закончить отмечать людей":
		return personalAccountPhoto, nil, nil
	case "Назад":
		return editCountPeopleParamsPhoto, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Выберите из предложенных вариантов")
		return editMarkedPeopleParamsPhoto, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state EditMarkedPeopleParamsPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
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

func (state EditMarkedPeopleParamsPhotoState) Name() stateName {
	return editMarkedPeopleParamsPhoto
}

// EditIsTeacherParamsPhotoState пользователь редактирует ответ на вопрос учитель ли это
type EditIsTeacherParamsPhotoState struct {
	postgres *postrgres.Repo
}

func (state EditIsTeacherParamsPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Да":
		return editTeacherNameParamsPhoto, nil, nil
	case "Нет":
		return editStudentNameParamsPhoto, nil, nil
	case "Назад":
		return editMarkedPeopleParamsPhoto, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Выберите из предложенных вариантов")
		return editIsTeacherParamsPhoto, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state EditIsTeacherParamsPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
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

func (state EditIsTeacherParamsPhotoState) Name() stateName {
	return editIsTeacherParamsPhoto
}

// EditTeacherNameParamsPhotoState пользователь выбирает преподавателя из списка
type EditTeacherNameParamsPhotoState struct {
	postgres *postrgres.Repo
}

func (state EditTeacherNameParamsPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhoto.GetPhotoRequestID(ctx, msg.PeerID)
	if err != nil {
		return editTeacherNameParamsPhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Ввести ФИО вручную":
		return editUserTeacherNameParamsPhoto, nil, nil
	case "Назад":
		err = state.postgres.RequestPhoto.DeletePointer(msg.PeerID)
		if err != nil {
			return editTeacherNameParamsPhoto, nil, err
		}
		return editIsTeacherParamsPhoto, nil, nil
	case "⬅️":
		err = state.postgres.RequestPhoto.ChangePointer(msg.PeerID, false)
		if err != nil {
			return editTeacherNameParamsPhoto, nil, err
		}
		return editTeacherNameParamsPhoto, nil, nil
	case "➡️":
		err = state.postgres.RequestPhoto.ChangePointer(msg.PeerID, true)
		if err != nil {
			return editTeacherNameParamsPhoto, nil, err
		}
		return editTeacherNameParamsPhoto, nil, nil
	default:
		teacherID, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введите номер преподавателя числом")
			return editTeacherNameParamsPhoto, []*params.MessagesSendBuilder{b}, nil
		}

		maxID, err := state.postgres.RequestPhoto.GetTeacherMaxID()
		if err != nil {
			return editTeacherNameParamsPhoto, []*params.MessagesSendBuilder{}, err
		}

		if !(teacherID >= 1 && teacherID <= maxID) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Такого преподавателя нет в списке")
			return editTeacherNameParamsPhoto, []*params.MessagesSendBuilder{b}, nil
		}

		teacherName, err := state.postgres.RequestPhoto.GetTeacherName(ctx, teacherID)
		if err != nil {
			return editTeacherNameRequestPhoto, []*params.MessagesSendBuilder{}, err
		}

		err = state.postgres.RequestPhoto.UpdateTeachers(ctx, photoID, teacherName)
		if err != nil {
			return editTeacherNameParamsPhoto, []*params.MessagesSendBuilder{}, err
		}

		allMarked, err := state.postgres.RequestPhoto.UpdateMarkedPeople(ctx, photoID, teacherName)
		if err != nil {
			return editTeacherNameParamsPhoto, []*params.MessagesSendBuilder{}, err
		}

		err = state.postgres.RequestPhoto.DeletePointer(msg.PeerID)
		if err != nil {
			return editTeacherNameParamsPhoto, nil, err
		}

		if allMarked {
			return personalAccountPhoto, nil, nil
		}

		return editMarkedPeopleParamsPhoto, nil, nil
	}
}

func (state EditTeacherNameParamsPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
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

func (state EditTeacherNameParamsPhotoState) Name() stateName {
	return editTeacherNameParamsPhoto
}

// EditUserTeacherNameParamsPhotoState пользователь вводит ФИО преподавателя
type EditUserTeacherNameParamsPhotoState struct {
	postgres *postrgres.Repo
}

func (state EditUserTeacherNameParamsPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return editUserTeacherNameParamsPhoto, nil, nil
	}

	photoID, err := state.postgres.RequestPhoto.GetPhotoRequestID(ctx, msg.PeerID)
	if err != nil {
		return editUserTeacherNameParamsPhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return editTeacherNameParamsPhoto, nil, nil
	default:
		err = state.postgres.RequestPhoto.UpdateTeachers(ctx, photoID, messageText)
		if err != nil {
			return editUserTeacherNameParamsPhoto, []*params.MessagesSendBuilder{}, err
		}

		allMarked, err := state.postgres.RequestPhoto.UpdateMarkedPeople(ctx, photoID, messageText)
		if err != nil {
			return editUserTeacherNameParamsPhoto, []*params.MessagesSendBuilder{}, err
		}

		if allMarked {
			return personalAccountPhoto, nil, nil
		}

		return editMarkedPeopleParamsPhoto, nil, nil
	}
}

func (state EditUserTeacherNameParamsPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите ФИО преподавателя")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditUserTeacherNameParamsPhotoState) Name() stateName {
	return editUserTeacherNameParamsPhoto
}

// EditStudentNameParamsPhotoState пользователь вводит ФИО студента
type EditStudentNameParamsPhotoState struct {
	postgres *postrgres.Repo
}

func (state EditStudentNameParamsPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return editStudentNameParamsPhoto, nil, nil
	}

	photoID, err := state.postgres.RequestPhoto.GetPhotoRequestID(ctx, msg.PeerID)
	if err != nil {
		return editStudentNameParamsPhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return editTeacherNameParamsPhoto, nil, nil
	default:
		allMarked, err := state.postgres.RequestPhoto.UpdateMarkedPeople(ctx, photoID, messageText)
		if err != nil {
			return editStudentNameParamsPhoto, []*params.MessagesSendBuilder{}, err
		}

		if allMarked {
			return personalAccountPhoto, nil, nil
		}

		return editMarkedPeopleParamsPhoto, nil, nil
	}
}

func (state EditStudentNameParamsPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите фамилию и имя студента")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditStudentNameParamsPhotoState) Name() stateName {
	return editStudentNameParamsPhoto
}
