package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

const (
	host = "localhost"
	port = "5432"
)

type video struct {
	id   string
	path string
	time string
}

type DatabaseController struct {
	db *sql.DB
}

func NewDatabaseController() *DatabaseController {
	return &DatabaseController{}
}

// TODO: change database user, right now user dont have password
func (dbc *DatabaseController) MakeConnection() {
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")
	// psqlInfo := fmt.Sprintf("host=%s port=%s user=%s"+
	// 	" password=%s dbname=%s sslmode=disable",
	// 	host, port, user, password, dbname)
	// log.Printf("psqlInfo = \"%s\"", psqlInfo)
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, host, dbname))
	if err != nil {
		log.Printf("Error connecting to database: ")
	}
	log.Println("Database connected by path: " + fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, host, dbname))
	dbc.db = db
}

func (dbc *DatabaseController) AddVideoToQueue(path, data string) error {
	_, err := dbc.db.Exec(fmt.Sprintf("INSERT INTO queue VALUES (%s, %s);"), path, data)
	if err != nil {
		return err
	}
	return nil
}

func (dbc *DatabaseController) CreateUser(name, second_name, username, email, login, password string) error {
	_, err := dbc.db.Exec(fmt.Sprintf("INSERT INTO users VALUES (%s, %s, %s, %s, %s, %s);",
		name, second_name, username, email, login, password))
	if err != nil {
		return err
	}
	return nil
}

func (dbc *DatabaseController) GetSoonerVideo() (video, error) {
	videoRow, err := dbc.db.Query("SELECT * FROM queue ORDER BY broadcast_time ASC LIMIT 1;")
	if err != nil {
		fmt.Printf("Error in getting sooner video: %s", videoRow)
		return video{}, err
	}
	defer videoRow.Close()
	video := video{}
	for videoRow.Next() {
		err = videoRow.Scan(&video.id, &video.path, &video.time)
		if err != nil {
			log.Printf("Error due scanning: %s", err)
		}
	}
	return video, nil
}
