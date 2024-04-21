package state

import (
	"context"
	"fmt"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	"strconv"
)

// IsPeoplePresentPhotoState пользователь отвечает да или нет на вопрос есть ли на фото люди
type IsPeoplePresentPhotoState struct {
	postgres *postrgres.Repo
}

func (state IsPeoplePresentPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhoto.GetPhotoLastID(ctx, msg.PeerID)
	if err != nil {
		return isPeoplePresentPhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Да":
		return countPeoplePhoto, nil, nil
	case "Нет":
		err := state.postgres.RequestPhoto.DeleteMarksOnPhoto(ctx, photoID)
		if err != nil {
			return isPeoplePresentPhoto, []*params.MessagesSendBuilder{}, err
		}
		return eventYearPhoto, nil, nil
	case "Пропустить":
		return eventYearPhoto, nil, nil
	case "Назад":
		err := state.postgres.RequestPhoto.DeletePhoto(ctx, photoID)
		if err != nil {
			return isPeoplePresentPhoto, []*params.MessagesSendBuilder{}, err
		}
		return loadPhoto, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Выберите из предложенных вариантов")
		return isPeoplePresentPhoto, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state IsPeoplePresentPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("На фото есть люди?")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Да", "", "secondary")
	k.AddTextButton("Нет", "", "secondary")
	k.AddRow()
	k.AddTextButton("Пропустить", "", "secondary")
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state IsPeoplePresentPhotoState) Name() stateName {
	return isPeoplePresentPhoto
}

// CountPeoplePhotoState пользователь вводит количество людей на фото
type CountPeoplePhotoState struct {
	postgres *postrgres.Repo
}

func (state CountPeoplePhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhoto.GetPhotoLastID(ctx, msg.PeerID)
	if err != nil {
		return countPeoplePhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return isPeoplePresentPhoto, nil, nil
	default:
		count, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Количество людей должно быть целым числом")
			return countPeoplePhoto, []*params.MessagesSendBuilder{b}, nil
		}

		if count < 1 {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Количество людей должно быть больше нуля")
			return countPeoplePhoto, []*params.MessagesSendBuilder{b}, nil
		}
		err = state.postgres.RequestPhoto.UpdateCountPeople(ctx, photoID, count)
		if err != nil {
			return countPeoplePhoto, []*params.MessagesSendBuilder{}, err
		}
		return markedPeoplePhoto, nil, nil
	}
}

func (state CountPeoplePhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите количество людей на фото")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state CountPeoplePhotoState) Name() stateName {
	return countPeoplePhoto
}

// MarkedPeoplePhotoState пользователь отмечает людей на фото
type MarkedPeoplePhotoState struct {
	postgres *postrgres.Repo
}

