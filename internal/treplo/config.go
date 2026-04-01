package treplo

type Config struct {
	DatabaseDSN string `env:"DATABASE_DSN" flag:"d" jsonConfig:"database_dsn" default:"" description:"Database connection string"`
	TgToken     string `env:"TG_TOKEN" flag:"t" jsonConfig:"tg_token" default:"" description:"Telegram bot token"`
	// SaluteSpeechClientID         string `env:"SALUTE_SPEECH_CLIENT_ID" flag:"s" jsonConfig:"salute_speech_client_id" default:"" description:"Salute speech client id"`
	SaluteSpeechAuthorizationKey string `env:"SALUTE_SPEECH_AUTHORIZATION_KEY" flag:"a" jsonConfig:"salute_speech_authorization_key" default:"" description:"Salute speech authorization key"`
	GigachatAuthorizationKey     string `env:"GIGACHAT_AUTHORIZATION_KEY" flag:"g" jsonConfig:"gigachat_authorization_key" default:"" description:"Gigachat authorization key"`
	StoragePath                  string `env:"STORAGE_PATH" flag:"s" jsonConfig:"storage" default:"./storage" description:"Path for files"`
}
