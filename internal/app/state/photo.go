package state

import (
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
)

type PhotoStubState struct {
	postgres *postrgres.Repo
}

func (state PhotoStubState) Handler(msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Назад":
		return selectArchive, nil, nil
	default:
		return photoStub, nil, nil
	}
}

func (state PhotoStubState) Show() ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Заглушка для фото")
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state PhotoStubState) Name() stateName {
	return photoStub
}
