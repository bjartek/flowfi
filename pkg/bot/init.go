package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type FlowFi struct {
	Logger               *zap.Logger
	Tgbot                *tgbotapi.BotAPI
	UpdateConfig         tgbotapi.UpdateConfig
	Config               Config
	BaseUrl              string
	Subscriptions        *Subscriptions
	ScreenshotUrlPattern string
	Store                SubscriptionStore
}

// TODO: error handling
func NewBot() *FlowFi {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	var config Config
	err = envdecode.Decode(&config)
	if err != nil {
		panic(err)
	}

	logger := CreateLogger(config)

	// Replace "YOUR_TELEGRAM_BOT_TOKEN" with the API token provided by the BotFather.
	tgbot, err := tgbotapi.NewBotAPI(config.Telegram.Token)
	if err != nil {
		panic(err)
	}

	logger = logger.With(zap.String("user", tgbot.Self.UserName))
	// Set the bot to use debug mode (verbose logging).
	tgbot.Debug = config.Telegram.Debug

	logAdapter := ZapBotLogger{
		Logger: logger,
	}
	logger.Info("Authorized")
	err = tgbotapi.SetLogger(logAdapter)
	if err != nil {
		panic(err)
	}

	updateConfig := tgbotapi.NewUpdate(0)

	updateConfig.Timeout = config.Telegram.Timeout

	// or use db store here, whatever
	store := NewFileSubscriptionStore("subscriptions.json")

	// Load existing subscriptions
	subData, err := store.LoadSubscriptions()
	if err != nil {
		logger.Error("Failed loading store", zap.Error(err))
	}

	subscriptions := &Subscriptions{
		pairs: subData,
	}
	return &FlowFi{
		Logger:               logger,
		Config:               config,
		Tgbot:                tgbot,
		UpdateConfig:         updateConfig,
		BaseUrl:              "https://api.geckoterminal.com/api/v2/networks/flow-evm/pools",
		ScreenshotUrlPattern: "https://www.geckoterminal.com/flow-evm/pools/%s?embed=1&info=0&swaps=0",
		Subscriptions:        subscriptions,
		Store:                store,
	}
}

type Config struct {
	Hostname    string `env:"HOSTNAME,default=kernel"`
	Application string `env:"APPLICATION,required"`

	Logger struct {
		Level zap.AtomicLevel `env:"LOGGING_LEVEL,default=info"`
		Local bool            `env:"LOGGING_LOCAL,default=true,strict"`
	}

	Telegram struct {
		Token   string `env:"TOKEN,required=true"`
		Debug   bool   `env:"DEBUG,default=true"`
		Timeout int    `env:"TIMEOUT,default=60"`
	}

	// repeat this amount of emoticon for each increase in binary bucket so 8/16/32/64/128 aso
	EmoticonStep int `env:"EMOTICON_STEP,default=4"`
}

func CreateLogger(cfg Config) *zap.Logger {
	loggingConfig := zap.Config{
		Level:            cfg.Logger.Level,
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	if !cfg.Logger.Local {
		// in kubernetes we do not want logs to stderr
		loggingConfig.OutputPaths = []string{"stdout"}
		loggingConfig.EncoderConfig = zap.NewProductionEncoderConfig()
		loggingConfig.Encoding = "json"
		loggingConfig.Development = false
	} else {
		loggingConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	logger, err := loggingConfig.Build()
	if err != nil {
		log.Fatal(err)
	}
	return logger
}
