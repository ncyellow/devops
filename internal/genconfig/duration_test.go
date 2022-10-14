package genconfig

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDuration_UnmarshalJSON(t *testing.T) {

	type Message struct {
		Elapsed Duration `json:"elapsed"`
	}
	// проверяем кейс - передана строка все ок
	var msg Message
	err := json.Unmarshal([]byte(`{"elapsed": "1h"}`), &msg)
	assert.NoError(t, err)
	assert.Equal(t, time.Hour*1, msg.Elapsed.Duration)

	// проверяем кейс - передано число все ок
	err = json.Unmarshal([]byte(`{"elapsed": 10}`), &msg)
	assert.NoError(t, err)
	assert.Equal(t, time.Nanosecond*10, msg.Elapsed.Duration)

	// проверяем кейс - кривое значение длительности
	err = json.Unmarshal([]byte(`{"elapsed": "adadaad1h"}`), &msg)
	assert.Error(t, err)
	assert.Equal(t, time.Nanosecond*0, msg.Elapsed.Duration)

	// проверяем кейс - кривой тип для даты
	err = json.Unmarshal([]byte(`{"elapsed": [1,2,3]}`), &msg)
	assert.Error(t, err)
	assert.Equal(t, time.Nanosecond*0, msg.Elapsed.Duration)
}
