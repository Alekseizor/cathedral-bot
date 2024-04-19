package postrgres

import (
	"fmt"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres/documents"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres/object_admin"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres/request_photo"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres/request_photo_archive"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres/search_album"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres/search_document"
	"github.com/jmoiron/sqlx"

	"github.com/Alekseizor/cathedral-bot/internal/app/config"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres/admin"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres/requests_documents"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres/state"
)

const (
	postgres = "postgres"
	sslMode  = "disable"
)

// Repo инстанс репо для работы с postgres
type Repo struct {
	repoCfg             config.PostgresConfig
	State               *state.Repo
	Document            *documents.Repo
	RequestsDocuments   *requests_documents.Repo
	RequestPhoto        *request_photo.Repo
	RequestPhotoArchive *request_photo_archive.Repo
	SearchAlbum         *search_album.Repo
	Admin               *admin.Repo
	Documents           *documents.Repo
	SearchDocument      *search_document.Repo
	ObjectAdmin         *object_admin.Repo
}

// New - создаем новое объект репо, подключения к бд еще нет!
func New(postgreCfg config.PostgresConfig) *Repo {
	return &Repo{
		repoCfg: postgreCfg,
	}
}

// Init - образует коннект к базе данных
func (r *Repo) Init() error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", r.repoCfg.Host, r.repoCfg.Port, r.repoCfg.User, r.repoCfg.Password, r.repoCfg.Name, sslMode)

	db, err := sqlx.Connect(postgres, dsn)
	if err != nil {
		return fmt.Errorf("[sqlx.Connect]: %w", err)
	}

	r.State = state.New(db)
	r.RequestPhoto = request_photo.New(db)
	r.RequestPhotoArchive = request_photo_archive.New(db)
	r.SearchAlbum = search_album.New(db)
	r.Document = documents.New(db)
	r.RequestsDocuments = requests_documents.New(db)
	r.Admin = admin.New(db)
	r.Documents = documents.New(db)
	r.ObjectAdmin = object_admin.New(db)
	r.SearchDocument = search_document.New(db)

	return nil
}
