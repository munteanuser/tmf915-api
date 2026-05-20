package events

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type CallbackFetcher interface {
	AllCallbacks(ctx context.Context) ([]string, error)
}

type Dispatcher struct {
	fetcher    CallbackFetcher
	timeout    time.Duration
	maxRetry   int
	httpClient *http.Client
	log        zerolog.Logger
}

func NewDispatcher(fetcher CallbackFetcher, timeout time.Duration, maxRetry int, log zerolog.Logger) *Dispatcher {
	return &Dispatcher{
		fetcher:    fetcher,
		timeout:    timeout,
		maxRetry:   maxRetry,
		httpClient: &http.Client{Timeout: timeout},
		log:        log,
	}
}

func (d *Dispatcher) Dispatch(eventType string, payload interface{}) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		callbacks, err := d.fetcher.AllCallbacks(ctx)
		if err != nil {
			d.log.Error().Err(err).Msg("dispatcher: fetch callbacks")
			return
		}

		body, err := json.Marshal(map[string]interface{}{
			"eventType": eventType,
			"eventTime": time.Now().UTC().Format(time.RFC3339),
			"event":     payload,
		})
		if err != nil {
			d.log.Error().Err(err).Msg("dispatcher: marshal payload")
			return
		}

		for _, cb := range callbacks {
			go d.post(cb, body)
		}
	}()
}

func (d *Dispatcher) post(url string, body []byte) {
	for attempt := 1; attempt <= d.maxRetry; attempt++ {
		resp, err := d.httpClient.Post(url, "application/json;charset=utf-8", bytes.NewReader(body))
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode < 300 {
				return
			}
			d.log.Warn().Str("url", url).Int("status", resp.StatusCode).
				Int("attempt", attempt).Msg("dispatcher: non-2xx response")
		} else {
			d.log.Warn().Err(err).Str("url", url).Int("attempt", attempt).Msg("dispatcher: post failed")
		}
		if attempt < d.maxRetry {
			time.Sleep(time.Duration(attempt) * 2 * time.Second)
		}
	}
	d.log.Error().Str("url", url).Int("retries", d.maxRetry).Msg("dispatcher: giving up")
}
