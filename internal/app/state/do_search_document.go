package state

import (
	"context"
	"fmt"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
)

// DoSearchDocumentState производится поиск документа
type DoSearchDocumentState struct {
	postgres *postrgres.Repo
}

func (state DoSearchDocumentState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Показать документы":
		return doSearchDocument, nil, nil
	case "Редактировать заявку":
		return doSearchDocument, nil, nil
	default:
		return doSearchDocument, nil, nil
	}
}

func (state DoSearchDocumentState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	searchParams, err := state.postgres.SearchDocument.ParseSearch(ctx, vkID)
	if err != nil {
		return []*params.MessagesSendBuilder{}, err
	}
	docNumber, err := state.postgres.Documents.SearchDocuments(ctx, searchParams)
	if err != nil {
		return []*params.MessagesSendBuilder{}, err
	}
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message(fmt.Sprintf("По вашему запросу нашлось %v документов. Показать документы или уточнить параметры?", docNumber))
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Показать документы", "", "secondary")
	k.AddRow()
	k.AddTextButton("Редактировать параметры", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state DoSearchDocumentState) Name() stateName {
	return doSearchDocument
}
