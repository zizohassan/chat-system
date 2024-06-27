package utils

import (
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"

	"log"
	"os"
	"time"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	logger *logrus.Logger
)

type LokiHook struct {
	URL    string
	Client *http.Client
}

type LokiEntry struct {
	Labels  string       `json:"labels"`
	Entries []LokiStream `json:"entries"`
}

type LokiStream struct {
	Timestamp string `json:"ts"`
	Line      string `json:"line"`
}


func InitLogger() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	Init()
}

func NewLokiHook(url string) *LokiHook {
	return &LokiHook{
		URL:    url,
		Client: &http.Client{},
	}
}

func (hook *LokiHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}

	stream := LokiEntry{
		Labels: `{job="chat-system"}`,
		Entries: []LokiStream{
			{
				Timestamp: time.Now().Format(time.RFC3339Nano),
				Line:      line,
			},
		},
	}

	payload, err := json.Marshal(stream)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", hook.URL, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := hook.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return err
	}

	return nil
}

func (hook *LokiHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func Init() {
	logger = logrus.New()
	lokiURL := "http://loki:3100/loki/api/v1/push"

	hook := NewLokiHook(lokiURL)
	logger.AddHook(hook)
}

func GetLogger() *logrus.Logger {
	return logger
}