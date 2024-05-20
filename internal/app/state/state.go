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
	start                                   = stateName("start")
	selectArchive                           = stateName("selectArchive")
	documentStart                           = stateName("documentStart")
	photoStub                               = stateName("photoStub")
	loadDocument                            = stateName("loadDocument")
	nameDocument                            = stateName("nameDocument")
	authorDocument                          = stateName("authorDocument")
	yearDocument                            = stateName("yearDocument")
	categoryDocument                        = stateName("categoryDocument")
	userCategoryDocument                    = stateName("userCategoryDocument")
	descriptionDocument                     = stateName("descriptionDocument")
	hashtagDocument                         = stateName("hashtagDocument")
	checkDocument                           = stateName("checkDocument")
	editDocument                            = stateName("editDocument")
	editNameDocument                        = stateName("editNameDocument")
	editAuthorDocument                      = stateName("editAuthorDocument")
	editYearDocument                        = stateName("editYearDocument")
	editCategoryDocument                    = stateName("editCategoryDocument")
	editUserCategoryDocument                = stateName("editUserCategoryDocument")
	editDescriptionDocument                 = stateName("editDescriptionDocument")
	editHashtagDocument                     = stateName("editHashtagDocument")
	loadArchive                             = stateName("loadArchive")
	nameArchive                             = stateName("nameArchive")
	authorArchive                           = stateName("authorArchive")
	yearArchive                             = stateName("yearArchive")
	categoryArchive                         = stateName("categoryArchive")
	userCategoryArchive                     = stateName("userCategoryArchive")
	descriptionArchive                      = stateName("descriptionArchive")
	hashtagArchive                          = stateName("hashtagArchive")
	checkArchive                            = stateName("checkArchive")
	documentCabinet                         = stateName("documentCabinet")
	albumsCabinet                           = stateName("albumsCabinet")
	blocking                                = stateName("blocking")
	blockUser                               = stateName("blockUser")
	workingRequestDocument                  = stateName("workingRequestDocument")
	workingDocument                         = stateName("workingDocument")
	nameSearchDocument                      = stateName("nameSearchDocument")
	authorSearchDocument                    = stateName("authorSearchDocument")
	yearSearchDocument                      = stateName("yearSearchDocument")
	categoriesSearchDocument                = stateName("categoriesSearchDocument")
	hashtagSearchDocument                   = stateName("hashtagSearchDocument")
	checkSearchDocument                     = stateName("checkSearchDocument")
	editSearchDocument                      = stateName("editSearchDocument")
	doSearchDocument                        = stateName("doSearchDocument")
	showSearchDocument                      = stateName("showSearchDocument")
	showChosenDocument                      = stateName("showChosenDocument")
	editNameSearchDocument                  = stateName("editNameSearchDocument")
	editAuthorSearchDocument                = stateName("editAuthorSearchDocument")
	editYearSearchDocument                  = stateName("editYearSearchDocument")
	editCategoriesSearchDocument            = stateName("editCategoriesSearchDocument")
	editHashtagSearchDocument               = stateName("editHashtagSearchDocument")
	actionOnDocument                        = stateName("actionOnDocument")
	changeDocument                          = stateName("changeDocument")
	changeTitleDocument                     = stateName("changeTitleDocument")
	changeDescriptionDocument               = stateName("changeDescriptionDocument")
	changeAuthorDocument                    = stateName("changeAuthorDocument")
	changeYearDocument                      = stateName("changeYearDocument")
	changeCategoryDocument                  = stateName("changeCategoryDocument")
	changeHashtagsDocument                  = stateName("changeHashtagsDocument")
	addDocumentAdministrator                = stateName("addDocumentAdministrator")
	removeDocumentAdministrator             = stateName("removeDocumentAdministrator")
	requestDocumentFromQueue                = stateName("requestDocumentFromQueue")
	requestDocumentSpecificApplication      = stateName("requestDocumentSpecificApplication")
	requestDocumentEntrySpecificApplication = stateName("requestDocumentEntrySpecificApplication")
	editDocumentAdmin                       = stateName("editDocumentAdmin")
	editNameDocumentAdmin                   = stateName("editNameDocumentAdmin")
	editAuthorDocumentAdmin                 = stateName("editAuthorDocumentAdmin")
	editYearDocumentAdmin                   = stateName("editYearDocumentAdmin")
	editCategoryDocumentAdmin               = stateName("editCategoryDocumentAdmin")
	editUserCategoryDocumentAdmin           = stateName("editUserCategoryDocumentAdmin")
	editDescriptionDocumentAdmin            = stateName("editDescriptionDocumentAdmin")
	editHashtagDocumentAdmin                = stateName("editHashtagDocumentAdmin")

	photoStart               = stateName("photoStart")
	loadPhoto                = stateName("loadPhoto")
	isPeoplePresentPhoto     = stateName("isPeoplePresentPhoto")
	countPeoplePhoto         = stateName("countPeoplePhoto")
	markedPeoplePhoto        = stateName("markedPeoplePhoto")
	isTeacherPhoto           = stateName("isTeacherPhoto")
	teacherNamePhoto         = stateName("teacherNamePhoto")
	userTeacherNamePhoto     = stateName("userTeacherNamePhoto")
	studentNamePhoto         = stateName("studentNamePhoto")
	eventYearPhoto           = stateName("eventYearPhoto")
	studyProgramPhoto        = stateName("studyProgramPhoto")
	eventNamePhoto           = stateName("eventNamePhoto")
	userEventNamePhoto       = stateName("userEventNamePhoto")
	descriptionPhoto         = stateName("descriptionPhoto")
	checkPhoto               = stateName("checkPhoto")
	editPhoto                = stateName("editPhoto")
	editEventYearPhoto       = stateName("editEventYearPhoto")
	editStudyProgramPhoto    = stateName("editStudyProgramPhoto")
	editEventNamePhoto       = stateName("editEventNamePhoto")
	editUserEventNamePhoto   = stateName("editUserEventNamePhoto")
	editDescriptionPhoto     = stateName("editDescriptionPhoto")
	editIsPeoplePresentPhoto = stateName("editIsPeoplePresentPhoto")
	editCountPeoplePhoto     = stateName("editCountPeoplePhoto")
	editMarkedPeoplePhoto    = stateName("editMarkedPeoplePhoto")
	editIsTeacherPhoto       = stateName("editIsTeacherPhoto")
	editTeacherNamePhoto     = stateName("editTeacherNamePhoto")
	editUserTeacherNamePhoto = stateName("editUserTeacherNamePhoto")
	editStudentNamePhoto     = stateName("editStudentNamePhoto")

	loadPhotoArchive              = stateName("loadPhotoArchive")
	eventYearPhotoArchive         = stateName("eventYearPhotoArchive")
	studyProgramPhotoArchive      = stateName("studyProgramPhotoArchive")
	eventNamePhotoArchive         = stateName("eventNamePhotoArchive")
	userEventNamePhotoArchive     = stateName("userEventNamePhotoArchive")
	descriptionPhotoArchive       = stateName("descriptionPhotoArchive")
	checkPhotoArchive             = stateName("checkPhotoArchive")
	editPhotoArchive              = stateName("editPhotoArchive")
	editEventYearPhotoArchive     = stateName("editEventYearPhotoArchive")
	editStudyProgramPhotoArchive  = stateName("editStudyProgramPhotoArchive")
	editEventNamePhotoArchive     = stateName("editEventNamePhotoArchive")
	editUserEventNamePhotoArchive = stateName("editUserEventNamePhotoArchive")
	editDescriptionPhotoArchive   = stateName("editDescriptionPhotoArchive")

	categorySearchAlbum              = stateName("categorySearchAlbum")
	yearSearchAlbum                  = stateName("yearSearchAlbum")
	findYearSearchAlbum              = stateName("findYearSearchAlbum")
	findYearLess2SearchAlbum         = stateName("findYearLess2SearchAlbum")
	showListYearSearchAlbum          = stateName("showListYearSearchAlbum")
	studyProgramSearchAlbum          = stateName("studyProgramSearchAlbum")
	findStudyProgramSearchAlbum      = stateName("findStudyProgramSearchAlbum")
	findStudyProgramLess2SearchAlbum = stateName("findStudyProgramLess2SearchAlbum")
	showListStudyProgramSearchAlbum  = stateName("showListStudyProgramSearchAlbum")
	eventSearchAlbum                 = stateName("eventSearchAlbum")
	findEventSearchAlbum             = stateName("findEventSearchAlbum")
	surnameTeacherSearchAlbum        = stateName("surnameTeacherSearchAlbum")
	teacherSearchAlbum               = stateName("teacherSearchAlbum")

	personalAccountPhoto = stateName("personalAccountPhoto")

	viewRequestsPhoto = stateName("viewRequestsPhoto")

	editRequestPhoto                = stateName("editRequestPhoto")
	editEventYearRequestPhoto       = stateName("editEventYearRequestPhoto")
	editStudyProgramRequestPhoto    = stateName("editStudyProgramRequestPhoto")
	editEventNameRequestPhoto       = stateName("editEventNameRequestPhoto")
	editUserEventNameRequestPhoto   = stateName("editUserEventNameRequestPhoto")
	editDescriptionRequestPhoto     = stateName("editDescriptionRequestPhoto")
	editIsPeoplePresentRequestPhoto = stateName("editIsPeoplePresentRequestPhoto")
	editCountPeopleRequestPhoto     = stateName("editCountPeopleRequestPhoto")
	editMarkedPeopleRequestPhoto    = stateName("editMarkedPeopleRequestPhoto")
	editIsTeacherRequestPhoto       = stateName("editIsTeacherRequestPhoto")
	editTeacherNameRequestPhoto     = stateName("editTeacherNameRequestPhoto")
	editUserTeacherNameRequestPhoto = stateName("editUserTeacherNameRequestPhoto")
	editStudentNameRequestPhoto     = stateName("editStudentNameRequestPhoto")
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

