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

type VendorHandler struct {
	VendorService *service.VendorService
	Router        *chi.Mux
}

type StatusMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (h *VendorHandler) GetAllVendorsHandler(w http.ResponseWriter, r *http.Request) {
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

	totalVendors, err := h.VendorService.GetTotalVendorsCount()
	if err != nil {
		slog.Error("Error getting total vendors count: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		return
	}

	totalPages := int(math.Ceil(float64(totalVendors) / float64(pageSize)))

	vendors, err := h.VendorService.GetAllVendors(page, pageSize)
	if err != nil {
		slog.Error("Error getting vendors: ", utils.Err(err))
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
	if len(vendors) == pageSize {
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
		"vendors":    vendors,
		"pagination": pagination,
	}

	utils.RespondWithJSON(w, status.OK, responseData)
}

func (h *VendorHandler) GetVendorByIDHandler(w http.ResponseWriter, r *http.Request) {
	vendorID := chi.URLParam(r, "id")

	objectID, err := primitive.ObjectIDFromHex(vendorID)
	if err != nil {
		slog.Error("Invalid vendor ID: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.BadRequest, errs.InvalidVendorID)
		return
	}

	vendor, err := h.VendorService.GetVendorByID(objectID)
	if err != nil {
		slog.Error("Error getting vendor by ID: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		return
	}

	if vendor == nil {
		utils.RespondWithErrorJSON(w, status.NotFound, errs.VendorNotFound)
		return
	}

	utils.RespondWithJSON(w, status.OK, vendor)
}

func (h *VendorHandler) CreateVendorHandler(w http.ResponseWriter, r *http.Request) {
	var createVendorRequest domain.CreateVendorRequest
	err := json.NewDecoder(r.Body).Decode(&createVendorRequest)
	if err != nil {
		utils.RespondWithErrorJSON(w, status.BadRequest, errs.InvalidRequestBody)
		return
	}

	vendor, err := h.VendorService.CreateVendor(&createVendorRequest)
	if err != nil {
		slog.Error("Error creating vendor: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, fmt.Sprintf("Error creating vendor: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(vendor)
}

func (h *VendorHandler) UpdateVendorHandler(w http.ResponseWriter, r *http.Request) {
	vendorID := chi.URLParam(r, "id")

	objectID, err := primitive.ObjectIDFromHex(vendorID)
	if err != nil {
		slog.Error("Invalid vendor ID: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.BadRequest, errs.InvalidVendorID)
		return
	}

	existingVendor, err := h.VendorService.GetVendorByID(objectID)
	if err != nil {
		slog.Error("Error checking if vendor exists: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		return
	}
	if existingVendor == nil {
		utils.RespondWithErrorJSON(w, status.NotFound, errs.VendorNotFound)
		return
	}

	var updateVendorRequest domain.UpdateVendorRequest
	err = json.NewDecoder(r.Body).Decode(&updateVendorRequest)
	if err != nil {
		utils.RespondWithErrorJSON(w, status.BadRequest, errs.InvalidRequestBody)
		return
	}

	vendor, err := h.VendorService.UpdateVendor(objectID, &updateVendorRequest)
	if err != nil {
		slog.Error("Error updating vendor: ", utils.Err(err))
		if err.Error() == "vendor not found" {
			utils.RespondWithErrorJSON(w, status.NotFound, errs.VendorNotFound)
		} else {
			utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		}
		return
	}

	utils.RespondWithJSON(w, status.OK, vendor)
}

func (h *VendorHandler) DeleteVendor(w http.ResponseWriter, r *http.Request) {
	vendorID := chi.URLParam(r, "id")

	objectID, err := primitive.ObjectIDFromHex(vendorID)
	if err != nil {
		slog.Error("Invalid vendor ID: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.BadRequest, errs.InvalidVendorID)
		return
	}

	err = h.VendorService.DeleteVendor(objectID)
	if err != nil {
		if err.Error() == "vendor not found" {
			utils.RespondWithErrorJSON(w, status.NotFound, errs.VendorNotFound)
		} else {
			slog.Error("Error deleting vendor:", utils.Err(err))
			utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		}
		return
	}

	response := StatusMessage{
		Code:    200,
		Message: "Vendor deleted successfully",
	}

	utils.RespondWithJSON(w, status.OK, response)
}

func (h *VendorHandler) SearchVendorsHandler(w http.ResponseWriter, r *http.Request) {
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

	totalVendors, err := h.VendorService.GetTotalVendorsCount()
	if err != nil {
		slog.Error("Error getting total vendor count: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		return
	}

	totalPages := int(math.Ceil(float64(totalVendors) / float64(pageSize)))

	query := r.URL.Query().Get("query")

	vendors, err := h.VendorService.SearchVendors(query, page, pageSize)
	if err != nil {
		slog.Error("Error searching vendors: ", utils.Err(err))
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
	if len(vendors) == pageSize {
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
		"vendors":    vendors,
		"pagination": pagination,
	}

	utils.RespondWithJSON(w, status.OK, responseData)
}

func (h *VendorHandler) FilterVendorsByTagsHandler(w http.ResponseWriter, r *http.Request) {
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

	totalVendors, err := h.VendorService.GetTotalVendorsCount()
	if err != nil {
		slog.Error("Error getting total vendors count: ", utils.Err(err))
		utils.RespondWithErrorJSON(w, status.InternalServerError, errs.InternalServerError)
		return
	}

	totalPages := int(math.Ceil(float64(totalVendors) / float64(pageSize)))

	if len(queryTags) == 0 {
		utils.RespondWithErrorJSON(w, status.BadRequest, errs.MissingTags)
		return
	}

	vendors, err := h.VendorService.FilterVendorsByTags(queryTags, page, pageSize)
	if err != nil {
		slog.Error("Error filtering vendors by tags: ", utils.Err(err))
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
	if len(vendors) == pageSize {
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
		"vendors":    vendors,
		"pagination": pagination,
	}

	utils.RespondWithJSON(w, status.OK, responseData)
}
