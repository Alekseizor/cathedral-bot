package state

import (
	"context"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
)

// CreateAlbumState пользователь выбирает чей альбом создать
type CreateAlbumState struct {
	postgres *postrgres.Repo
	vk       *api.VK
}

func (state CreateAlbumState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return createAlbum, nil, nil
	}

	switch messageText {
	case "Назад":
		return photoStart, nil, nil
	case "Студентов":
		return createStudentAlbum, nil, nil
	case "Преподавателя":
		return createTeacherAlbum, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Такого действия нет в предложенных вариантах")
		return createAlbum, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state CreateAlbumState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Чей альбом создать?")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Студентов", "", "secondary")
	k.AddTextButton("Преподавателя", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state CreateAlbumState) Name() stateName {
	return createAlbum
}

// CreateStudentAlbumState пользователь создаёт альбом студентов
type CreateStudentAlbumState struct {
	postgres *postrgres.Repo
	vk       *api.VK
	vkUser   *api.VK
	groupID  int
}

func (state CreateStudentAlbumState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return createStudentAlbum, nil, nil
	}

	switch messageText {
	case "Назад":
		return createAlbum, nil, nil
	default:
		url, flag, err := state.postgres.CreateAlbum.CreateStudentAlbum(ctx, state.vkUser, state.groupID, messageText)
		if err != nil {
			return createStudentAlbum, nil, nil
		}
		if flag {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Такой альбом уже существует")
			return createTeacherAlbum, []*params.MessagesSendBuilder{b}, nil
		}
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Альбом успешно создан\n" + url)
		return photoStart, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state CreateStudentAlbumState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напиши название альбома в формате:\nГод//Программа обучения//Событие\nНапример:\n2024//Бакалавриат//Защита диплома")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state CreateStudentAlbumState) Name() stateName {
	return createStudentAlbum
}

// CreateTeacherAlbumState пользователь создаёт альбом преподавателя
type CreateTeacherAlbumState struct {
	postgres *postrgres.Repo
	vk       *api.VK
	vkUser   *api.VK
	groupID  int
}

func (state CreateTeacherAlbumState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return createTeacherAlbum, nil, nil
	}

	switch messageText {
	case "Назад":
		return createAlbum, nil, nil
	default:
		url, flag, err := state.postgres.CreateAlbum.CreateTeacherAlbum(ctx, state.vkUser, state.groupID, messageText)
		if err != nil {
			return createTeacherAlbum, nil, nil
		}
		if flag {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Такой альбом уже существует")
			return createTeacherAlbum, []*params.MessagesSendBuilder{b}, nil
		}
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Альбом успешно создан\n" + url)
		return photoStart, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state CreateTeacherAlbumState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напиши полностью ФИО преподавателя\nНапример: Филиппович Анна Юрьевна")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state CreateTeacherAlbumState) Name() stateName {
	return createTeacherAlbum
}
