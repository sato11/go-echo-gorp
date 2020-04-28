package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-gorp/gorp"

	"github.com/labstack/echo"
)

func truncateTable(dbmap *gorp.DbMap) error {
	err := dbmap.TruncateTables()
	if err != nil {
		return err
	}
	return nil
}

func init() {
	dsn = "host=db user=postgres dbname=postgres password=password sslmode=disable"
}

func TestMain(m *testing.M) {
	exitCode := m.Run()

	dbmap, err := setupDB()
	if err != nil {
		log.Fatal(err)
	}

	err = truncateTable(dbmap)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(exitCode)
}

func TestListComments(t *testing.T) {
	dbmap, err := setupDB()
	if err != nil {
		log.Fatal(err)
	}

	controller := &Controller{dbmap: dbmap}

	req := httptest.NewRequest(http.MethodGet, "/api/comments", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := setupEcho()
	c := e.NewContext(req, rec)

	err = controller.ListComments(c)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Code != 200 {
		t.Fatalf("status code should be 200, got: %d", rec.Code)
	}

	var comments []Comment
	err = json.NewDecoder(rec.Body).Decode(&comments)
	if err != nil {
		t.Fatal(err)
	}
	if len(comments) > 0 {
		t.Fatalf("len(comments) should be empty, got: %d", len(comments))
	}
}

func TestInsertComment(t *testing.T) {
	dbmap, err := setupDB()
	if err != nil {
		log.Fatal(err)
	}

	controller := &Controller{dbmap: dbmap}

	req := httptest.NewRequest(http.MethodPost, "/api/comments", strings.NewReader(`
{
	"name": "sato11",
	"text": "Hello World!"
}
	`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := setupEcho()
	c := e.NewContext(req, rec)

	err = controller.InsertComment(c)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Code != 201 {
		t.Fatalf("code should be 201, got: %d", rec.Code)
	}
}
