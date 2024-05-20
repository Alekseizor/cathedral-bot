package state

import (
	"context"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	"strconv"
	"time"
)

// CategorySearchAlbumState пользователь выбирает категорию для поиска альбома
type CategorySearchAlbumState struct {
	postgres *postrgres.Repo
}

func (state CategorySearchAlbumState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return categorySearchAlbum, nil, nil
	}

	switch messageText {
	case "Студентов":
		return yearSearchAlbum, nil, nil
	case "Преподавателя":
		return surnameTeacherSearchAlbum, nil, nil
	case "Назад":
		err := state.postgres.SearchAlbum.DeleteSearchAlbum(ctx, msg.PeerID)
		if err != nil {
			return categorySearchAlbum, nil, err
		}
		return photoStart, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Такой категории нет в предложенных вариантах")
		return categorySearchAlbum, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state CategorySearchAlbumState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Выберите чей альбом искать")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Студентов", "", "secondary")
	k.AddTextButton("Преподавателя", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state CategorySearchAlbumState) Name() stateName {
	return categorySearchAlbum
}

// YearSearchAlbumState пользователь указывает год для поиска альбома
type YearSearchAlbumState struct {
	postgres *postrgres.Repo
}

func (state YearSearchAlbumState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Пропустить":
	case "Назад":
		return categorySearchAlbum, nil, nil
	default:
		year, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введён год недопустимого формата")
			return yearSearchAlbum, []*params.MessagesSendBuilder{b}, nil
		}

		currentYear := time.Now().Year()
		if !(year >= 1900 && year <= currentYear) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введён несуществующий год")
			return yearSearchAlbum, []*params.MessagesSendBuilder{b}, nil
		}

		err = state.postgres.SearchAlbum.UpdateYear(ctx, msg.PeerID, year)
		if err != nil {
			return yearSearchAlbum, nil, err
		}
	}

	count, err := state.postgres.SearchAlbum.CountAlbums(ctx, msg.PeerID)
	if err != nil {
		return yearSearchAlbum, nil, err
	}

	if count < 2 {
		return findYearLess2SearchAlbum, nil, nil
	}
	return findYearSearchAlbum, nil, nil
}

func (state YearSearchAlbumState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите год события в формате YYYY")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Пропустить", "", "secondary")
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state YearSearchAlbumState) Name() stateName {
	return yearSearchAlbum
}

// FindYearSearchAlbumState выводится количество найденных альбомов, если их два или больше
type FindYearSearchAlbumState struct {
	postgres *postrgres.Repo
}

