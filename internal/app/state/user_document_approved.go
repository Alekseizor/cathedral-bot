package state

import (
	"context"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
)

// ShowUserDocumentApprovedState вывод пользовательских заявок на публикацию документа
type ShowUserDocumentApprovedState struct {
	postgres *postrgres.Repo
}

func (state ShowUserDocumentApprovedState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return showUserDocumentApproved, nil, nil
	}

	switch messageText {
	case "Завершить просмотр документов":
		err := state.postgres.UserDocumentPublication.DeletePointer(msg.PeerID)
		if err != nil {
			return showUserDocumentApproved, nil, err
		}
		return documentStart, nil, nil
	case "⬅️":
		err := state.postgres.UserDocumentApproved.ChangePointer(msg.PeerID, false)
		if err != nil {
			return showUserDocumentApproved, nil, err
		}
		return showUserDocumentApproved, nil, nil
	case "➡️":
		err := state.postgres.UserDocumentApproved.ChangePointer(msg.PeerID, true)
		if err != nil {
			return showUserDocumentApproved, nil, err
		}
		return showUserDocumentApproved, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Просматривайте документы с помощью кнопок или завершите просмотр")
		return showUserDocumentApproved, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state ShowUserDocumentApprovedState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	message, attachment, pointer, count, err := state.postgres.UserDocumentApproved.GetApprovedDocument(vkID)
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
	k.AddTextButton("Завершить просмотр документов", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state ShowUserDocumentApprovedState) Name() stateName {
	return showUserDocumentApproved
}
