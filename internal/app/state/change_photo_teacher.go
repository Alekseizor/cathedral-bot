package state

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
)

type ChangeNameTeacherPhotoState struct {
	postgres *postrgres.Repo
	vk       *api.VK
	vkUser   *api.VK
}

func (state ChangeNameTeacherPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.ObjectAdmin.Get(ctx, msg.PeerID)
	if err != nil {
		return changeNameTeacherPhoto, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Новый преподаватель":
		return changeUserNameTeacherPhoto, nil, nil
	case "Назад":
		return workingAlbums, nil, nil
	default:
		maxID, err := state.postgres.RequestPhoto.GetTeacherMaxID()
		if err != nil {
			return changeNameTeacherPhoto, []*params.MessagesSendBuilder{}, err
		}

		eventNumber, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введите номер преподавателя числом")
			return changeNameTeacherPhoto, []*params.MessagesSendBuilder{b}, nil
		}

		if !(eventNumber >= 1 && eventNumber <= maxID) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Такого события нет в списке")
			return changeNameTeacherPhoto, []*params.MessagesSendBuilder{b}, nil
		}

		err = state.postgres.TeacherAlbums.UpdateName(ctx, photoID, eventNumber)
		if err != nil {
			return changeNameTeacherPhoto, nil, err
		}

		albumVKID, err := state.postgres.TeacherAlbums.GetVKID(ctx, photoID)
		if err != nil {
			return changeNameTeacherPhoto, nil, err
		}

		// достали id группы
		groupID, err := state.vk.GroupsGetByID(nil)
		if err != nil {
			return changeNameTeacherPhoto, nil, err
		}

		newAlbumName, err := state.postgres.RequestPhoto.GetTeacherName(ctx, eventNumber)
		if err != nil {
			return changeNameTeacherPhoto, nil, err
		}

		_, err = state.vkUser.PhotosEditAlbum(api.Params{"album_id": albumVKID, "title": newAlbumName, "owner_id": groupID[0].ID * (-1)})
		if err != nil {
			return changeNameTeacherPhoto, nil, err
		}

		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message(fmt.Sprintf("Для альбома №%d изменено название на - %s", photoID, newAlbumName))
		return workingAlbums, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state ChangeNameTeacherPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	teachers, err := state.postgres.RequestPhoto.GetAllTeacherNames()
	if err != nil {
		return []*params.MessagesSendBuilder{}, err
	}
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите номер нового преподавателя из списка ниже:\n" + teachers)
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Новый преподаватель", "", "secondary")
	addBackButton(k)
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state ChangeNameTeacherPhotoState) Name() stateName {
	return changeNameTeacherPhoto
}

type ChangeDescriptionPhotoTeacherState struct {
	postgres *postrgres.Repo
	vk       *api.VK
	vkUser   *api.VK
}

func (state ChangeDescriptionPhotoTeacherState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Назад":
		return workingAlbums, nil, nil
	default:
		photoID, err := state.postgres.ObjectAdmin.Get(ctx, msg.PeerID)
		if err != nil {
			return changeDescriptionPhotoTeacher, []*params.MessagesSendBuilder{}, err
		}

		err = state.postgres.TeacherAlbums.UpdateDescription(ctx, photoID, messageText)
		if err != nil {
			return changeDescriptionPhotoTeacher, []*params.MessagesSendBuilder{}, err
		}

		albumVKID, err := state.postgres.TeacherAlbums.GetVKID(ctx, photoID)
		if err != nil {
			return changeDescriptionPhotoTeacher, []*params.MessagesSendBuilder{}, err
		}

		// достали id группы
		groupID, err := state.vk.GroupsGetByID(nil)
		if err != nil {
			return changeDescriptionPhotoTeacher, []*params.MessagesSendBuilder{}, err
		}

		_, err = state.vkUser.PhotosEditAlbum(api.Params{"album_id": albumVKID, "description": messageText, "owner_id": groupID[0].ID * (-1)})
		if err != nil {
			return changeDescriptionPhotoTeacher, []*params.MessagesSendBuilder{}, err
		}

		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message(fmt.Sprintf("Для альбома №%d изменено описание на - %s", photoID, messageText))

		return workingAlbums, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state ChangeDescriptionPhotoTeacherState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите описание альбома")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state ChangeDescriptionPhotoTeacherState) Name() stateName {
	return changeDescriptionPhotoTeacher
}

type ChangeUserNameTeacherPhotoState struct {
	postgres *postrgres.Repo
	vk       *api.VK
	vkUser   *api.VK
}

func (state ChangeUserNameTeacherPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	switch messageText {
	case "Назад":
		return changeNameTeacherPhoto, nil, nil
	default:
		photoID, err := state.postgres.ObjectAdmin.Get(ctx, msg.PeerID)
		if err != nil {
			return changeUserNameTeacherPhoto, []*params.MessagesSendBuilder{}, err
		}

		err = state.postgres.TeacherAlbums.UpdateNewName(ctx, photoID, messageText)
		if err != nil {
			return changeUserNameTeacherPhoto, []*params.MessagesSendBuilder{}, err
		}

		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message(fmt.Sprintf("Для альбома №%d задано новое название - %s", photoID, messageText))
		return workingAlbums, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state ChangeUserNameTeacherPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите ФИО нового преподавателя")
	k := object.NewMessagesKeyboard(true)
	addBackButton(k)
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state ChangeUserNameTeacherPhotoState) Name() stateName {
	return changeUserNameTeacherPhoto
}
