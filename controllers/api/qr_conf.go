package api

import (
	"encoding/json"
	"net/http"

	"github.com/gophish/gophish/models"
)

// QRConf handles requests for the /api/qrconf/ endpoint
func (as *Server) QRConf(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		qr, err := models.GetQR()
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		JSONResponse(w, qr, http.StatusOK)

	// POST: Update database
	case r.Method == "POST":
		qr := models.QR{}
		err := json.NewDecoder(r.Body).Decode(&qr)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Error decoding QR configs."}, http.StatusBadRequest)
			return
		}
		err = qr.Validate()
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
		}
		err = models.UpdateQR(&qr)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		JSONResponse(w, models.Response{Success: true, Message: "Successfully saved QR settings."}, http.StatusCreated)
	}
}
