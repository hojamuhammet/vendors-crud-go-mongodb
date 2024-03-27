package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"vendors/internal/domain"
	"vendors/internal/service"
	"vendors/pkg/lib/errs"
	"vendors/pkg/lib/status"
	"vendors/pkg/lib/utils"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TheatreHandler struct {
	TheatreService *service.TheatreService
	Router         *chi.Mux
}

func (h *TheatreHandler) GetAllTheatresHandler(w http.ResponseWriter, r *http.Request) {
	page := 1      // Default page if not provided
	pageSize := 10 // Default page size, adjust as needed

	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		pageNum, err := strconv.Atoi(pageStr)
		if err != nil || pageNum < 1 {
			utils.RespondWithErrorJSON(w, status.BadRequest, errs.InvalidRequestFormat)
			return
		}
		page = pageNum
	}

	totalTheatres, err := h.TheatreService.GetTotalTheatresCount()
	if err != nil {
		slog.Error("Error getting total Theatres count: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		return
	}

	totalPages := int(math.Ceil(float64(totalTheatres) / float64(pageSize)))

	theatres, err := h.TheatreService.GetAllTheatres(page, pageSize)
	if err != nil {
		slog.Error("Error getting theatres: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		return
	}

	var prevPage interface{}
	if page > 1 {
		prevPage = page - 1
	} else {
		prevPage = nil
	}

	var nextPage interface{}
	if len(theatres) == pageSize {
		nextPage = page + 1
	} else {
		nextPage = nil
	}

	var firstPage interface{}
	if totalPages > 0 {
		firstPage = 1
	} else {
		firstPage = nil
	}

	var lastPage interface{}
	if totalPages >= 1 {
		lastPage = totalPages
	} else {
		lastPage = firstPage
	}

	pagination := map[string]interface{}{
		"current_page": page,
		"prev_page":    prevPage,
		"next_page":    nextPage,
		"first_page":   firstPage,
		"last_page":    lastPage,
	}

	responseData := map[string]interface{}{
		"theatres":   theatres,
		"pagination": pagination,
	}

	utils.RespondWithJSON(w, status.OK, responseData)
}

func (h *TheatreHandler) GetTheatreByIDHandler(w http.ResponseWriter, r *http.Request) {
	theatreID := chi.URLParam(r, "id")

	objectID, err := primitive.ObjectIDFromHex(theatreID)
	if err != nil {
		slog.Error("Invalid theatre ID: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.BadRequest, errs.InvalidTheatreID)
		return
	}

	theatre, err := h.TheatreService.GetTheatreByID(objectID)
	if err != nil {
		slog.Error("Error getting theatre by ID: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		return
	}

	if theatre == nil {
		utils.RespondWithErrorJSON(w, status.NotFound, errs.TheatreNotFound)
		return
	}

	utils.RespondWithJSON(w, status.OK, theatre)
}

func (h *TheatreHandler) CreateTheatreHandler(w http.ResponseWriter, r *http.Request) {
	var createTheatreRequest domain.CreateTheatreRequest
	err := json.NewDecoder(r.Body).Decode(&createTheatreRequest)
	if err != nil {
		utils.RespondWithErrorJSON(w, status.BadRequest, errs.InvalidRequestBody)
		return
	}

	theatre, err := h.TheatreService.CreateTheatre(&createTheatreRequest)
	if err != nil {
		slog.Error("Error creating theatre: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, fmt.Sprintf("Error creating theatre: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(theatre)
}

func (h *TheatreHandler) UpdateTheatreHandler(w http.ResponseWriter, r *http.Request) {
	theatreID := chi.URLParam(r, "id")

	objectID, err := primitive.ObjectIDFromHex(theatreID)
	if err != nil {
		slog.Error("Invalid theatre ID: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.BadRequest, errs.InvalidTheatreID)
		return
	}

	existingTheatre, err := h.TheatreService.GetTheatreByID(objectID)
	if err != nil {
		slog.Error("Error checking if theatre exists: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		return
	}
	if existingTheatre == nil {
		utils.RespondWithErrorJSON(w, status.NotFound, errs.TheatreNotFound)
		return
	}

	var updateTheatreRequest domain.UpdateTheatreRequest
	err = json.NewDecoder(r.Body).Decode(&updateTheatreRequest)
	if err != nil {
		utils.RespondWithErrorJSON(w, status.BadRequest, errs.InvalidRequestBody)
		return
	}

	theatre, err := h.TheatreService.UpdateTheatre(objectID, &updateTheatreRequest)
	if err != nil {
		slog.Error("Error updating theatre: ", utils.Err(err))
		if err.Error() == "theatre not found" {
			utils.RespondWithErrorJSON(w, status.NotFound, errs.TheatreNotFound)
		} else {
			utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		}
		return
	}

	utils.RespondWithJSON(w, status.OK, theatre)
}

func (h *TheatreHandler) DeleteTheatre(w http.ResponseWriter, r *http.Request) {
	theatreID := chi.URLParam(r, "id")

	objectID, err := primitive.ObjectIDFromHex(theatreID)
	if err != nil {
		slog.Error("Invalid theatre ID: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.BadRequest, errs.InvalidTheatreID)
		return
	}

	err = h.TheatreService.DeleteTheatre(objectID)
	if err != nil {
		if err.Error() == "theatre not found" {
			utils.RespondWithErrorJSON(w, status.NotFound, errs.TheatreNotFound)
		} else {
			slog.Error("Error deleting theatre:", utils.Err(err))
			utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		}
		return
	}

	response := StatusMessage{
		Code:    200,
		Message: "Theatre deleted successfully",
	}

	utils.RespondWithJSON(w, status.OK, response)
}

func (h *TheatreHandler) SearchTheatresHandler(w http.ResponseWriter, r *http.Request) {
	page := 1      // Default page if not provided
	pageSize := 10 // Default page size, adjust as needed

	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		pageNum, err := strconv.Atoi(pageStr)
		if err != nil || pageNum < 1 {
			utils.RespondWithErrorJSON(w, status.BadRequest, errs.InvalidRequestFormat)
			return
		}
		page = pageNum
	}

	totalTheatres, err := h.TheatreService.GetTotalTheatresCount()
	if err != nil {
		slog.Error("Error getting total theatres count: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		return
	}

	totalPages := int(math.Ceil(float64(totalTheatres) / float64(pageSize)))

	query := r.URL.Query().Get("query")

	theatres, err := h.TheatreService.SearchTheatres(query, page, pageSize)
	if err != nil {
		slog.Error("Error searching theatres: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		return
	}

	var prevPage interface{}
	if page > 1 {
		prevPage = page - 1
	} else {
		prevPage = nil
	}

	var nextPage interface{}
	if len(theatres) == pageSize {
		nextPage = page + 1
	} else {
		nextPage = nil
	}

	var firstPage interface{}
	if totalPages > 0 {
		firstPage = 1
	} else {
		firstPage = nil
	}

	var lastPage interface{}
	if totalPages >= 1 {
		lastPage = totalPages
	} else {
		lastPage = firstPage
	}

	pagination := map[string]interface{}{
		"current_page": page,
		"prev_page":    prevPage,
		"next_page":    nextPage,
		"first_page":   firstPage,
		"last_page":    lastPage,
	}

	responseData := map[string]interface{}{
		"theatres":   theatres,
		"pagination": pagination,
	}

	utils.RespondWithJSON(w, status.OK, responseData)
}

func (h *TheatreHandler) FilterTheatresByTagsHandler(w http.ResponseWriter, r *http.Request) {
	page := 1      // Default page if not provided
	pageSize := 10 // Default page size, adjust as needed
	queryTags := r.URL.Query()["tags"]

	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		pageNum, err := strconv.Atoi(pageStr)
		if err != nil || pageNum < 1 {
			utils.RespondWithErrorJSON(w, status.BadRequest, errs.InvalidRequestFormat)
			return
		}
		page = pageNum
	}

	totalTheatres, err := h.TheatreService.GetTotalTheatresCount()
	if err != nil {
		slog.Error("Error getting total theatres count: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		return
	}

	totalPages := int(math.Ceil(float64(totalTheatres) / float64(pageSize)))

	if len(queryTags) == 0 {
		utils.RespondWithErrorJSON(w, status.BadRequest, errs.MissingTags)
		return
	}

	theatres, err := h.TheatreService.FilterTheatresByTags(queryTags, page, pageSize)
	if err != nil {
		slog.Error("Error filtering theatres by tags: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		return
	}

	var prevPage interface{}
	if page > 1 {
		prevPage = page - 1
	} else {
		prevPage = nil
	}

	var nextPage interface{}
	if len(theatres) == pageSize {
		nextPage = page + 1
	} else {
		nextPage = nil
	}

	var firstPage interface{}
	if totalPages > 0 {
		firstPage = 1
	} else {
		firstPage = nil
	}

	var lastPage interface{}
	if totalPages >= 1 {
		lastPage = totalPages
	} else {
		lastPage = firstPage
	}

	pagination := map[string]interface{}{
		"current_page": page,
		"prev_page":    prevPage,
		"next_page":    nextPage,
		"first_page":   firstPage,
		"last_page":    lastPage,
	}

	responseData := map[string]interface{}{
		"theatres":   theatres,
		"pagination": pagination,
	}

	utils.RespondWithJSON(w, status.OK, responseData)
}
