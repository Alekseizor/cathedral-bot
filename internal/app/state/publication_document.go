package state

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Alekseizor/cathedral-bot/internal/app/ds"
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

	requestID, err := state.postgres.ObjectAdmin.Get(ctx, msg.PeerID)
	if err != nil {
		return "", nil, fmt.Errorf("[object_admin.Get]: %w", err)
	}

	switch messageText {
	case "Изменить заявку":
		return editDocumentAdmin, nil, nil
	case "Принять заявку":
		reqDocument, err := state.postgres.RequestsDocuments.GetByID(ctx, requestID)
		if err != nil {
			return "", nil, fmt.Errorf("[requests_documents.GetByID]: %w", err)
		}

		if reqDocument.IsCategoryNew {
			err = state.postgres.Documents.NewCategory(ctx, reqDocument.Category.String)
			if err != nil {
				return "", nil, fmt.Errorf("[documents.NewCategory]: %w", err)
			}
		}

		err = state.postgres.Documents.UploadDocument(ctx, *reqDocument)
		if err != nil {
			return "", nil, fmt.Errorf("[documents.UploadDocument]: %w", err)
		}

		err = state.postgres.RequestsDocuments.DeleteDocumentRequestByID(ctx, reqDocument.ID)
		if err != nil {
			return "", nil, fmt.Errorf("[requests_documents.DeleteDocumentRequestByID]: %w", err)
		}
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message(fmt.Sprintf("Файл c названием - %s - успешно добавлен в документоархив", reqDocument.Title))
		k := object.NewMessagesKeyboard(true)
		addBackButton(k)
		b.Keyboard(k)
		return workingRequestDocument, []*params.MessagesSendBuilder{b}, nil
	case "Отклонить заявку":
		err = state.postgres.RequestsDocuments.UpdateStatus(ctx, ds.StatusAdminDeclined, requestID)
		if err != nil {
			return "", nil, fmt.Errorf("[requests_documents.UpdateStatus]: %w", err)
		}

		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message(fmt.Sprintf("Заявка c ID - %d отклонена", requestID))
		k := object.NewMessagesKeyboard(true)
		addBackButton(k)
		b.Keyboard(k)
		return workingRequestDocument, []*params.MessagesSendBuilder{b}, nil
	case "Назад":
		status, err := state.postgres.RequestsDocuments.GetStatus(ctx, requestID)
		if status == ds.StatusAdminWorking {
			err = state.postgres.RequestsDocuments.UpdateStatus(ctx, ds.StatusUserConfirmed, requestID)
			if err != nil {
				return "", nil, fmt.Errorf("[requests_documents.UpdateStatus]: %w", err)
			}
		}
		return workingRequestDocument, nil, nil
	default:
		return workingRequestDocument, nil, nil
	}
}

func (state RequestDocumentFromQueueState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	output, attachment, requestID, err := state.postgres.RequestsDocuments.GetRequestFromQueue(ctx)
	if err != nil {
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

	err = state.postgres.RequestsDocuments.UpdateStatus(ctx, ds.StatusAdminWorking, requestID)
	if err != nil {
		return nil, fmt.Errorf("[requests_documents.UpdateStatus]: %w", err)
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

type RequestDocumentEntrySpecificApplicationState struct {
	postgres *postrgres.Repo
}

func (state RequestDocumentEntrySpecificApplicationState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Назад":
		return workingRequestDocument, nil, nil
	default:
		requestID, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("ID заявки должно быть числом!")
			k := object.NewMessagesKeyboard(true)
			addBackButton(k)
			b.Keyboard(k)
			return requestDocumentEntrySpecificApplication, []*params.MessagesSendBuilder{b}, nil
		}

		requestDocument, err := state.postgres.RequestsDocuments.GetByID(ctx, requestID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				b := params.NewMessagesSendBuilder()
				b.RandomID(0)
				b.Message(fmt.Sprintf("Заяки с ID - %d - не найдено. Возможно, файл из этой заявки уже опубликован!", requestID))
				k := object.NewMessagesKeyboard(true)
				addBackButton(k)
				b.Keyboard(k)
				return requestDocumentEntrySpecificApplication, []*params.MessagesSendBuilder{b}, nil
			}
			return requestDocumentEntrySpecificApplication, nil, fmt.Errorf("[requests_documents.GetByID]: %w", err)
		}

		switch requestDocument.Status {
		case ds.StatusAdminDeclined:
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message(fmt.Sprintf("Заяки с ID - %d - уже отклонена", requestID))
			k := object.NewMessagesKeyboard(true)
			addBackButton(k)
			b.Keyboard(k)
			return requestDocumentEntrySpecificApplication, []*params.MessagesSendBuilder{b}, nil
		case ds.StatusAdminWorking:
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message(fmt.Sprintf("С заякой №%d сейчас работает другой администратор", requestID))
			k := object.NewMessagesKeyboard(true)
			addBackButton(k)
			b.Keyboard(k)
			return requestDocumentEntrySpecificApplication, []*params.MessagesSendBuilder{b}, nil
		}

		err = state.postgres.ObjectAdmin.Update(ctx, requestID, msg.PeerID)
		if err != nil {
			return requestDocumentEntrySpecificApplication, nil, fmt.Errorf("[object_admin.Update]: %w", err)
		}

		err = state.postgres.RequestsDocuments.UpdateStatus(ctx, ds.StatusAdminWorking, requestID)
		if err != nil {
			return requestDocumentEntrySpecificApplication, nil, fmt.Errorf("[requests_documents.UpdateStatus]: %w", err)
		}

		return requestDocumentSpecificApplication, nil, nil
	}
}

