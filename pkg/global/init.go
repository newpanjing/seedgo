package global

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func InitConfig(configPath string) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	Config = &Configuration{}
	if err := v.Unmarshal(Config); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	v.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed: %s", e.Name)
		if err := v.Unmarshal(Config); err != nil {
			log.Printf("Unable to decode into struct: %v", err)
		}
	})
	v.WatchConfig()

	log.Printf("Config loaded from %s", configPath)
}
