package config

type Config struct {
	Bot    Bot
	Logger Logger
	VK     VK
}

type Bot struct {
	Token                string
	Debug                bool `default:"false"`
	Timeout              int  `default:"5"`
	EventLongpollTimeout int  `default:"60"`
}

type Logger struct {
	Level string `default:"error"`
}

type VK struct {
	Token string
}