func (state RequestDocumentEntrySpecificApplicationState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите ID заявки, над которой хотите поработать")
	k := object.NewMessagesKeyboard(true)
	addBackButton(k)
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state RequestDocumentEntrySpecificApplicationState) Name() stateName {
	return requestDocumentEntrySpecificApplication
}

type RequestDocumentSpecificApplicationState struct {
	postgres *postrgres.Repo
}

func (state RequestDocumentSpecificApplicationState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	requestID, err := state.postgres.ObjectAdmin.Get(ctx, msg.PeerID)
	if err != nil {
		return "", nil, fmt.Errorf("[object_admin.Get]: %w", err)
	}

	switch messageText {
	case "Изменить заявку":
		return editDocumentAdmin, nil, nil
	case "Принять заявку":
		reqDocument, err := state.postgres.RequestsDocuments.GetByID(ctx, requestID)
		if err != nil {
			return "", nil, fmt.Errorf("[requests_documents.GetByID]: %w", err)
		}

		if reqDocument.IsCategoryNew {
			err = state.postgres.Documents.NewCategory(ctx, reqDocument.Category.String)
			if err != nil {
				return "", nil, fmt.Errorf("[documents.NewCategory]: %w", err)
			}
		}

		err = state.postgres.Documents.UploadDocument(ctx, *reqDocument)
		if err != nil {
			return "", nil, fmt.Errorf("[documents.UploadDocument]: %w", err)
		}

		err = state.postgres.RequestsDocuments.DeleteDocumentRequestByID(ctx, reqDocument.ID)
		if err != nil {
			return "", nil, fmt.Errorf("[requests_documents.DeleteDocumentRequestByID]: %w", err)
		}
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message(fmt.Sprintf("Файл c названием - %s - успешно добавлен в документоархив", reqDocument.Title))
		k := object.NewMessagesKeyboard(true)
		addBackButton(k)
		b.Keyboard(k)
		return workingRequestDocument, []*params.MessagesSendBuilder{b}, nil
	case "Отклонить заявку":
		err = state.postgres.RequestsDocuments.UpdateStatus(ctx, ds.StatusAdminDeclined, requestID)
		if err != nil {
			return "", nil, fmt.Errorf("[requests_documents.UpdateStatus]: %w", err)
		}

		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message(fmt.Sprintf("Заявка c ID - %d отклонена", requestID))
		k := object.NewMessagesKeyboard(true)
		addBackButton(k)
		b.Keyboard(k)
		return workingRequestDocument, []*params.MessagesSendBuilder{b}, nil
	case "Назад":
		status, err := state.postgres.RequestsDocuments.GetStatus(ctx, requestID)
		if status == ds.StatusAdminWorking {
			err = state.postgres.RequestsDocuments.UpdateStatus(ctx, ds.StatusUserConfirmed, requestID)
			if err != nil {
				return "", nil, fmt.Errorf("[requests_documents.UpdateStatus]: %w", err)
			}
		}
		return workingRequestDocument, nil, nil
	default:
		return workingRequestDocument, nil, nil
	}
}

