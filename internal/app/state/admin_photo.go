package state

import (
	"context"
	"fmt"
	"strconv"

	"github.com/SevereCloud/vksdk/v2/api"
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
	postgres *postrgres.Repo
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
}

func (state WorkingAlbumsState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	switch messageText {
	case "Студенческие":
		return workingAlbumsFromStudents, nil, nil
	case "Преподавательские":
		return workingAlbumsFromTeachers, nil, nil
	case "Назад":
		return albumsCabinet, nil, nil
	default:
		return workingAlbums, nil, nil
	}
}

func (state WorkingAlbumsState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Вы хотите работать с альбомами студенческими или преподавательскими?")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Студенческие", "", "secondary")
	k.AddTextButton("Преподавательские", "", "secondary")
	addBackButton(k)
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state WorkingAlbumsState) Name() stateName {
	return workingAlbums
}

///////////

type WorkingAlbumsFromStudentsState struct {
	postgres *postrgres.Repo
}

func (state WorkingAlbumsFromStudentsState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	if messageText == "Назад" {
		return workingAlbums, nil, nil
	}

	documentID, err := strconv.Atoi(messageText)
	if err != nil {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("ID должно быть числом, например, 12")
		return workingAlbumsFromStudents, []*params.MessagesSendBuilder{b}, nil
	}

	exists, err := state.postgres.StudentAlbums.CheckExistence(ctx, documentID)
	if err != nil {
		return workingAlbumsFromStudents, nil, err
	}

	if !exists {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Альбома студентов с таким ID не найдено")
		return workingAlbumsFromStudents, []*params.MessagesSendBuilder{b}, nil
	}

	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	output, err := state.postgres.StudentAlbums.GetAlbum(ctx, documentID)
	if err != nil {
		return workingAlbumsFromStudents, nil, err
	}
	b.Message(output)
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Изменить", documentID, "secondary")
	k.AddTextButton("Удалить", documentID, "secondary")
	addBackButton(k)
	b.Keyboard(k)

	return actionOnPhoto, []*params.MessagesSendBuilder{b}, nil
}

func (state WorkingAlbumsFromStudentsState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	output, err := state.postgres.StudentAlbums.GetAllAlbumsOutput(ctx)
	if err != nil {
		return nil, fmt.Errorf("[student_albums.GetAllAlbumsOutput]: %w", err)
	}
	messages := make([]*params.MessagesSendBuilder, 0)

	if len(output) > 0 {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message(fmt.Sprintf("Существующие студенческие альбомы:"))

		messages = append(messages, b)
	}

	for _, outputElem := range output {
		messages = append(messages, params.NewMessagesSendBuilder().RandomID(0).Message(outputElem))
	}

	b1 := params.NewMessagesSendBuilder()
	b1.RandomID(0)
	b1.Message("Введите ID альбома, над которым хотите поработать. Например: 12")
	k := object.NewMessagesKeyboard(true)
	addBackButton(k)
	b1.Keyboard(k)

	return append(messages, b1), nil
}

func (state WorkingAlbumsFromStudentsState) Name() stateName {
	return workingAlbumsFromStudents
}

///////////

