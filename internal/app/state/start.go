package state

import (
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"

	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
)

type StartState struct {
	postgres *postrgres.Repo
}

func (state StartState) Handler(msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "Начать" {
		return selectArchive, nil, nil
	}

	return start, nil, nil
}

func (state StartState) Show() ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.PeerID(236322856)
	b.Message("Привет! Для того, чтобы начать работу с кафедральным ботом, нажми кнопку 'Начать'")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Начать", "", "")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state StartState) Name() stateName {
	return start
}

/////////////////////////////////////////////
