package api

import (
	"net/http"

	"github.com/gophish/gophish/models"
)

// ClientHistory handles requests for the /api/client_history endpoint.
func (as *Server) ClientHistory(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		// Get all client history records.
		client, err := models.GetAllClientHistory()
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		// Respond with the client data
		JSONResponse(w, client, http.StatusOK)
	}
}
