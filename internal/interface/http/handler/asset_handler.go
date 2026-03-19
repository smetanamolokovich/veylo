package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	appasset "github.com/smetanamolokovich/veylo/internal/application/asset"
	"github.com/smetanamolokovich/veylo/internal/domain/asset"
	authmiddleware "github.com/smetanamolokovich/veylo/internal/interface/http/middleware"
)

type AssetHandler struct {
	createVehicleUseCase *appasset.CreateVehicleAssetUseCase
	getAssetUseCase      *appasset.GetAssetUseCase
}

func NewAssetHandler(createVehicleUC *appasset.CreateVehicleAssetUseCase, getAssetUC *appasset.GetAssetUseCase) *AssetHandler {
	return &AssetHandler{
		createVehicleUseCase: createVehicleUC,
		getAssetUseCase:      getAssetUC,
	}
}

type createVehicleRequest struct {
	VIN             string `json:"vin"`
	LicensePlate    string `json:"license_plate"`
	Brand           string `json:"brand"`
	Model           string `json:"model"`
	BodyType        string `json:"body_type"`
	FuelType        string `json:"fuel_type"`
	Transmission    string `json:"transmission"`
	OdometerReading int    `json:"odometer_reading"`
	Color           string `json:"color"`
	EnginePower     int    `json:"engine_power"`
}

func (h *AssetHandler) CreateVehicle(w http.ResponseWriter, r *http.Request) {
	orgID, ok := authmiddleware.OrganizationIDFromCtx(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req createVehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.createVehicleUseCase.Execute(r.Context(), appasset.CreateVehicleAssetRequest{
		OrganizationID:  orgID,
		VIN:             req.VIN,
		LicensePlate:    req.LicensePlate,
		Brand:           req.Brand,
		Model:           req.Model,
		BodyType:        req.BodyType,
		FuelType:        req.FuelType,
		Transmission:    req.Transmission,
		OdometerReading: req.OdometerReading,
		Color:           req.Color,
		EnginePower:     req.EnginePower,
	})
	if err != nil {
		if errors.Is(err, asset.ErrAlreadyExists) {
			writeError(w, http.StatusConflict, "asset with this VIN or license plate already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *AssetHandler) Get(w http.ResponseWriter, r *http.Request) {
	orgID, ok := authmiddleware.OrganizationIDFromCtx(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id := chi.URLParam(r, "id")

	resp, err := h.getAssetUseCase.Execute(r.Context(), appasset.GetAssetRequest{
		ID:             id,
		OrganizationID: orgID,
	})
	if err != nil {
		if errors.Is(err, asset.ErrNotFound) {
			writeError(w, http.StatusNotFound, "asset not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
