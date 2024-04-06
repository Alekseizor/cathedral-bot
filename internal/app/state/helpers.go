package state

import "github.com/SevereCloud/vksdk/v2/object"

func contains(slice []int64, value int64) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func addWBackButton(k *object.MessagesKeyboard) {
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
}
