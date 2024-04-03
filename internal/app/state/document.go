package state

import (
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"

	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
)

type DocumentStubState struct {
	postgres *postrgres.Repo
}

func (state DocumentStubState) Handler(msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	state.postgres.Admin.GetDocumentsAdmins()
	switch messageText {
	case "Назад":
		return selectArchive, nil, nil
	default:
		return documentStub, nil, nil
	}
}

func (state DocumentStubState) Show() ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Заглушка для документов")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state DocumentStubState) Name() stateName {
	return documentStub
}
