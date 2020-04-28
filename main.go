package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-gorp/gorp"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
)

// Comment represents schema of the comments table
type Comment struct {
	ID      int64     `json:"id" db:"id,primarykey,autoincrement"`
	Name    string    `json:"name" db:"name,notnull,default:'noname',size:200"`
	Text    string    `json:"text" db:"text,notnull,size:399"`
	Created time.Time `json:"created" db:"created,notnull"`
	Updated time.Time `json:"updated" db:"updated,notnull"`
}

var dsn = os.Getenv("DSN")
var port = os.Getenv("PORT")

func main() {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	dbmap.AddTableWithName(Comment{}, "comments")
	err = dbmap.CreateTablesIfNotExists()
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.GET("/api/comments", func(c echo.Context) error {
		var comments []Comment
		_, err := dbmap.Select(&comments, "SELECT * FROM comments ORDER BY created DESC LIMIT 10")
		if err != nil {
			c.Logger().Error("Select: ", err)
			return c.String(http.StatusBadRequest, "Select: "+err.Error())
		}
		return c.JSON(http.StatusOK, comments)
	})
	e.POST("/api/comments", func(c echo.Context) error {
		var comment Comment
		if err = c.Bind(&comment); err != nil {
			c.Logger().Error("Bind: ", err)
			return c.String(http.StatusBadRequest, "Bind "+err.Error())
		}
		if err = dbmap.Insert(&comment); err != nil {
			c.Logger().Error("Insert: ", err)
			return c.String(http.StatusBadRequest, "Insert: "+err.Error())
		}
		c.Logger().Infof("ADDED: %v", comment.ID)
		return c.JSON(http.StatusCreated, "")
	})
	e.Static("/", "static/")
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
