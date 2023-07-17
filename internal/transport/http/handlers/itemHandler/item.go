package itemHandler

import (
	"encoding/json"
	"hezzl/internal/core"
	"log"
	"net/http"
	"strconv"
)

type ItemService interface {
	CreateItem(name string) (core.Item, error)
	GetAllItems() ([]core.Item, error)
	UpdateItem(update core.Item) (core.Item, error)
	RemoveItem(id, campaignId int) (core.Item, error)
}

type ItemHandler struct {
	ItemService
	Logger *log.Logger
}

func NewItemHandler(i ItemService, log *log.Logger) *ItemHandler {
	return &ItemHandler{
		ItemService: i,
		Logger:      log,
	}
}

func (i ItemHandler) RegisterItemHandlers(router *http.ServeMux) {
	router.HandleFunc("/item/create", i.createItem)
	router.HandleFunc("/items/list", i.getItemsList)
	router.HandleFunc("/item/update", i.updateItem)
	router.HandleFunc("/item/remove", i.removeItem)
}

func (i ItemHandler) createItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Write(i.responseWithError(http.StatusMethodNotAllowed, MethodNotAllowed))
		return
	}

	var newItem core.Item
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		i.Logger.Printf("error: %v", err)
		w.Write(i.responseWithError(http.StatusBadRequest, BadJSON))
		return
	}

	newItem, err := i.ItemService.CreateItem(newItem.Name)
	if err != nil {
		i.Logger.Printf("error: %v", err)
		w.Write(i.responseWithError(http.StatusInternalServerError, UnexpectedInternalError))
		return
	}

	responseJSON, err := json.Marshal(&newItem)
	if err != nil {
		i.Logger.Printf("error: %v", err)
		w.Write(i.responseWithError(http.StatusInternalServerError, UnexpectedInternalError))
		return
	}
	w.Write(responseJSON)
}

func (i ItemHandler) getItemsList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Write(i.responseWithError(http.StatusMethodNotAllowed, MethodNotAllowed))
		return
	}

	items, err := i.ItemService.GetAllItems()
	if err != nil {
		i.Logger.Printf("error: %v", err)
		w.Write(i.responseWithError(http.StatusInternalServerError, UnexpectedInternalError))
		return
	}

	responseJSON, err := json.Marshal(&items)
	if err != nil {
		i.Logger.Printf("error: %v", err)
		w.Write(i.responseWithError(http.StatusInternalServerError, UnexpectedInternalError))
		return
	}
	w.Write(responseJSON)
}

func (i ItemHandler) updateItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		w.Write(i.responseWithError(http.StatusMethodNotAllowed, MethodNotAllowed))
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		i.Logger.Printf("error: %v", err)
		w.Write(i.responseWithError(http.StatusBadRequest, BadQuery))
		return
	}

	campaignId, err := strconv.Atoi(r.URL.Query().Get("campaignId"))
	if err != nil {
		i.Logger.Printf("error: %v", err)
		w.Write(i.responseWithError(http.StatusBadRequest, BadQuery))
		return
	}

	var update core.Item
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		i.Logger.Printf("error: %v", err)
		w.Write(i.responseWithError(http.StatusBadRequest, BadJSON))
		return
	}
	update.Id = id
	update.CampaignId = campaignId

	if !i.validate(update) {
		w.Write(i.responseWithError(http.StatusBadRequest, BadJSON))
		return
	}

	update, err = i.ItemService.UpdateItem(update)
	if err != nil {
		i.Logger.Printf("error: %v", err)
		if unwrapErr(err, core.ItemNotFound) {
			w.Write(i.responseWithError(http.StatusNotFound, NotFound))
			return
		}
		w.Write(i.responseWithError(http.StatusInternalServerError, UnexpectedInternalError))
		return
	}

	responseJSON, err := json.Marshal(&update)
	if err != nil {
		i.Logger.Printf("error: %v", err)
		w.Write(i.responseWithError(http.StatusInternalServerError, UnexpectedInternalError))
		return
	}
	w.Write(responseJSON)
}

func (i ItemHandler) removeItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		w.Write(i.responseWithError(http.StatusMethodNotAllowed, MethodNotAllowed))
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		i.Logger.Printf("error: %v", err)
		w.Write(i.responseWithError(http.StatusBadRequest, BadQuery))
		return
	}

	campaignId, err := strconv.Atoi(r.URL.Query().Get("campaignId"))
	if err != nil {
		i.Logger.Printf("error: %v", err)
		w.Write(i.responseWithError(http.StatusBadRequest, BadQuery))
		return
	}

	item, err := i.ItemService.RemoveItem(id, campaignId)
	if err != nil {
		i.Logger.Printf("error: %v", err)
		if unwrapErr(err, core.ItemNotFound) {
			w.Write(i.responseWithError(http.StatusNotFound, NotFound))
			return
		}
		w.Write(i.responseWithError(http.StatusInternalServerError, UnexpectedInternalError))
		return
	}

	responseJSON, err := json.Marshal(&item)
	if err != nil {
		i.Logger.Printf("error: %v", err)
		w.Write(i.responseWithError(http.StatusInternalServerError, UnexpectedInternalError))
		return
	}
	w.Write(responseJSON)
}

func (i ItemHandler) validate(in core.Item) bool {
	if in.Name == "" {
		return false
	}

	return true
}
