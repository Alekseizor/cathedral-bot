package state

import (
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
)

type SelectArchiveState struct {
	postgres *postrgres.Repo
}

func (state SelectArchiveState) Handler(msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Фото":
		return photoStub, nil, nil
	case "Документы":
		return documentStub, nil, nil
	default:
		return selectArchive, nil, nil
	}
}

func (state SelectArchiveState) Show() ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("В нашем боте есть 2 архива: фотоархив и документоархив. С каким хочешь поработать?")
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Фото", "", "secondary")
	k.AddTextButton("Документы", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state SelectArchiveState) Name() stateName {
	return selectArchive
}
