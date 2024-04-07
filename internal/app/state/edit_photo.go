package state

import (
	"context"
	"fmt"
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
	k.AddTextButton("Назад", "", "secondary")
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
	k.AddTextButton("Назад", "", "secondary")
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
	k.AddTextButton("Назад", "", "secondary")
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

		err = state.postgres.RequestPhoto.UpdateEvent(ctx, photoID, eventNumber)
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
	k.AddTextButton("Назад", "", "secondary")
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
	b.Message("Напишите название своего события. Оно будет рассмотрено администратором")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
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
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditDescriptionPhotoState) Name() stateName {
	return editDescriptionPhoto
}

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
		return editIsPeoplePresentPhoto, nil, nil
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
	k.AddTextButton("Назад", "", "secondary")
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
	k.AddTextButton("Назад", "", "secondary")
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
	k.AddTextButton("Назад", "", "secondary")
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
	k.AddTextButton("Назад", "", "secondary")
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
		return teacherNamePhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Ввести ФИО вручную":
		return editUserTeacherNamePhoto, nil, nil
	case "Назад":
		return editIsTeacherPhoto, nil, nil
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

		if allMarked {
			return checkPhoto, nil, nil
		}

		return editMarkedPeoplePhoto, nil, nil
	}
}

func (state EditTeacherNamePhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
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
	k.AddTextButton("Назад", "", "secondary")
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
	k.AddTextButton("Назад", "", "secondary")
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
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditStudentNamePhotoState) Name() stateName {
	return editStudentNamePhoto
}
