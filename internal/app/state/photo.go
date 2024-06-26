package state

import (
	"context"
	"fmt"

	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"

	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
)

type PhotoStartState struct {
	postgres *postrgres.Repo
}

func (state PhotoStartState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Загрузка фото":
		return loadPhoto, nil, nil
	case "Загрузка архива":
		return loadPhotoArchive, nil, nil
	case "Поиск альбома":
		err := state.postgres.SearchAlbum.CreateSearchAlbum(ctx, msg.PeerID)
		if err != nil {
			return photoStart, []*params.MessagesSendBuilder{}, err
		}
		return categorySearchAlbum, nil, nil
	case "Личный кабинет":
		err := state.postgres.PersonalAccountPhoto.CreatePersonalAccountPhoto(ctx, msg.PeerID)
		if err != nil {
			return photoStart, []*params.MessagesSendBuilder{}, err
		}
		return personalAccountPhoto, nil, nil
	case "Кабинет администратора фотоархива":
		return albumsCabinet, nil, nil
	case "Назад":
		return selectArchive, nil, nil
	default:
		return photoStart, nil, nil
	}
}

func (state PhotoStartState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Добро пожаловать в архив фотографий. Выберите нужное действие")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Загрузка фото", "", "secondary")
	k.AddTextButton("Загрузка архива", "", "secondary")
	k.AddRow()
	k.AddTextButton("Поиск альбома", "", "secondary")
	k.AddTextButton("Личный кабинет", "", "secondary")

	albumsAdmins, err := state.postgres.Admin.GetAlbumsAdmins(ctx)
	if err != nil {
		return nil, fmt.Errorf("[admin.GetDocumentsAdmins]: %w", err)
	}

	if contains(albumsAdmins, int64(vkID)) {
		k.AddRow()
		k.AddTextButton("Кабинет администратора фотоархива", "", "secondary")
	}

	addBackButton(k)
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state PhotoStartState) Name() stateName {
	return photoStart
}
