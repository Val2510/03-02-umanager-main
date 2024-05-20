package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"gitlab.com/robotomize/gb-golang/homework/03-02-umanager/pkg/pb"
)

func newLinksHandler(linksClient linksClient) *linksHandler {
	return &linksHandler{client: linksClient}
}

type linksHandler struct {
	client linksClient
}

func (h *linksHandler) GetLinks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	resp, err := h.client.ListLinks(ctx, &pb.Empty{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp.Links); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *linksHandler) PostLinks(w http.ResponseWriter, r *http.Request) {
	var req pb.CreateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	_, err := h.client.CreateLink(ctx, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *linksHandler) DeleteLinksId(w http.ResponseWriter, r *http.Request, id string) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	_, err := h.client.DeleteLink(ctx, &pb.DeleteLinkRequest{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h Handler) GetLinksId(w http.ResponseWriter, r *http.Request, id string) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	resp, err := h.linksHandler.client.GetLink(ctx, &pb.GetLinkRequest{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *linksHandler) PutLinksId(w http.ResponseWriter, r *http.Request, id string) {
	var req pb.UpdateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.Id = id

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	_, err := h.client.UpdateLink(ctx, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *linksHandler) GetLinksUserUserID(w http.ResponseWriter, r *http.Request, userID string) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	resp, err := h.client.GetLinkByUserID(ctx, &pb.GetLinksByUserId{UserId: userID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp.Links); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
