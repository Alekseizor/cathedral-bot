package state

import (
	"context"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
)

var validExtension = map[string]struct{}{
	"jpg":  struct{}{},
	"jpeg": struct{}{},
	"png":  struct{}{},
	"tiff": struct{}{},
}

// LoadPhotoState пользователь загружает документ
type LoadPhotoState struct {
	postgres *postrgres.Repo
	vk       *api.VK
}

func (state LoadPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "Назад" {
		return photoStub, nil, nil
	}
	attachment := msg.Attachments

	if len(attachment) == 0 {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Загрузите вашу фотографию, прикрепив её к сообщению")
		return loadPhoto, []*params.MessagesSendBuilder{b}, nil
	}

	if len(attachment) > 1 {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Можно загрузить лишь одну фотографию, для загрузки множества фотографий воспользуйтесь загрузкой архива")
		return loadPhoto, []*params.MessagesSendBuilder{b}, nil
	}

	if attachment[0].Photo.AccessKey != "" {
		err := state.postgres.RequestPhoto.UploadPhoto(ctx, state.vk, attachment[0].Photo, msg.PeerID)
		if err != nil {
			return loadPhoto, []*params.MessagesSendBuilder{}, err
		}

		switch messageText {
		default:
			return loadPhoto, nil, nil
		}
	}

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

	switch messageText {
	default:
		return loadPhoto, nil, nil
	}
}

func (state LoadPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Загрузите фото. Допустимые  форматы фото: jpg, jpeg, png, tiff")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state LoadPhotoState) Name() stateName {
	return loadPhoto
}
