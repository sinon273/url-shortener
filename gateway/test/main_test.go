//go:build integration

package gateway_test

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tckafka "github.com/testcontainers/testcontainers-go/modules/kafka"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

type testEnv struct {
	pgDSN       string
	redisHost   string
	redisPort   string
	mongoURI    string
	kafkaBroker string
}

func setupInfra(t *testing.T) *testEnv {
	t.Helper()
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("url_shortener"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(90*time.Second)),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = pgContainer.Terminate(ctx) })

	pgHost, _ := pgContainer.Host(ctx)
	pgPort, _ := pgContainer.MappedPort(ctx, "5432")
	pgDSN := fmt.Sprintf("postgres://postgres:postgres@%s:%s/url_shortener?sslmode=disable", pgHost, pgPort.Port())

	redisContainer, err := redis.Run(ctx,
		"redis:7-alpine",
		testcontainers.WithWaitStrategy(
			wait.ForLog("Ready to accept connections").
				WithStartupTimeout(60*time.Second)),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = redisContainer.Terminate(ctx) })

	redisHost, _ := redisContainer.Host(ctx)
	redisPort, _ := redisContainer.MappedPort(ctx, "6379")

	mongoContainer, err := mongodb.Run(ctx,
		"mongo:7",
		testcontainers.WithWaitStrategy(
			wait.ForLog("Waiting for connections").
				WithStartupTimeout(90*time.Second)),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = mongoContainer.Terminate(ctx) })

	mongoURI, _ := mongoContainer.ConnectionString(ctx)

	kafkaContainer, err := tckafka.Run(ctx,
		"confluentinc/confluent-local:7.6.0",
		tckafka.WithClusterID("test-cluster"),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = kafkaContainer.Terminate(ctx) })

	kafkaBrokers, _ := kafkaContainer.Brokers(ctx)

	return &testEnv{
		pgDSN:       pgDSN,
		redisHost:   redisHost,
		redisPort:   redisPort.Port(),
		mongoURI:    mongoURI,
		kafkaBroker: kafkaBrokers[0],
	}
}

func startService(t *testing.T, dir string, env map[string]string) {
	t.Helper()
	cmd := exec.Command("go", "run", "./cmd")
	cmd.Dir = dir
	cmd.Env = os.Environ()
	for k, v := range env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	require.NoError(t, cmd.Start())

	t.Cleanup(func() {
		if cmd.Process == nil {
			return
		}
		kill := exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprint(cmd.Process.Pid))
		_ = kill.Run()
	})
}

func waitForPort(t *testing.T, addr string, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, time.Second)
		if err == nil {
			_ = conn.Close()
			return
		}
		time.Sleep(300 * time.Millisecond)
	}
	t.Fatalf("timeout waiting for %s", addr)
}

func createKafkaTopic(t *testing.T, broker, topic string) {
	t.Helper()
	conn, err := kafka.Dial("tcp", broker)
	require.NoError(t, err)
	defer conn.Close()

	controller, err := conn.Controller()
	require.NoError(t, err)

	controllerConn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
	require.NoError(t, err)
	defer controllerConn.Close()

	err = controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})
	require.NoError(t, err)
}

func repoRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs("../..")
	require.NoError(t, err)
	return root
}
