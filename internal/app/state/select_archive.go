package state

import (
	"context"

	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"

	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
)

type SelectArchiveState struct {
	postgres *postrgres.Repo
}

func (state SelectArchiveState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Фото":
		return photoStub, nil, nil
	case "Документы":
		return documentStart, nil, nil
	default:
		return selectArchive, nil, nil
	}
}

func (state SelectArchiveState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("В нашем боте есть 2 архива: фотоархив и документоархив. С каким хочешь поработать?")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Фото", "", "secondary")
	k.AddTextButton("Документы", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state SelectArchiveState) Name() stateName {
	return selectArchive
}
