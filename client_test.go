package serpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_Search(t *testing.T) {
	apiKey := "test_key"
	mockResponse := map[string]interface{}{
		"search_metadata": map[string]interface{}{"status": "Success"},
		"organic_results": []interface{}{
			map[string]interface{}{"title": "Result 1"},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/search" {
			t.Errorf("expected path /search, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("api_key") != apiKey {
			t.Errorf("expected api_key %s, got %s", apiKey, r.URL.Query().Get("api_key"))
		}
		if r.Header.Get("User-Agent") == "" {
			t.Errorf("expected User-Agent header")
		}
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient(apiKey, WithBaseURL(server.URL))
	results, err := client.Search(context.Background(), map[string]string{"q": "test"})
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if results["search_metadata"].(map[string]interface{})["status"] != "Success" {
		t.Errorf("expected status Success")
	}
}

func TestClient_PaginationHelper(t *testing.T) {
	resp := Response{
		"serpapi_pagination": map[string]interface{}{
			"next": "https://serpapi.com/search.json?engine=google&q=coffee&start=10",
		},
	}

	params := resp.NextPageParams()
	if params["start"] != "10" {
		t.Errorf("expected start=10, got %s", params["start"])
	}
	if params["engine"] != "google" {
		t.Errorf("expected engine=google, got %s", params["engine"])
	}
}

func TestClient_ErrorHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Missing query"}`))
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	_, err := client.Search(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expectedErr := "serpapi error (400): Missing query"
	if err.Error() != expectedErr {
		t.Errorf("expected error %q, got %q", expectedErr, err.Error())
	}
}
