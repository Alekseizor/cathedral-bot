package state

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"

	"github.com/Alekseizor/cathedral-bot/internal/app/config"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
)

type stateName string

const (
	start                        = stateName("start")
	selectArchive                = stateName("selectArchive")
	documentStart                = stateName("documentStart")
	photoStub                    = stateName("photoStub")
	loadDocument                 = stateName("loadDocument")
	nameDocument                 = stateName("nameDocument")
	authorDocument               = stateName("authorDocument")
	yearDocument                 = stateName("yearDocument")
	categoryDocument             = stateName("categoryDocument")
	userCategoryDocument         = stateName("userCategoryDocument")
	descriptionDocument          = stateName("descriptionDocument")
	hashtagDocument              = stateName("hashtagDocument")
	checkDocument                = stateName("checkDocument")
	editDocument                 = stateName("editDocument")
	editNameDocument             = stateName("editNameDocument")
	editAuthorDocument           = stateName("editAuthorDocument")
	editYearDocument             = stateName("editYearDocument")
	editCategoryDocument         = stateName("editCategoryDocument")
	editUserCategoryDocument     = stateName("editUserCategoryDocument")
	editDescriptionDocument      = stateName("editDescriptionDocument")
	editHashtagDocument          = stateName("editHashtagDocument")
	loadArchive                  = stateName("loadArchive")
	nameArchive                  = stateName("nameArchive")
	authorArchive                = stateName("authorArchive")
	yearArchive                  = stateName("yearArchive")
	categoryArchive              = stateName("categoryArchive")
	userCategoryArchive          = stateName("userCategoryArchive")
	descriptionArchive           = stateName("descriptionArchive")
	hashtagArchive               = stateName("hashtagArchive")
	checkArchive                 = stateName("checkArchive")
	documentCabinet              = stateName("documentCabinet")
	albumsCabinet                = stateName("albumsCabinet")
	blocking                     = stateName("blocking")
	blockUser                    = stateName("blockUser")
	workingRequestDocument       = stateName("workingRequestDocument")
	workingDocument              = stateName("workingDocument")
	nameSearchDocument           = stateName("nameSearchDocument")
	authorSearchDocument         = stateName("authorSearchDocument")
	yearSearchDocument           = stateName("yearSearchDocument")
	categoriesSearchDocument     = stateName("categoriesSearchDocument")
	hashtagSearchDocument        = stateName("hashtagSearchDocument")
	checkSearchDocument          = stateName("checkSearchDocument")
	editSearchDocument           = stateName("editSearchDocument")
	doSearchDocument             = stateName("doSearchDocument")
	showSearchDocument           = stateName("showSearchDocument")
	editNameSearchDocument       = stateName("editNameSearchDocument")
	editAuthorSearchDocument     = stateName("editAuthorSearchDocument")
	editYearSearchDocument       = stateName("editYearSearchDocument")
	editCategoriesSearchDocument = stateName("editCategoriesSearchDocument")
	editHashtagSearchDocument    = stateName("editHashtagSearchDocument")
)

type State interface {
	Name() stateName
	Handler(context.Context, object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error)
	Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error)
}

type States struct {
	cfg        config.Config
	statesList map[stateName]State
	postgres   *postrgres.Repo
}

func New(cfg config.Config) *States {
	return &States{
		cfg: cfg,
	}
}

