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

type CinemaHandler struct {
	CinemaService *service.CinemaService
	Router        *chi.Mux
}

func (h *CinemaHandler) GetAllCinemasHandler(w http.ResponseWriter, r *http.Request) {
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

	totalCinemas, err := h.CinemaService.GetTotalCinemasCount()
	if err != nil {
		slog.Error("Error getting total cinema count: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		return
	}

	totalPages := int(math.Ceil(float64(totalCinemas) / float64(pageSize)))

	cinemas, err := h.CinemaService.GetAllCinemas(page, pageSize)
	if err != nil {
		slog.Error("Error getting cinemas: ", utils.Err(err))
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
	if len(cinemas) == pageSize {
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
		"cinemas":    cinemas,
		"pagination": pagination,
	}

	utils.RespondWithJSON(w, status.OK, responseData)
}

func (h *CinemaHandler) GetCinemaByIDHandler(w http.ResponseWriter, r *http.Request) {
	cinemaID := chi.URLParam(r, "id")

	objectID, err := primitive.ObjectIDFromHex(cinemaID)
	if err != nil {
		slog.Error("Invalid cinema ID: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.BadRequest, errs.InvalidCinemaID)
		return
	}

	cinema, err := h.CinemaService.GetCinemaByID(objectID)
	if err != nil {
		slog.Error("Error getting cinema by ID: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		return
	}

	if cinema == nil {
		utils.RespondWithErrorJSON(w, status.NotFound, errs.CinemaNotFound)
		return
	}

	utils.RespondWithJSON(w, status.OK, cinema)
}

func (h *CinemaHandler) CreateCinemaHandler(w http.ResponseWriter, r *http.Request) {
	var createCinemaRequest domain.CreateCinemaRequest
	err := json.NewDecoder(r.Body).Decode(&createCinemaRequest)
	if err != nil {
		utils.RespondWithErrorJSON(w, status.BadRequest, errs.InvalidRequestBody)
		return
	}

	cinema, err := h.CinemaService.CreateCinema(&createCinemaRequest)
	if err != nil {
		slog.Error("Error creating cinema: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, fmt.Sprintf("Error creating cinema: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cinema)
}

func (h *CinemaHandler) UpdateCinemaHandler(w http.ResponseWriter, r *http.Request) {
	cinemaID := chi.URLParam(r, "id")

	objectID, err := primitive.ObjectIDFromHex(cinemaID)
	if err != nil {
		slog.Error("Invalid cinema ID: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.BadRequest, errs.InvalidCinemaID)
		return
	}

	existingCinema, err := h.CinemaService.GetCinemaByID(objectID)
	if err != nil {
		slog.Error("Error checking if cinema exists: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		return
	}
	if existingCinema == nil {
		utils.RespondWithErrorJSON(w, status.NotFound, errs.CinemaNotFound)
		return
	}

	var updateCinemaRequest domain.UpdateCinemaRequest
	err = json.NewDecoder(r.Body).Decode(&updateCinemaRequest)
	if err != nil {
		utils.RespondWithErrorJSON(w, status.BadRequest, errs.InvalidRequestBody)
		return
	}

	cinema, err := h.CinemaService.UpdateCinema(objectID, &updateCinemaRequest)
	if err != nil {
		slog.Error("Error updating cinema: ", utils.Err(err))
		if err.Error() == "cinema not found" {
			utils.RespondWithErrorJSON(w, status.NotFound, errs.CinemaNotFound)
		} else {
			utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		}
		return
	}

	utils.RespondWithJSON(w, status.OK, cinema)
}

func (h *CinemaHandler) DeleteCinema(w http.ResponseWriter, r *http.Request) {
	cinemaID := chi.URLParam(r, "id")

	objectID, err := primitive.ObjectIDFromHex(cinemaID)
	if err != nil {
		slog.Error("Invalid cinema ID: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.BadRequest, errs.InvalidCafeID)
		return
	}

	err = h.CinemaService.DeleteCinema(objectID)
	if err != nil {
		if err.Error() == "cinema not found" {
			utils.RespondWithErrorJSON(w, status.NotFound, errs.CafeNotFound)
		} else {
			slog.Error("Error deleting cinema:", utils.Err(err))
			utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		}
		return
	}

	response := StatusMessage{
		Code:    200,
		Message: "Cinema deleted successfully",
	}

	utils.RespondWithJSON(w, status.OK, response)
}

func (h *CinemaHandler) SearchCinemasHandler(w http.ResponseWriter, r *http.Request) {
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

	totalCinemas, err := h.CinemaService.GetTotalCinemasCount()
	if err != nil {
		slog.Error("Error getting total cinemas count: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		return
	}

	totalPages := int(math.Ceil(float64(totalCinemas) / float64(pageSize)))

	query := r.URL.Query().Get("query")

	cinemas, err := h.CinemaService.SearchCinemas(query, page, pageSize)
	if err != nil {
		slog.Error("Error searching cinemas: ", utils.Err(err))
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
	if len(cinemas) == pageSize {
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
		"cinemas":    cinemas,
		"pagination": pagination,
	}

	utils.RespondWithJSON(w, status.OK, responseData)
}

func (h *CinemaHandler) FilterCinemasByTagsHandler(w http.ResponseWriter, r *http.Request) {
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

	totalCinemas, err := h.CinemaService.GetTotalCinemasCount()
	if err != nil {
		slog.Error("Error getting total cinemas count: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		return
	}

	totalPages := int(math.Ceil(float64(totalCinemas) / float64(pageSize)))

	if len(queryTags) == 0 {
		utils.RespondWithErrorJSON(w, status.BadRequest, errs.MissingTags)
		return
	}

	cinemas, err := h.CinemaService.FilterCinemasByTags(queryTags, page, pageSize)
	if err != nil {
		slog.Error("Error filtering cinemas by tags: ", utils.Err(err))
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
	if len(cinemas) == pageSize {
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
		"cinemas":    cinemas,
		"pagination": pagination,
	}

	utils.RespondWithJSON(w, status.OK, responseData)
}
