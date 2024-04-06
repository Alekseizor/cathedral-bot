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

// UploadDocument загружает документ
func (r *Repo) UploadPhotoAsFile(ctx context.Context, VK *api.VK, doc object.DocsDoc, vkID int) error {
	resp, err := http.Get(doc.URL)
	if err != nil {
		log.Println(err)
	}
	upload, _ := VK.DocsGetMessagesUploadServer(api.Params{
		"type":    "doc",
		"peer_id": vkID,
	})
	file, err := io.ReadAll(resp.Body)
	fileBody := bytes.NewReader(file)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", doc.Title)
	io.Copy(part, fileBody)
	writer.Close()
	req, _ := http.NewRequest("POST", upload.UploadURL, bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	response, _ := client.Do(req)
	docs := &docsDoc{}
	json.NewDecoder(response.Body).Decode(docs)
	log.Println(docs.File)
	savedDoc, _ := VK.DocsSave(api.Params{
		"file":  docs.File,
		"title": doc.Title,
	})

	attachment := "doc" + strconv.Itoa(savedDoc.Doc.OwnerID) + "_" + strconv.Itoa(savedDoc.Doc.ID)

	_, err = r.db.ExecContext(ctx, "INSERT INTO request_photo(attachment, user_id) VALUES ($1, $2)", attachment, vkID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// UploadDocument загружает документ
func (r *Repo) UploadPhoto(ctx context.Context, VK *api.VK, doc object.PhotosPhoto, vkID int) error {
	resp, err := http.Get(doc.Sizes[9].URL)
	if err != nil {
		log.Println(err)
	}
	upload, _ := VK.DocsGetMessagesUploadServer(api.Params{
		"type":    "doc",
		"peer_id": vkID,
	})
	file, err := io.ReadAll(resp.Body)
	fileBody := bytes.NewReader(file)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", doc.Title)
	io.Copy(part, fileBody)
	writer.Close()
	req, _ := http.NewRequest("POST", upload.UploadURL, bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	response, _ := client.Do(req)
	docs := &docsDoc{}
	json.NewDecoder(response.Body).Decode(docs)
	log.Println(docs.File)
	savedDoc, _ := VK.DocsSave(api.Params{
		"file":  docs.File,
		"title": doc.Title,
	})

	attachment := "doc" + strconv.Itoa(savedDoc.Doc.OwnerID) + "_" + strconv.Itoa(savedDoc.Doc.ID)

	_, err = r.db.ExecContext(ctx, "INSERT INTO request_photo(attachment, user_id) VALUES ($1, $2)", attachment, vkID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}
