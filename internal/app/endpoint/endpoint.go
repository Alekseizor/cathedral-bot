package endpoint

import (
	"context"
	"fmt"
	"time"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"github.com/rs/zerolog/log"

	"github.com/Alekseizor/cathedral-bot/internal/app/config"
	"github.com/Alekseizor/cathedral-bot/internal/app/metrics/prometheus"
	"github.com/Alekseizor/cathedral-bot/internal/app/state"
)

type Endpoint struct {
	cfg     config.Config
	lp      *longpoll.LongPoll
	vk      *api.VK
	states  *state.States
	metrics *prometheus.Metrics
}

func New(cfg config.Config) *Endpoint {
	return &Endpoint{
		cfg:     cfg,
		states:  state.New(cfg),
		metrics: prometheus.New(cfg.MonitoringConfig.Address),
	}
}

func (e *Endpoint) Init(ctx context.Context) error {
	err := e.states.Init()
	if err != nil {
		return fmt.Errorf("[state.Init]: %w", err)
	}

	vk := api.NewVK(e.cfg.BotConfig.Token)
	e.vk = vk

	group, err := vk.GroupsGetByID(nil)
	if err != nil {
		return fmt.Errorf("[vk.GroupsGetByID]: %w", err)
	}

	lp, err := longpoll.NewLongPoll(vk, group[0].ID)
	if err != nil {
		return fmt.Errorf("[longpoll.NewLongPoll]: %w", err)
	}

	// ждем новые сообщения
	lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		// будем считать время выполнения запроса
		startTime := time.Now()

		// делаем обработку от паники
		defer func() {
			if err := recover(); err != nil {
				log.Ctx(ctx).Error().Msgf("[Endpoint.Init:MessageNew:recover]: vkID -%d ,error - %v", obj.Message.PeerID, err)
				e.metrics.Panic()
			}
		}()

		log.Ctx(ctx).Info().Msgf("User - %d: Message - %s", obj.Message.PeerID, obj.Message.Text)

		// обрабатываем сообщения и подготавливаем ответ
		respMessages, oldState, err := e.states.Handler(ctx, obj.Message)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msgf("[Endpoint.Init:MessageNew:states.Handler]: vkID -%d", obj.Message.PeerID)
			e.metrics.HandlerError(oldState)

			// произошла ошибка при работе, напишем об этом пользователю
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("При обработке сообщения произошла ошибка. Мы уже разбираемся")
			respMessages = []*params.MessagesSendBuilder{b}
		}

		// считаем время обработки запроса для каждого из стейтов
		defer func() {
			e.metrics.Percentiles(time.Since(startTime), oldState)
		}()

		// происходит отправка сообщений
		for _, message := range respMessages {
			message.PeerID(obj.Message.PeerID)

			_, err := e.vk.MessagesSend(message.Params)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msgf("[Endpoint: vk.MessagesSend] vkID - %d", obj.Message.PeerID)
				e.metrics.SendError(oldState)
			}
		}
	})

	e.lp = lp

	return nil
}

func (e *Endpoint) Run(ctx context.Context) error {
	go func() {
		err := e.metrics.Run()
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("[prometheus.Listen]")
		}
	}()

	err := e.lp.Run()
	if err != nil {
		return fmt.Errorf("[lp.Run]: %w", err)
	}

	return nil
}
