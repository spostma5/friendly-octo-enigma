package risk

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

// This is not at all comprehensive, but just an example of how testing could be done
// If this was production, I'd either set a lot of these cases up in smaller tests,
// or make a map of "test runs" that get run in each case with a path, value(s), expected status,
// and expected output to cut down on duplication
// Would also need a lot more tests for happy/sad paths and edge cases, as this is
// pretty bare bones for now

func TestGetRisks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(HandleGetRisks))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Incorrect status: expected StatusOk - got %d", resp.StatusCode)
	}

	expectedType := "application/json"
	if typ := resp.Header.Get("Content-Type"); typ != expectedType {
		t.Errorf("Incorrect type: expected %s - got %s", expectedType, typ)
	}

	expected := "{}\n"

	bt, err := io.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		t.Error(err)
	}

	if string(bt) != expected {
		t.Errorf("Incorrect body: expected %s - got %s", expected, string(bt))
	}

}

// Wasn't initially planning on using a response recorder here, but ran
// into issues with go 1.22s path values.
// Not entirely sure why, but it didn't like it being just added to the URL
// and instead I had to manually set the value in the request
func TestGetRisk(t *testing.T) {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Error(err)
	}

	// Invalid id
	req.SetPathValue("id", "1")
	HandleGetRisk(rr, req)
	defer rr.Result().Body.Close()

	if rr.Result().StatusCode != http.StatusNotFound {
		t.Errorf("Incorrect status: expected StatusNotFound - got %d", rr.Result().StatusCode)
	}

	expected := "Unable to find requested risk\n"

	bt, err := io.ReadAll(rr.Result().Body)

	if err != nil {
		t.Error(err)
	}

	if string(bt) != expected {
		t.Errorf("Incorrect body: expected %s - got %s", expected, string(bt))
	}

	// Matching id
	tRisk := Risk{
		State: "open",
		Title: "Test-risk",
	}

	createRisk(&tRisk)

	rr = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Error(err)
	}
	req.SetPathValue("id", tRisk.Id)

	HandleGetRisk(rr, req)
	defer rr.Result().Body.Close()

	if rr.Result().StatusCode != http.StatusOK {
		t.Errorf("Incorrect status: expected StatusOK - got %d", rr.Result().StatusCode)
	}

	// Could use the actual uuid regex, but this is easier to read and fine for demo purposes
	expected = `\{"id":".*","state":"open","title":"Test-risk","description":""\}\n`

	bt, err = io.ReadAll(rr.Result().Body)

	if err != nil {
		t.Error(err)
	}

	matched, _ := regexp.MatchString(expected, string(bt))

	if !matched {
		t.Errorf("Incorrect body: expected %s - got %s", expected, string(bt))
	}
}

func TestPostRisk(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(HandlePostRisk))
	defer server.Close()

	slog.SetLogLoggerLevel(slog.LevelDebug)

	// No payload
	resp, err := http.Post(server.URL, "application/json", nil)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Incorrect status: expected StatusBadRequest - got %d", resp.StatusCode)
	}

	expected := "Invalid JSON\n"

	bt, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	if string(bt) != expected {
		t.Errorf("Incorrect body: expected %s - got %s", expected, string(bt))
	}

	// Correct payload
	tRisk := Risk{
		State: "open",
		Title: "Test-risk",
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(tRisk)

	resp, err = http.Post(server.URL, "application/json", b)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Incorrect status: expected StatusCreated - got %d", resp.StatusCode)
	}

	expectedType := "application/json"
	if typ := resp.Header.Get("Content-Type"); typ != expectedType {
		t.Errorf("Incorrect type: expected %s - got %s", expectedType, typ)
	}

	expected = `\{"id":".*","state":"open","title":"Test-risk","description":""\}\n`

	bt, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	matched, _ := regexp.MatchString(expected, string(bt))

	if !matched {
		t.Errorf("Incorrect body: expected to match regex %s - got %s", expected, string(bt))
	}

	// Invalid payload
	tRisk = Risk{
		State: "openandshut",
		Title: "Test-risk2",
	}

	b = new(bytes.Buffer)
	json.NewEncoder(b).Encode(tRisk)

	resp, err = http.Post(server.URL, "application/json", b)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Incorrect status: expected StatusBadRequest - got %d", resp.StatusCode)
	}

	expected = "Validation error: Key: 'Risk.State' Error:Field validation for 'State' failed on the 'oneof' tag\n"

	bt, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	if string(bt) != expected {
		t.Errorf("Incorrect body: expected %s - got %s", expected, string(bt))
	}
}
