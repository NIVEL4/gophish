package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	ctx "github.com/gophish/gophish/context"
	log "github.com/gophish/gophish/logger"
	"github.com/gophish/gophish/models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// SendingProfiles handles requests for the /api/whatsapp/ endpoint
func (as *Server) WhatsappProfiles(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		ss, err := models.GetWhatsapps(ctx.Get(r, "user_id").(int64))
		if err != nil {
			log.Error(err)
		}
		JSONResponse(w, ss, http.StatusOK)
	//POST: Create a new Whatsapp and return it as JSON
	case r.Method == "POST":
		wsp := models.Whatsapp{}
		// Put the request into a page
		err := json.NewDecoder(r.Body).Decode(&wsp)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Invalid request"}, http.StatusBadRequest)
			return
		}
		// Check to make sure the name is unique
		_, err = models.GetWhatsappByName(wsp.Name, ctx.Get(r, "user_id").(int64))
		if err != gorm.ErrRecordNotFound {
			JSONResponse(w, models.Response{Success: false, Message: "Whatsapp name already in use"}, http.StatusConflict)
			log.Error(err)
			return
		}
		wsp.ModifiedDate = time.Now().UTC()
		wsp.UserId = ctx.Get(r, "user_id").(int64)
		err = models.PostWhatsapp(&wsp)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		JSONResponse(w, wsp, http.StatusCreated)
	}
}

// SendingProfile contains functions to handle the GET'ing, DELETE'ing, and PUT'ing
// of a Whatsapp object
func (as *Server) WhatsappProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 0, 64)
	wsp, err := models.GetWhatsapp(id, ctx.Get(r, "user_id").(int64))
	if err != nil {
		JSONResponse(w, models.Response{Success: false, Message: "Whatsapp not found"}, http.StatusNotFound)
		return
	}
	switch {
	case r.Method == "GET":
		JSONResponse(w, wsp, http.StatusOK)
	case r.Method == "DELETE":
		err = models.DeleteWhatsapp(id, ctx.Get(r, "user_id").(int64))
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Error deleting Whatsapp"}, http.StatusInternalServerError)
			return
		}
		JSONResponse(w, models.Response{Success: true, Message: "Whatsapp Deleted Successfully"}, http.StatusOK)
	case r.Method == "PUT":
		wsp = models.Whatsapp{}
		err = json.NewDecoder(r.Body).Decode(&wsp)
		if err != nil {
			log.Error(err)
		}
		if wsp.Id != id {
			JSONResponse(w, models.Response{Success: false, Message: "/:id and /:whatsapp_id mismatch"}, http.StatusBadRequest)
			return
		}
		err = wsp.Validate()
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusBadRequest)
			return
		}
		wsp.ModifiedDate = time.Now().UTC()
		wsp.UserId = ctx.Get(r, "user_id").(int64)
		err = models.PutWhatsapp(&wsp)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Error updating page"}, http.StatusInternalServerError)
			return
		}
		JSONResponse(w, wsp, http.StatusOK)
	}
}