func (state MarkedPeoplePhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhoto.GetPhotoLastID(ctx, msg.PeerID)
	if err != nil {
		return markedPeoplePhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Знаю ФИО человека":
		return isTeacherPhoto, nil, nil
	case "Не знаю ФИО человека":
		allMarked, err := state.postgres.RequestPhoto.UpdateMarkedPeople(ctx, photoID, "???")
		if err != nil {
			return markedPeoplePhoto, []*params.MessagesSendBuilder{}, err
		}
		if allMarked {
			return eventYearPhoto, nil, nil
		}
		return markedPeoplePhoto, nil, nil
	case "Закончить отмечать людей":
		return eventYearPhoto, nil, nil
	case "Назад":
		return countPeoplePhoto, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Выберите из предложенных вариантов")
		return markedPeoplePhoto, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state MarkedPeoplePhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
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

func (state MarkedPeoplePhotoState) Name() stateName {
	return markedPeoplePhoto
}

// IsTeacherPhotoState пользователь отмечает человека, отвечая на вопрос учитель ли это
type IsTeacherPhotoState struct {
	postgres *postrgres.Repo
}

func (state IsTeacherPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Да":
		return teacherNamePhoto, nil, nil
	case "Нет":
		return studentNamePhoto, nil, nil
	case "Назад":
		return markedPeoplePhoto, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Выберите из предложенных вариантов")
		return isTeacherPhoto, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state IsTeacherPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Это преподаватель?")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Да", "", "secondary")
	k.AddTextButton("Нет", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state IsTeacherPhotoState) Name() stateName {
	return isTeacherPhoto
}

// TeacherNamePhotoState пользователь выбирает преподавателя из списка
type TeacherNamePhotoState struct {
	postgres *postrgres.Repo
}

func (state TeacherNamePhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhoto.GetPhotoLastID(ctx, msg.PeerID)
	if err != nil {
		return teacherNamePhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Ввести ФИО вручную":
		return userTeacherNamePhoto, nil, nil
	case "Назад":
		return isTeacherPhoto, nil, nil
	default:
		teacherID, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введите номер преподавателя числом")
			return teacherNamePhoto, []*params.MessagesSendBuilder{b}, nil
		}

		maxID, err := state.postgres.RequestPhoto.GetTeacherMaxID()
		if err != nil {
			return teacherNamePhoto, []*params.MessagesSendBuilder{}, err
		}

		if !(teacherID >= 1 && teacherID <= maxID) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Такого преподавателя нет в списке")
			return teacherNamePhoto, []*params.MessagesSendBuilder{b}, nil
		}

		teacherName, err := state.postgres.RequestPhoto.GetTeacherName(ctx, teacherID)
		if err != nil {
			return teacherNamePhoto, []*params.MessagesSendBuilder{}, err
		}

		err = state.postgres.RequestPhoto.UpdateTeachers(ctx, photoID, teacherName)
		if err != nil {
			return teacherNamePhoto, []*params.MessagesSendBuilder{}, err
		}

		allMarked, err := state.postgres.RequestPhoto.UpdateMarkedPeople(ctx, photoID, teacherName)
		if err != nil {
			return teacherNamePhoto, []*params.MessagesSendBuilder{}, err
		}

		if allMarked {
			return eventYearPhoto, nil, nil
		}

		return markedPeoplePhoto, nil, nil
	}
}

func (state TeacherNamePhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	teacherNames, err := state.postgres.RequestPhoto.GetTeacherNames()
	if err != nil {
		return []*params.MessagesSendBuilder{}, err
	}

	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите номер преподавателя из списка ниже:\n" + teacherNames)
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Ввести ФИО вручную", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state TeacherNamePhotoState) Name() stateName {
	return teacherNamePhoto
}

// UserTeacherNamePhotoState пользователь вводит ФИО преподавателя
type UserTeacherNamePhotoState struct {
	postgres *postrgres.Repo
}

func (state UserTeacherNamePhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return userTeacherNamePhoto, nil, nil
	}

	photoID, err := state.postgres.RequestPhoto.GetPhotoLastID(ctx, msg.PeerID)
	if err != nil {
		return userTeacherNamePhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return teacherNamePhoto, nil, nil
	default:
		err = state.postgres.RequestPhoto.UpdateTeachers(ctx, photoID, messageText)
		if err != nil {
			return userTeacherNamePhoto, []*params.MessagesSendBuilder{}, err
		}

		allMarked, err := state.postgres.RequestPhoto.UpdateMarkedPeople(ctx, photoID, messageText)
		if err != nil {
			return userTeacherNamePhoto, []*params.MessagesSendBuilder{}, err
		}

		if allMarked {
			return eventYearPhoto, nil, nil
		}

		return markedPeoplePhoto, nil, nil
	}
}

func (state UserTeacherNamePhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите ФИО преподавателя")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state UserTeacherNamePhotoState) Name() stateName {
	return userTeacherNamePhoto
}

// StudentNamePhotoState пользователь вводит ФИО студента
type StudentNamePhotoState struct {
	postgres *postrgres.Repo
}

func (state StudentNamePhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return studentNamePhoto, nil, nil
	}

	photoID, err := state.postgres.RequestPhoto.GetPhotoLastID(ctx, msg.PeerID)
	if err != nil {
		return studentNamePhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return teacherNamePhoto, nil, nil
	default:
		allMarked, err := state.postgres.RequestPhoto.UpdateMarkedPeople(ctx, photoID, messageText)
		if err != nil {
			return studentNamePhoto, []*params.MessagesSendBuilder{}, err
		}

		if allMarked {
			return eventYearPhoto, nil, nil
		}

		return markedPeoplePhoto, nil, nil
	}
}

func (state StudentNamePhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите фамилию и имя человека")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state StudentNamePhotoState) Name() stateName {
	return studentNamePhoto
}
