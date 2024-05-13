package state

import (
	"bytes"
	"context"
	"github.com/Alekseizor/cathedral-bot/internal/app/ds"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	"github.com/mholt/archiver"
	"io"
	"net/http"
	"strconv"
	"time"
)

var validExtensionPhotoArchive = map[string]struct{}{
	"rar": struct{}{},
}

// LoadPhotoArchiveState пользователь загружает архив
type LoadPhotoArchiveState struct {
	postgres *postrgres.Repo
	vk       *api.VK
}

func (state LoadPhotoArchiveState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "Назад" {
		return photoStart, nil, nil
	}

	attachment := msg.Attachments
	if len(attachment) == 0 {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Загрузите архив, прикрепив его к сообщению")
		return loadPhotoArchive, []*params.MessagesSendBuilder{b}, nil
	}

	if len(attachment) > 1 {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Можно загрузить только один архив")
		return loadPhotoArchive, []*params.MessagesSendBuilder{b}, nil
	}

	if attachment[0].Type == "doc" {
		if _, ok := validExtensionPhotoArchive[attachment[0].Doc.Ext]; !ok {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Данный архив недопустимого формата")
			return loadPhotoArchive, []*params.MessagesSendBuilder{b}, nil
		}

		resp, err := http.Get(attachment[0].Doc.URL)
		if err != nil {
			return loadPhotoArchive, nil, err
		}
		defer resp.Body.Close()

		archive, err := io.ReadAll(resp.Body)
		if err != nil {
			return loadPhotoArchive, nil, err
		}

		archiveBody := bytes.NewReader(archive)

		r := &archiver.Rar{}
		err = r.Open(archiveBody, int64(archiveBody.Len()))
		if err != nil {
			return loadPhotoArchive, nil, err
		}

		attachments := make([]string, 0)
		for {
			filePhoto, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				return loadPhotoArchive, nil, err
			}

			photo, err := io.ReadAll(filePhoto)
			if err != nil {
				return loadPhotoArchive, nil, err
			}

			contentType := http.DetectContentType(photo)
			if _, ok := validExtensionPhoto[contentType]; !ok {
				b := params.NewMessagesSendBuilder()
				b.RandomID(0)
				b.Message("Недопустимый тип изображения в архиве: " + contentType)
				return loadPhotoArchive, []*params.MessagesSendBuilder{b}, nil
			}

			if contentType == "application/octet-stream" {
				photo, err = convertTiffToJpg(photo)
				if err != nil {
					return loadPhotoArchive, nil, err
				}
			}

			attach, err := state.postgres.RequestPhotoArchive.GetAttachmentPhoto(ctx, state.vk, photo, msg.PeerID)
			if err != nil {
				return loadPhotoArchive, nil, err
			}

			attachments = append(attachments, attach)
		}

		err = state.postgres.RequestPhotoArchive.UploadArchivePhoto(ctx, state.vk, attachments, msg.PeerID)
		if err != nil {
			return loadPhotoArchive, nil, err
		}

		return eventYearPhotoArchive, nil, nil
	}

	return loadPhotoArchive, nil, nil
}

func (state LoadPhotoArchiveState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Загрузите архив rar. Фото в архиве должны быть одной категории. Допустимые форматы фото: jpg, jpeg, png, tiff")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state LoadPhotoArchiveState) Name() stateName {
	return loadPhotoArchive
}

// EventYearPhotoArchiveState пользователь указывает год создания фотографий в архиве
type EventYearPhotoArchiveState struct {
	postgres *postrgres.Repo
}

func (state EventYearPhotoArchiveState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhotoArchive.GetPhotoLastID(ctx, msg.PeerID)
	if err != nil {
		return eventYearPhotoArchive, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Назад":
		err = state.postgres.RequestPhotoArchive.DeletePhoto(ctx, photoID)
		if err != nil {
			return eventYearPhotoArchive, []*params.MessagesSendBuilder{}, err
		}
		return loadPhotoArchive, nil, nil
	case "Пропустить":
		return studyProgramPhotoArchive, nil, nil
	default:
		year, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введён год недопустимого формата")
			return eventYearPhotoArchive, []*params.MessagesSendBuilder{b}, nil
		}

		currentYear := time.Now().Year()
		if !(year >= 1900 && year <= currentYear) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введён несуществующий год")
			return eventYearPhotoArchive, []*params.MessagesSendBuilder{b}, nil
		}
		err = state.postgres.RequestPhotoArchive.UpdateYear(ctx, photoID, year)
		if err != nil {
			return eventYearPhotoArchive, []*params.MessagesSendBuilder{}, err
		}
		return studyProgramPhotoArchive, nil, nil
	}
}

