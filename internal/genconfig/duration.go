package genconfig

import (
	"encoding/json"
	"fmt"
	"time"
)

// Duration приходится делать так, как стандартный парсинг json в struct не поддерживает time.Duration из строки
type Duration struct {
	time.Duration
}

func (duration *Duration) UnmarshalJSON(b []byte) error {
	var unmarshalledJson interface{}

	err := json.Unmarshal(b, &unmarshalledJson)
	if err != nil {
		return err
	}

	switch value := unmarshalledJson.(type) {
	case float64:
		duration.Duration = time.Duration(value)
	case string:
		duration.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid duration: %#v", unmarshalledJson)
	}

	return nil
}
