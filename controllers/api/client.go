package api

import (
	"encoding/json"
	"net/http"

	"github.com/gophish/gophish/models"
)

// Client handles requests for the /api/client/ endpoint
func (as *Server) Client(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		client, err := models.GetClient()
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		JSONResponse(w, client, http.StatusOK)

	// POST: Update database
	case r.Method == "POST":
		client := models.Client{}
		err := json.NewDecoder(r.Body).Decode(&client)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Error decoding client data."}, http.StatusBadRequest)
			return
		}
		err = client.Validate()
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
		}
		err = models.UpdateClient(&client)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		JSONResponse(w, models.Response{Success: true, Message: "Successfully saved client data."}, http.StatusCreated)
	}
}
