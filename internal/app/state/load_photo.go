package state

import (
	"context"
	"fmt"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	"strconv"
	"time"
)

var validExtension = map[string]struct{}{
	"jpg":  struct{}{},
	"jpeg": struct{}{},
	"png":  struct{}{},
	"tif":  struct{}{},
	"tiff": struct{}{},
}

// LoadPhotoState пользователь загружает фотографию
type LoadPhotoState struct {
	postgres *postrgres.Repo
	vk       *api.VK
}

func (state LoadPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "Назад" {
		return photoStart, nil, nil
	}
	attachment := msg.Attachments

	if len(attachment) == 0 {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Загрузите фотографию, прикрепив её к сообщению")
		return loadPhoto, []*params.MessagesSendBuilder{b}, nil
	}

	if len(attachment) > 1 {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Можно загрузить лишь одну фотографию, для загрузки множества фотографий воспользуйтесь загрузкой архива")
		return loadPhoto, []*params.MessagesSendBuilder{b}, nil
	}

	if attachment[0].Type == "photo" {
		err := state.postgres.RequestPhoto.UploadPhoto(ctx, state.vk, attachment[0].Photo, msg.PeerID)
		if err != nil {
			return loadPhoto, []*params.MessagesSendBuilder{}, err
		}

		return isPeoplePresentPhoto, nil, nil
	}

	if attachment[0].Type == "doc" {
		if _, ok := validExtension[attachment[0].Doc.Ext]; !ok {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Данная фотография недопустимого формата")
			return loadPhoto, []*params.MessagesSendBuilder{b}, nil
		}

		err := state.postgres.RequestPhoto.UploadPhotoAsFile(ctx, state.vk, attachment[0].Doc, msg.PeerID)
		if err != nil {
			return loadPhoto, []*params.MessagesSendBuilder{}, err
		}

		return isPeoplePresentPhoto, nil, nil
	}

	return loadPhoto, nil, nil
}

func (state LoadPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Загрузите фото. Допустимые  форматы фото: jpg, jpeg, png, tif, tiff")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state LoadPhotoState) Name() stateName {
	return loadPhoto
}

// IsPeoplePresentPhotoState пользователь отвечает да или нет на вопрос есть ли на фото люди
type IsPeoplePresentPhotoState struct {
	postgres *postrgres.Repo
}

func (state IsPeoplePresentPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Да":
		return countPeoplePhoto, nil, nil
	case "Нет":
		return eventYearPhoto, nil, nil
	case "Пропустить":
		return eventYearPhoto, nil, nil
	case "Назад":
		err := state.postgres.RequestPhoto.DeletePhotoRequest(ctx, msg.PeerID)
		if err != nil {
			return isPeoplePresentPhoto, []*params.MessagesSendBuilder{}, err
		}
		return loadPhoto, nil, nil
	default:
		return isPeoplePresentPhoto, nil, nil
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
	k.AddTextButton("Назад", "", "secondary")
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
		return editEventYearPhoto, []*params.MessagesSendBuilder{}, err
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
	k.AddTextButton("Назад", "", "secondary")
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
		return editEventYearPhoto, []*params.MessagesSendBuilder{}, err
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
	k.AddTextButton("Назад", "", "secondary")
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
		return markedPeoplePhoto, nil, nil
	case "Нет":
		return markedPeoplePhoto, nil, nil
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
	k.AddRow()
	k.AddTextButton("Нет", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state IsTeacherPhotoState) Name() stateName {
	return isTeacherPhoto
}

// TeacherPhotoState пользователь выбирает преподавателя из списка
type TeacherPhotoState struct {
	postgres *postrgres.Repo
}

func (state TeacherPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhoto.GetPhotoLastID(ctx, msg.PeerID)
	if err != nil {
		return teacherPhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Ввести ФИО вручную":
		return teacherPhoto, nil, nil
	case "Назад":
		return isTeacherPhoto, nil, nil
	default:
		teacherID, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введите номер преподавателя числом")
			return teacherPhoto, []*params.MessagesSendBuilder{b}, nil
		}

		maxID, err := state.postgres.RequestPhoto.GetTeacherMaxID()
		if err != nil {
			return teacherPhoto, []*params.MessagesSendBuilder{}, err
		}

		if !(teacherID >= 1 && teacherID <= maxID) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Такого преподавателя нет в списке")
			return teacherPhoto, []*params.MessagesSendBuilder{b}, nil
		}

		teacherName, err := state.postgres.RequestPhoto.GetTeacherName(ctx, teacherID)
		if err != nil {
			return teacherPhoto, []*params.MessagesSendBuilder{}, err
		}

		allMarked, err := state.postgres.RequestPhoto.UpdateMarkedPeople(ctx, photoID, teacherName)
		if err != nil {
			return teacherPhoto, []*params.MessagesSendBuilder{}, err
		}

		if allMarked {
			return eventYearPhoto, nil, nil
		}

		return markedPeoplePhoto, nil, nil
	}
}