func (s *States) Init(vk *api.VK) error {
	postgresRepo := postrgres.New(s.cfg.ClientsConfig.PostgresConfig)

	err := postgresRepo.Init()
	if err != nil {
		return fmt.Errorf("[postrgres.Init]: %w", err)
	}

	s.postgres = postgresRepo

	//здесь инициализируются все стейты
	startState := &StartState{postgres: postgresRepo}
	selectArchiveState := &SelectArchiveState{postgres: postgresRepo}
	documentStartState := &DocumentStartState{postgres: postgresRepo}
	photoStubState := &PhotoStubState{postgres: postgresRepo}
	loadDocumentState := &LoadDocumentState{postgres: postgresRepo, vk: vk}
	nameDocumentState := &NameDocumentState{postgres: postgresRepo}
	authorDocumentState := &AuthorDocumentState{postgres: postgresRepo}
	yearDocumentState := &YearDocumentState{postgres: postgresRepo}
	categoryDocumentState := &CategoryDocumentState{postgres: postgresRepo}
	userCategoryDocumentState := &UserCategoryDocumentState{postgres: postgresRepo}
	descriptionDocumentState := &DescriptionDocumentState{postgres: postgresRepo}
	hashtagDocumentState := &HashtagDocumentState{postgres: postgresRepo}
	checkDocumentState := &CheckDocumentState{postgres: postgresRepo}
	editDocumentState := &EditDocumentState{postgres: postgresRepo}
	editNameDocumentState := &EditNameDocumentState{postgres: postgresRepo}
	editAuthorDocumentState := &EditAuthorDocumentState{postgres: postgresRepo}
	editYearDocumentState := &EditYearDocumentState{postgres: postgresRepo}
	editCategoryDocumentState := &EditCategoryDocumentState{postgres: postgresRepo}
	editUserCategoryDocumentState := &EditUserCategoryDocumentState{postgres: postgresRepo}
	editDescriptionDocumentState := &EditDescriptionDocumentState{postgres: postgresRepo}
	editHashtagDocumentState := &EditHashtagDocumentState{postgres: postgresRepo}
	loadArchiveState := &LoadArchiveState{postgres: postgresRepo, vk: vk}
	nameArchiveState := &NameArchiveState{postgres: postgresRepo}
	authorArchiveState := &AuthorArchiveState{postgres: postgresRepo}
	yearArchiveState := &YearArchiveState{postgres: postgresRepo}
	categoryArchiveState := &CategoryArchiveState{postgres: postgresRepo}
	userCategoryArchiveState := &UserCategoryArchiveState{postgres: postgresRepo}
	descriptionArchiveState := &DescriptionArchiveState{postgres: postgresRepo}
	hashtagArchiveState := &HashtagArchiveState{postgres: postgresRepo}
	checkArchiveState := &CheckArchiveState{postgres: postgresRepo}
	albumsCabinetState := &AlbumsCabinetState{postgres: postgresRepo}
	documentCabinetState := &DocumentCabinetState{postgres: postgresRepo}
	blockUserState := &BlockUserState{postgres: postgresRepo}
	blockingState := &BlockingState{}
	workingRequestDocumentState := &WorkingRequestDocumentState{}
	nameSearchDocumentState := &NameSearchDocumentState{postgres: postgresRepo}
	authorSearchDocumentState := &AuthorSearchDocumentState{postgres: postgresRepo}
	yearSearchDocumentState := &YearSearchDocumentState{postgres: postgresRepo}
	categoriesSearchDocumentState := &CategoriesSearchDocumentState{postgres: postgresRepo}
	hashtagSearchDocumentState := &HashtagSearchDocumentState{postgres: postgresRepo}
	checkSearchDocumentState := &CheckSearchDocumentState{postgres: postgresRepo}
	doSearchDocumentState := &DoSearchDocumentState{postgres: postgresRepo}
	editSearchDocumentState := &EditSearchDocumentState{postgres: postgresRepo}
	editNameSearchDocumentState := &EditNameSearchDocumentState{postgres: postgresRepo}
	editAuthorSearchDocumentState := &EditAuthorSearchDocumentState{postgres: postgresRepo}
	editYearSearchDocumentState := &EditYearSearchDocumentState{postgres: postgresRepo}
	editCategoriesSearchDocumentState := &EditCategoriesSearchDocumentState{postgres: postgresRepo}
	editHashtagSearchDocumentState := &EditHashtagSearchDocumentState{postgres: postgresRepo}
	showSearchDocumentState := &ShowSearchDocumentState{postgres: postgresRepo}

	//мапаем все стейты
	s.statesList = map[stateName]State{
		startState.Name():                        startState,
		selectArchiveState.Name():                selectArchiveState,
		documentStartState.Name():                documentStartState,
		photoStubState.Name():                    photoStubState,
		loadDocumentState.Name():                 loadDocumentState,
		nameDocumentState.Name():                 nameDocumentState,
		authorDocumentState.Name():               authorDocumentState,
		yearDocumentState.Name():                 yearDocumentState,
		categoryDocumentState.Name():             categoryDocumentState,
		userCategoryDocumentState.Name():         userCategoryDocumentState,
		descriptionDocumentState.Name():          descriptionDocumentState,
		hashtagDocumentState.Name():              hashtagDocumentState,
		checkDocumentState.Name():                checkDocumentState,
		editDocumentState.Name():                 editDocumentState,
		editNameDocumentState.Name():             editNameDocumentState,
		editAuthorDocumentState.Name():           editAuthorDocumentState,
		editYearDocumentState.Name():             editYearDocumentState,
		editCategoryDocumentState.Name():         editCategoryDocumentState,
		editUserCategoryDocumentState.Name():     editUserCategoryDocumentState,
		editDescriptionDocumentState.Name():      editDescriptionDocumentState,
		editHashtagDocumentState.Name():          editHashtagDocumentState,
		loadArchiveState.Name():                  loadArchiveState,
		nameArchiveState.Name():                  nameArchiveState,
		authorArchiveState.Name():                authorArchiveState,
		yearArchiveState.Name():                  yearArchiveState,
		categoryArchiveState.Name():              categoryArchiveState,
		userCategoryArchiveState.Name():          userCategoryArchiveState,
		descriptionArchiveState.Name():           descriptionArchiveState,
		hashtagArchiveState.Name():               hashtagArchiveState,
		checkArchiveState.Name():                 checkArchiveState,
		albumsCabinetState.Name():                albumsCabinetState,
		documentCabinetState.Name():              documentCabinetState,
		blockUserState.Name():                    blockUserState,
		blockingState.Name():                     blockingState,
		workingRequestDocumentState.Name():       workingRequestDocumentState,
		nameSearchDocumentState.Name():           nameSearchDocumentState,
		authorSearchDocumentState.Name():         authorSearchDocumentState,
		yearSearchDocumentState.Name():           yearSearchDocumentState,
		categoriesSearchDocumentState.Name():     categoriesSearchDocumentState,
		hashtagSearchDocumentState.Name():        hashtagSearchDocumentState,
		checkSearchDocumentState.Name():          checkSearchDocumentState,
		doSearchDocumentState.Name():             doSearchDocumentState,
		editSearchDocumentState.Name():           editSearchDocumentState,
		editNameSearchDocumentState.Name():       editNameSearchDocumentState,
		editAuthorSearchDocumentState.Name():     editAuthorSearchDocumentState,
		editYearSearchDocumentState.Name():       editYearSearchDocumentState,
		editCategoriesSearchDocumentState.Name(): editCategoriesSearchDocumentState,
		editHashtagSearchDocumentState.Name():    editHashtagSearchDocumentState,
		showSearchDocumentState.Name():           showSearchDocumentState,
	}

	return nil
}

