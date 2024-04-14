package state

import (
	"bytes"
	"context"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	"github.com/mholt/archiver"
	"io"
	"net/http"
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

			// Проверка типа изображения
			contentType := http.DetectContentType(photo)
			if contentType != "image/jpeg" && contentType != "image/png" && contentType != "application/octet-stream" {
				b := params.NewMessagesSendBuilder()
				b.RandomID(0)
				b.Message("Недопустимый тип изображения в архиве")
				return loadPhotoArchive, []*params.MessagesSendBuilder{b}, nil
			}

			if contentType == "application/octet-stream" {
				photo, err = convertTiffToJpg(photo)
				if err != nil {
					return loadPhotoArchive, nil, err
				}
			}

			attach, err := state.postgres.RequestArchivePhoto.GetAttachmentPhoto(ctx, state.vk, photo, msg.PeerID)
			if err != nil {
				return loadPhotoArchive, nil, err
			}

			attachments = append(attachments, attach)
		}

		err = state.postgres.RequestArchivePhoto.UploadArchivePhoto(ctx, state.vk, attachments, msg.PeerID)
		if err != nil {
			return loadPhotoArchive, nil, err
		}
		return loadPhotoArchive, nil, nil
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
