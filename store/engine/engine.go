package engine

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/mahingaRodin/carzone-go/models"
)


type EngineStore struct {
	db *sql.DB
}

func New(db *sql.DB) *EngineStore {
	return &EngineStore{db: db}
}

func (e EngineStore) EngineById(ctx context.Context, id string) (models.Engine, error) {
	var engine models.Engine

	tx,err := e.db.BeginTx(ctx,nil)
	if err != nil {
		return engine,nil
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("Transaction rollback error: %v\n", rbErr)
			}
		} else {
			if cmErr := tx.Commit(); cmErr != nil {
				fmt.Printf("Transaction commit error: %v\n", cmErr)
			}
		}
	}()

	err = tx.QueryRowContext(ctx,"SELECT id, type, horsepower, torque, created_at, updated_at FROM engine WHERE id=$1", id).Scan(
		&engine.EngineID,&engine.Displacement,&engine.NoOfCyclinders,&engine.CarRange,)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return engine,nil
		}
		return engine,err
	}
	return engine,err
}

func (e EngineStore) CreateEngine(ctx context.Context, engineReq *models.EngineRequest) (models.Engine,error) {
	tx,err := e.db.BeginTx(ctx,nil)
	if err != nil {
		return models.Engine{},err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("Transaction rollback error: %v\n", rbErr)
			}
		} else {
			if cmErr := tx.Commit(); cmErr != nil {
				fmt.Printf("Transaction commit error: %v\n", cmErr)
			}
		}
	}()

	engineID := uuid.New()

	_,err = tx.ExecContext(ctx,
		"INSERT INTO engine (id, displacement, no_of_cylinders, car_range) VALUES ($1, $2, $3, $4)",engineID, engineReq.Displacement, engineReq.NoOfCyclinders, engineReq.CarRange)

	if err != nil {
		return models.Engine{},err
	}
	engine := models.Engine{
		EngineID: engineID,
		Displacement: engineReq.Displacement,
		NoOfCyclinders: engineReq.NoOfCyclinders,
		CarRange: engineReq.CarRange,
	}
	return engine,nil
	
}

func (e EngineStore) EngineUpdate(ctx context.Context, id string, engineReq *models.EngineRequest) (models.Engine, error) {
	engineID, err := uuid.Parse(id)
	if err != nil {
		return models.Engine{}, err
	}

	tx,err := e.db.BeginTx(ctx,nil)
	if err != nil {
		return models.Engine{},err
	}

		defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("Transaction rollback error: %v\n", rbErr)
			}
		} else {
			if cmErr := tx.Commit(); cmErr != nil {
				fmt.Printf("Transaction commit error: %v\n", cmErr)
			}
		}
	}()

	results, err := tx.ExecContext(ctx,
	"UPDATE engine SET displacement=$1, no_of_cylinders=$2, car_range=$3 WHERE id=$4",
	engineReq.Displacement,
	engineReq.NoOfCyclinders,
	engineReq.CarRange,
	engineID,
	)

	if err != nil {
		return models.Engine{}, err
	}

	    rowsAffected, err := results.RowsAffected()
    if err != nil {
        return models.Engine{}, err
    }

    if rowsAffected == 0 {
        return models.Engine{}, errors.New("no rows were deleted")
    }

	engine := models.Engine{
		EngineID: engineID,
		Displacement: engineReq.Displacement,
		NoOfCyclinders: engineReq.NoOfCyclinders,
		CarRange: engineReq.CarRange,
	}

	return engine,nil

}

func (e EngineStore) EngineDelete(ctx context.Context, id string) (models.Engine, error) {
    var engine models.Engine
    var err error

    engineID, err := uuid.Parse(id)
    if err != nil {
        return models.Engine{}, err
    }

    tx, err := e.db.BeginTx(ctx, nil)
    if err != nil {
        return models.Engine{}, err
    }

    defer func() {
        if err != nil {
            if rbErr := tx.Rollback(); rbErr != nil {
                fmt.Printf("Transaction rollback error: %v\n", rbErr)
            }
        } else {
            if cmErr := tx.Commit(); cmErr != nil {
                fmt.Printf("Transaction commit error: %v\n", cmErr)
            }
        }
    }()

    row := tx.QueryRowContext(ctx,
        "SELECT id, displacement, no_of_cylinders, car_range, created_at, updated_at FROM engine WHERE id=$1",
        engineID,
    )

    err = row.Scan(
        &engine.EngineID,
        &engine.Displacement,
        &engine.NoOfCyclinders,
        &engine.CarRange,
    )

    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return models.Engine{}, errors.New("engine not found")
        }
        return models.Engine{}, err
    }

    var result sql.Result
    result, err = tx.ExecContext(ctx, "DELETE FROM engine WHERE id=$1", engineID)
    if err != nil {
        return models.Engine{}, err
    }

    var rowsAffected int64
    rowsAffected, err = result.RowsAffected()
    if err != nil {
        return models.Engine{}, err
    }

    if rowsAffected == 0 {
        return models.Engine{}, errors.New("engine deletion failed")
    }

    return engine, nil
}


