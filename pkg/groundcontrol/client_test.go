package groundcontrol_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/robherley/groundcontrol-go/pkg/groundcontrol"
)

const (
	testProjectID = "test-project-id"
	testAPIKey    = "gcp_test-api-key"
	enabledFlag   = "enabled-flag"
	disabledFlag  = "disabled-flag"
	enabledActor  = groundcontrol.Actor("enabled-actor")
	disabledActor = groundcontrol.Actor("disabled-actor")
)

func flagCheckPath(projectID, flag string) string {
	return fmt.Sprintf("/projects/%s/flags/%s/check", projectID, flag)
}

func flagCheck(enabled bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"enabled": %t}`, enabled)
	}
}

func TestIsFeatureFlagEnabled(t *testing.T) {
	t.Run("non-200 response", func(t *testing.T) {
		testServer := httptest.NewServer(http.NotFoundHandler())
		defer testServer.Close()

		client := groundcontrol.New(testProjectID, testAPIKey, groundcontrol.WithBaseURL(testServer.URL))

		isEnabled, err := client.IsFeatureFlagEnabled(context.Background(), "this-flag-does-not-exist")
		if err == nil {
			t.Fatal("expected error")
		}

		if isEnabled {
			t.Fatal("expected flag to be disabled")
		}
	})

	t.Run("everything enabled", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(flagCheck(true)))
		defer testServer.Close()

		client := groundcontrol.New(testProjectID, testAPIKey, groundcontrol.WithBaseURL(testServer.URL))

		isEnabled, err := client.IsFeatureFlagEnabled(context.Background(), enabledFlag)
		if err != nil {
			t.Fatal(err)
		}

		if !isEnabled {
			t.Fatal("expected flag to be enabled")
		}
	})

	t.Run("everything disabled", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(flagCheck(false)))
		defer testServer.Close()

		client := groundcontrol.New(testProjectID, testAPIKey, groundcontrol.WithBaseURL(testServer.URL))

		isEnabled, err := client.IsFeatureFlagEnabled(context.Background(), disabledFlag)
		if err != nil {
			t.Fatal(err)
		}

		if isEnabled {
			t.Fatal("expected flag to be disabled")
		}
	})

	t.Run("specific flag enabled", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			flagCheck(r.URL.Path == flagCheckPath(testProjectID, enabledFlag))(w, r)
		}))

		defer testServer.Close()

		client := groundcontrol.New(testProjectID, testAPIKey, groundcontrol.WithBaseURL(testServer.URL))

		isEnabled, err := client.IsFeatureFlagEnabled(context.Background(), enabledFlag)
		if err != nil {
			t.Fatal(err)
		}

		if !isEnabled {
			t.Fatal("expected flag to be enabled")
		}

		isEnabled, err = client.IsFeatureFlagEnabled(context.Background(), disabledFlag)
		if err != nil {
			t.Fatal(err)
		}

		if isEnabled {
			t.Fatal("expected flag to be disabled")
		}
	})

	t.Run("flag enabled for actors", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != flagCheckPath(testProjectID, enabledFlag) {
				flagCheck(false)(w, r)
				return
			}

			query := r.URL.Query()
			actors := query["actorIds"]

			for _, actor := range actors {
				if actor == enabledActor.Identifier() {
					flagCheck(true)(w, r)
					return
				}
			}

			flagCheck(false)(w, r)
		}))

		defer testServer.Close()

		client := groundcontrol.New(testProjectID, testAPIKey, groundcontrol.WithBaseURL(testServer.URL))

		// single actor, enabled == enabled
		isEnabled, err := client.IsFeatureFlagEnabled(context.Background(), enabledFlag, enabledActor)
		if err != nil {
			t.Fatal(err)
		}

		if !isEnabled {
			t.Fatal("expected flag to be enabled for actor")
		}

		// single actor, disabled == disabled
		isEnabled, err = client.IsFeatureFlagEnabled(context.Background(), enabledFlag, disabledActor)
		if err != nil {
			t.Fatal(err)
		}

		if isEnabled {
			t.Fatal("expected flag to be disabled for actor")
		}

		// multiple actors, mixed enablement == enabled
		isEnabled, err = client.IsFeatureFlagEnabled(context.Background(), enabledFlag, disabledActor, enabledActor)
		if err != nil {
			t.Fatal(err)
		}

		if !isEnabled {
			t.Fatal("expected flag to be enabled for actors")
		}
	})
}
