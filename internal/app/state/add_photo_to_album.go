package state

import (
	"context"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
)

// AddPhotoToAlbumState пользователь добавляет фотографию в альбом
type AddPhotoToAlbumState struct {
	postgres *postrgres.Repo
	vk       *api.VK
	vkUser   *api.VK
	groupID  int
}

func (state AddPhotoToAlbumState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return addPhotoToAlbum, nil, nil
	}

	switch messageText {
	case "Назад":
		return photoStart, nil, nil
	default:
		err := state.postgres.AddPhotoToAlbum.AddPhotoToAlbum(ctx, state.vkUser, 302847033, state.groupID, "https://sun9-67.userapi.com/impg/uIPQ5qaDjyIcN_EGbevzRraknKJQcePdxqSLWA/ioBpykqF1Js.jpg?size=1280x960&quality=96&sign=0b6ab11a40e0936ac1d73f15943705c3&c_uniq_tag=GM0R-7W9pSjmsyXznk9ABjIBKKLQ8p7yHg4HPyAuzv4&type=album")
		if err != nil {
			return photoStart, nil, err
		}

		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Фотография успешно добавлена в альбом")
		return photoStart, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state AddPhotoToAlbumState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напиши 1")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state AddPhotoToAlbumState) Name() stateName {
	return addPhotoToAlbum
}
