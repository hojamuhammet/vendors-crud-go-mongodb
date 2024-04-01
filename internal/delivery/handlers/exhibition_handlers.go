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

type ExhibitionHandler struct {
	ExhibitionService *service.ExhibitionService
	Router            *chi.Mux
}

func (h *ExhibitionHandler) GetAllExhibitionsHandler(w http.ResponseWriter, r *http.Request) {
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

	totalExhibitions, err := h.ExhibitionService.GetTotalExhibitionsCount()
	if err != nil {
		slog.Error("Error getting total exhibitions count: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		return
	}

	totalPages := int(math.Ceil(float64(totalExhibitions) / float64(pageSize)))

	exhibitions, err := h.ExhibitionService.GetAllExhibitions(page, pageSize)
	if err != nil {
		slog.Error("Error getting exhibitions: ", utils.Err(err))
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
	if len(exhibitions) == pageSize {
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
		"exhibitions": exhibitions,
		"pagination":  pagination,
	}

	utils.RespondWithJSON(w, status.OK, responseData)
}

func (h *ExhibitionHandler) GetExhibitionByIDHandler(w http.ResponseWriter, r *http.Request) {
	exhibitionID := chi.URLParam(r, "id")

	objectID, err := primitive.ObjectIDFromHex(exhibitionID)
	if err != nil {
		slog.Error("Invalid exhibition ID: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.BadRequest, errs.InvalidExhibitionID)
		return
	}

	exhibition, err := h.ExhibitionService.GetExhibitionByID(objectID)
	if err != nil {
		slog.Error("Error getting exhibition by ID: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		return
	}

	if exhibition == nil {
		utils.RespondWithErrorJSON(w, status.NotFound, errs.ExhibitionNotFound)
		return
	}

	utils.RespondWithJSON(w, status.OK, exhibition)
}

func (h *ExhibitionHandler) CreateExhibitionHandler(w http.ResponseWriter, r *http.Request) {
	var createExhibitionRequest domain.CreateExhibitionRequest
	err := json.NewDecoder(r.Body).Decode(&createExhibitionRequest)
	if err != nil {
		utils.RespondWithErrorJSON(w, status.BadRequest, errs.InvalidRequestBody)
		return
	}

	exhibition, err := h.ExhibitionService.CreateExhibition(&createExhibitionRequest)
	if err != nil {
		slog.Error("Error creating exhibition: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, fmt.Sprintf("Error creating exhibition: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(exhibition)
}

func (h *ExhibitionHandler) UpdateExhibitionHandler(w http.ResponseWriter, r *http.Request) {
	ExhibitionID := chi.URLParam(r, "id")

	objectID, err := primitive.ObjectIDFromHex(ExhibitionID)
	if err != nil {
		slog.Error("Invalid exhibition ID: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.BadRequest, errs.InvalidExhibitionID)
		return
	}

	existingExhibition, err := h.ExhibitionService.GetExhibitionByID(objectID)
	if err != nil {
		slog.Error("Error checking if exhibition exists: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		return
	}
	if existingExhibition == nil {
		utils.RespondWithErrorJSON(w, status.NotFound, errs.ExhibitionNotFound)
		return
	}

	var updateExhibitionRequest domain.UpdateExhibitionRequest
	err = json.NewDecoder(r.Body).Decode(&updateExhibitionRequest)
	if err != nil {
		utils.RespondWithErrorJSON(w, status.BadRequest, errs.InvalidRequestBody)
		return
	}

	exhibition, err := h.ExhibitionService.UpdateExhibition(objectID, &updateExhibitionRequest)
	if err != nil {
		slog.Error("Error updating exhibition: ", utils.Err(err))
		if err.Error() == "exhibition not found" {
			utils.RespondWithErrorJSON(w, status.NotFound, errs.ExhibitionNotFound)
		} else {
			utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		}
		return
	}

	utils.RespondWithJSON(w, status.OK, exhibition)
}

func (h *ExhibitionHandler) DeleteExhibition(w http.ResponseWriter, r *http.Request) {
	exhibitionID := chi.URLParam(r, "id")

	objectID, err := primitive.ObjectIDFromHex(exhibitionID)
	if err != nil {
		slog.Error("Invalid exhibition ID: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.BadRequest, errs.InvalidExhibitionID)
		return
	}

	err = h.ExhibitionService.DeleteExhibition(objectID)
	if err != nil {
		if err.Error() == "exhibition not found" {
			utils.RespondWithErrorJSON(w, status.NotFound, errs.ExhibitionNotFound)
		} else {
			slog.Error("Error deleting exhibition:", utils.Err(err))
			utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		}
		return
	}

	response := StatusMessage{
		Code:    200,
		Message: "Exhibition deleted successfully",
	}

	utils.RespondWithJSON(w, status.OK, response)
}

func (h *ExhibitionHandler) SearchExhibitionsHandler(w http.ResponseWriter, r *http.Request) {
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

	totalExhibitions, err := h.ExhibitionService.GetTotalExhibitionsCount()
	if err != nil {
		slog.Error("Error getting total exhibitions count: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		return
	}

	totalPages := int(math.Ceil(float64(totalExhibitions) / float64(pageSize)))

	query := r.URL.Query().Get("query")

	exhibitions, err := h.ExhibitionService.SearchExhibitions(query, page, pageSize)
	if err != nil {
		slog.Error("Error searching exhibitions: ", utils.Err(err))
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
	if len(exhibitions) == pageSize {
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
		"exhibitions": exhibitions,
		"pagination":  pagination,
	}

	utils.RespondWithJSON(w, status.OK, responseData)
}

func (h *ExhibitionHandler) FilterExhibitionsByTagsHandler(w http.ResponseWriter, r *http.Request) {
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

	totalExhibitions, err := h.ExhibitionService.GetTotalExhibitionsCount()
	if err != nil {
		slog.Error("Error getting total exhibitions count: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		return
	}

	totalPages := int(math.Ceil(float64(totalExhibitions) / float64(pageSize)))

	if len(queryTags) == 0 {
		utils.RespondWithErrorJSON(w, status.BadRequest, errs.MissingTags)
		return
	}

	exhibitions, err := h.ExhibitionService.FilterExhibitionsByTags(queryTags, page, pageSize)
	if err != nil {
		slog.Error("Error filtering exhibitions by tags: ", utils.Err(err))
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
	if len(exhibitions) == pageSize {
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
		"exhibitions": exhibitions,
		"pagination":  pagination,
	}

	utils.RespondWithJSON(w, status.OK, responseData)
}
