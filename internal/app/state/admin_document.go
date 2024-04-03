package state

import (
	"context"

	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"

	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
)

type DocumentCabinetState struct {
	postgres *postrgres.Repo
}

func (state DocumentCabinetState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Работа с заявкой":
		return selectArchive, nil, nil
	case "Работа с файлом":
		return documentCabinet, nil, nil
	case "Заблокировать пользователя":
		return blockUser, nil, nil
	case "Добавить администратора":
		return documentCabinet, nil, nil
	case "Удалить администратора":
		return documentCabinet, nil, nil
	case "Выйти из кабинета администратора":
		return selectArchive, nil, nil
	default:
		return documentCabinet, nil, nil
	}
}

func (state DocumentCabinetState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Вы в кабинете администратора документоархива, выберите действие")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Работа с заявкой", "", "secondary")
	k.AddTextButton("Работа с файлом", "", "secondary")
	k.AddRow()
	k.AddTextButton("Заблокировать пользователя", "", "secondary")
	k.AddRow()
	k.AddTextButton("Добавить администратора", "", "secondary")
	k.AddTextButton("Удалить администратора", "", "secondary")
	k.AddRow()
	k.AddTextButton("Выйти из кабинета администратора", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state DocumentCabinetState) Name() stateName {
	return documentCabinet
}
