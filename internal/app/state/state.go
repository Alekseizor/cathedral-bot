package state

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"

	"github.com/Alekseizor/cathedral-bot/internal/app/config"
	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
)

type stateName string

const (
	start         = stateName("start")
	selectArchive = stateName("selectArchive")
	documentStub  = stateName("documentStub")
	photoStub     = stateName("photoStub")
)

type State interface {
	Name() stateName
	Handler(object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error)
	Show() ([]*params.MessagesSendBuilder, error)
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

func (s *States) Init() error {
	postgresRepo := postrgres.New(s.cfg.ClientsConfig.PostgresConfig)

	err := postgresRepo.Init()
	if err != nil {
		return fmt.Errorf("[postrgres.Init]: %w", err)
	}

	s.postgres = postgresRepo

	//здесь инициализируются все стейты
	startState := &StartState{postgres: postgresRepo}
	selectArchiveState := &SelectArchiveState{postgres: postgresRepo}
	documentStubState := &DocumentStubState{postgres: postgresRepo}
	photoStubState := &PhotoStubState{postgres: postgresRepo}

	//мапаем все стейты
	s.statesList = map[stateName]State{
		startState.Name():         startState,
		selectArchiveState.Name(): selectArchiveState,
		documentStubState.Name():  documentStubState,
		photoStubState.Name():     photoStubState,
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
	newState, respMessage, err := state.Handler(obj)
	if err != nil {
		return nil, stateStr, fmt.Errorf("[state.Handler]: %w", err)
	}

	// достали нужную структуру стейта
	// далее берем сообщения, которые надо отправить, для этого стейта
	state = s.statesList[newState]
	newStateMessage, err := state.Show()
	if err != nil {
		return nil, stateStr, fmt.Errorf("[state.Handler]: %w", err)
	}

	respMessage = append(respMessage, newStateMessage...)

	//обновляем стейт пользователя
	err = s.postgres.State.Update(ctx, vkID, string(newState))
	if err != nil {
		return nil, stateStr, fmt.Errorf("[State.Update]: %w", err)
	}

	return respMessage, stateStr, nil
}