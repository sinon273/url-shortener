//go:build integration

package gateway_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
)

func TestFullScenario(t *testing.T) {
	env := setupInfra(t)
	root := repoRoot(t)

	startService(t, filepath.Join(root, "user"), map[string]string{
		"DATABASE_URL":      env.pgDSN,
		"REDIS_HOST":        env.redisHost,
		"REDIS_PORT":        env.redisPort,
		"JWT_SECRET":        "test-secret",
		"ACCESS_TOKEN_TTL":  "15m",
		"REFRESH_TOKEN_TTL": "168h",
		"USER_SERVICE_PORT": "50051",
	})

	startService(t, filepath.Join(root, "url"), map[string]string{
		"MONGO_URI":         env.mongoURI,
		"MONGO_DB":          "url_shortener",
		"REDIS_HOST":        env.redisHost,
		"REDIS_PORT":        env.redisPort,
		"USER_SERVICE_ADDR": "localhost:50051",
		"URL_SERVICE_PORT":  "50052",
		"SHORT_URL_BASE":    "https://short.url",
		"KAFKA_BROKERS":     env.kafkaBroker,
	})

	startService(t, filepath.Join(root, "gateway"), map[string]string{
		"GATEWAY_PORT":      "8080",
		"USER_SERVICE_ADDR": "localhost:50051",
		"URL_SERVICE_ADDR":  "localhost:50052",
	})

	waitForPort(t, "localhost:50051", 40*time.Second)
	waitForPort(t, "localhost:50052", 40*time.Second)
	waitForPort(t, "localhost:8080", 40*time.Second)

	time.Sleep(2 * time.Second)

	createKafkaTopic(t, env.kafkaBroker, "click_events")

	regBody := `{"email":"test@example.com","password":"password123"}`
	regResp, err := http.Post(
		"http://localhost:8080/api/v1/auth/register",
		"application/json",
		bytes.NewBufferString(regBody),
	)
	require.NoError(t, err)
	defer regResp.Body.Close()
	if regResp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(regResp.Body)
		t.Fatalf("register failed: status=%d body=%s", regResp.StatusCode, string(b))
	}

	loginResp, err := http.Post(
		"http://localhost:8080/api/v1/auth/login",
		"application/json",
		bytes.NewBufferString(regBody),
	)
	require.NoError(t, err)
	defer loginResp.Body.Close()
	if loginResp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(loginResp.Body)
		t.Fatalf("login failed: status=%d body=%s", loginResp.StatusCode, string(b))
	}

	var loginBody struct {
		AccessToken string `json:"access_token"`
	}
	require.NoError(t, json.NewDecoder(loginResp.Body).Decode(&loginBody))
	require.NotEmpty(t, loginBody.AccessToken)

	createReq, _ := http.NewRequest(
		http.MethodPost,
		"http://localhost:8080/api/v1/urls",
		bytes.NewBufferString(`{"original_url":"https://example.com"}`),
	)
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+loginBody.AccessToken)

	createResp, err := http.DefaultClient.Do(createReq)
	require.NoError(t, err)
	defer createResp.Body.Close()
	if createResp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(createResp.Body)
		t.Fatalf("create failed: status=%d body=%s", createResp.StatusCode, string(b))
	}

	var createBody struct {
		ShortCode string `json:"short_code"`
	}
	require.NoError(t, json.NewDecoder(createResp.Body).Decode(&createBody))
	require.NotEmpty(t, createBody.ShortCode)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	redirectResp, err := client.Get("http://localhost:8080/" + createBody.ShortCode)
	require.NoError(t, err)
	defer redirectResp.Body.Close()
	require.Equal(t, http.StatusFound, redirectResp.StatusCode)
	require.Equal(t, "https://example.com", redirectResp.Header.Get("Location"))

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{env.kafkaBroker},
		Topic:   "click_events",
		GroupID: "test-reader",
	})
	defer reader.Close()

	readCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	msg, err := reader.ReadMessage(readCtx)
	require.NoError(t, err)

	var event struct {
		ShortCode string `json:"short_code"`
		UserID    string `json:"user_id"`
		ClickedAt string `json:"clicked_at"`
		IP        string `json:"ip"`
		UserAgent string `json:"user_agent"`
		Referer   string `json:"referer"`
	}
	require.NoError(t, json.Unmarshal(msg.Value, &event))
	require.Equal(t, createBody.ShortCode, event.ShortCode)
	require.NotEmpty(t, event.UserID)
	require.NotEmpty(t, event.ClickedAt)
}
