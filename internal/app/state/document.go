package state

import (
	"context"
	"fmt"

	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"

	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
)

type DocumentStubState struct {
	postgres *postrgres.Repo
}

func (state DocumentStubState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Кабинет администратора документоархива":
		return documentCabinet, nil, nil
	case "Назад":
		return selectArchive, nil, nil
	default:
		return documentStub, nil, nil
	}
}

func (state DocumentStubState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Заглушка для документов")
	k := object.NewMessagesKeyboard(true)
	documentsAdmins, err := state.postgres.Admin.GetDocumentsAdmins(ctx)
	if err != nil {
		return nil, fmt.Errorf("[admin.GetDocumentsAdmins]: %w", err)
	}

	if contains(documentsAdmins, int64(vkID)) {
		k.AddRow()
		k.AddTextButton("Кабинет администратора документоархива", "", "secondary")
	}

	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state DocumentStubState) Name() stateName {
	return documentStub
}
