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
		return showSearchDocument, nil, nil
	case "Редактировать параметры":
		return editSearchDocument, nil, nil
	default:
		return doSearchDocument, nil, nil
	}
}

func (state DoSearchDocumentState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	searchParams, err := state.postgres.SearchDocument.ParseSearch(ctx, vkID)
	if err != nil {
		return []*params.MessagesSendBuilder{}, err
	}
	docNumber, err := state.postgres.Documents.SearchDocuments(ctx, searchParams, vkID)
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

// ShowSearchDocumentState вывод найденных документов
type ShowSearchDocumentState struct {
	postgres *postrgres.Repo
}

func (state ShowSearchDocumentState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "⬅️":
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Ваши документы:\n1.---\n2.---\n3.---\n4.---\n5.---")
		return showSearchDocument, []*params.MessagesSendBuilder{b}, nil
	case "➡️":
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Ваши документы:\n1.---\n2.---\n3.---\n4.---\n5.---")
		return showSearchDocument, []*params.MessagesSendBuilder{b}, nil
	default:
		return showSearchDocument, nil, nil
	}
}

func (state ShowSearchDocumentState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("________________________________")
	k := object.NewMessagesKeyboardInline()
	k.AddRow()
	k.AddTextButton("⬅️", "", "secondary")
	k.AddTextButton("➡️", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state ShowSearchDocumentState) Name() stateName {
	return showSearchDocument
}