func (state FindYearSearchAlbumState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return findYearSearchAlbum, nil, nil
	}

	switch messageText {
	case "Показать найденные альбомы":
		return showListYearSearchAlbum, nil, nil
	case "Добавить фильтр по программе обучения":
		return studyProgramSearchAlbum, nil, nil
	case "Назад":
		err := state.postgres.SearchAlbum.DeleteYear(ctx, msg.PeerID)
		if err != nil {
			return findYearSearchAlbum, nil, err
		}
		return yearSearchAlbum, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Такого действия нет в предложенных вариантах")
		return findYearSearchAlbum, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state FindYearSearchAlbumState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	count, err := state.postgres.SearchAlbum.CountAlbums(ctx, vkID)
	if err != nil {
		return nil, err
	}
	countString := strconv.Itoa(count)

	searchParams, err := state.postgres.SearchAlbum.GetSearchParams(vkID)
	if err != nil {
		return nil, err
	}

	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message(searchParams + "\n" + "Найдено альбомов: " + countString)
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Показать найденные альбомы", "", "secondary")
	k.AddRow()
	k.AddTextButton("Добавить фильтр по программе обучения", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state FindYearSearchAlbumState) Name() stateName {
	return findYearSearchAlbum
}

// FindYearLess2SearchAlbumState выводится количество найденных альбомов, если их меньше двух
type FindYearLess2SearchAlbumState struct {
	postgres *postrgres.Repo
}

func (state FindYearLess2SearchAlbumState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return findYearLess2SearchAlbum, nil, nil
	}

	switch messageText {
	case "Завершить поиск":
		err := state.postgres.SearchAlbum.DeleteSearchAlbum(ctx, msg.PeerID)
		if err != nil {
			return findYearLess2SearchAlbum, nil, err
		}
		return photoStart, nil, nil
	case "Назад":
		err := state.postgres.SearchAlbum.DeleteYear(ctx, msg.PeerID)
		if err != nil {
			return findYearLess2SearchAlbum, nil, err
		}
		return yearSearchAlbum, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Такого действия нет в предложенных вариантах")
		return findYearLess2SearchAlbum, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state FindYearLess2SearchAlbumState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	count, err := state.postgres.SearchAlbum.CountAlbums(ctx, vkID)
	if err != nil {
		return nil, err
	}

	countString := strconv.Itoa(count)
	albums, _, _, err := state.postgres.SearchAlbum.ShowList(ctx, vkID)
	if err != nil {
		return nil, err
	}

	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Найдено альбомов: " + countString + "\n" + albums)
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Завершить поиск", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state FindYearLess2SearchAlbumState) Name() stateName {
	return findYearLess2SearchAlbum
}

// ShowListYearSearchAlbumState пользователь получает список найденных альбомов
type ShowListYearSearchAlbumState struct {
	postgres *postrgres.Repo
}

func (state ShowListYearSearchAlbumState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return showListYearSearchAlbum, nil, nil
	}

	switch messageText {
	case "Завершить поиск":
		err := state.postgres.SearchAlbum.DeleteSearchAlbum(ctx, msg.PeerID)
		if err != nil {
			return showListYearSearchAlbum, nil, err
		}
		return photoStart, nil, nil
	case "Назад":
		err := state.postgres.SearchAlbum.DeletePointer(msg.PeerID)
		if err != nil {
			return showListYearSearchAlbum, nil, err
		}
		return findYearSearchAlbum, nil, nil
	case "⬅️":
		err := state.postgres.SearchAlbum.ChangePointerStudents(msg.PeerID, false)
		if err != nil {
			return showListYearSearchAlbum, nil, err
		}
		return showListYearSearchAlbum, nil, nil
	case "➡️":
		err := state.postgres.SearchAlbum.ChangePointerStudents(msg.PeerID, true)
		if err != nil {
			return showListYearSearchAlbum, nil, err
		}
		return showListYearSearchAlbum, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Такого действия нет в предложенных вариантах")
		return showListYearSearchAlbum, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state ShowListYearSearchAlbumState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	albums, pointer, count, err := state.postgres.SearchAlbum.ShowList(ctx, vkID)
	if err != nil {
		return nil, err
	}
	countString := strconv.Itoa(count)

	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Найдено альбомов: " + countString + "\n" + albums)
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
	k.AddTextButton("Завершить поиск", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state ShowListYearSearchAlbumState) Name() stateName {
	return showListYearSearchAlbum
}

// StudyProgramSearchAlbumState пользователь указывает программу обучения для поиска альбома
type StudyProgramSearchAlbumState struct {
	postgres *postrgres.Repo
}

func (state StudyProgramSearchAlbumState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

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
	case "Назад":
		return findYearSearchAlbum, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Такой программы обучения нет в предложенных вариантах")
		return studyProgramSearchAlbum, []*params.MessagesSendBuilder{b}, nil
	}

	if educationProgram != "" {
		err := state.postgres.SearchAlbum.UpdateStudyProgram(ctx, msg.PeerID, educationProgram)
		if err != nil {
			return studyProgramSearchAlbum, nil, err
		}
	}

	count, err := state.postgres.SearchAlbum.CountAlbums(ctx, msg.PeerID)
	if err != nil {
		return studyProgramSearchAlbum, nil, err
	}

	if count < 2 {
		return findStudyProgramLess2SearchAlbum, nil, nil
	}
	return findStudyProgramSearchAlbum, nil, nil
}

func (state StudyProgramSearchAlbumState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
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
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state StudyProgramSearchAlbumState) Name() stateName {
	return studyProgramSearchAlbum
}

// FindStudyProgramSearchAlbumState выводится количество найденных альбомов, если их два или больше
type FindStudyProgramSearchAlbumState struct {
	postgres *postrgres.Repo
}

func (state FindStudyProgramSearchAlbumState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return findStudyProgramSearchAlbum, nil, nil
	}

	switch messageText {
	case "Показать найденные альбомы":
		return showListStudyProgramSearchAlbum, nil, nil
	case "Добавить фильтр по названию события":
		return eventSearchAlbum, nil, nil
	case "Назад":
		err := state.postgres.SearchAlbum.DeleteStudyProgram(ctx, msg.PeerID)
		if err != nil {
			return findStudyProgramSearchAlbum, nil, err
		}
		return studyProgramSearchAlbum, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Такого действия нет в предложенных вариантах")
		return findStudyProgramSearchAlbum, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state FindStudyProgramSearchAlbumState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	count, err := state.postgres.SearchAlbum.CountAlbums(ctx, vkID)
	if err != nil {
		return nil, err
	}
	countString := strconv.Itoa(count)

	searchParams, err := state.postgres.SearchAlbum.GetSearchParams(vkID)
	if err != nil {
		return nil, err
	}

	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message(searchParams + "\n" + "Найдено альбомов: " + countString)
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Показать найденные альбомы", "", "secondary")
	k.AddRow()
	k.AddTextButton("Добавить фильтр по названию события", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state FindStudyProgramSearchAlbumState) Name() stateName {
	return findStudyProgramSearchAlbum
}

// FindStudyProgramLess2SearchAlbumState выводится количество найденных альбомов, если их меньше двух
type FindStudyProgramLess2SearchAlbumState struct {
	postgres *postrgres.Repo
}

func (state FindStudyProgramLess2SearchAlbumState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return findStudyProgramLess2SearchAlbum, nil, nil
	}

	switch messageText {
	case "Завершить поиск":
		err := state.postgres.SearchAlbum.DeleteSearchAlbum(ctx, msg.PeerID)
		if err != nil {
			return findStudyProgramLess2SearchAlbum, nil, err
		}
		return photoStart, nil, nil
	case "Назад":
		err := state.postgres.SearchAlbum.DeleteStudyProgram(ctx, msg.PeerID)
		if err != nil {
			return findStudyProgramLess2SearchAlbum, nil, err
		}
		return studyProgramSearchAlbum, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Такого действия нет в предложенных вариантах")
		return findStudyProgramLess2SearchAlbum, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state FindStudyProgramLess2SearchAlbumState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	count, err := state.postgres.SearchAlbum.CountAlbums(ctx, vkID)
	if err != nil {
		return nil, err
	}

	countString := strconv.Itoa(count)
	albums, _, _, err := state.postgres.SearchAlbum.ShowList(ctx, vkID)
	if err != nil {
		return nil, err
	}

	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Найдено альбомов: " + countString + "\n" + albums)
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Завершить поиск", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state FindStudyProgramLess2SearchAlbumState) Name() stateName {
	return findStudyProgramLess2SearchAlbum
}

// ShowListStudyProgramSearchAlbumState пользователь получает список найденных альбомов
type ShowListStudyProgramSearchAlbumState struct {
	postgres *postrgres.Repo
}

func (state ShowListStudyProgramSearchAlbumState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return showListStudyProgramSearchAlbum, nil, nil
	}

	switch messageText {
	case "Завершить поиск":
		err := state.postgres.SearchAlbum.DeleteSearchAlbum(ctx, msg.PeerID)
		if err != nil {
			return showListStudyProgramSearchAlbum, nil, err
		}
		return photoStart, nil, nil
	case "Назад":
		err := state.postgres.SearchAlbum.DeletePointer(msg.PeerID)
		if err != nil {
			return showListStudyProgramSearchAlbum, nil, err
		}
		return findStudyProgramSearchAlbum, nil, nil
	case "⬅️":
		err := state.postgres.SearchAlbum.ChangePointerStudents(msg.PeerID, false)
		if err != nil {
			return showListStudyProgramSearchAlbum, nil, err
		}
		return showListStudyProgramSearchAlbum, nil, nil
	case "➡️":
		err := state.postgres.SearchAlbum.ChangePointerStudents(msg.PeerID, true)
		if err != nil {
			return showListStudyProgramSearchAlbum, nil, err
		}
		return showListStudyProgramSearchAlbum, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Такого действия нет в предложенных вариантах")
		return showListStudyProgramSearchAlbum, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state ShowListStudyProgramSearchAlbumState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	albums, pointer, count, err := state.postgres.SearchAlbum.ShowList(ctx, vkID)
	if err != nil {
		return nil, err
	}
	countString := strconv.Itoa(count)

	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Найдено альбомов: " + countString + "\n" + albums)
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
	k.AddTextButton("Завершить поиск", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state ShowListStudyProgramSearchAlbumState) Name() stateName {
	return showListStudyProgramSearchAlbum
}

// EventSearchAlbumState пользователь указывает название события для поиска альбома
type EventSearchAlbumState struct {
	postgres *postrgres.Repo
}

func (state EventSearchAlbumState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return eventSearchAlbum, nil, nil
	}

	switch messageText {
	case "Пропустить":
	case "Назад":
		return findStudyProgramSearchAlbum, nil, nil
	default:
		maxID, err := state.postgres.SearchAlbum.GetEventMaxID()
		if err != nil {
			return eventSearchAlbum, []*params.MessagesSendBuilder{}, err
		}

		eventNumber, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введена не цифра")
			return eventSearchAlbum, []*params.MessagesSendBuilder{b}, nil
		}

		if !(eventNumber >= 1 && eventNumber <= maxID) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Такого события нет в списке")
			return eventSearchAlbum, []*params.MessagesSendBuilder{b}, nil
		}

		err = state.postgres.SearchAlbum.UpdateEvent(ctx, msg.PeerID, eventNumber-1)
		if err != nil {
			return eventSearchAlbum, []*params.MessagesSendBuilder{}, err
		}
	}

	return findEventSearchAlbum, nil, nil
}

