package vkcomaudio

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os/exec"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kurnosovmak/tb-bot-music/internal/config"
	"github.com/kurnosovmak/tb-bot-music/internal/text"
)

var (
	startMessage string
)

func init() {
	startMessage = text.MustLoad("start.html")
}

type Handler struct {
	botApi *tgbotapi.BotAPI
	cfg    config.VK
	logger *slog.Logger
}

func NewHandler(cfg config.VK, botApi *tgbotapi.BotAPI, logger *slog.Logger) *Handler {
	return &Handler{
		botApi: botApi,
		cfg:    cfg,
		logger: logger,
	}
}

// VK API response struct
type vkAudioResponse struct {
	Response []struct {
		URL    string `json:"url"`
		Title  string `json:"title"`
		Artist string `json:"artist"`
		Thumb  struct {
			Photo string `json:"photo_300"`
		} `json:"thumb"`
	} `json:"response"`
}

func (h *Handler) Handle(ctx context.Context, newMessage tgbotapi.Message) error {
	vkAudioId, err := extractVKAudioID(newMessage.Text)
	if err != nil {
		return fmt.Errorf("error parce audio id")
	}

	h.logger.Info("audio id", slog.Any("id", vkAudioId))

	vkApiUrl := fmt.Sprintf("https://api.vk.com/method/audio.getById?v=5.269&client_id=6287487&audios=%s&access_token=%s", vkAudioId, h.cfg.Token)

	resp, err := http.Get(vkApiUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("error bad status: %d", resp.StatusCode))
	}

	var vkResp vkAudioResponse
	if err := json.NewDecoder(resp.Body).Decode(&vkResp); err != nil {
		return err
	}

	if len(vkResp.Response) == 0 || vkResp.Response[0].URL == "" {
		return fmt.Errorf("не удалось получить аудио")
	}

	m3u8URL := vkResp.Response[0].URL // тут m3u8 файл

	// начинаем скачивать
	msg := tgbotapi.NewMessage(newMessage.Chat.ID, fmt.Sprintf("Начинаем скачивать: %s (%s)", vkResp.Response[0].Title, vkResp.Response[0].Artist))
	_, err = h.botApi.Send(msg)

	cmd := exec.Command("ffmpeg", "-i", m3u8URL, "-f", "mp3", "pipe:1")
	audioReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	defer audioReader.Close()

	msgWithPhoto := tgbotapi.NewPhoto(newMessage.Chat.ID, tgbotapi.FileURL(vkResp.Response[0].Thumb.Photo))
	_, err = h.botApi.Send(msgWithPhoto)

	msgWithAudio := tgbotapi.NewAudio(newMessage.Chat.ID, tgbotapi.FileReader{
		Name:   vkResp.Response[0].Title, //+ ".mp3",
		Reader: audioReader,
	})
	_, err = h.botApi.Send(msgWithAudio)
	return err
}

func extractVKAudioID(url string) (string, error) {
	parts := strings.Split(url, "/audio")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid vk audio url")
	}
	idPart := parts[1]

	// если есть GET параметры, убираем их
	id := strings.SplitN(idPart, "?", 2)[0]

	return id, nil
}
