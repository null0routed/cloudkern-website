package main

import (
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// newTestServer builds the same router main() serves, behind an httptest server.
func newTestServer() *httptest.Server {
	app := &app{
		Templates: template.Must(template.ParseFS(TmplFiles, "ui/tmpl/*.html")),
	}
	return httptest.NewServer(newRouter(app))
}

func get(t *testing.T, ts *httptest.Server, path string) (int, string) {
	t.Helper()

	resp, err := http.Get(ts.URL + path)
	if err != nil {
		t.Fatalf("GET %s: %v", path, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading body of GET %s: %v", path, err)
	}
	return resp.StatusCode, string(body)
}

func TestIndexServesHomepage(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	status, body := get(t, ts, "/")
	if status != http.StatusOK {
		t.Fatalf("GET / status = %d, want %d", status, http.StatusOK)
	}
	if !strings.Contains(body, "CloudKern") {
		t.Fatal("GET / body does not contain site wordmark")
	}
}

func TestFaviconIsServed(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	status, body := get(t, ts, "/favicon.ico")
	if status != http.StatusOK {
		t.Fatalf("GET /favicon.ico status = %d, want %d", status, http.StatusOK)
	}
	if body == "" {
		t.Fatal("GET /favicon.ico returned an empty body")
	}
}

func TestStaticAssetIsServed(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	status, body := get(t, ts, "/static/css/site.css")
	if status != http.StatusOK {
		t.Fatalf("GET /static/css/site.css status = %d, want %d", status, http.StatusOK)
	}
	if !strings.Contains(body, "--bg") {
		t.Fatal("GET /static/css/site.css body does not look like the stylesheet")
	}
}

func TestStaticAssetMissingIsNotFound(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	status, _ := get(t, ts, "/static/css/missing.css")
	if status != http.StatusNotFound {
		t.Fatalf("GET /static/css/missing.css status = %d, want %d", status, http.StatusNotFound)
	}
}

func TestHealthcheck(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	status, _ := get(t, ts, "/healthcheck")
	if status != http.StatusOK {
		t.Fatalf("GET /healthcheck status = %d, want %d", status, http.StatusOK)
	}
}
