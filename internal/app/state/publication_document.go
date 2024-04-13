package state

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
)

///////////

type RequestDocumentFromQueueState struct {
	postgres *postrgres.Repo
}

func (state RequestDocumentFromQueueState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Заявки из очереди":
		return workingRequestDocument, nil, nil
	case "Конкретная заявка":
		return workingRequestDocument, nil, nil
	case "Назад":
		return documentCabinet, nil, nil
	default:
		return workingRequestDocument, nil, nil
	}
}

func (state RequestDocumentFromQueueState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	output, attachment, requestID, err := state.postgres.RequestsDocuments.GetRequestFromQueue(ctx)
	if err != nil {
		//if err == sql.ErrNoRows {
		if errors.Is(err, sql.ErrNoRows) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Заявок, ожидающих проверки, нет!")
			k := object.NewMessagesKeyboard(true)
			addBackButton(k)
			b.Keyboard(k)
			return []*params.MessagesSendBuilder{b}, nil
		}
		return nil, fmt.Errorf("[requests_documents.GetRequestFromQueue]: %w", err)
	}

	err = state.postgres.ObjectAdmin.Update(ctx, requestID, vkID)
	if err != nil {
		return nil, fmt.Errorf("[object_admin.Update]: %w", err)
	}

	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message(output)
	b.Attachment(attachment)
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Изменить заявку", "", "secondary")
	k.AddTextButton("Принять заявку", "", "secondary")
	k.AddTextButton("Отклонить заявку", "", "secondary")
	addBackButton(k)
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state RequestDocumentFromQueueState) Name() stateName {
	return requestDocumentFromQueue
}
