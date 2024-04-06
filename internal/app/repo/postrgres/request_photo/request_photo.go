package request_photo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/object"
	"github.com/jmoiron/sqlx"
	"io"
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

type docsDoc struct {
	File string `json:"file"`
}

// UploadPhotoAsFile загружает фотографию как документ
func (r *Repo) UploadPhotoAsFile(ctx context.Context, VK *api.VK, doc object.DocsDoc, vkID int) error {
	resp, err := http.Get(doc.URL)
	if err != nil {
		return fmt.Errorf("failed to retrieve document URL: %w", err)
	}
	defer resp.Body.Close()

	upload, err := VK.DocsGetMessagesUploadServer(api.Params{
		"type":    "doc",
		"peer_id": vkID,
	})
	if err != nil {
		return fmt.Errorf("failed to get upload server URL: %w", err)
	}

	file, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read document content: %w", err)
	}
	fileBody := bytes.NewReader(file)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", doc.Title)
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}
	_, err = io.Copy(part, fileBody)
	if err != nil {
		return fmt.Errorf("failed to copy file content to form file: %w", err)
	}
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	req, err := http.NewRequest("POST", upload.UploadURL, body)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}
	defer response.Body.Close()

	docs := &docsDoc{}
	err = json.NewDecoder(response.Body).Decode(docs)
	if err != nil {
		return fmt.Errorf("failed to decode response body: %w", err)
	}

	log.Println(docs.File)

	savedDoc, err := VK.DocsSave(api.Params{
		"file":  docs.File,
		"title": doc.Title,
	})
	if err != nil {
		return fmt.Errorf("failed to save document: %w", err)
	}

	attachment := "doc" + strconv.Itoa(savedDoc.Doc.OwnerID) + "_" + strconv.Itoa(savedDoc.Doc.ID)

	_, err = r.db.ExecContext(ctx, "INSERT INTO request_photo(attachment, user_id) VALUES ($1, $2)", attachment, vkID)
	if err != nil {
		return fmt.Errorf("failed to insert into database: %w", err)
	}

	return nil
}

type photosPhoto struct {
	Server int    `json:"server"`
	Photo  string `json:"photo"`
	Hash   string `json:"hash"`
}

// UploadPhoto загружает фотографию
func (r *Repo) UploadPhoto(ctx context.Context, VK *api.VK, photo object.PhotosPhoto, vkID int) error {
	resp, err := http.Get(photo.Sizes[4].URL)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	uploadServer, err := VK.PhotosGetMessagesUploadServer(api.Params{
		"peer_id": vkID,
	})
	if err != nil {
		log.Println(err)
		return err
	}

	file, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	fileBody := bytes.NewReader(file)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("photo", "photo.jpg")
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = io.Copy(part, fileBody)
	if err != nil {
		log.Println(err)
		return err
	}
	writer.Close()

	req, err := http.NewRequest("POST", uploadServer.UploadURL, bytes.NewReader(body.Bytes()))
	if err != nil {
		log.Println(err)
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(req)

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
		return err
	}

	attachment := "photo" + strconv.Itoa(savedPhoto[0].OwnerID) + "_" + strconv.Itoa(savedPhoto[0].ID)

	_, err = r.db.ExecContext(ctx, "INSERT INTO request_photo(attachment, user_id) VALUES ($1, $2)", attachment, vkID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}