func (state EventYearPhotoArchiveState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите год события в формате YYYY")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Пропустить", "", "secondary")
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EventYearPhotoArchiveState) Name() stateName {
	return eventYearPhotoArchive
}

// StudyProgramPhotoArchiveState пользователь указывает программу обучения
type StudyProgramPhotoArchiveState struct {
	postgres *postrgres.Repo
}

func (state StudyProgramPhotoArchiveState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhotoArchive.GetPhotoLastID(ctx, msg.PeerID)
	if err != nil {
		return studyProgramPhotoArchive, []*params.MessagesSendBuilder{}, err
	}

	var educationProgram string

	switch messageText {
	case "Бакалавриат":
		educationProgram = "Бакалавриат"
	case "Магистратура":
		educationProgram = "Магистратура"
	case "Специалитет":
		educationProgram = "Специалитет"
	case "Аспирантура":
		educationProgram = "Аспирантура"
	case "Пропустить":
		return eventNamePhotoArchive, nil, nil
	case "Назад":
		return eventYearPhotoArchive, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Такой программы обучения нет в предложенных вариантах")
		return studyProgramPhotoArchive, []*params.MessagesSendBuilder{b}, nil
	}

	err = state.postgres.RequestPhotoArchive.UpdateStudyProgram(ctx, photoID, educationProgram)
	if err != nil {
		return studyProgramPhotoArchive, []*params.MessagesSendBuilder{}, err
	}
	return eventNamePhotoArchive, nil, nil
}

func (state StudyProgramPhotoArchiveState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Выберите программу обучения")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Бакалавриат", "", "secondary")
	k.AddTextButton("Магистратура", "", "secondary")
	k.AddRow()
	k.AddTextButton("Специалитет", "", "secondary")
	k.AddTextButton("Аспирантура", "", "secondary")
	k.AddRow()
	k.AddTextButton("Пропустить", "", "secondary")
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state StudyProgramPhotoArchiveState) Name() stateName {
	return studyProgramPhotoArchive
}

// EventNamePhotoArchiveState пользователь указывает существующее название события
type EventNamePhotoArchiveState struct {
	postgres *postrgres.Repo
}

func (state EventNamePhotoArchiveState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	photoID, err := state.postgres.RequestPhotoArchive.GetPhotoLastID(ctx, msg.PeerID)
	if err != nil {
		return eventNamePhotoArchive, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Своё событие":
		return userEventNamePhotoArchive, nil, nil
	case "Пропустить":
		return descriptionPhotoArchive, nil, nil
	case "Назад":
		return studyProgramPhotoArchive, nil, nil
	default:
		maxID, err := state.postgres.RequestPhotoArchive.GetEventMaxID()
		if err != nil {
			return eventNamePhotoArchive, []*params.MessagesSendBuilder{}, err
		}

		eventNumber, err := strconv.Atoi(messageText)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Введена не цифра")
			return eventNamePhotoArchive, []*params.MessagesSendBuilder{b}, nil
		}

		if !(eventNumber >= 1 && eventNumber <= maxID) {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Такого события нет в списке")
			return eventNamePhotoArchive, []*params.MessagesSendBuilder{b}, nil
		}

		err = state.postgres.RequestPhotoArchive.UpdateEvent(ctx, photoID, eventNumber)
		if err != nil {
			return eventNamePhotoArchive, []*params.MessagesSendBuilder{}, err
		}
		return descriptionPhotoArchive, nil, nil
	}
}

func (state EventNamePhotoArchiveState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	categories, err := state.postgres.RequestPhotoArchive.GetEventNames()
	if err != nil {
		return []*params.MessagesSendBuilder{}, err
	}
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите номер события из списка ниже:\n" + categories)
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Своё событие", "", "secondary")
	k.AddRow()
	k.AddTextButton("Пропустить", "", "secondary")
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state EventNamePhotoArchiveState) Name() stateName {
	return eventNamePhotoArchive
}

