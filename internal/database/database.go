package database

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Cirqach/GoStream/cmd/logger"
	_ "github.com/lib/pq"
)

const (
	host = "localhost"
	port = "5432"
)

type Video struct {
	Id   string
	Path string
	Time string
}

// DatabaseController struct    allow access to dabase
type DatabaseController struct {
	db *sql.DB
}

// NewDatabaseController function    create new DataBaseController obj
func NewDatabaseController() *DatabaseController {
	return &DatabaseController{}
}

// MakeConnection    create connection between DatabaseController and database
func (dbc *DatabaseController) MakeConnection() {
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")
	if user == "" && password == "" && dbname == "" {
		logger.LogError("MakeConnection",
			fmt.Sprintf("user = \"%s\", password = \"%s\", dbname = \"%s\"", user, password, dbname))
	}
	db, err := sql.Open("postgres",
		fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, host, dbname))
	if err != nil {
		logger.LogError("MakeConnection",
			fmt.Sprintf("Error connecting to database: %s", err))
	}
	logger.LogMessage("MakeConnection", "Database connected by path: "+
		fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, host, dbname))
	dbc.db = db
}

// AddVideoToQueue method    add new video to Queue table
func (dbc *DatabaseController) AddVideoToQueue(path, data string) error {
	_, err := dbc.db.Exec(fmt.Sprintf("INSERT INTO queue VALUES ('%s', '%s');", path, data))
	if err != nil {
		return err
	}
	return nil
}

// CreateUser method     create new row in User table in database
func (dbc *DatabaseController) CreateUser(name, second_name, username, email, login, password string) error {
	_, err := dbc.db.Exec(fmt.Sprintf("INSERT INTO users (name, second_name, username, email, login, password) VALUES ('%s', '%s', '%s', '%s', '%s', '%s');",
		name, second_name, username, email, login, password))
	if err != nil {
		return err
	}
	return nil
}

// NOTE: mayby I should return sooner video in comparison with time.Now()
// GetSoonerVideo method    return sooner video from database
func (dbc *DatabaseController) GetSoonerVideo() (Video, error) {
	videoRow, err := dbc.db.Query("SELECT * FROM queue ORDER BY broadcast_time ASC LIMIT 1;")
	if err != nil {
		logger.LogError("GetSoonerVideo",
			fmt.Sprintf("Error in getting sooner video: %s", videoRow))
		return Video{}, err
	}
	defer videoRow.Close()
	video := Video{}

	for videoRow.Next() {
		if err = videoRow.Scan(&video.Id, &video.Path, &video.Time); err != nil {
			logger.LogError("GetSoonerVideo", fmt.Sprintf("Error due scanning: %s", err))
		}
	}
	return video, nil
}

// ClearQueue method    delete rows from database which older than time.Now()
func (dbc *DatabaseController) ClearQueue() error {
	t := fmt.Sprintf(strings.Split(time.Now().Local().String(), ".")[0])
	_, err := dbc.db.Exec(
		fmt.Sprintf("DELETE FROM queue WHERE broadcast_time < '%s'", t))
	if err != nil {
		return err
	}
	return nil

}
