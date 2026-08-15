package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/andreshungbz/cmps4191-gatekeeper/internal/data"
	"github.com/andreshungbz/cmps4191-gatekeeper/internal/validator"
	"github.com/google/uuid"
)

// createAPIKeyHandler reads JSON input to create an API key record, then returns JSON of the created record.
func (app *application) createAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ConsumerID uuid.UUID  `json:"consumer_id"`
		KeyHash    string     `json:"key_hash"`
		KeyPrefix  string     `json:"key_prefix"`
		Status     string     `json:"status"`
		ExpiresAt  *time.Time `json:"expires_at"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	apiKey := &data.APIKey{
		ConsumerID: input.ConsumerID,
		KeyHash:    input.KeyHash,
		KeyPrefix:  input.KeyPrefix,
		Status:     data.KeyStatus(input.Status),
		ExpiresAt:  input.ExpiresAt,
	}

	v := validator.New()
	if data.ValidateAPIKey(v, apiKey); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.APIKey.Insert(apiKey)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"api_key": apiKey}, nil)
}

// showAPIKeyHandler reads the UUID in the URL path, then returns JSON of the matching API key record.
func (app *application) showAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readUUIDParam("id", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	apiKey, err := app.models.APIKey.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"api_key": apiKey}, nil)
}

// updateAPIKeyHandler updates the API key record matching UUID in the URL path, then returns JSON of the updated record.
func (app *application) updateAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readUUIDParam("id", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	apiKey, err := app.models.APIKey.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Status     *string    `json:"status"`
		LastUsedAt *time.Time `json:"last_used_at"`
		ExpiresAt  *time.Time `json:"expires_at"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Status != nil {
		apiKey.Status = data.KeyStatus(*input.Status)
	}
	if input.LastUsedAt != nil {
		apiKey.LastUsedAt = input.LastUsedAt
	}
	if input.ExpiresAt != nil {
		apiKey.ExpiresAt = input.ExpiresAt
	}

	v := validator.New()
	if data.ValidateAPIKey(v, apiKey); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.APIKey.Update(apiKey)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"api_key": apiKey}, nil)
}

// deleteAPIKeyHandler deletes the API key record matching UUID in the URL path.
func (app *application) deleteAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readUUIDParam("id", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.APIKey.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "API key successfully deleted"}, nil)
}
