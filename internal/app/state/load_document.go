package state

import (
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
)

var validExtension = map[string]struct{}{
	"doc":  struct{}{},
	"docx": struct{}{},
	"ppt":  struct{}{},
	"pptx": struct{}{},
	"txt":  struct{}{},
	"pdf":  struct{}{},
}

type LoadDocumentState struct {
	postgres *postrgres.Repo
}

func (state LoadDocumentState) Handler(msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	attachment := msg.Attachments

	if len(attachment) > 1 {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Можно загрузить лишь один документ, для загрузки множества документов воспользуйтесь загрузкой архива")
		return loadDocument, []*params.MessagesSendBuilder{b}, nil
	}

	if _, ok := validExtension[attachment[0].Doc.Ext]; !ok {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Поддерживаются только документы формата doc/docx, pdf, txt, ppt/pptx.")
		return loadDocument, []*params.MessagesSendBuilder{b}, nil
	}

	switch messageText {
	case "Назад":
		return documentStart, nil, nil
	default:
		return loadDocument, nil, nil
	}
}

func (state LoadDocumentState) Show() ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Загрузите ваш документ")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state LoadDocumentState) Name() stateName {
	return loadDocument
}