func (state TeacherPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
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

func (state TeacherPhotoState) Name() stateName {
	return teacherPhoto
}

// EventYearPhotoState пользователь указывает год создания документа
type EventYearPhotoState struct {
	postgres *postrgres.Repo
}

func (state EventYearPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhoto.GetPhotoLastID(ctx, msg.PeerID)
	if err != nil {
		return editEventYearPhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return isPeoplePresentPhoto, nil, nil
	case "Пропустить":
		return studyProgramPhoto, nil, nil
	default:
		year, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введён год недопустимого формата")
			return eventYearPhoto, []*params.MessagesSendBuilder{b}, nil
		}

		currentYear := time.Now().Year()
		if !(year >= 1900 && year <= currentYear) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введён несуществующий год")
			return eventYearPhoto, []*params.MessagesSendBuilder{b}, nil
		}
		err = state.postgres.RequestPhoto.UpdateYear(ctx, photoID, year)
		if err != nil {
			return eventYearPhoto, []*params.MessagesSendBuilder{}, err
		}
		return studyProgramPhoto, nil, nil
	}
}

func (state EventYearPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите год события в формате YYYY")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Пропустить", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EventYearPhotoState) Name() stateName {
	return eventYearPhoto
}

// StudyProgramPhotoState пользователь указывает программу обучения
type StudyProgramPhotoState struct {
	postgres *postrgres.Repo
}

func (state StudyProgramPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhoto.GetPhotoLastID(ctx, msg.PeerID)
	if err != nil {
		return editEventYearPhoto, []*params.MessagesSendBuilder{}, err
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
	case "Пропустить":
		return eventNamePhoto, nil, nil
	case "Назад":
		return eventYearPhoto, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Такой программы обучения нет в предложенных вариантах")
		return eventYearPhoto, []*params.MessagesSendBuilder{b}, nil
	}

	err = state.postgres.RequestPhoto.UpdateStudyProgram(ctx, photoID, educationProgram)
	if err != nil {
		return studyProgramPhoto, []*params.MessagesSendBuilder{}, err
	}
	return eventNamePhoto, nil, nil
}

func (state StudyProgramPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
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
	k.AddTextButton("Пропустить", "", "secondary")
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state StudyProgramPhotoState) Name() stateName {
	return studyProgramPhoto
}

// EventNamePhotoState пользователь указывает существующее название события
type EventNamePhotoState struct {
	postgres *postrgres.Repo
}

