package add_photo_to_album

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/jmoiron/sqlx"
	"io"
	"mime/multipart"
	"net/http"
)

// Repo инстанс репо для работы с параметрами для поиска альбома
type Repo struct {
	db *sqlx.DB
}

// New - создаем новое объект репо для работы с параметрами для поиска альбома
func New(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

type UploadResponse struct {
	Server     int    `json:"server"`
	PhotosList string `json:"photos_list"`
	AID        int    `json:"aid"`
	Hash       string `json:"hash"`
}

func (r *Repo) AddPhotoToAlbum(ctx context.Context, vk *api.VK, albumID int, groupID int, photoURL string) error {
	// Получаем URL сервера для загрузки фотографии
	uploadServer, err := vk.PhotosGetUploadServer(api.Params{"album_id": albumID, "group_id": groupID})
	if err != nil {
		return fmt.Errorf("ошибка при получении сервера загрузки: %w", err)
	}

	// Получаем содержимое фотографии по URL
	resp, err := http.Get(photoURL)
	if err != nil {
		return fmt.Errorf("ошибка при загрузке фотографии: %w", err)
	}
	defer resp.Body.Close()

	// Создаем тело запроса
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Добавляем файл в тело запроса
	part, err := writer.CreateFormFile("file1", "photo.jpg")
	if err != nil {
		return fmt.Errorf("ошибка при создании формы для файла: %w", err)
	}
	_, err = io.Copy(part, resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка при записи фотографии в форму: %w", err)
	}
	writer.Close()

	// Отправляем POST-запрос на сервер загрузки
	resp, err = http.Post(uploadServer.UploadURL, writer.FormDataContentType(), body)
	if err != nil {
		return fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	// Декодируем ответ
	var uploadResponse UploadResponse
	err = json.NewDecoder(resp.Body).Decode(&uploadResponse)
	if err != nil {
		return fmt.Errorf("ошибка при декодировании ответа: %w", err)
	}

	// Сохраняем фотографию
	_, err = vk.PhotosSave(api.Params{
		"album_id":    albumID,
		"group_id":    groupID,
		"server":      uploadResponse.Server,
		"photos_list": uploadResponse.PhotosList,
		"hash":        uploadResponse.Hash,
	})
	if err != nil {
		return fmt.Errorf("ошибка при сохранении фотографии: %w", err)
	}

	return nil
}
