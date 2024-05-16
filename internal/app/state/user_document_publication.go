package state

import (
	"context"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
)

// ShowUserDocumentPublicationsState вывод пользовательских заявок на публикацию документа
type ShowUserDocumentPublicationsState struct {
	postgres *postrgres.Repo
}

func (state ShowUserDocumentPublicationsState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return showUserDocumentPublication, nil, nil
	}

	switch messageText {
	case "Завершить просмотр заявок":
		err := state.postgres.UserDocumentPublication.DeletePointer(msg.PeerID)
		if err != nil {
			return showUserDocumentPublication, nil, err
		}
		return documentStart, nil, nil
	case "⬅️":
		err := state.postgres.UserDocumentPublication.ChangePointer(msg.PeerID, false)
		if err != nil {
			return showUserDocumentPublication, nil, err
		}
		return showUserDocumentPublication, nil, nil
	case "➡️":
		err := state.postgres.UserDocumentPublication.ChangePointer(msg.PeerID, true)
		if err != nil {
			return showUserDocumentPublication, nil, err
		}
		return showUserDocumentPublication, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Просматривайте заявки с помощью кнопок или завершите просмотр")
		return showUserDocumentPublication, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state ShowUserDocumentPublicationsState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	message, attachment, pointer, count, err := state.postgres.UserDocumentPublication.GetRequestDocument(vkID)
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

func (state ShowUserDocumentPublicationsState) Name() stateName {
	return showUserDocumentPublication
}
