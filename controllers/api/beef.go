package api

import (
	"encoding/json"
	"net/http"

	"github.com/gophish/gophish/models"
)

// BeEF handles requests for the /api/beef/ endpoint
func (as *Server) BeEF(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		beef, err := models.GetBeEF()
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		JSONResponse(w, beef, http.StatusOK)

	// POST: Update database
	case r.Method == "POST":
		beef := models.BeEF{}
		err := json.NewDecoder(r.Body).Decode(&beef)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Error decoding beef data."}, http.StatusBadRequest)
			return
		}
		err = beef.Validate()
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
		}
		err = models.UpdateBeEF(&beef)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		JSONResponse(w, models.Response{Success: true, Message: "Successfully saved beef data."}, http.StatusCreated)
	}
}
