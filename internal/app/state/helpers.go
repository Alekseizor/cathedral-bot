package state

import (
	"bytes"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/object"
	"github.com/disintegration/imaging"
)

func contains(slice []int64, value int64) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func addBackButton(k *object.MessagesKeyboard) {
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
}

func convertTiffToJpg(tiffImage []byte) ([]byte, error) {
	src, err := imaging.Decode(bytes.NewReader(tiffImage))
	if err != nil {
		return nil, fmt.Errorf("ошибка при декодировании изображения: %v", err)
	}

	var jpgImage []byte
	buf := new(bytes.Buffer)
	err = imaging.Encode(buf, src, imaging.JPEG)
	if err != nil {
		return nil, fmt.Errorf("ошибка при кодировании изображения в JPEG: %v", err)
	}
	jpgImage = buf.Bytes()

	return jpgImage, nil
}