func (state RequestDocumentSpecificApplicationState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	requestID, err := state.postgres.ObjectAdmin.Get(ctx, vkID)
	if err != nil {
		return nil, fmt.Errorf("[object_admin.Get]: %w", err)
	}

	output, attachment, err := state.postgres.RequestsDocuments.GetRequestByID(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("[requests_documents.GetRequestByID]: %w", err)
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

func (state RequestDocumentSpecificApplicationState) Name() stateName {
	return requestDocumentSpecificApplication
}

type EditDocumentAdminState struct {
	postgres *postrgres.Repo
}

func (state EditDocumentAdminState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Назад":
		return requestDocumentSpecificApplication, nil, nil
	case "Изменить название":
		return editNameDocumentAdmin, nil, nil
	case "Изменить ФИО автора":
		return editAuthorDocumentAdmin, nil, nil
	case "Изменить год":
		return editYearDocumentAdmin, nil, nil
	case "Изменить категорию":
		return editCategoryDocumentAdmin, nil, nil
	case "Изменить описание":
		return editDescriptionDocumentAdmin, nil, nil
	case "Изменить хэштеги":
		return editHashtagDocumentAdmin, nil, nil
	default:
		return editDocumentAdmin, nil, nil
	}
}

func (state EditDocumentAdminState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Выберите параметр для редактирования")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Изменить название", "", "secondary")
	k.AddTextButton("Изменить ФИО автора", "", "secondary")
	k.AddRow()
	k.AddTextButton("Изменить год", "", "secondary")
	k.AddTextButton("Изменить категорию", "", "secondary")
	k.AddRow()
	k.AddTextButton("Изменить описание", "", "secondary")
	k.AddTextButton("Изменить хэштеги", "", "secondary")
	addBackButton(k)
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditDocumentAdminState) Name() stateName {
	return editDocumentAdmin
}

type EditNameDocumentAdminState struct {
	postgres *postrgres.Repo
}

func (state EditNameDocumentAdminState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	reqID, err := state.postgres.ObjectAdmin.Get(ctx, msg.PeerID)
	if err != nil {
		return editNameDocumentAdmin, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return editDocumentAdmin, nil, nil
	default:
		err = state.postgres.RequestsDocuments.EditName(ctx, messageText, reqID)
		if err != nil {
			return editNameDocumentAdmin, []*params.MessagesSendBuilder{}, err
		}
		return requestDocumentSpecificApplication, nil, nil
	}
}

func (state EditNameDocumentAdminState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите новое название загружаемого документа")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditNameDocumentAdminState) Name() stateName {
	return editNameDocumentAdmin
}

type EditAuthorDocumentAdminState struct {
	postgres *postrgres.Repo
}

func (state EditAuthorDocumentAdminState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	reqID, err := state.postgres.ObjectAdmin.Get(ctx, msg.PeerID)
	if err != nil {
		return editAuthorDocumentAdmin, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return editDocumentAdmin, nil, nil
	default:
		if len(messageText) > 60 {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("ФИО автора слишком длинное, повторите ввод")
			return editAuthorDocumentAdmin, []*params.MessagesSendBuilder{b}, nil
		}
		russianRegex := regexp.MustCompile("^[а-яА-Я\\s]+$")
		if !russianRegex.MatchString(messageText) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("ФИО автора должно состоять из русских букв, повторите ввод")
			return editAuthorDocumentAdmin, []*params.MessagesSendBuilder{b}, nil
		}
		err = state.postgres.RequestsDocuments.EditAuthor(ctx, messageText, reqID)
		if err != nil {
			return editAuthorDocumentAdmin, []*params.MessagesSendBuilder{}, err
		}
		return requestDocumentSpecificApplication, nil, nil
	}
}

func (state EditAuthorDocumentAdminState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите новое ФИО автора загружаемого документа")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditAuthorDocumentAdminState) Name() stateName {
	return editAuthorDocumentAdmin
}

type EditYearDocumentAdminState struct {
	postgres *postrgres.Repo
}

func (state EditYearDocumentAdminState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	reqID, err := state.postgres.ObjectAdmin.Get(ctx, msg.PeerID)
	if err != nil {
		return editYearDocumentAdmin, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return editDocumentAdmin, nil, nil
	default:
		year, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введите год числом в формате YYYY")
			return editYearDocumentAdmin, []*params.MessagesSendBuilder{b}, nil
		}
		currentYear := time.Now().Year()
		if !(year >= 1800 && year <= currentYear) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введите существующий год в формате YYYY")
			return editYearDocumentAdmin, []*params.MessagesSendBuilder{b}, nil
		}
		err = state.postgres.RequestsDocuments.EditYear(ctx, year, reqID)
		if err != nil {
			return editYearDocumentAdmin, []*params.MessagesSendBuilder{}, err
		}
		return requestDocumentSpecificApplication, nil, nil
	}
}

