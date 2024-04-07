package state

import (
	"context"
	"fmt"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	"github.com/lib/pq"
	"strconv"
	"strings"
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
		cleanedOffset := strings.Replace(msg.Payload, "\"", "", -1)
		offset, err := strconv.Atoi(cleanedOffset)
		if err != nil {
			return showSearchDocument, []*params.MessagesSendBuilder{}, err
		}
		err = state.postgres.SearchDocument.UpdatePointer(ctx, offset, msg.PeerID)
		if err != nil {
			return showSearchDocument, []*params.MessagesSendBuilder{}, err
		}
		return showSearchDocument, []*params.MessagesSendBuilder{}, nil
	case "➡️":
		cleanedOffset := strings.Replace(msg.Payload, "\"", "", -1)
		offset, err := strconv.Atoi(cleanedOffset)
		if err != nil {
			return showSearchDocument, []*params.MessagesSendBuilder{}, err
		}
		err = state.postgres.SearchDocument.UpdatePointer(ctx, offset, msg.PeerID)
		if err != nil {
			return showSearchDocument, []*params.MessagesSendBuilder{}, err
		}
		return showSearchDocument, []*params.MessagesSendBuilder{}, nil
	default:
		searchParams, err := state.postgres.SearchDocument.ParseSearchButtons(ctx, msg.PeerID)
		if err != nil {
			return showSearchDocument, []*params.MessagesSendBuilder{}, err
		}
		var documentIDs pq.Int64Array
		if len(searchParams.Documents) >= searchParams.PointerDoc+5 {
			documentIDs = searchParams.Documents[searchParams.PointerDoc : searchParams.PointerDoc+5]
		} else {
			documentIDs = searchParams.Documents[searchParams.PointerDoc:]
		}
		docNumber, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введите номер документа из списка числом без лишних символов")
			return showSearchDocument, []*params.MessagesSendBuilder{b}, nil
		}
		if docNumber >= searchParams.PointerDoc+1 && docNumber <= searchParams.PointerDoc+len(documentIDs) {
			err = state.postgres.SearchDocument.UpdateChosenDocSearch(ctx, int(searchParams.Documents[docNumber-1]), msg.PeerID)
			return showChosenDocument, []*params.MessagesSendBuilder{}, nil
		} else {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введите номер документа из списка числом без лишних символов")
			return showSearchDocument, []*params.MessagesSendBuilder{b}, nil
		}
	}
}

func (state ShowSearchDocumentState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	output, err := state.postgres.Documents.GetSearchDocuments(ctx, vkID)
	if err != nil {
		return []*params.MessagesSendBuilder{}, err
	}
	searchParams, err := state.postgres.SearchDocument.ParseSearchButtons(ctx, vkID)
	if err != nil {
		return []*params.MessagesSendBuilder{}, err
	}
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите номер нужного документа из списка:\n" + output)
	if len(searchParams.Documents) >= 6 {
		k := object.NewMessagesKeyboardInline()
		k.AddRow()
		if !(searchParams.PointerDoc == 0) {
			k.AddTextButton("⬅️", -5, "secondary")
		}
		if !(len(searchParams.Documents)-searchParams.PointerDoc < 6) {
			k.AddTextButton("➡️", 5, "secondary")
		}
		b.Keyboard(k)
	}

	return []*params.MessagesSendBuilder{b}, nil
}

func (state ShowSearchDocumentState) Name() stateName {
	return showSearchDocument
}

// ShowChosenDocumentState пользователь смотрит выбранный им документ
type ShowChosenDocumentState struct {
	postgres *postrgres.Repo
}

func (state ShowChosenDocumentState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Назад":
		return showSearchDocument, nil, nil
	case "Закончить поиск":
		return documentStart, nil, nil
	default:
		return showSearchDocument, nil, nil
	}
}

func (state ShowChosenDocumentState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	searchParams, err := state.postgres.SearchDocument.ParseSearchButtons(ctx, vkID)
	if err != nil {
		return []*params.MessagesSendBuilder{}, err
	}
	output, attachment, err := state.postgres.Documents.GetOutput(ctx, *searchParams.ChosenDoc)
	if err != nil {
		return []*params.MessagesSendBuilder{}, err
	}
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Выбранный документ:\n" + output)
	b.Attachment(attachment)
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	k.AddRow()
	k.AddTextButton("Закончить поиск", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state ShowChosenDocumentState) Name() stateName {
	return showChosenDocument
}
