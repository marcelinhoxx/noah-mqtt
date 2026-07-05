package growatt_app

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

func (h *Client) postForm(url string, data url.Values, responseBody any) (*http.Response, error) {

	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if h.token != "" {
		req.Header.Set("Authorization", "Bearer "+h.token)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10) Noah-MQTT")

	slog.Info("HTTP POST",
		slog.String("url", url))

	resp, err := h.client.Do(req)
	if err != nil {
		slog.Error("HTTP request failed",
			slog.String("url", url),
			slog.String("error", err.Error()))
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	slog.Info("HTTP Response",
		slog.String("url", url),
		slog.Int("status", resp.StatusCode),
		slog.Int("bytes", len(body)))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http %d: %s", resp.StatusCode, string(body))
	}

	if responseBody != nil {

		if err := json.Unmarshal(body, responseBody); err != nil {

			if strings.Contains(string(body), "<html") ||
				strings.Contains(string(body), "<!DOCTYPE") {

				slog.Error("Growatt returned HTML instead of JSON")

				return nil, fmt.Errorf("growatt returned html instead of json")
			}

			return nil, err
		}
	}

	return resp, nil
}
