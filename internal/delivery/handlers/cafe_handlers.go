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

type CafeHandler struct {
	CafeService *service.CafeService
	Router      *chi.Mux
}

type StatusMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (h *CafeHandler) GetAllCafesHandler(w http.ResponseWriter, r *http.Request) {
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

	totalCafes, err := h.CafeService.GetTotalCafesCount()
	if err != nil {
		slog.Error("Error getting total cafes count: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		return
	}

	totalPages := int(math.Ceil(float64(totalCafes) / float64(pageSize)))

	cafes, err := h.CafeService.GetAllCafes(page, pageSize)
	if err != nil {
		slog.Error("Error getting cafes: ", utils.Err(err))
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
	if len(cafes) == pageSize {
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
		"cafes":      cafes,
		"pagination": pagination,
	}

	utils.RespondWithJSON(w, status.OK, responseData)
}

func (h *CafeHandler) GetCafeByIDHandler(w http.ResponseWriter, r *http.Request) {
	cafeID := chi.URLParam(r, "id")

	objectID, err := primitive.ObjectIDFromHex(cafeID)
	if err != nil {
		slog.Error("Invalid cafe ID: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.BadRequest, errs.InvalidCafeID)
		return
	}

	cafe, err := h.CafeService.GetCafeByID(objectID)
	if err != nil {
		slog.Error("Error getting cafe by ID: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		return
	}

	if cafe == nil {
		utils.RespondWithErrorJSON(w, status.NotFound, errs.CafeNotFound)
		return
	}

	utils.RespondWithJSON(w, status.OK, cafe)
}

func (h *CafeHandler) CreateCafeHandler(w http.ResponseWriter, r *http.Request) {
	var createCafeRequest domain.CreateCafeRequest
	err := json.NewDecoder(r.Body).Decode(&createCafeRequest)
	if err != nil {
		utils.RespondWithErrorJSON(w, status.BadRequest, errs.InvalidRequestBody)
		return
	}

	cafe, err := h.CafeService.CreateCafe(&createCafeRequest)
	if err != nil {
		slog.Error("Error creating cafe: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, fmt.Sprintf("Error creating Cafe: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cafe)
}

func (h *CafeHandler) UpdateCafeHandler(w http.ResponseWriter, r *http.Request) {
	cafeID := chi.URLParam(r, "id")

	objectID, err := primitive.ObjectIDFromHex(cafeID)
	if err != nil {
		slog.Error("Invalid cafe ID: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.BadRequest, errs.InvalidCafeID)
		return
	}

	existingCafe, err := h.CafeService.GetCafeByID(objectID)
	if err != nil {
		slog.Error("Error checking if cafe exists: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		return
	}
	if existingCafe == nil {
		utils.RespondWithErrorJSON(w, status.NotFound, errs.CafeNotFound)
		return
	}

	var updateCafeRequest domain.UpdateCafeRequest
	err = json.NewDecoder(r.Body).Decode(&updateCafeRequest)
	if err != nil {
		utils.RespondWithErrorJSON(w, status.BadRequest, errs.InvalidRequestBody)
		return
	}

	cafe, err := h.CafeService.UpdateCafe(objectID, &updateCafeRequest)
	if err != nil {
		slog.Error("Error updating cafe: ", utils.Err(err))
		if err.Error() == "cafe not found" {
			utils.RespondWithErrorJSON(w, status.NotFound, errs.CafeNotFound)
		} else {
			utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		}
		return
	}

	utils.RespondWithJSON(w, status.OK, cafe)
}

func (h *CafeHandler) DeleteCafe(w http.ResponseWriter, r *http.Request) {
	CafeID := chi.URLParam(r, "id")

	objectID, err := primitive.ObjectIDFromHex(CafeID)
	if err != nil {
		slog.Error("Invalid cafe ID: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.BadRequest, errs.InvalidCafeID)
		return
	}

	err = h.CafeService.DeleteCafe(objectID)
	if err != nil {
		if err.Error() == "cafe not found" {
			utils.RespondWithErrorJSON(w, status.NotFound, errs.CafeNotFound)
		} else {
			slog.Error("Error deleting cafe:", utils.Err(err))
			utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		}
		return
	}

	response := StatusMessage{
		Code:    200,
		Message: "Cafe deleted successfully",
	}

	utils.RespondWithJSON(w, status.OK, response)
}

func (h *CafeHandler) SearchCafesHandler(w http.ResponseWriter, r *http.Request) {
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

	totalCafes, err := h.CafeService.GetTotalCafesCount()
	if err != nil {
		slog.Error("Error getting total cafes count: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		return
	}

	totalPages := int(math.Ceil(float64(totalCafes) / float64(pageSize)))

	query := r.URL.Query().Get("query")

	cafes, err := h.CafeService.SearchCafes(query, page, pageSize)
	if err != nil {
		slog.Error("Error searching cafes: ", utils.Err(err))
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
	if len(cafes) == pageSize {
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
		"cafes":      cafes,
		"pagination": pagination,
	}

	utils.RespondWithJSON(w, status.OK, responseData)
}

func (h *CafeHandler) FilterCafesByTagsHandler(w http.ResponseWriter, r *http.Request) {
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

	totalCafes, err := h.CafeService.GetTotalCafesCount()
	if err != nil {
		slog.Error("Error getting total cafes count: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		return
	}

	totalPages := int(math.Ceil(float64(totalCafes) / float64(pageSize)))

	if len(queryTags) == 0 {
		utils.RespondWithErrorJSON(w, status.BadRequest, errs.MissingTags)
		return
	}

	cafes, err := h.CafeService.FilterCafesByTags(queryTags, page, pageSize)
	if err != nil {
		slog.Error("Error filtering cafes by tags: ", utils.Err(err))
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
	if len(cafes) == pageSize {
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
		"cafes":      cafes,
		"pagination": pagination,
	}

	utils.RespondWithJSON(w, status.OK, responseData)
}
