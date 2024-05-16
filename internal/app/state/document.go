package state

import (
	"context"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"

	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
)

type DocumentStartState struct {
	postgres *postrgres.Repo
}

func (state DocumentStartState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Загрузка документа":
		return loadDocument, nil, nil
	case "Загрузка архива":
		return loadArchive, nil, nil
	case "Поиск документа":
		err := state.postgres.SearchDocument.CreateSearch(ctx, msg.PeerID)
		if err != nil {
			return documentStart, []*params.MessagesSendBuilder{}, err
		}
		return nameSearchDocument, nil, nil
	case "Кабинет администратора документоархива":
		return documentCabinet, nil, nil
	case "Мои заявки":
		err := state.postgres.UserDocumentPublication.CreateUserDocumentPublication(ctx, msg.PeerID)
		if err != nil {
			return documentStart, []*params.MessagesSendBuilder{}, err
		}
		return showUserDocumentPublication, nil, nil
	case "Мои документы":
		err := state.postgres.UserDocumentApproved.CreateUserDocumentApproved(ctx, msg.PeerID)
		if err != nil {
			return documentStart, []*params.MessagesSendBuilder{}, err
		}
		return showUserDocumentApproved, nil, nil
	case "Назад":
		return selectArchive, nil, nil
	default:
		return documentStart, nil, nil
	}
}

func (state DocumentStartState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Добро пожаловать в архив документов. Выберите нужный пункт из списка ниже:")
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
	k.AddTextButton("Загрузка документа", "", "secondary")
	k.AddRow()
	k.AddTextButton("Загрузка архива", "", "secondary")
	k.AddRow()
	k.AddTextButton("Поиск документа", "", "secondary")
	k.AddRow()
	k.AddTextButton("Мои заявки", "", "secondary")
	k.AddTextButton("Мои документы", "", "secondary")
	addBackButton(k)
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state DocumentStartState) Name() stateName {
	return documentStart
}