func (state EventSearchAlbumState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	categories, err := state.postgres.SearchAlbum.GetEventNames()
	if err != nil {
		return []*params.MessagesSendBuilder{}, err
	}
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите номер события из списка ниже:\n" + categories)
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Пропустить", "", "secondary")
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EventSearchAlbumState) Name() stateName {
	return eventSearchAlbum
}

// FindEventSearchAlbumState выводится количество найденных альбомов
type FindEventSearchAlbumState struct {
	postgres *postrgres.Repo
}

func (state FindEventSearchAlbumState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return findEventSearchAlbum, nil, nil
	}

	switch messageText {
	case "Завершить поиск":
		err := state.postgres.SearchAlbum.DeleteSearchAlbum(ctx, msg.PeerID)
		if err != nil {
			return findEventSearchAlbum, nil, err
		}
		return photoStart, nil, nil
	case "Назад":
		err := state.postgres.SearchAlbum.DeleteEvent(ctx, msg.PeerID)
		if err != nil {
			return findEventSearchAlbum, nil, err
		}

		err = state.postgres.SearchAlbum.DeletePointer(msg.PeerID)
		if err != nil {
			return findEventSearchAlbum, nil, err
		}

		return eventSearchAlbum, nil, nil
	case "⬅️":
		err := state.postgres.SearchAlbum.ChangePointerStudents(msg.PeerID, false)
		if err != nil {
			return findEventSearchAlbum, nil, err
		}
		return findEventSearchAlbum, nil, nil
	case "➡️":
		err := state.postgres.SearchAlbum.ChangePointerStudents(msg.PeerID, true)
		if err != nil {
			return findEventSearchAlbum, nil, err
		}
		return findEventSearchAlbum, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Такого действия нет в предложенных вариантах")
		return findEventSearchAlbum, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state FindEventSearchAlbumState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	albums, pointer, count, err := state.postgres.SearchAlbum.ShowList(ctx, vkID)
	if err != nil {
		return nil, err
	}
	countString := strconv.Itoa(count)

	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Найдено альбомов: " + countString + "\n" + albums)
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
	k.AddTextButton("Завершить поиск", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state FindEventSearchAlbumState) Name() stateName {
	return findEventSearchAlbum
}

// SurnameTeacherSearchAlbumState пользователь вводит первые буквы фамилии преподавателя
type SurnameTeacherSearchAlbumState struct {
	postgres *postrgres.Repo
}

func (state SurnameTeacherSearchAlbumState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return surnameTeacherSearchAlbum, nil, nil
	}

	switch messageText {
	case "Пропустить":
		return teacherSearchAlbum, nil, nil
	case "Назад":
		return categorySearchAlbum, nil, nil
	default:
		name := messageText

		if containsDigits(name) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Не используйте цифры")
			return surnameTeacherSearchAlbum, []*params.MessagesSendBuilder{b}, nil
		}

		if containsNonRussianChars(name) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Используйте только русские символы")
			return surnameTeacherSearchAlbum, []*params.MessagesSendBuilder{b}, nil
		}

		err := state.postgres.SearchAlbum.UpdateName(ctx, msg.PeerID, name)
		if err != nil {
			return surnameTeacherSearchAlbum, nil, err
		}
	}

	return teacherSearchAlbum, nil, nil
}

