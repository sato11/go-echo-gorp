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
	"gopkg.in/go-playground/validator.v9"
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

func setupDB() (*gorp.DbMap, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	dbmap.AddTableWithName(Comment{}, "comments")
	err = dbmap.CreateTablesIfNotExists()
	if err != nil {
		return nil, err
	}

	return dbmap, nil
}

func setupEcho() *echo.Echo {
	e := echo.New()
	e.Validator = &Validator{validator: validator.New()}
	return e
}

// PreInsert updates timestamp columns before insert
func (c *Comment) PreInsert(s gorp.SqlExecutor) error {
	c.Created = time.Now()
	c.Updated = c.Created
	return nil
}

// PreUpdate updates Updated column before update
func (c *Comment) PreUpdate(s gorp.SqlExecutor) error {
	c.Updated = time.Now()
	return nil
}

// Controller groups api functions by route
type Controller struct {
	dbmap *gorp.DbMap
}

// ListComments returns an array of comments
func (controller *Controller) ListComments(c echo.Context) error {
	var comments []Comment
	_, err := controller.dbmap.Select(&comments, "SELECT * FROM comments ORDER BY created DESC LIMIT 10")
	if err != nil {
		c.Logger().Error("Select: ", err)
		return c.String(http.StatusBadRequest, "Select: "+err.Error())
	}
	return c.JSON(http.StatusOK, comments)
}

// InsertComment creates a record in table and returns empty string
func (controller *Controller) InsertComment(c echo.Context) error {
	var comment Comment
	if err := c.Bind(&comment); err != nil {
		c.Logger().Error("Bind: ", err)
		return c.String(http.StatusBadRequest, "Bind "+err.Error())
	}
	if err := c.Validate(&comment); err != nil {
		c.Logger().Error("Validate: ", err)
		return c.String(http.StatusBadRequest, "Validate "+err.Error())
	}
	if err := controller.dbmap.Insert(&comment); err != nil {
		c.Logger().Error("Insert: ", err)
		return c.String(http.StatusBadRequest, "Insert: "+err.Error())
	}
	c.Logger().Infof("ADDED: %v", comment.ID)
	return c.JSON(http.StatusCreated, "")
}

// Validator represents validator
type Validator struct {
	validator *validator.Validate
}

// Validate validates parameters by checking struct tags
func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func main() {
	dbmap, err := setupDB()
	if err != nil {
		log.Fatal(err)
	}

	controller := &Controller{dbmap: dbmap}

	e := setupEcho()
	e.GET("/api/comments", controller.ListComments)
	e.POST("/api/comments", controller.InsertComment)
	e.Static("/", "static/")
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