func (state EditYearDocumentAdminState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите новый год создания документа")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditYearDocumentAdminState) Name() stateName {
	return editYearDocumentAdmin
}

type EditCategoryDocumentAdminState struct {
	postgres *postrgres.Repo
}

func (state EditCategoryDocumentAdminState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	reqID, err := state.postgres.ObjectAdmin.Get(ctx, msg.PeerID)
	if err != nil {
		return editCategoryDocumentAdmin, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return editDocumentAdmin, nil, nil
	case "Своя категория":
		return editUserCategoryDocumentAdmin, nil, nil
	default:
		maxID, err := state.postgres.RequestsDocuments.GetCategoryMaxID()
		if err != nil {
			return editCategoryDocumentAdmin, []*params.MessagesSendBuilder{}, err
		}
		categoryNumber, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введите номер категории числом, повторите ввод")
			return editCategoryDocumentAdmin, []*params.MessagesSendBuilder{b}, nil
		}
		if !(categoryNumber >= 1 && categoryNumber <= maxID) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Категории с таким номером нет в списке, повторите ввод")
			return editCategoryDocumentAdmin, []*params.MessagesSendBuilder{b}, nil
		}
		err = state.postgres.RequestsDocuments.EditCategory(ctx, categoryNumber, reqID)
		if err != nil {
			return editCategoryDocumentAdmin, []*params.MessagesSendBuilder{}, err
		}
		return requestDocumentSpecificApplication, nil, nil
	}
}

func (state EditCategoryDocumentAdminState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	categories, err := state.postgres.RequestsDocuments.GetCategoryNames()
	if err != nil {
		return []*params.MessagesSendBuilder{}, err
	}
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите новый номер категории документа из списка ниже:\n" + categories)
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	k.AddRow()
	k.AddTextButton("Своя категория", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditCategoryDocumentAdminState) Name() stateName {
	return editCategoryDocumentAdmin
}

type EditUserCategoryDocumentAdminState struct {
	postgres *postrgres.Repo
}

func (state EditUserCategoryDocumentAdminState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	reqID, err := state.postgres.ObjectAdmin.Get(ctx, msg.PeerID)
	if err != nil {
		return editUserCategoryDocumentAdmin, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return editCategoryDocumentAdmin, nil, nil
	default:
		err = state.postgres.RequestsDocuments.EditUserCategory(ctx, messageText, reqID)
		if err != nil {
			return editUserCategoryDocumentAdmin, []*params.MessagesSendBuilder{}, err
		}

		return requestDocumentSpecificApplication, nil, nil
	}
}

func (state EditUserCategoryDocumentAdminState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите новое название категории. Категория будет автоматически создана после принятия заявки")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditUserCategoryDocumentAdminState) Name() stateName {
	return editUserCategoryDocumentAdmin
}

type EditDescriptionDocumentAdminState struct {
	postgres *postrgres.Repo
}

func (state EditDescriptionDocumentAdminState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	reqID, err := state.postgres.ObjectAdmin.Get(ctx, msg.PeerID)
	if err != nil {
		return editDescriptionDocumentAdmin, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return editDocumentAdmin, nil, nil
	default:
		err = state.postgres.RequestsDocuments.EditDescription(ctx, messageText, reqID)
		if err != nil {
			return editDescriptionDocumentAdmin, []*params.MessagesSendBuilder{}, err
		}
		return requestDocumentSpecificApplication, nil, nil
	}
}

func (state EditDescriptionDocumentAdminState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите новое описание документа")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditDescriptionDocumentAdminState) Name() stateName {
	return editDescriptionDocumentAdmin
}

type EditHashtagDocumentAdminState struct {
	postgres *postrgres.Repo
}

func (state EditHashtagDocumentAdminState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	reqID, err := state.postgres.ObjectAdmin.Get(ctx, msg.PeerID)
	if err != nil {
		return editHashtagDocumentAdmin, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		return editDocumentAdmin, nil, nil
	default:
		hashtags := strings.Split(messageText, " ")
		err = state.postgres.RequestsDocuments.EditHashtags(ctx, hashtags, reqID)
		if err != nil {
			return editHashtagDocumentAdmin, []*params.MessagesSendBuilder{}, err
		}
		return requestDocumentSpecificApplication, nil, nil
	}
}

func (state EditHashtagDocumentAdminState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите новые названия хештегов через пробел (например, фамилия преподавателя или название предмета)")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EditHashtagDocumentAdminState) Name() stateName {
	return editHashtagDocumentAdmin
}
