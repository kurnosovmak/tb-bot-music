package router

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	tgevents "github.com/kurnosovmak/tb-bot-music/internal/listener/tg-events"
)

var (
	ErrComandNotFound = fmt.Errorf("command not found")
)

// Command — структура команды
type Command struct {
	handler tgevents.HandleFunc
}

// Handler — роутер команд
type Router struct {
	staticRoutes    map[string]Command
	pregmatchRoutes map[*regexp.Regexp]Command
	notFoundHanle   tgevents.HandleFunc
}

// NewHandler — конструктор роутера с дефолтными командами
func NewRouter() *Router {
	return &Router{
		staticRoutes:    make(map[string]Command),
		pregmatchRoutes: make(map[*regexp.Regexp]Command),
	}
}

// RegisterStatic — регистрация статической команды
func (h *Router) RegisterStatic(cmd string, fn tgevents.HandleFunc) {
	h.staticRoutes[cmd] = Command{
		handler: fn,
	}
}

// RegisterStatic — регистрация статической команды
func (h *Router) RegisterNotFound(fn tgevents.HandleFunc) {
	h.notFoundHanle = fn
}

// RegisterPregmatch — регистрация команды по регулярке
func (h *Router) RegisterPregmatch(pattern string, fn tgevents.HandleFunc) {
	re := regexp.MustCompile(pattern)
	h.pregmatchRoutes[re] = Command{
		handler: fn,
	}
}

// Handle — основная функция маршрутизации
func (h *Router) Handle(ctx context.Context, newMessage tgbotapi.Message) error {
	// Разделяем команду и payload по пробелу
	parts := strings.Fields(newMessage.Text)
	if len(parts) == 0 {
		return fmt.Errorf("empty input")
	}
	cmd := parts[0]
	// Сначала ищем в статических маршрутах
	if c, ok := h.staticRoutes[cmd]; ok {
		return c.handler(ctx, newMessage)
	}

	// Потом проверяем pregmatch маршруты
	for re, c := range h.pregmatchRoutes {
		if re.MatchString(newMessage.Text) {
			return c.handler(ctx, newMessage)
		}
	}

	if h.notFoundHanle != nil {
		return h.notFoundHanle(ctx, newMessage)
	}

	return fmt.Errorf("command not found text = %s : %w", newMessage.Text, ErrComandNotFound)
}