func (s *States) Init(vk *api.VK, vkUser *api.VK, groupID int) error {
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
	showChosenDocumentState := &ShowChosenDocumentState{postgres: postgresRepo}
	workingDocumentState := &WorkingDocumentState{postgres: postgresRepo}
	actionOnDocumentState := &ActionOnDocumentState{postgres: postgresRepo}
	changeDocumentState := &ChangeDocumentState{postgres: postgresRepo}
	changeTitleDocumentState := &ChangeTitleDocumentState{postgres: postgresRepo}
	changeDescriptionDocumentState := &ChangeDescriptionDocumentState{postgres: postgresRepo}
	changeAuthorDocumentState := &ChangeAuthorDocumentState{postgres: postgresRepo}
	changeYearDocumentState := &ChangeYearDocumentState{postgres: postgresRepo}
	changeCategoryDocumentState := &ChangeCategoryDocumentState{postgres: postgresRepo}
	changeHashtagsDocumentState := &ChangeHashtagsDocumentState{postgres: postgresRepo}
	addDocumentAdministratorState := &AddDocumentAdministratorState{postgres: postgresRepo}
	removeDocumentAdministratorState := &RemoveDocumentAdministratorState{postgres: postgresRepo}
	requestDocumentFromQueueState := &RequestDocumentFromQueueState{postgres: postgresRepo}
	editDocumentAdminState := &EditDocumentAdminState{postgres: postgresRepo}
	requestDocumentSpecificApplicationState := &RequestDocumentSpecificApplicationState{postgres: postgresRepo}
	editNameDocumentAdminState := &EditNameDocumentAdminState{postgres: postgresRepo}
	editAuthorDocumentAdminState := &EditAuthorDocumentAdminState{postgres: postgresRepo}
	editYearDocumentAdminState := &EditYearDocumentAdminState{postgres: postgresRepo}
	editCategoryDocumentAdminState := &EditCategoryDocumentAdminState{postgres: postgresRepo}
	editUserCategoryDocumentAdminState := &EditUserCategoryDocumentAdminState{postgres: postgresRepo}
	editDescriptionDocumentAdminState := &EditDescriptionDocumentAdminState{postgres: postgresRepo}
	editHashtagDocumentAdminState := &EditHashtagDocumentAdminState{postgres: postgresRepo}
	requestDocumentEntrySpecificApplicationState := &RequestDocumentEntrySpecificApplicationState{postgres: postgresRepo}

	photoStartState := &PhotoStartState{postgres: postgresRepo}
	loadPhotoState := &LoadPhotoState{postgres: postgresRepo, vk: vk}
	isPeoplePresentPhotoState := &IsPeoplePresentPhotoState{postgres: postgresRepo}
	countPeoplePhotoState := &CountPeoplePhotoState{postgres: postgresRepo}
	markedPeoplePhotoState := &MarkedPeoplePhotoState{postgres: postgresRepo}
	isTeacherPhotoState := &IsTeacherPhotoState{postgres: postgresRepo}
	teacherNamePhotoState := &TeacherNamePhotoState{postgres: postgresRepo}
	userTeacherNamePhotoState := &UserTeacherNamePhotoState{postgres: postgresRepo}
	studentNamePhotoState := &StudentNamePhotoState{postgres: postgresRepo}
	eventYearPhotoState := &EventYearPhotoState{postgres: postgresRepo}
	studyProgramPhotoState := &StudyProgramPhotoState{postgres: postgresRepo}
	eventNamePhotoState := &EventNamePhotoState{postgres: postgresRepo}
	userEventNamePhotoState := &UserEventNamePhotoState{postgres: postgresRepo}
	descriptionPhotoState := &DescriptionPhotoState{postgres: postgresRepo}
	checkPhotoState := &CheckPhotoState{postgres: postgresRepo}
	editPhotoState := &EditPhotoState{postgres: postgresRepo}
	editEventYearPhotoState := &EditEventYearPhotoState{postgres: postgresRepo}
	editStudyProgramPhotoState := &EditStudyProgramPhotoState{postgres: postgresRepo}
	editEventNamePhotoState := &EditEventNamePhotoState{postgres: postgresRepo}
	editUserEventNamePhotoState := &EditUserEventNamePhotoState{postgres: postgresRepo}
	editDescriptionPhotoState := &EditDescriptionPhotoState{postgres: postgresRepo}
	editIsPeoplePresentPhotoState := &EditIsPeoplePresentPhotoState{postgres: postgresRepo}
	editCountPeoplePhotoState := &EditCountPeoplePhotoState{postgres: postgresRepo}
	editMarkedPeoplePhotoState := &EditMarkedPeoplePhotoState{postgres: postgresRepo}
	editIsTeacherPhotoState := &EditIsTeacherPhotoState{postgres: postgresRepo}
	editTeacherNamePhotoState := &EditTeacherNamePhotoState{postgres: postgresRepo}
	editUserTeacherNamePhotoState := &EditUserTeacherNamePhotoState{postgres: postgresRepo}
	editStudentNamePhotoState := &EditStudentNamePhotoState{postgres: postgresRepo}

	loadPhotoArchiveState := &LoadPhotoArchiveState{postgres: postgresRepo, vk: vk}
	eventYearPhotoArchiveState := &EventYearPhotoArchiveState{postgres: postgresRepo}
	studyProgramPhotoArchiveState := &StudyProgramPhotoArchiveState{postgres: postgresRepo}
	eventNamePhotoArchiveState := &EventNamePhotoArchiveState{postgres: postgresRepo}
	userEventNamePhotoArchiveState := &UserEventNamePhotoArchiveState{postgres: postgresRepo}
	descriptionPhotoArchiveState := &DescriptionPhotoArchiveState{postgres: postgresRepo}
	checkPhotoArchiveState := &CheckPhotoArchiveState{postgres: postgresRepo}
	editPhotoArchiveState := &EditPhotoArchiveState{postgres: postgresRepo}
	editEventYearPhotoArchiveState := &EditEventYearPhotoArchiveState{postgres: postgresRepo}
	editStudyProgramPhotoArchiveState := &EditStudyProgramPhotoArchiveState{postgres: postgresRepo}
	editEventNamePhotoArchiveState := &EditEventNamePhotoArchiveState{postgres: postgresRepo}
	editUserEventNamePhotoArchiveState := &EditUserEventNamePhotoArchiveState{postgres: postgresRepo}
	editDescriptionPhotoArchiveState := &EditDescriptionPhotoArchiveState{postgres: postgresRepo}

	categorySearchAlbumState := &CategorySearchAlbumState{postgres: postgresRepo}
	yearSearchAlbumState := &YearSearchAlbumState{postgres: postgresRepo}
	findYearSearchAlbumState := &FindYearSearchAlbumState{postgres: postgresRepo}
	findYearLess2SearchAlbumState := &FindYearLess2SearchAlbumState{postgres: postgresRepo}
	showListYearSearchAlbumState := &ShowListYearSearchAlbumState{postgres: postgresRepo}
	studyProgramSearchAlbumState := &StudyProgramSearchAlbumState{postgres: postgresRepo}
	findStudyProgramSearchAlbumState := &FindStudyProgramSearchAlbumState{postgres: postgresRepo}
	findStudyProgramLess2SearchAlbumState := &FindStudyProgramLess2SearchAlbumState{postgres: postgresRepo}
	showListStudyProgramSearchAlbumState := &ShowListStudyProgramSearchAlbumState{postgres: postgresRepo}
	eventSearchAlbumState := &EventSearchAlbumState{postgres: postgresRepo}
	findEventSearchAlbumState := &FindEventSearchAlbumState{postgres: postgresRepo}
	surnameTeacherSearchAlbumState := &SurnameTeacherSearchAlbumState{postgres: postgresRepo}
	teacherSearchAlbumState := &TeacherSearchAlbumState{postgres: postgresRepo}

	personalAccountPhotoState := &PersonalAccountPhotoState{postgres: postgresRepo}

	viewRequestsPhotoState := &ViewRequestsPhotoState{postgres: postgresRepo, vk: vk, vkUser: vkUser, groupID: groupID}

	editRequestPhotoState := &EditRequestPhotoState{postgres: postgresRepo}
	editEventYearRequestPhotoState := &EditEventYearRequestPhotoState{postgres: postgresRepo}
	editStudyProgramRequestPhotoState := &EditStudyProgramRequestPhotoState{postgres: postgresRepo}
	editEventNameRequestPhotoState := &EditEventNameRequestPhotoState{postgres: postgresRepo}
	editUserEventNameRequestPhotoState := &EditUserEventNameRequestPhotoState{postgres: postgresRepo}
	editDescriptionRequestPhotoState := &EditDescriptionRequestPhotoState{postgres: postgresRepo}
	editIsPeoplePresentRequestPhotoState := &EditIsPeoplePresentRequestPhotoState{postgres: postgresRepo}
	editCountPeopleRequestPhotoState := &EditCountPeopleRequestPhotoState{postgres: postgresRepo}
	editMarkedPeopleRequestPhotoState := &EditMarkedPeopleRequestPhotoState{postgres: postgresRepo}
	editIsTeacherRequestPhotoState := &EditIsTeacherRequestPhotoState{postgres: postgresRepo}
	editTeacherNameRequestPhotoState := &EditTeacherNameRequestPhotoState{postgres: postgresRepo}
	editUserTeacherNameRequestPhotoState := &EditUserTeacherNameRequestPhotoState{postgres: postgresRepo}
	editStudentNameRequestPhotoState := &EditStudentNameRequestPhotoState{postgres: postgresRepo}

	//мапаем все стейты
	s.statesList = map[stateName]State{
		startState.Name():                                   startState,
		selectArchiveState.Name():                           selectArchiveState,
		documentStartState.Name():                           documentStartState,
		loadDocumentState.Name():                            loadDocumentState,
		nameDocumentState.Name():                            nameDocumentState,
		authorDocumentState.Name():                          authorDocumentState,
		yearDocumentState.Name():                            yearDocumentState,
		categoryDocumentState.Name():                        categoryDocumentState,
		userCategoryDocumentState.Name():                    userCategoryDocumentState,
		descriptionDocumentState.Name():                     descriptionDocumentState,
		hashtagDocumentState.Name():                         hashtagDocumentState,
		checkDocumentState.Name():                           checkDocumentState,
		editDocumentState.Name():                            editDocumentState,
		editNameDocumentState.Name():                        editNameDocumentState,
		editAuthorDocumentState.Name():                      editAuthorDocumentState,
		editYearDocumentState.Name():                        editYearDocumentState,
		editCategoryDocumentState.Name():                    editCategoryDocumentState,
		editUserCategoryDocumentState.Name():                editUserCategoryDocumentState,
		editDescriptionDocumentState.Name():                 editDescriptionDocumentState,
		editHashtagDocumentState.Name():                     editHashtagDocumentState,
		loadArchiveState.Name():                             loadArchiveState,
		nameArchiveState.Name():                             nameArchiveState,
		authorArchiveState.Name():                           authorArchiveState,
		yearArchiveState.Name():                             yearArchiveState,
		categoryArchiveState.Name():                         categoryArchiveState,
		userCategoryArchiveState.Name():                     userCategoryArchiveState,
		descriptionArchiveState.Name():                      descriptionArchiveState,
		hashtagArchiveState.Name():                          hashtagArchiveState,
		checkArchiveState.Name():                            checkArchiveState,
		albumsCabinetState.Name():                           albumsCabinetState,
		documentCabinetState.Name():                         documentCabinetState,
		blockUserState.Name():                               blockUserState,
		blockingState.Name():                                blockingState,
		workingRequestDocumentState.Name():                  workingRequestDocumentState,
		nameSearchDocumentState.Name():                      nameSearchDocumentState,
		authorSearchDocumentState.Name():                    authorSearchDocumentState,
		yearSearchDocumentState.Name():                      yearSearchDocumentState,
		categoriesSearchDocumentState.Name():                categoriesSearchDocumentState,
		hashtagSearchDocumentState.Name():                   hashtagSearchDocumentState,
		checkSearchDocumentState.Name():                     checkSearchDocumentState,
		doSearchDocumentState.Name():                        doSearchDocumentState,
		editSearchDocumentState.Name():                      editSearchDocumentState,
		editNameSearchDocumentState.Name():                  editNameSearchDocumentState,
		editAuthorSearchDocumentState.Name():                editAuthorSearchDocumentState,
		editYearSearchDocumentState.Name():                  editYearSearchDocumentState,
		editCategoriesSearchDocumentState.Name():            editCategoriesSearchDocumentState,
		editHashtagSearchDocumentState.Name():               editHashtagSearchDocumentState,
		showSearchDocumentState.Name():                      showSearchDocumentState,
		showChosenDocumentState.Name():                      showChosenDocumentState,
		workingDocumentState.Name():                         workingDocumentState,
		actionOnDocumentState.Name():                        actionOnDocumentState,
		changeDocumentState.Name():                          changeDocumentState,
		changeTitleDocumentState.Name():                     changeTitleDocumentState,
		changeDescriptionDocumentState.Name():               changeDescriptionDocumentState,
		changeAuthorDocumentState.Name():                    changeAuthorDocumentState,
		changeYearDocumentState.Name():                      changeYearDocumentState,
		changeCategoryDocumentState.Name():                  changeCategoryDocumentState,
		changeHashtagsDocumentState.Name():                  changeHashtagsDocumentState,
		addDocumentAdministratorState.Name():                addDocumentAdministratorState,
		removeDocumentAdministratorState.Name():             removeDocumentAdministratorState,
		requestDocumentFromQueueState.Name():                requestDocumentFromQueueState,
		editDocumentAdminState.Name():                       editDocumentAdminState,
		editNameDocumentAdminState.Name():                   editNameDocumentAdminState,
		editAuthorDocumentAdminState.Name():                 editAuthorDocumentAdminState,
		editYearDocumentAdminState.Name():                   editYearDocumentAdminState,
		editCategoryDocumentAdminState.Name():               editCategoryDocumentAdminState,
		editUserCategoryDocumentAdminState.Name():           editUserCategoryDocumentAdminState,
		editDescriptionDocumentAdminState.Name():            editDescriptionDocumentAdminState,
		editHashtagDocumentAdminState.Name():                editHashtagDocumentAdminState,
		requestDocumentSpecificApplicationState.Name():      requestDocumentSpecificApplicationState,
		requestDocumentEntrySpecificApplicationState.Name(): requestDocumentEntrySpecificApplicationState,

		photoStartState.Name():               photoStartState,
		loadPhotoState.Name():                loadPhotoState,
		isPeoplePresentPhotoState.Name():     isPeoplePresentPhotoState,
		countPeoplePhotoState.Name():         countPeoplePhotoState,
		markedPeoplePhotoState.Name():        markedPeoplePhotoState,
		isTeacherPhotoState.Name():           isTeacherPhotoState,
		teacherNamePhotoState.Name():         teacherNamePhotoState,
		userTeacherNamePhotoState.Name():     userTeacherNamePhotoState,
		studentNamePhotoState.Name():         studentNamePhotoState,
		eventYearPhotoState.Name():           eventYearPhotoState,
		studyProgramPhotoState.Name():        studyProgramPhotoState,
		eventNamePhotoState.Name():           eventNamePhotoState,
		userEventNamePhotoState.Name():       userEventNamePhotoState,
		descriptionPhotoState.Name():         descriptionPhotoState,
		checkPhotoState.Name():               checkPhotoState,
		editPhotoState.Name():                editPhotoState,
		editEventYearPhotoState.Name():       editEventYearPhotoState,
		editStudyProgramPhotoState.Name():    editStudyProgramPhotoState,
		editEventNamePhotoState.Name():       editEventNamePhotoState,
		editUserEventNamePhotoState.Name():   editUserEventNamePhotoState,
		editDescriptionPhotoState.Name():     editDescriptionPhotoState,
		editIsPeoplePresentPhotoState.Name(): editIsPeoplePresentPhotoState,
		editCountPeoplePhotoState.Name():     editCountPeoplePhotoState,
		editMarkedPeoplePhotoState.Name():    editMarkedPeoplePhotoState,
		editIsTeacherPhotoState.Name():       editIsTeacherPhotoState,
		editTeacherNamePhotoState.Name():     editTeacherNamePhotoState,
		editUserTeacherNamePhotoState.Name(): editUserTeacherNamePhotoState,
		editStudentNamePhotoState.Name():     editStudentNamePhotoState,

		loadPhotoArchiveState.Name():              loadPhotoArchiveState,
		eventYearPhotoArchiveState.Name():         eventYearPhotoArchiveState,
		studyProgramPhotoArchiveState.Name():      studyProgramPhotoArchiveState,
		eventNamePhotoArchiveState.Name():         eventNamePhotoArchiveState,
		userEventNamePhotoArchiveState.Name():     userEventNamePhotoArchiveState,
		descriptionPhotoArchiveState.Name():       descriptionPhotoArchiveState,
		checkPhotoArchiveState.Name():             checkPhotoArchiveState,
		editPhotoArchiveState.Name():              editPhotoArchiveState,
		editEventYearPhotoArchiveState.Name():     editEventYearPhotoArchiveState,
		editStudyProgramPhotoArchiveState.Name():  editStudyProgramPhotoArchiveState,
		editEventNamePhotoArchiveState.Name():     editEventNamePhotoArchiveState,
		editUserEventNamePhotoArchiveState.Name(): editUserEventNamePhotoArchiveState,
		editDescriptionPhotoArchiveState.Name():   editDescriptionPhotoArchiveState,

		categorySearchAlbumState.Name():              categorySearchAlbumState,
		yearSearchAlbumState.Name():                  yearSearchAlbumState,
		findYearSearchAlbumState.Name():              findYearSearchAlbumState,
		findYearLess2SearchAlbumState.Name():         findYearLess2SearchAlbumState,
		showListYearSearchAlbumState.Name():          showListYearSearchAlbumState,
		studyProgramSearchAlbumState.Name():          studyProgramSearchAlbumState,
		findStudyProgramSearchAlbumState.Name():      findStudyProgramSearchAlbumState,
		findStudyProgramLess2SearchAlbumState.Name(): findStudyProgramLess2SearchAlbumState,
		showListStudyProgramSearchAlbumState.Name():  showListStudyProgramSearchAlbumState,
		eventSearchAlbumState.Name():                 eventSearchAlbumState,
		findEventSearchAlbumState.Name():             findEventSearchAlbumState,
		surnameTeacherSearchAlbumState.Name():        surnameTeacherSearchAlbumState,
		teacherSearchAlbumState.Name():               teacherSearchAlbumState,

		personalAccountPhotoState.Name(): personalAccountPhotoState,

		viewRequestsPhotoState.Name(): viewRequestsPhotoState,

		editRequestPhotoState.Name():                editRequestPhotoState,
		editEventYearRequestPhotoState.Name():       editEventYearRequestPhotoState,
		editStudyProgramRequestPhotoState.Name():    editStudyProgramRequestPhotoState,
		editEventNameRequestPhotoState.Name():       editEventNameRequestPhotoState,
		editUserEventNameRequestPhotoState.Name():   editUserEventNameRequestPhotoState,
		editDescriptionRequestPhotoState.Name():     editDescriptionRequestPhotoState,
		editIsPeoplePresentRequestPhotoState.Name(): editIsPeoplePresentRequestPhotoState,
		editCountPeopleRequestPhotoState.Name():     editCountPeopleRequestPhotoState,
		editMarkedPeopleRequestPhotoState.Name():    editMarkedPeopleRequestPhotoState,
		editIsTeacherRequestPhotoState.Name():       editIsTeacherRequestPhotoState,
		editTeacherNameRequestPhotoState.Name():     editTeacherNameRequestPhotoState,
		editUserTeacherNameRequestPhotoState.Name(): editUserTeacherNameRequestPhotoState,
		editStudentNameRequestPhotoState.Name():     editStudentNameRequestPhotoState,
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
