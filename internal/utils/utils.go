package utils

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/maxuanquang/social-network/configs"
	"go.uber.org/zap"
)

func NewLogger(cfg *configs.LoggerConfig) (*zap.Logger, error) {
	loggerCfg := zap.NewDevelopmentConfig()

	switch cfg.Level {
	case "debug":
		loggerCfg.Level.SetLevel(zap.DebugLevel)
	case "info":
		loggerCfg.Level.SetLevel(zap.InfoLevel)
	case "warn":
		loggerCfg.Level.SetLevel(zap.WarnLevel)
	case "error":
		loggerCfg.Level.SetLevel(zap.ErrorLevel)
	case "dpanic":
		loggerCfg.Level.SetLevel(zap.DPanicLevel)
	case "panic":
		loggerCfg.Level.SetLevel(zap.PanicLevel)
	case "fatal":
		loggerCfg.Level.SetLevel(zap.FatalLevel)
	}

	logger := zap.Must(loggerCfg.Build())
	return logger, nil
}

// Unmarshal converts map[string]string to a struct.
// It takes name of each field in map[string]string and maps with json tags in target struct
func Unmarshal(sourceMap map[string]string, objectPointer interface{}) error {
	objValue := reflect.ValueOf(objectPointer)
	if objValue.Kind() != reflect.Ptr || objValue.IsNil() {
		return fmt.Errorf("obj must be a non-nil pointer to a struct")
	}

	// Iterate over struct fields
	for i := 0; i < objValue.Elem().NumField(); i++ {
		field := objValue.Elem().Field(i)

		jsonTag := objValue.Elem().Type().Field(i).Tag.Get("json")
		mapValue, ok := sourceMap[jsonTag]
		if !ok || len(mapValue) == 0 {
			continue
		}

		switch field.Kind() {
		case reflect.Int64:
			intValue, err := strconv.ParseInt(mapValue, 10, 64)
			if err != nil {
				return err
			}
			field.SetInt(intValue)
		case reflect.String:
			field.SetString(mapValue)
		case reflect.Bool:
			boolValue, err := strconv.ParseBool(mapValue)
			if err != nil {
				return err
			}
			field.SetBool(boolValue)
		}
	}
	return nil
}
