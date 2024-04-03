package state

import (
	"context"
	"strconv"

	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"

	"github.com/Alekseizor/cathedral-bot/internal/app/repo/postrgres"
)

type BlockUserState struct {
	postgres *postrgres.Repo
}

func (state BlockUserState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	messageText := msg.Text
	if messageText == "Назад" {
		return documentCabinet, nil, nil
	}
	vkID, err := strconv.Atoi(messageText)
	if err != nil || vkID < 100000000 || vkID >= 1000000000 {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("VK ID должно быть числом, например, 221486551")
		return blockUser, []*params.MessagesSendBuilder{b}, nil
	}

	err = state.postgres.State.Update(ctx, vkID, string(blocking))
	if err != nil {
		return blockUser, nil, err
	}

	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Пользователь заблокирован, он больше не сможет пользоваться сервисом")
	return blockUser, []*params.MessagesSendBuilder{b}, nil
}

func (state BlockUserState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите VK ID пользователя, которого нужно заблокировать")
	k := object.NewMessagesKeyboard(true)
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	return []*params.MessagesSendBuilder{b}, nil
}

func (state BlockUserState) Name() stateName {
	return blockUser
}

////////////////////

type BlockingState struct {
}

func (state BlockingState) Handler(ctx context.Context, msg object.MessagesMessage) (stateName, []*params.MessagesSendBuilder, error) {
	return blocking, nil, nil
}

func (state BlockingState) Show(ctx context.Context, vkID int) ([]*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Администратор заблокировал Вас навсегда, доступ к ресурсу запрещен")
	return []*params.MessagesSendBuilder{b}, nil
}

func (state BlockingState) Name() stateName {
	return blocking
}
