package config

import (
	"flag"
	"fmt"
	"reflect"
	"strings"

	"github.com/creasty/defaults"
	"github.com/spf13/viper"
)

var (
	envPrefix string
)

func init() {
	flag.StringVar(&envPrefix, "env-prefix", "APP", "Env prefix (default=APP)")
	flag.Parse()
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)

	var cfg Config
	defaults.Set(&cfg)
	autoSetDefaults("", &cfg)

	// Чтение файла, игнорируем ошибку так как мы попробоуем по итогу причитать из env системы или вернем дефолтнуе поля
	viper.ReadInConfig()
	// Автозагрузка из ENV
	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Маппинг в struct
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &cfg, nil
}

// autoSetDefaults рекурсивно ставит дефолты из структуры в Viper
func autoSetDefaults(prefix string, s interface{}) {
	val := reflect.ValueOf(s).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		key := fieldType.Tag.Get("mapstructure")
		if key == "" {
			key = strings.ToLower(fieldType.Name)
		}
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		switch field.Kind() {
		case reflect.Struct:
			autoSetDefaults(fullKey, field.Addr().Interface())
		default:
			viper.SetDefault(fullKey, field.Interface())
		}
	}
}