type ActionOnPhotoState struct {
	postgres *postrgres.Repo
	vk       *api.VK
	vkUser   *api.VK
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
		albumVKID, err := state.postgres.StudentAlbums.GetVKID(ctx, albumID)
		if err != nil {
			return "", nil, fmt.Errorf("[StudentAlbums.GetVKID]: %w", err)
		}

		// достали id группы
		groupID, err := state.vk.GroupsGetByID(nil)
		if err != nil {
			return "", nil, fmt.Errorf("[vk.GroupsGetByID]: %w", err)
		}

		resp, err := state.vkUser.PhotosDeleteAlbum(api.Params{"album_id": albumVKID, "group_id": groupID[0].ID})
		if err != nil || resp != 1 {
			return "", nil, fmt.Errorf("[vkUser.PhotosDeleteAlbum]: %w", err)
		}

		err = state.postgres.StudentAlbums.Delete(ctx, albumID)
		if err != nil {
			return "", nil, fmt.Errorf("[TeacherAlbums.Delete]: %w", err)
		}

		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message(fmt.Sprintf("Студенческий альбом №%d успешно удален", albumID))
		return workingAlbums, []*params.MessagesSendBuilder{b}, nil
	case "Изменить":
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Что Вы хотите изменить в альбоме?")
		k := object.NewMessagesKeyboard(true)
		k.AddRow()
		k.AddTextButton("Год события", albumID, "secondary")
		k.AddTextButton("Программа обучения", albumID, "secondary")
		k.AddRow()
		k.AddTextButton("Название события", albumID, "secondary")
		k.AddTextButton("Описание", albumID, "secondary")
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
	case "Год события":
		return changeEventYearPhoto, nil, nil
	case "Программа обучения":
		return changeStudyProgramPhoto, nil, nil
	case "Название события":
		return changeEventNamePhoto, nil, nil
	case "Описание":
		return changeDescriptionPhoto, nil, nil
	case "Назад":
		return workingAlbums, nil, nil
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

/////

///////////

type WorkingAlbumsFromTeacherState struct {
	postgres *postrgres.Repo
}

func (state WorkingAlbumsFromTeacherState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text

	if messageText == "Назад" {
		return workingAlbums, nil, nil
	}

	documentID, err := strconv.Atoi(messageText)
	if err != nil {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("ID должно быть числом, например, 12")
		return workingAlbumsFromTeachers, []*params.MessagesSendBuilder{b}, nil
	}

	exists, err := state.postgres.TeacherAlbums.CheckExistence(ctx, documentID)
	if err != nil {
		return workingAlbumsFromTeachers, nil, err
	}

	if !exists {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Альбома студентов с таким ID не найдено")
		return workingAlbumsFromTeachers, []*params.MessagesSendBuilder{b}, nil
	}

	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	output, err := state.postgres.TeacherAlbums.GetAlbum(ctx, documentID)
	if err != nil {
		return workingAlbumsFromTeachers, nil, err
	}
	b.Message(output)
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Изменить", documentID, "secondary")
	k.AddTextButton("Удалить", documentID, "secondary")
	addBackButton(k)
	b.Keyboard(k)

	return actionOnPhotoTeacher, []*params.MessagesSendBuilder{b}, nil
}

func (state WorkingAlbumsFromTeacherState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	output, err := state.postgres.TeacherAlbums.GetAllAlbumsOutput(ctx)
	if err != nil {
		return nil, fmt.Errorf("[teacher_albums.GetAllAlbumsOutput]: %w", err)
	}

	messages := make([]*params.MessagesSendBuilder, 0)

	if len(output) > 0 {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message(fmt.Sprintf("Существующие альбомы преподавателей:"))

		messages = append(messages, b)
	}

	for _, outputElem := range output {
		messages = append(messages, params.NewMessagesSendBuilder().RandomID(0).Message(outputElem))
	}

	b1 := params.NewMessagesSendBuilder()
	b1.RandomID(0)
	b1.Message("Введите ID альбома, над которым хотите поработать. Например: 12")
	k := object.NewMessagesKeyboard(true)
	addBackButton(k)
	b1.Keyboard(k)
	return append(messages, b1), nil
}

func (state WorkingAlbumsFromTeacherState) Name() stateName {
	return workingAlbumsFromTeachers
}

///////////

type ActionOnPhotoTeacherState struct {
	postgres *postrgres.Repo
	vk       *api.VK
	vkUser   *api.VK
}

func (state ActionOnPhotoTeacherState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "Назад" {
		return workingAlbums, nil, nil
	}

	payload := msg.Payload

	albumID, err := strconv.Atoi(payload)
	if err != nil {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Введите одну из команд, которые предлагает бот")
		return workingAlbums, []*params.MessagesSendBuilder{b}, nil
	}

	switch messageText {
	case "Удалить":
		albumVKID, err := state.postgres.TeacherAlbums.GetVKID(ctx, albumID)
		if err != nil {
			return "", nil, err
		}

		// достали id группы
		groupID, err := state.vk.GroupsGetByID(nil)
		if err != nil {
			return "", nil, fmt.Errorf("[vk.GroupsGetByID]: %w", err)
		}

		resp, err := state.vkUser.PhotosDeleteAlbum(api.Params{"album_id": albumVKID, "group_id": groupID[0].ID})
		if err != nil || resp != 1 {
			return "", nil, fmt.Errorf("[vkUser.PhotosDeleteAlbum]: %w", err)
		}

		err = state.postgres.TeacherAlbums.Delete(ctx, albumID)
		if err != nil {
			return "", nil, fmt.Errorf("[TeacherAlbums.Delete]: %w", err)
		}

		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message(fmt.Sprintf("Альбом учителя №%d успешно удален", albumID))
		return workingAlbums, []*params.MessagesSendBuilder{b}, nil
	case "Изменить":
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Что Вы хотите изменить в альбоме?")
		k := object.NewMessagesKeyboard(true)
		k.AddRow()
		k.AddTextButton("ФИО", albumID, "secondary")
		k.AddTextButton("Описание", albumID, "secondary")
		addBackButton(k)
		b.Keyboard(k)
		return changeAlbumsTeacher, []*params.MessagesSendBuilder{b}, nil
	default:
		return workingAlbums, nil, nil
	}
}

func (state ActionOnPhotoTeacherState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	return nil, nil
}

func (state ActionOnPhotoTeacherState) Name() stateName {
	return actionOnPhotoTeacher
}

///////////

type ChangeAlbumsTeacherState struct {
	postgres *postrgres.Repo
}

func (state ChangeAlbumsTeacherState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
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
	case "ФИО":
		return changeNameTeacherPhoto, nil, nil
	case "Описание":
		return changeDescriptionPhotoTeacher, nil, nil
	case "Назад":
		return workingAlbums, nil, nil
	default:
		return workingAlbums, nil, nil
	}
}

