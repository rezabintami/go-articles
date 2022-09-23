package config

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"sync"

	"github.com/spf13/viper"
)

var (
	config      *viper.Viper
	message     *viper.Viper
	onceConfig  sync.Once
	onceMessage sync.Once
)

func GetConfiguration(code string) string {
	if os.Getenv("app.env") == "" {
		onceConfig.Do(func() {
			config = viper.New()
			config.SetConfigName("config")
			config.SetConfigType("yaml")
			config.AddConfigPath("./")

			err := config.ReadInConfig()
			if err != nil {
				log.Fatalf("err GetConfiguration: %s, code: %s", err.Error(), code)
			}
		})

		value := config.GetString(code)
		if value == "" && !config.InConfig(code) {
			log.Fatalf("err GetConfiguration: not found, code: %s \n", code)
		}
		return value

	}

	value := os.Getenv(code)

	return value
}

// Message get message for message.json
func Message(code string, data map[string]interface{}) string {
	country := "en"

	onceMessage.Do(func() {
		message = viper.New()
		message.SetConfigType("json")
		message.SetConfigName("global.json")

		langFile := "./resources/lang/" + country
		message.AddConfigPath(langFile)

		err := message.ReadInConfig()
		if err != nil {
			log.Printf("err Message: %s, code: %s", err.Error(), code)
		} else {
			log.Printf("Successfully read language, file: %s", langFile)
		}
	})

	if message == nil {
		log.Printf("err Message: not initialized, code: %s", code)
	}

	text := message.GetString(code)
	if text == "" {
		if !message.InConfig(code) {
			log.Printf("err Message: not found, code: %s", code)
		}
		return ""
	}

	var tpl bytes.Buffer
	t, err := template.New("").Parse(text)
	if err != nil {
		log.Printf("err Message: %s, code: %s", err.Error(), code)
		return ""
	}

	err = t.Execute(&tpl, data)
	if err != nil {
		log.Printf("err Message: %s, code: %s", err.Error(), code)
		return ""
	}

	return tpl.String()
}
