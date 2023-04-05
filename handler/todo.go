package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		params := r.URL.Query()
		prevID := 0
		size := 5

		if len(params["prev_id"]) > 0 {
			prevID, _ = strconv.Atoi(params.Get("prev_id"))
		}
		if len(params["size"]) > 0 {
			size, _ = strconv.Atoi(params["size"][0])
		}
		readTODORequest := &model.ReadTODORequest{
			PrevID: int64(prevID),
			Size:   int64(size),
		}
		ctx := r.Context()
		res, err := h.Read(ctx, readTODORequest)
		if err != nil {
			log.Println(err)
			return
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "Failed to encode response as JSON", http.StatusInternalServerError)
			return
		}
	case "POST":
		createTODORequest := &model.CreateTODORequest{}
		err := json.NewDecoder(r.Body).Decode(createTODORequest)
		if err != nil {
			log.Println(err)
			return
		}

		if createTODORequest.Subject == "" {
			http.Error(w, "Subject cannot be empty", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		res, err := h.Create(ctx, createTODORequest)
		if err != nil {
			log.Println(err)
			return
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "Failed to encode response as JSON", http.StatusInternalServerError)
			return
		}
	case "PUT":
		updateTODORequest := &model.UpdateTODORequest{}
		err := json.NewDecoder(r.Body).Decode(updateTODORequest)
		if err != nil {
			log.Println(err)
			return
		}

		if updateTODORequest.ID == 0 {
			http.Error(w, "ID cannot be 0", http.StatusBadRequest)
			return
		}
		if updateTODORequest.Subject == "" {
			http.Error(w, "Subject cannot be empty", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		res, err := h.Update(ctx, updateTODORequest)
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "Failed to encode response as JSON", http.StatusInternalServerError)
			return
		}
	case "DELETE":
		deleteTODORequest := &model.DeleteTODORequest{}
		err := json.NewDecoder(r.Body).Decode(deleteTODORequest)

		if err != nil {
			log.Println(err)
			return
		}

		if len(deleteTODORequest.IDs) == 0 {
			http.Error(w, "IDs cannot be empty", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		res, err := h.Delete(ctx, deleteTODORequest)
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "Failed to encode response as JSON", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	todo, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}

	return &model.CreateTODOResponse{TODO: *todo}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	todos, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	if err != nil {
		return nil, err
	}
	return &model.ReadTODOResponse{TODOs: todos}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	todo, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.UpdateTODOResponse{TODO: *todo}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	err := h.svc.DeleteTODO(ctx, req.IDs)
	if err != nil {
		return nil, err
	}
	return &model.DeleteTODOResponse{}, nil
}
