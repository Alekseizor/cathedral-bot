package postrgres

import (
	"fmt"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres/document"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres/request_photo"

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
	repoCfg           config.PostgresConfig
	State             *state.Repo
	Document          *document.Repo
	RequestsDocuments *requests_documents.Repo
	RequestPhoto      *request_photo.Repo
	Admin             *admin.Repo
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
	r.Document = document.New(db)
	r.RequestsDocuments = requests_documents.New(db)
	r.Admin = admin.New(db)

	return nil
}
