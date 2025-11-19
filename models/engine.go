package models

import (
	"errors"

	"github.com/google/uuid"
)

type Engine struct {
	EngineID  uuid.UUID `json: "engine_id"`
	Displacement int64 `json: "displacement"`
	NoOfCyclinders int64 `json: "noOfCyclinders"`
	CarRange int64 `json: "carRange"`
}

type EngineRequest struct {
	EngineID  uuid.UUID `json: "engine_id"`
	Displacement int64 `json: "displacement"`
	NoOfCyclinders int64 `json: "noOfCyclinders"`
	CarRange int64 `json: "carRange"`
}

func ValidateEngineRequest(engineReq EngineRequest) error {
	if err := validateDisplacement(engineReq.Displacement); err != nil {
		return err
	}
	if err := validateNoOfCyclinders(engineReq.NoOfCyclinders); err != nil {
		return err
	}
	if err := validateCarRange(engineReq.CarRange); err != nil {
		return err
	}
}

func validateDisplacement(displacement int64) error {
	if displacement <= 0 {
		return errors.New("displacement must be greater than zero")
	}
	return nil
}

func validateNoOfCyclinders(noOfCyclinders int64) error {
	if noOfCyclinders <= 0 {
		return  errors.New("noOfCyclinders must be greater than zero")
	}
	return nil
}

func validateCarRange(carRange int64) error {
	if carRange <= 0 {
		return errors.New("Car range must be greater than zero")
	}
	return nil
}
