package util

import (
	"encoding/json"
	"time"
	"websocketbenchmark/model"
)

func GetEvent(c int64) ([]byte, error) {
	payload := model.Payload{
		Count:     c,
		Timestamp: time.Now().UnixMilli(),
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return []byte{}, err
	}

	return b, nil
}
