package state

import (
	"context"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
)

// ViewRequestsPhotoState админ смотрит заявки на добавление фото в альбом
type ViewRequestsPhotoState struct {
	postgres *postrgres.Repo
	vk       *api.VK
	vkUser   *api.VK
	groupID  int
}

func (state ViewRequestsPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return viewRequestsPhoto, nil, nil
	}

	switch messageText {
	case "Одобрить":
		comment, err := state.postgres.ViewRequestPhoto.ApprovePhoto(msg.PeerID, state.vkUser, state.groupID)
		if err != nil {
			return viewRequestsPhoto, nil, err
		}
		if comment != "" {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message(comment)
			return viewRequestsPhoto, []*params.MessagesSendBuilder{b}, nil
		}
		return viewRequestsPhoto, nil, nil
	case "Отклонить":
		err := state.postgres.ViewRequestPhoto.RejectPhoto(msg.PeerID)
		if err != nil {
			return viewRequestsPhoto, nil, err
		}
		return viewRequestsPhoto, nil, nil
	case "Редактировать":
		return editRequestPhoto, nil, nil
	case "Завершить просмотр заявок":
		err := state.postgres.ViewRequestPhoto.DeletePointer(msg.PeerID)
		if err != nil {
			return viewRequestsPhoto, nil, err
		}
		return photoStart, nil, nil
	case "⬅️":
		err := state.postgres.ViewRequestPhoto.ChangePointer(msg.PeerID, false)
		if err != nil {
			return viewRequestsPhoto, nil, err
		}
		return viewRequestsPhoto, nil, nil
	case "➡️":
		err := state.postgres.ViewRequestPhoto.ChangePointer(msg.PeerID, true)
		if err != nil {
			return viewRequestsPhoto, nil, err
		}
		return viewRequestsPhoto, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Такого действия нет в предложенных вариантах")
		return viewRequestsPhoto, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state ViewRequestsPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	message, attachment, pointer, count, err := state.postgres.ViewRequestPhoto.GetRequestPhoto(vkID)
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
	k.AddTextButton("Одобрить", "", "secondary")
	k.AddTextButton("Отклонить", "", "secondary")
	k.AddRow()
	k.AddTextButton("Редактировать", "", "secondary")
	k.AddTextButton("Завершить просмотр заявок", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state ViewRequestsPhotoState) Name() stateName {
	return viewRequestsPhoto
}
