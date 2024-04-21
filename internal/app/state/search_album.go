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
		return categorySearchAlbum, nil, nil
	case "Назад":
		err := state.postgres.SearchAlbum.DeleteSearchAlbum(ctx, msg.PeerID)
		if err != nil {
			return eventYearPhotoArchive, nil, err
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
	switch count {
	case 1:
		return findYearSearchAlbum, nil, nil
	case 0:
		return findYearSearchAlbum, nil, nil
	default:
		return findYearSearchAlbum, nil, nil
	}
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

// FindYearSearchAlbumState выводится количество найденных альбомов
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
		return showListSearchAlbum, nil, nil
	case "Добавить фильтр по программе обучения":
		return findYearSearchAlbum, nil, nil
	case "Назад":
		err := state.postgres.SearchAlbum.DeleteYear(ctx, msg.PeerID)
		if err != nil {
			return yearSearchAlbum, nil, err
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

	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Найдено альбомов: " + countString)
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

// ShowListSearchAlbumState пользователь получает список найденных альбомов
type ShowListSearchAlbumState struct {
	postgres *postrgres.Repo
}

func (state ShowListSearchAlbumState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return showListSearchAlbum, nil, nil
	}

	switch messageText {
	case "Завершить поиск":
		err := state.postgres.SearchAlbum.DeleteSearchAlbum(ctx, msg.PeerID)
		if err != nil {
			return eventYearPhotoArchive, nil, err
		}
		return photoStart, nil, nil
	case "Назад":
		return findYearSearchAlbum, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Такого действия нет в предложенных вариантах")
		return showListSearchAlbum, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state ShowListSearchAlbumState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	output, err := state.postgres.SearchAlbum.ShowList(ctx, vkID)
	if err != nil {
		return nil, err
	}

	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Найденные альбомы:\n" + output)
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Завершить поиск", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state ShowListSearchAlbumState) Name() stateName {
	return showListSearchAlbum
}
