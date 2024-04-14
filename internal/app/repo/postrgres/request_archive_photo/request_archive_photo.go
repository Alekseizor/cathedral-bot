package request_archive_photo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
)

// Repo инстанс репо для работы с фотографиями пользователя
type Repo struct {
	db *sqlx.DB
}

// New - создаем новое объект репо для работы с фотографиями пользователя
func New(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

type photosPhoto struct {
	Server int    `json:"server"`
	Photo  string `json:"photo"`
	Hash   string `json:"hash"`
}

func (r *Repo) GetAttachmentPhoto(ctx context.Context, VK *api.VK, photoData []byte, vkID int) (string, error) {
	uploadServer, err := VK.PhotosGetMessagesUploadServer(api.Params{
		"peer_id": vkID,
	})
	if err != nil {
		log.Println(err)
		return "", err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("photo", "photo.jpg")
	if err != nil {
		log.Println(err)
		return "", err
	}
	_, err = part.Write(photoData)
	if err != nil {
		log.Println(err)
		return "", err
	}
	writer.Close()

	req, err := http.NewRequest("POST", uploadServer.UploadURL, body)
	if err != nil {
		log.Println(err)
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer response.Body.Close()

	uploadResult := &photosPhoto{}
	err = json.NewDecoder(response.Body).Decode(uploadResult)
	if err != nil {
		log.Fatal(err)
	}

	savedPhoto, err := VK.PhotosSaveMessagesPhoto(api.Params{
		"photo":  uploadResult.Photo,
		"server": uploadResult.Server,
		"hash":   uploadResult.Hash,
	})
	if err != nil {
		log.Println(err)
		return "", err
	}

	attachment := "photo" + strconv.Itoa(savedPhoto[0].OwnerID) + "_" + strconv.Itoa(savedPhoto[0].ID)

	return attachment, nil
}

func (r *Repo) UploadArchivePhoto(ctx context.Context, VK *api.VK, attachments []string, vkID int) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO request_photo(attachments, user_id) VALUES ($1, $2)", pq.Array(attachments), vkID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}
	return nil
}