// Handler - вся бизнес логика приложения выполняется здесь
func (s *States) Handler(ctx context.Context, obj object.MessagesMessage) ([]*params.MessagesSendBuilder, string, error) {
	message := obj
	vkID := message.PeerID

	//достаем стейт пользователя
	stateStr, err := s.postgres.State.Get(ctx, vkID)
	if err != nil {
		// пользователь впервые пришел к нам, добавляем в базу
		if err != sql.ErrNoRows {
			return nil, stateStr, fmt.Errorf("[State.Get]: %w", err)
		} else {
			err = s.postgres.State.Insert(ctx, vkID)
			stateStr = string(start)
			if err != nil {
				return nil, stateStr, fmt.Errorf("[State.Insert]: %w", err)
			}
		}
	}

	// достали нужную структуру стейта
	userState := stateName(stateStr)
	state := s.statesList[userState]

	//выполняем обработку сообщения согласно стейту
	newState, respMessage, err := state.Handler(ctx, obj)
	if err != nil {
		return nil, stateStr, fmt.Errorf("[state.Handler]: %w", err)
	}

	// достали нужную структуру стейта
	// далее берем сообщения, которые надо отправить, для этого стейта
	state = s.statesList[newState]
	newStateMessage, err := state.Show(ctx, vkID)
	if err != nil {
		return nil, stateStr, fmt.Errorf("[state.Show]: %w", err)
	}

	respMessage = append(respMessage, newStateMessage...)

	//обновляем стейт пользователя
	err = s.postgres.State.Update(ctx, vkID, string(newState))
	if err != nil {
		return nil, stateStr, fmt.Errorf("[State.Update]: %w", err)
	}

	return respMessage, stateStr, nil
}
