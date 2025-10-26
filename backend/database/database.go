package database

import (
	"errors"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/kube-dash/kube-dash-backend/models"
)

func InitDB() (*gorm.DB, error) {

	// open local sqlite database
	// TODO: change to more powerful database
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// migrate the schema and create tables
	db.AutoMigrate(
		&models.DBContainerMetricsModel{},
		&models.DBPodMetricsModel{},
		&models.DBClusterMetricsModel{},
	)

	return db, nil

}

// remove pod metrics stored in the database
func DBDeletePodMetrics(
	db *gorm.DB,
	starttime *time.Time, endtime *time.Time,
) error {

	tx := db.Begin() // start a transaction

	var whereClause string
	var args []interface{}
	if starttime != nil && endtime != nil {

		// validate if starttime is not after than endtime
		if starttime.After(*endtime) {
			return errors.New("start_time cannot be after end_time")
		}

		// delete records created within the specified time range
		whereClause = "created_at BETWEEN ? AND ?"
		args = []interface{}{starttime, endtime}

	} else if starttime != nil {
		// delete records created after the specified start time
		whereClause = "created_at >= ?"
		args = []interface{}{starttime}
	} else if endtime != nil {
		// delete records created before the specified end time
		whereClause = "created_at <= ?"
		args = []interface{}{endtime}
	} else {
		// no time range specified, delete all records
		whereClause = "1=1"
	}

	// define the models to be deleted
	modelRows := []interface{}{
		&models.DBContainerMetricsModel{},
		&models.DBPodMetricsModel{},
		&models.DBClusterMetricsModel{},
	}

	// delete records from each model using the set query
	for _, mr := range modelRows {

		// hard-delete records of the current model
		err := tx.Unscoped().Where(whereClause, args...).Delete(mr).Error
		if err != nil {
			// rollback on error
			tx.Rollback()
			return err
		}
	}

	// commit the transaction if everything was successful
	if err := tx.Commit().Error; err != nil {
		return err
	}

	// run a VACUUM command to clean up the database
	db.Exec("VACUUM")

	return nil
}

func cleanPodMetricsDB(db *gorm.DB) error {
	// remove everything before 168 hours ago (7 days)
	endTime := time.Now().Add(-168 * time.Hour)
	return DBDeletePodMetrics(db, nil, &endTime)
}

// start a job that removes old metrics records from the database
// and runs every 1 hour
func StartDBPodMetricsCleaner(
	db *gorm.DB,
) error {

	ticker := time.NewTicker(3600 * time.Second)
	go func() {
		for range ticker.C {
			err := cleanPodMetricsDB(db)
			if err != nil {
				// TODO: find a way to error handle the ticker
				return
			}
		}
	}()

	return nil
}
