package dbRepo

import (
	"database/sql"
	"errors"
	"goTemp/model"
)

type TempDb interface {
	CreateValue(v *model.Value) error
	GetValue(id int) ([]*model.Value, error)
	ListValues() ([]*model.Value, error)
}

type PostgesTempDb struct {
	Db *sql.DB
}

func (t *PostgesTempDb) ListValue() ([]*model.Value, error) {
	rows, err := t.Db.Query("SELECT * FROM value")

	if err != nil {
		return nil, errors.New("Error load data from database")
	}

	valueList := make([]*model.Value, 0)

	var index int
	var hum string
	var temp string

	if rows != nil {

		for rows.Next() {
			err := rows.Scan(&index, &temp, &hum)

			if err != nil {
				return nil, errors.New("Error load data from database")
			}

			valueList = append(valueList, &model.Value{index, temp, hum})
		}

		err = rows.Err()

		if err != nil {
			return nil, errors.New("Error load data from database")
		}
	}

	return valueList, nil
}

func (t *PostgesTempDb) GetValue(id int) ([]*model.Value, error) {

	rows, err := t.Db.Query("SELECT * FROM value WHERE sensorid=$1", id)
	valueList := make([]*model.Value, 0)

	var index int
	var hum string
	var temp string

	if rows != nil {

		for rows.Next() {
			err := rows.Scan(&index, &temp, &hum)

			if err != nil {
				return nil, errors.New("Error load data from database")
			}

			valueList = append(valueList, &model.Value{index, temp, hum})
		}

		err = rows.Err()

		if err != nil {
			return nil, errors.New("Error load data from database")
		}

	}

	return valueList, err
}

func (t *PostgesTempDb) CreateValue(v *model.Value) error {

	_, err := t.Db.Exec("CREATE TABLE IF NOT EXISTS value(SensorId INT, Temp VARCHAR(50), Hum VARCHAR(50))")

	if err != nil {
		return errors.New("Error creating new table for value")
	}

	var sensorId int

	err = t.Db.QueryRow(`INSERT INTO value
	VALUES($1, $2, $3) RETURNING SensorId`, v.Index, v.Temp, v.Hum).Scan(&sensorId)

	if err != nil {
		return errors.New("Error writing data to database")
	} else {
		return nil
	}
}