func (state SurnameTeacherSearchAlbumState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите первую(-ые) букву(-ы) фамилии преподавателя")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Пропустить", "", "secondary")
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state SurnameTeacherSearchAlbumState) Name() stateName {
	return surnameTeacherSearchAlbum
}

// TeacherSearchAlbumState пользователь ищет альбом преподавателя
type TeacherSearchAlbumState struct {
	postgres *postrgres.Repo
}

func (state TeacherSearchAlbumState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return teacherSearchAlbum, nil, nil
	}

	switch messageText {
	case "Завершить поиск":
		err := state.postgres.SearchAlbum.DeleteSearchAlbum(ctx, msg.PeerID)
		if err != nil {
			return teacherSearchAlbum, nil, err
		}
		return photoStart, nil, nil
	case "Назад":
		err := state.postgres.SearchAlbum.DeletePointer(msg.PeerID)
		if err != nil {
			return teacherSearchAlbum, nil, err
		}

		err = state.postgres.SearchAlbum.DeleteSurname(msg.PeerID)
		if err != nil {
			return teacherSearchAlbum, nil, err
		}

		return surnameTeacherSearchAlbum, nil, nil
	case "⬅️":
		err := state.postgres.SearchAlbum.ChangePointerTeacher(msg.PeerID, false)
		if err != nil {
			return teacherSearchAlbum, nil, err
		}
		return teacherSearchAlbum, nil, nil
	case "➡️":
		err := state.postgres.SearchAlbum.ChangePointerTeacher(msg.PeerID, true)
		if err != nil {
			return teacherSearchAlbum, nil, err
		}
		return teacherSearchAlbum, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Такого действия нет в предложенных вариантах")
		return teacherSearchAlbum, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state TeacherSearchAlbumState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	teacherNames, pointer, count, err := state.postgres.SearchAlbum.GetTeacherNames(vkID)
	if err != nil {
		return []*params.MessagesSendBuilder{}, err
	}
	countString := strconv.Itoa(count)

	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Найдено альбомов: " + countString + "\n" + teacherNames)
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
	k.AddTextButton("Завершить поиск", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state TeacherSearchAlbumState) Name() stateName {
	return teacherSearchAlbum
}
