//go:build integration

package gateway_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"testing"
	"time"
)

func TestDeleteURL(t *testing.T) {
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

	body := `{"email":"del@example.com","password":"password123"}`
	resp, err := http.Post("http://localhost:8080/api/v1/auth/register", "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("register failed: status=%d body=%s", resp.StatusCode, string(b))
	}

	loginResp, err := http.Post("http://localhost:8080/api/v1/auth/login", "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatal(err)
	}
	defer loginResp.Body.Close()
	if loginResp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(loginResp.Body)
		t.Fatalf("login failed: status=%d body=%s", loginResp.StatusCode, string(b))
	}

	var login struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(loginResp.Body).Decode(&login); err != nil {
		t.Fatal(err)
	}

	createReq, _ := http.NewRequest("POST", "http://localhost:8080/api/v1/urls",
		bytes.NewBufferString(`{"original_url":"https://example.com"}`))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+login.AccessToken)

	createResp, err := http.DefaultClient.Do(createReq)
	if err != nil {
		t.Fatal(err)
	}
	defer createResp.Body.Close()
	if createResp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(createResp.Body)
		t.Fatalf("create failed: status=%d body=%s", createResp.StatusCode, string(b))
	}

	var created struct {
		ShortCode string `json:"short_code"`
	}
	if err := json.NewDecoder(createResp.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}

	delReq, _ := http.NewRequest("DELETE", "http://localhost:8080/api/v1/urls/"+created.ShortCode, nil)
	delReq.Header.Set("Authorization", "Bearer "+login.AccessToken)

	delResp, err := http.DefaultClient.Do(delReq)
	if err != nil {
		t.Fatal(err)
	}
	defer delResp.Body.Close()
	if delResp.StatusCode != http.StatusNoContent {
		b, _ := io.ReadAll(delResp.Body)
		t.Fatalf("delete failed: status=%d body=%s", delResp.StatusCode, string(b))
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	afterDelResp, err := client.Get("http://localhost:8080/" + created.ShortCode)
	if err != nil {
		t.Fatal(err)
	}
	defer afterDelResp.Body.Close()
	if afterDelResp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", afterDelResp.StatusCode)
	}
}