func (state EventNamePhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhoto.GetPhotoLastID(ctx, msg.PeerID)
	if err != nil {
		return editEventYearPhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Своё событие":
		return userEventNamePhoto, nil, nil
	case "Пропустить":
		return descriptionPhoto, nil, nil
	case "Назад":
		return studyProgramPhoto, nil, nil
	default:
		maxID, err := state.postgres.RequestPhoto.GetEventMaxID()
		if err != nil {
			return eventNamePhoto, []*params.MessagesSendBuilder{}, err
		}

		eventNumber, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введите номер события числом")
			return eventNamePhoto, []*params.MessagesSendBuilder{b}, nil
		}

		if !(eventNumber >= 1 && eventNumber <= maxID) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Такого события нет в списке")
			return eventNamePhoto, []*params.MessagesSendBuilder{b}, nil
		}

		err = state.postgres.RequestPhoto.UpdateEvent(ctx, photoID, eventNumber)
		if err != nil {
			return eventNamePhoto, []*params.MessagesSendBuilder{}, err
		}
		return descriptionPhoto, nil, nil
	}
}

func (state EventNamePhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	categories, err := state.postgres.RequestPhoto.GetTeacherNames()
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
	k.AddTextButton("Пропустить", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EventNamePhotoState) Name() stateName {
	return eventNamePhoto
}

// UserEventNamePhotoState пользователь указывает своё название события
type UserEventNamePhotoState struct {
	postgres *postrgres.Repo
}

func (state UserEventNamePhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhoto.GetPhotoLastID(ctx, msg.PeerID)
	if err != nil {
		return editEventYearPhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Пропустить":
		return descriptionPhoto, nil, nil
	case "Назад":
		return eventNamePhoto, nil, nil
	default:
		err := state.postgres.RequestPhoto.UpdateUserEvent(ctx, photoID, messageText)
		if err != nil {
			return userEventNamePhoto, []*params.MessagesSendBuilder{}, err
		}
		return descriptionPhoto, nil, nil
	}
}

func (state UserEventNamePhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите название своего события. Оно будет рассмотрено администратором")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Пропустить", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state UserEventNamePhotoState) Name() stateName {
	return userEventNamePhoto
}

// DescriptionPhotoState пользователь вводит описание фотографии
type DescriptionPhotoState struct {
	postgres *postrgres.Repo
}

func (state DescriptionPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhoto.GetPhotoLastID(ctx, msg.PeerID)
	if err != nil {
		return editEventYearPhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Пропустить":
		return checkPhoto, nil, nil
	case "Назад":
		return eventNamePhoto, nil, nil
	default:
		err := state.postgres.RequestPhoto.UpdateDescription(ctx, photoID, messageText)
		if err != nil {
			return descriptionPhoto, []*params.MessagesSendBuilder{}, err
		}
		return checkPhoto, nil, nil
	}
}

func (state DescriptionPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите описание фотографии")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Пропустить", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state DescriptionPhotoState) Name() stateName {
	return descriptionPhoto
}

// CheckPhotoState пользователь проверяет заявку на загрузку фотографии
type CheckPhotoState struct {
	postgres *postrgres.Repo
}

func (state CheckPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Отправить":
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Фотография отправлена администратору на рассмотрение. Вы можете отслеживать статус своей заявки в личном кабинете")
		return photoStart, []*params.MessagesSendBuilder{b}, nil
	case "Редактировать заявку":
		return editPhoto, nil, nil
	case "Назад":
		return descriptionPhoto, nil, nil
	default:
		return checkPhoto, nil, nil
	}
}

func (state CheckPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	photoID, err := state.postgres.RequestPhoto.GetPhotoLastID(ctx, vkID)
	if err != nil {
		return []*params.MessagesSendBuilder{}, err
	}

	output, attachment, err := state.postgres.RequestPhoto.CheckParams(ctx, photoID)
	if err != nil {
		return []*params.MessagesSendBuilder{}, err
	}

	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Проверьте правильность введённых параметров:\n" + output)
	b.Attachment(attachment)
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Отправить", "", "secondary")
	k.AddRow()
	k.AddTextButton("Редактировать заявку", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state CheckPhotoState) Name() stateName {
	return checkPhoto
}
