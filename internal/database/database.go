package database

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Cirqach/GoStream/cmd/logger"
	_ "github.com/lib/pq"
)

const (
	host = "localhost"
	port = "5432"
)

// Video struct    allow access to queue table with video schedule
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
	logger.LogMessage(logger.GetFuncName(0), "Creating new database controller")
	return &DatabaseController{}
}

// MakeConnection    create connection between DatabaseController and database
func (dbc *DatabaseController) MakeConnection() {
	logger.LogMessage(logger.GetFuncName(0), "Connecting to database")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_USER_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")
	if user == "" && password == "" && dbname == "" {
		logger.LogError(logger.GetFuncName(0),
			fmt.Sprintf("user = \"%s\", password = \"%s\", dbname = \"%s\"", user, password, dbname))
	}
	db, err := sql.Open("postgres",
		fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, host, dbname))
	if err != nil {
		logger.LogError(logger.GetFuncName(0),
			fmt.Sprintf("Error connecting to database: %s", err))
	}
	logger.LogMessage(logger.GetFuncName(0), "Database connected by path: "+
		fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, host, dbname))
	dbc.db = db
}

// AddVideoToQueue method    add new video to Queue table
// Take video file name and date in format "2006-01-02 15:04:05"
func (dbc *DatabaseController) AddVideoToQueue(fileName, date string) error {
	logger.LogMessage(logger.GetFuncName(0), "Adding new video record to queue")
	_, err := dbc.db.Exec(fmt.Sprintf("INSERT INTO queue (path, broadcast_time) VALUES ('%s', '%s');", fileName, date))
	if err != nil {
		return err
	}
	return nil
}

// CreateUser method     create new row in User table in database
func (dbc *DatabaseController) CreateUser(name, second_name, username, email, login, password string) error {
	logger.LogMessage(logger.GetFuncName(0), "Adding new user record")
	_, err := dbc.db.Exec(
		fmt.Sprintf("INSERT INTO users (name, second_name, username, email, login, password)"+
			" VALUES ('%s', '%s', '%s', '%s', '%s', '%s');",
			name, second_name, username, email, login, password))
	if err != nil {
		return err
	}
	return nil
}

// GetSoonerVideo method    return sooner video from database
// NOTE: mayby I should return sooner video in comparison with time.Now()
func (dbc *DatabaseController) GetSoonerVideo() (Video, error) {
	logger.LogMessage(logger.GetFuncName(0), "Getting sooner video")
	videoRow, err := dbc.db.Query("SELECT * FROM queue ORDER BY broadcast_time ASC LIMIT 1;")
	if err != nil {
		logger.LogError(logger.GetFuncName(0), err.Error())
		return Video{}, err
	}
	defer videoRow.Close()
	video := Video{}

	for videoRow.Next() {
		if err = videoRow.Scan(&video.Id, &video.Path, &video.Time); err != nil {
			logger.LogError(logger.GetFuncName(0), fmt.Sprintf("Error due scanning: %s", err))
		}
	}
	logger.LogMessage(logger.GetFuncName(0), fmt.Sprintf("Sooner video: %s", video))
	return video, nil
}

// ClearQueue method    delete rows from database which older than time.Now()
func (dbc *DatabaseController) ClearQueue() error {
	logger.LogMessage(logger.GetFuncName(0), "Clearing record older than "+time.Now().String())
	t := fmt.Sprintf(strings.Split(time.Now().Local().String(), ".")[0])
	result, err := dbc.db.Exec(
		fmt.Sprintf("DELETE FROM queue WHERE broadcast_time < '%s'", t))
	if err != nil {
		logger.LogError(logger.GetFuncName(0), err.Error())
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.LogError(logger.GetFuncName(0), err.Error())
		return err
	}
	logger.LogMessage(logger.GetFuncName(0), fmt.Sprintf("deleted %d rows", rowsAffected))
	return nil
}

// TODO: check password
// VerifyUser method    check if user exists in database
func (dbc *DatabaseController) VerifyUser(username, password string) bool {
	logger.LogMessage(logger.GetFuncName(0), fmt.Sprintf("trying verify user: username: %s, password: %s", username, password))
	_, err := dbc.db.Exec(fmt.Sprintf("SELECT * FROM users WHERE username='%s';", username))
	if err != nil {
		logger.LogError(logger.GetFuncName(0), err.Error())
		return false
	}
	logger.LogMessage(logger.GetFuncName(0), "user verified")
	return true
}

func (dbc *DatabaseController) CheckTimeOverlap(wantedTime time.Time, duration string) bool {
	logger.LogMessage(logger.GetFuncName(0), "Checking time overlap")
	// Parse the video duration into hours, minutes, and seconds
	durationParts := strings.Split(duration, ":")
	if len(durationParts) != 3 {
		logger.LogError(logger.GetFuncName(0), "Invalid duration format")
		return false // Invalid duration format
	}
	durationHours, _ := strconv.Atoi(durationParts[0])
	durationMinutes, _ := strconv.Atoi(durationParts[1])
	durationSeconds, _ := strconv.Atoi(durationParts[2])
	logger.LogMessage(logger.GetFuncName(0), fmt.Sprintf("Duration: %d:%d:%d", durationHours, durationMinutes, durationSeconds))
	// Calculate the end time of the wanted broadcast
	wantedEndTime := wantedTime.Add(
		time.Duration(durationHours*60+durationMinutes)*
			(time.Duration(wantedTime.Minute())*time.Minute) +
			time.Duration(durationSeconds)*
				(time.Duration(wantedTime.Second())*time.Second))
	logger.LogMessage(logger.GetFuncName(0), fmt.Sprintf("Wanted time: %s, end time: %s", wantedTime, wantedEndTime))
	// Query the database for existing broadcasts that overlap with the wanted time
	rows, err := dbc.db.Query("SELECT broadcast_time FROM queue")
	if err != nil {
		logger.LogError(logger.GetFuncName(0), "Error querying database: "+err.Error())
		return false
	}
	defer rows.Close()

	for rows.Next() {
		var existingBroadcastTime time.Time
		err := rows.Scan(&existingBroadcastTime)
		logger.LogMessage(logger.GetFuncName(0), fmt.Sprintf("Existing broadcast time: %s", existingBroadcastTime))
		if err != nil {
			logger.LogError(logger.GetFuncName(0), "Error scanning row: "+err.Error())
			return false
		}

		// Check if the wanted time range overlaps with the existing broadcast time
		if wantedTime.Before(existingBroadcastTime.Add(time.Minute)) && wantedEndTime.After(existingBroadcastTime) {
			logger.LogMessage(logger.GetFuncName(0), fmt.Sprintf("Wanted time overlaps with existing broadcast: %s", existingBroadcastTime))
			return true // Overlap detected
		}
	}
	logger.LogMessage(logger.GetFuncName(0), "No overlap found")
	return false // No overlap found

}
