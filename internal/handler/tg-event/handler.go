package tgevent

import (
	"github.com/kurnosovmak/tb-bot-music/internal/handler/tg-event/start"
	vkcomaudio "github.com/kurnosovmak/tb-bot-music/internal/handler/tg-event/vkcom-audio"
	"github.com/kurnosovmak/tb-bot-music/internal/listener/tg-events/router"
)

func InitRoutes(
	router *router.Router,

	startHandle *start.Handler,
	vkcomAudioHandle *vkcomaudio.Handler,

) {
	router.RegisterStatic("/start", startHandle.Handle)
	router.RegisterPregmatch(`^https?:\/\/vk\.com\/audio-?\d+_-?\d+$`, vkcomAudioHandle.Handle)
}