// UserEventNamePhotoArchiveState пользователь указывает своё название события
type UserEventNamePhotoArchiveState struct {
	postgres *postrgres.Repo
}

func (state UserEventNamePhotoArchiveState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return userEventNamePhotoArchive, nil, nil
	}

	photoID, err := state.postgres.RequestPhotoArchive.GetPhotoLastID(ctx, msg.PeerID)
	if err != nil {
		return userEventNamePhotoArchive, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Пропустить":
		return descriptionPhotoArchive, nil, nil
	case "Назад":
		return eventNamePhotoArchive, nil, nil
	default:
		err := state.postgres.RequestPhotoArchive.UpdateUserEvent(ctx, photoID, messageText)
		if err != nil {
			return userEventNamePhotoArchive, []*params.MessagesSendBuilder{}, err
		}
		return descriptionPhotoArchive, nil, nil
	}
}

func (state UserEventNamePhotoArchiveState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите название своего события")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Пропустить", "", "secondary")
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state UserEventNamePhotoArchiveState) Name() stateName {
	return userEventNamePhotoArchive
}

// DescriptionPhotoArchiveState пользователь вводит описание фотографии
type DescriptionPhotoArchiveState struct {
	postgres *postrgres.Repo
}

func (state DescriptionPhotoArchiveState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return descriptionPhotoArchive, nil, nil
	}

	photoID, err := state.postgres.RequestPhotoArchive.GetPhotoLastID(ctx, msg.PeerID)
	if err != nil {
		return descriptionPhotoArchive, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Пропустить":
		return checkPhotoArchive, nil, nil
	case "Назад":
		return eventNamePhotoArchive, nil, nil
	default:
		err := state.postgres.RequestPhotoArchive.UpdateDescription(ctx, photoID, messageText)
		if err != nil {
			return descriptionPhotoArchive, []*params.MessagesSendBuilder{}, err
		}
		return checkPhotoArchive, nil, nil
	}
}

func (state DescriptionPhotoArchiveState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите общее описание фотографий в архиве")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Пропустить", "", "secondary")
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state DescriptionPhotoArchiveState) Name() stateName {
	return descriptionPhotoArchive
}

// CheckPhotoArchiveState пользователь проверяет заявку на загрузку фотографии
type CheckPhotoArchiveState struct {
	postgres *postrgres.Repo
}

func (state CheckPhotoArchiveState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return checkPhotoArchive, nil, nil
	}

	photoID, err := state.postgres.RequestPhotoArchive.GetPhotoLastID(ctx, msg.PeerID)
	if err != nil {
		return checkPhotoArchive, []*params.MessagesSendBuilder{}, err
	}

	switch messageText {
	case "Отправить":
		err = state.postgres.RequestPhotoArchive.UpdateStatus(ctx, ds.StatusUserConfirmed, photoID)
		if err != nil {
			return checkPhotoArchive, []*params.MessagesSendBuilder{}, err
		}

		err = state.postgres.RequestPhotoArchive.ChangeArchiveToPhotos(ctx, photoID)
		if err != nil {
			return checkPhotoArchive, []*params.MessagesSendBuilder{}, err
		}

		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Фотография отправлена администратору на рассмотрение. Вы можете отслеживать статус своей заявки в личном кабинете")
		return photoStart, []*params.MessagesSendBuilder{b}, nil
	case "Редактировать заявку":
		return editPhotoArchive, nil, nil
	case "Назад":
		return descriptionPhotoArchive, nil, nil
	default:
		return checkPhotoArchive, nil, nil
	}
}

func (state CheckPhotoArchiveState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	photoID, err := state.postgres.RequestPhotoArchive.GetPhotoLastID(ctx, vkID)
	if err != nil {
		return []*params.MessagesSendBuilder{}, err
	}

	output, err := state.postgres.RequestPhotoArchive.CheckParams(ctx, photoID)
	if err != nil {
		return []*params.MessagesSendBuilder{}, err
	}

	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Проверьте правильность введённых параметров:\n" + output)
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Редактировать заявку", "", "secondary")
	k.AddRow()
	k.AddTextButton("Отправить", "", "secondary")
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state CheckPhotoArchiveState) Name() stateName {
	return checkPhotoArchive
}