func (state ChangeAlbumsTeacherState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	return nil, nil
}

func (state ChangeAlbumsTeacherState) Name() stateName {
	return changeAlbumsTeacher
}

// RequestPhotoFromQueueState админ смотрит заявки на добавление фото в альбом
type RequestPhotoFromQueueState struct {
	postgres *postrgres.Repo
	vk       *api.VK
	vkUser   *api.VK
	groupID  int
}

func (state RequestPhotoFromQueueState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "" {
		return requestPhotoFromQueue, nil, nil
	}

	switch messageText {
	case "Одобрить":
		comment, err := state.postgres.ViewRequestPhoto.ApprovePhoto(msg.PeerID, state.vkUser, state.groupID)
		if err != nil {
			return requestPhotoFromQueue, nil, err
		}
		if comment != "" {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message(comment)
			return requestPhotoFromQueue, []*params.MessagesSendBuilder{b}, nil
		}
		return requestPhotoFromQueue, nil, nil
	case "Отклонить":
		err := state.postgres.ViewRequestPhoto.RejectPhoto(msg.PeerID)
		if err != nil {
			return requestPhotoFromQueue, nil, err
		}
		return requestPhotoFromQueue, nil, nil
	case "Редактировать":
		return editRequestPhoto, nil, nil
	case "Завершить просмотр заявок":
		err := state.postgres.ViewRequestPhoto.DeletePointer(msg.PeerID)
		if err != nil {
			return requestPhotoFromQueue, nil, err
		}
		return photoStart, nil, nil
	case "⬅️":
		err := state.postgres.ViewRequestPhoto.ChangePointer(msg.PeerID, false)
		if err != nil {
			return requestPhotoFromQueue, nil, err
		}
		return requestPhotoFromQueue, nil, nil
	case "➡️":
		err := state.postgres.ViewRequestPhoto.ChangePointer(msg.PeerID, true)
		if err != nil {
			return requestPhotoFromQueue, nil, err
		}
		return requestPhotoFromQueue, nil, nil
	default:
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Такого действия нет в предложенных вариантах")
		return requestPhotoFromQueue, []*params.MessagesSendBuilder{b}, nil
	}
}

func (state RequestPhotoFromQueueState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	err := state.postgres.ViewRequestPhoto.CreatePersonalAccountPhoto(vkID)
	if err != nil {
		return []*params.MessagesSendBuilder{}, err
	}
	message, attachment, pointer, count, err := state.postgres.ViewRequestPhoto.GetRequestPhoto(vkID)
	if err != nil {
		return []*params.MessagesSendBuilder{}, err
	}
	if message == "" {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Заявок нет")
		k := object.NewMessagesKeyboard(true)
		k.AddRow()
		k.AddTextButton("Завершить просмотр заявок", "", "negative")
		b.Keyboard(k)
		return []*params.MessagesSendBuilder{b}, err
	}
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message(message)
	b.Attachment(attachment)
	k := object.NewMessagesKeyboard(true)
	if count > 1 {
		k.AddRow()
		if pointer != 0 {
			k.AddTextButton("⬅️", "", "secondary")
		}
		if count-pointer > 1 {
			k.AddTextButton("➡️", "", "secondary")
		}
	}
	k.AddRow()
	k.AddTextButton("Одобрить", "", "secondary")
	k.AddTextButton("Отклонить", "", "secondary")
	k.AddRow()
	k.AddTextButton("Редактировать", "", "secondary")
	k.AddTextButton("Завершить просмотр заявок", "", "negative")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state RequestPhotoFromQueueState) Name() stateName {
	return requestPhotoFromQueue
}
