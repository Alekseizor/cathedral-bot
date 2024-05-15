package create_album

import (
	"context"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/jmoiron/sqlx"
	"strings"
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

// CreateStudentAlbum создает альбом студентов
func (r *Repo) CreateStudentAlbum(ctx context.Context, vk *api.VK, groupID int, nameAlbum string) (string, bool, error) {
	name := strings.Split(nameAlbum, "//")
	year := name[0]
	studyProgram := name[1]
	event := name[2]
	resultName := fmt.Sprintf("%s // %s // %s", year, studyProgram, event)

	var flag bool
	err := r.db.Get(&flag, "SELECT EXISTS(SELECT FROM student_albums WHERE year = $1 AND study_program = $2 AND event = $3)", year, studyProgram, event)
	if err != nil {
		return "", false, fmt.Errorf("[db.Get]: %w", err)
	}
	if flag {
		return "", true, nil
	}

	album, err := vk.PhotosCreateAlbum(api.Params{
		"title":    resultName,
		"group_id": groupID,
	})

	albumURL := fmt.Sprintf("https://vk.com/album-%d_%d", groupID, album.ID)

	_, err = r.db.ExecContext(ctx, "INSERT INTO student_albums (year, study_program, event, url, vk_id) VALUES ($1, $2, $3, $4, $5)", year, studyProgram, event, albumURL, album.ID)
	if err != nil {
		return "", false, fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return albumURL, false, nil
}

// CreateTeacherAlbum создает альбом преподавателя
func (r *Repo) CreateTeacherAlbum(ctx context.Context, vk *api.VK, groupID int, nameAlbum string) (string, bool, error) {
	var flag bool
	err := r.db.Get(&flag, "SELECT EXISTS(SELECT FROM teacher_albums WHERE name = $1)", nameAlbum)
	if err != nil {
		return "", false, fmt.Errorf("[db.Get]: %w", err)
	}
	if flag {
		return "", true, nil
	}

	album, err := vk.PhotosCreateAlbum(api.Params{
		"title":    nameAlbum,
		"group_id": groupID,
	})

	albumURL := fmt.Sprintf("https://vk.com/album-%d_%d", groupID, album.ID)

	_, err = r.db.ExecContext(ctx, "INSERT INTO teacher_albums (name, url, vk_id) VALUES ($1, $2, $3)", nameAlbum, albumURL, album.ID)
	if err != nil {
		return "", false, fmt.Errorf("[db.ExecContext]: %w", err)
	}

	fmt.Println(album)
	return albumURL, false, nil
}
