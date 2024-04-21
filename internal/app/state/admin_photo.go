package state

import (
	"context"
	"fmt"
	"strconv"

	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"

	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
)

type AlbumsCabinetState struct {
	postgres *postrgres.Repo
}

func (state AlbumsCabinetState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Работа с заявкой":
		return workingRequestPhoto, nil, nil
	case "Работа с альбомом":
		return workingAlbums, nil, nil
	case "Заблокировать пользователя":
		return blockUser, nil, nil
	case "Добавить администратора":
		return addPhotoAdministrator, nil, nil
	case "Удалить администратора":
		return removePhotoAdministrator, nil, nil
	case "Выйти из кабинета администратора":
		return selectArchive, nil, nil
	default:
		return albumsCabinet, nil, nil
	}
}

func (state AlbumsCabinetState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Вы в кабинете администратора фотоархива, выберите действие")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Работа с заявкой", "", "secondary")
	k.AddTextButton("Работа с альбомом", "", "secondary")
	k.AddRow()
	k.AddTextButton("Заблокировать пользователя", "", "secondary")
	k.AddRow()
	k.AddTextButton("Добавить администратора", "", "secondary")
	k.AddTextButton("Удалить администратора", "", "secondary")
	k.AddRow()
	k.AddTextButton("Выйти из кабинета администратора", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state AlbumsCabinetState) Name() stateName {
	return albumsCabinet
}

///////////

type WorkingRequestPhotoState struct {
}

func (state WorkingRequestPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Заявки из очереди":
		return requestPhotoFromQueue, nil, nil
	case "Конкретная заявка":
		return requestPhotoSpecificApplication, nil, nil
	case "Назад":
		return albumsCabinet, nil, nil
	default:
		return workingRequestPhoto, nil, nil
	}
}

func (state WorkingRequestPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Вы хотите работать с заявками из очереди или с конкретной заявкой?")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Заявки из очереди", "", "secondary")
	k.AddTextButton("Конкретная заявка", "", "secondary")
	addBackButton(k)
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state WorkingRequestPhotoState) Name() stateName {
	return workingRequestPhoto
}

///////////

type WorkingAlbumsState struct {
	postgres *postrgres.Repo
}

func (state WorkingAlbumsState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	if messageText == "Назад" {
		return albumsCabinet, nil, nil
	}

	documentID, err := strconv.Atoi(messageText)
	if err != nil {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("ID должно быть числом, например, 12")
		return workingAlbums, []*params.MessagesSendBuilder{b}, nil
	}

	exists, err := state.postgres.Documents.CheckExistence(ctx, documentID)
	if err != nil {
		return workingAlbums, nil, err
	}
	if !exists {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("ID c таким фото не найдено")
		return workingAlbums, []*params.MessagesSendBuilder{b}, nil
	}

	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	output, attachment, err := state.postgres.Documents.GetOutput(ctx, documentID)
	b.Attachment(attachment)
	b.Message(output)
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Изменить", documentID, "secondary")
	k.AddTextButton("Удалить", documentID, "secondary")
	addBackButton(k)
	b.Keyboard(k)

	return workingAlbums, []*params.MessagesSendBuilder{b}, nil
}

func (state WorkingAlbumsState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	//TODO: показываем существующие альбомы как id и название
	b.Message("Показываем существующие альбомы как id и название")

	b1 := params.NewMessagesSendBuilder()
	b1.RandomID(0)
	b1.Message("Введите ID альбома, над которым хотите поработать. Например: 12")
	k := object.NewMessagesKeyboard(true)
	addBackButton(k)
	b1.Keyboard(k)
	return []*params.MessagesSendBuilder{b, b1}, nil
}

func (state WorkingAlbumsState) Name() stateName {
	return workingAlbums
}

///////////

type ActionOnPhotoState struct {
	postgres *postrgres.Repo
}

func (state ActionOnPhotoState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "Назад" {
		return workingAlbums, nil, nil
	}

	payload := msg.Payload

	albumID, err := strconv.Atoi(payload)
	if err != nil {
		return "", nil, err
	}

	switch messageText {
	case "Удалить":
		err = state.postgres.Documents.Delete(ctx, albumID)
		if err != nil {
			return "", nil, err
		}
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Документ успешно удален")
		return workingAlbums, []*params.MessagesSendBuilder{b}, nil
	case "Изменить":
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Что Вы хотите изменить в документе?")
		k := object.NewMessagesKeyboard(true)
		k.AddRow()
		k.AddTextButton("Название", albumID, "secondary")
		k.AddTextButton("Описание", albumID, "secondary")
		k.AddRow()
		k.AddTextButton("Автор", albumID, "secondary")
		k.AddTextButton("Год", albumID, "secondary")
		k.AddRow()
		k.AddTextButton("Категория", albumID, "secondary")
		k.AddTextButton("Хештеги", albumID, "secondary")
		addBackButton(k)
		b.Keyboard(k)
		return changeAlbums, []*params.MessagesSendBuilder{b}, nil
	default:
		return workingAlbums, nil, nil
	}
}

func (state ActionOnPhotoState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	return nil, nil
}

func (state ActionOnPhotoState) Name() stateName {
	return actionOnPhoto
}

///////////

type ChangeAlbumsState struct {
	postgres *postrgres.Repo
}

func (state ChangeAlbumsState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "Назад" {
		return workingAlbums, nil, nil
	}

	payload := msg.Payload

	documentID, err := strconv.Atoi(payload)
	if err != nil {
		return "", nil, err
	}

	err = state.postgres.ObjectAdmin.Update(ctx, documentID, msg.PeerID)
	if err != nil {
		return "", nil, fmt.Errorf("[object_admin.Update]: %w", err)
	}

	switch messageText {
	case "Название":
		return changeTitleDocument, nil, nil
	case "Описание":
		return changeDescriptionDocument, nil, nil
	case "Автор":
		return changeAuthorDocument, nil, nil
	case "Год":
		return changeYearDocument, nil, nil
	case "Категория":
		return changeCategoryDocument, nil, nil
	case "Хештеги":
		return changeHashtagsDocument, nil, nil
	default:
		return workingAlbums, nil, nil
	}
}

func (state ChangeAlbumsState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	return nil, nil
}

func (state ChangeAlbumsState) Name() stateName {
	return changeAlbums
}
