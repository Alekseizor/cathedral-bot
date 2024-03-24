package state

import (
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
)

type DocumentStubState struct {
	postgres *postrgres.Repo
}

func (state DocumentStubState) Handler(msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
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

func (state DocumentStubState) Show() ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Заглушка для документов")
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state DocumentStubState) Name() stateName {
	return documentStub
}
