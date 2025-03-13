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
		// Get the most recent client
		client, err := models.GetClient()
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		// Check if the client is empty
		if client.Name == "" && client.Email == "" {
			JSONResponse(w, models.Response{Success: false, Message: "No clients registered."}, http.StatusOK)
			return
		}
		JSONResponse(w, client, http.StatusOK)
	case r.Method == "POST":
		// Handle client update
		client := models.Client{}
		err := json.NewDecoder(r.Body).Decode(&client)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Error decoding client data."}, http.StatusBadRequest)
			return
		}
		// Validate the incoming client data
		err = client.Validate()
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		// Save new client or update in case of doesnt exist and save the change in history records.
		err = models.UpdateClient(&client)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}

		// Respond with success message
		JSONResponse(w, models.Response{Success: true, Message: "Successfully saved and recorded client data."}, http.StatusCreated)
	}
}
