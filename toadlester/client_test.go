package toadlester_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	tl "github.com/maroda/terraform-provider-toadlester/toadlester"
)

func TestNewAPIClient(t *testing.T) {
	url := "http://test"
	want := url
	got := tl.NewAPIClient(url)
	if got.BaseURL != want {
		t.Errorf("got %T, want %q", got, want)
	}
}

func TestAPIClient_ReadType(t *testing.T) {
	t.Run("Reads API endpoint", func(t *testing.T) {
		returns := "Metric_int_up: 42"
		wwwServ := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := fmt.Fprintln(w, returns)
			assertError(t, err, nil)
		}))
		defer wwwServ.Close()

		want := returns + "\n"
		client := tl.NewAPIClient(wwwServ.URL)
		set := &tl.Setting{
			Name:  "INT_SIZE",
			Value: "100",
			Algo:  "up",
		}
		got, err := client.ReadType(set)
		assertError(t, err, nil)
		assertStringContains(t, got, want)
	})
}

func TestAPIClient_CRUD(t *testing.T) {
	returns := `Set new up value INT_SIZE for 100
Set new down value INT_SIZE for 100
`
	wwwServ := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintln(w, returns)
		assertError(t, err, nil)
	}))
	defer wwwServ.Close()

	want := returns + "\n"
	client := tl.NewAPIClient(wwwServ.URL)
	set := &tl.Setting{
		Name:  "INT_SIZE",
		Value: "100",
		Algo:  "up",
	}

	t.Run("Create endpoint works", func(t *testing.T) {
		got, err := client.CreateType(set)
		assertError(t, err, nil)
		assertStringContains(t, got, want)
	})

	t.Run("Update endpoint works", func(t *testing.T) {
		got, err := client.UpdateType(set)
		assertError(t, err, nil)
		assertStringContains(t, got, want)
	})

	t.Run("Delete endpoint works", func(t *testing.T) {
		got, err := client.DeleteType(set)
		assertError(t, err, nil)
		assertStringContains(t, got, want)
	})
}

func TestAPIClient_ReadTypeLC(t *testing.T) {
	t.Run("Local Container: Reads API endpoint", func(t *testing.T) {
		want := "INT_SIZE"
		stop, url := makeLocalToadLester(t)
		defer stop()

		client := tl.NewAPIClient(url)
		set := &tl.Setting{}
		got, err := client.ReadType(set)
		assertError(t, err, nil)
		assertStringContains(t, got, want)
	})
}

// Helpers //

// Start up a test container of the API for integration testing
func makeLocalToadLester(t *testing.T) (func(), string) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "ghcr.io/maroda/toadlester:latest",
		ExposedPorts: []string{"8899/tcp"},
		WaitingFor: wait.ForHTTP("/metrics").
			WithPort("8899/tcp").
			WithStartupTimeout(30 * time.Second),
	}
	container, err := testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
	if err != nil {
		t.Fatalf("could not start test container: %v", err)
	}

	mappedPort, err := container.MappedPort(ctx, "8899")
	if err != nil {
		t.Fatalf("could not map port to container: %v", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("could not get container host: %v", err)
	}

	endpoint := fmt.Sprintf("http://%s:%s", host, mappedPort.Port())
	resp, err := http.Get(endpoint + "/metrics")
	if err != nil {
		t.Fatalf("could not get metrics from container: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	assertStatus(t, resp.StatusCode, http.StatusOK)

	end := func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("could not stop test container: %v", err)
		}
	}

	fmt.Printf("Started test container with endpoint %s\n", endpoint)

	return end, endpoint
}

// Assertions //

func assertError(t testing.TB, got, want error) {
	t.Helper()
	if !errors.Is(got, want) {
		t.Errorf("got error %q want %q", got, want)
	}
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertStringContains(t *testing.T, full, want string) {
	t.Helper()
	if !strings.Contains(full, want) {
		t.Errorf("Did not find %q, expected string contains %q", want, full)
	}
}

/*
func assertGotError(t testing.TB, got error) {
	t.Helper()
	if got == nil {
		t.Errorf("Expected an error but got %q", got)
	}
}

func assertInt(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct value, got %d, want %d", got, want)
	}
}
*/
