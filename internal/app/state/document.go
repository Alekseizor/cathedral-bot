package state

import (
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"

	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
)

type DocumentStartState struct {
	postgres *postrgres.Repo
}

func (state DocumentStartState) Handler(msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Загрузка документа":
		return loadDocument, nil, nil
	case "Загрузка архива":
		return loadArchive, nil, nil
	case "Назад":
		return selectArchive, nil, nil
	default:
		return documentStart, nil, nil
	}
}

func (state DocumentStartState) Show(vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Добро пожаловать в архив документов. Выберите нужный пункт из списка ниже:")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Загрузка документа", "", "secondary")
	k.AddRow()
	k.AddTextButton("Загрузка архива", "", "secondary")
	k.AddRow()
	k.AddTextButton("Поиск документа", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state DocumentStartState) Name() stateName {
	return documentStart
}
