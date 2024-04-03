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

type PlaceHandler struct {
	PlaceService *service.PlaceService
	Router       *chi.Mux
}

type StatusMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (h *PlaceHandler) GetAllPlacesHandler(w http.ResponseWriter, r *http.Request) {
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

	totalPlaces, err := h.PlaceService.GetTotalPlacesCount()
	if err != nil {
		slog.Error("Error getting total places count: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		return
	}

	totalPages := int(math.Ceil(float64(totalPlaces) / float64(pageSize)))

	places, err := h.PlaceService.GetAllPlaces(page, pageSize)
	if err != nil {
		slog.Error("Error getting places: ", utils.Err(err))
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
	if len(places) == pageSize {
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
		"places":     places,
		"pagination": pagination,
	}

	utils.RespondWithJSON(w, status.OK, responseData)
}

func (h *PlaceHandler) GetPlaceByIDHandler(w http.ResponseWriter, r *http.Request) {
	placeID := chi.URLParam(r, "id")

	objectID, err := primitive.ObjectIDFromHex(placeID)
	if err != nil {
		slog.Error("Invalid place ID: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.BadRequest, errs.InvalidPlaceID)
		return
	}

	place, err := h.PlaceService.GetPlaceByID(objectID)
	if err != nil {
		slog.Error("Error getting place by ID: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		return
	}

	if place == nil {
		utils.RespondWithErrorJSON(w, status.NotFound, errs.PlaceNotFound)
		return
	}

	utils.RespondWithJSON(w, status.OK, place)
}

func (h *PlaceHandler) CreatePlaceHandler(w http.ResponseWriter, r *http.Request) {
	var createPlaceRequest domain.CreatePlaceRequest
	err := json.NewDecoder(r.Body).Decode(&createPlaceRequest)
	if err != nil {
		utils.RespondWithErrorJSON(w, status.BadRequest, errs.InvalidRequestBody)
		return
	}

	place, err := h.PlaceService.CreatePlace(&createPlaceRequest)
	if err != nil {
		slog.Error("Error creating place: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, fmt.Sprintf("Error creating place: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(place)
}

func (h *PlaceHandler) UpdatePlaceHandler(w http.ResponseWriter, r *http.Request) {
	placeID := chi.URLParam(r, "id")

	objectID, err := primitive.ObjectIDFromHex(placeID)
	if err != nil {
		slog.Error("Invalid place ID: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.BadRequest, errs.InvalidPlaceID)
		return
	}

	existingPlace, err := h.PlaceService.GetPlaceByID(objectID)
	if err != nil {
		slog.Error("Error checking if place exists: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		return
	}
	if existingPlace == nil {
		utils.RespondWithErrorJSON(w, status.NotFound, errs.PlaceNotFound)
		return
	}

	var updatePlaceRequest domain.UpdatePlaceRequest
	err = json.NewDecoder(r.Body).Decode(&updatePlaceRequest)
	if err != nil {
		utils.RespondWithErrorJSON(w, status.BadRequest, errs.InvalidRequestBody)
		return
	}

	place, err := h.PlaceService.UpdatePlace(objectID, &updatePlaceRequest)
	if err != nil {
		slog.Error("Error updating place: ", utils.Err(err))
		if err.Error() == "place not found" {
			utils.RespondWithErrorJSON(w, status.NotFound, errs.PlaceNotFound)
		} else {
			utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		}
		return
	}

	utils.RespondWithJSON(w, status.OK, place)
}

func (h *PlaceHandler) DeletePlace(w http.ResponseWriter, r *http.Request) {
	placeID := chi.URLParam(r, "id")

	objectID, err := primitive.ObjectIDFromHex(placeID)
	if err != nil {
		slog.Error("Invalid place ID: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.BadRequest, errs.InvalidPlaceID)
		return
	}

	err = h.PlaceService.DeletePlace(objectID)
	if err != nil {
		if err.Error() == "place not found" {
			utils.RespondWithErrorJSON(w, status.NotFound, errs.PlaceNotFound)
		} else {
			slog.Error("Error deleting place:", utils.Err(err))
			utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		}
		return
	}

	response := StatusMessage{
		Code:    200,
		Message: "Place deleted successfully",
	}

	utils.RespondWithJSON(w, status.OK, response)
}

func (h *PlaceHandler) SearchPlacesHandler(w http.ResponseWriter, r *http.Request) {
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

	totalPlaces, err := h.PlaceService.GetTotalPlacesCount()
	if err != nil {
		slog.Error("Error getting total places count: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		return
	}

	totalPages := int(math.Ceil(float64(totalPlaces) / float64(pageSize)))

	query := r.URL.Query().Get("query")

	places, err := h.PlaceService.SearchPlaces(query, page, pageSize)
	if err != nil {
		slog.Error("Error searching places: ", utils.Err(err))
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
	if len(places) == pageSize {
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
		"places":     places,
		"pagination": pagination,
	}

	utils.RespondWithJSON(w, status.OK, responseData)
}

func (h *PlaceHandler) FilterPlacesByTagsHandler(w http.ResponseWriter, r *http.Request) {
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

	totalPlaces, err := h.PlaceService.GetTotalPlacesCount()
	if err != nil {
		slog.Error("Error getting total places count: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		return
	}

	totalPages := int(math.Ceil(float64(totalPlaces) / float64(pageSize)))

	if len(queryTags) == 0 {
		utils.RespondWithErrorJSON(w, status.BadRequest, errs.MissingTags)
		return
	}

	places, err := h.PlaceService.FilterPlacesByTags(queryTags, page, pageSize)
	if err != nil {
		slog.Error("Error filtering places by tags: ", utils.Err(err))
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
	if len(places) == pageSize {
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
		"places":     places,
		"pagination": pagination,
	}

	utils.RespondWithJSON(w, status.OK, responseData)
}
