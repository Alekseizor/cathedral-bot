package state

import (
	"context"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
)

// PersonalAccountPhotoState пользователь открывает личный кабинет со своими заявками на загрузку фотографий
type PersonalAccountPhotoState struct {
	postgres *postrgres.Repo
}

func (state PersonalAccountPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return personalAccountPhoto, nil, nil
	}

	switch messageText {
	case "Завершить просмотр заявок":
		err := state.postgres.PersonalAccountPhoto.DeletePointer(msg.PeerID)
		if err != nil {
			return personalAccountPhoto, nil, err
		}
		return photoStart, nil, nil
	case "⬅️":
		err := state.postgres.PersonalAccountPhoto.ChangePointer(msg.PeerID, false)
		if err != nil {
			return personalAccountPhoto, nil, err
		}
		return personalAccountPhoto, nil, nil
	case "➡️":
		err := state.postgres.PersonalAccountPhoto.ChangePointer(msg.PeerID, true)
		if err != nil {
			return personalAccountPhoto, nil, err
		}
		return personalAccountPhoto, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Такого действия нет в предложенных вариантах")
		return personalAccountPhoto, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state PersonalAccountPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	message, attachment, pointer, count, err := state.postgres.PersonalAccountPhoto.GetRequestPhoto(vkID)
	if err != nil {
		return []*params.MessagesSendBuilder{}, err
	}

	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message(message)
	b.Attachment(attachment)
	k := object.NewMessagesKeyboard(true)
	if count > 1 {
		k.AddRow()
		if pointer != 0 {
			k.AddTextButton("⬅️", "", "secondary")
		}
		if count-pointer > 1 {
			k.AddTextButton("➡️", "", "secondary")
		}
	}
	k.AddRow()
	k.AddTextButton("Завершить просмотр заявок", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state PersonalAccountPhotoState) Name() stateName {
	return personalAccountPhoto
}
