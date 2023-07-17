package itemService

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"hezzl/internal/core"
)

const (
	AllItemsCacheKey = "allItems"
)

type ItemRepository interface {
	CreateItem(name string) (core.Item, error)
	GetItemById(id int) (core.Item, error)
	GetAllItems() ([]core.Item, error)
	UpdateItem(update core.Item) error
	RemoveItem(id, campaignId int) (core.Item, error)
}

type MessageBroker interface {
	Publish(subj string, data []byte) error
}

type Cache interface {
	Set(key string, value interface{}) error
	Get(key string) ([]byte, error)
	Delete(key string) error
}

type ItemService struct {
	ItemRepository
	MessageBroker
	Cache
}

func NewItemService(i ItemRepository, mb MessageBroker, c Cache) *ItemService {

	return &ItemService{
		ItemRepository: i,
		MessageBroker:  mb,
		Cache:          c,
	}
}

func (i ItemService) CreateItem(name string) (core.Item, error) {

	item, err := i.ItemRepository.CreateItem(name)
	if err != nil {
		return core.Item{}, fmt.Errorf("error while trying to create item: %w", err)
	}

	if err := i.sendEvent(item); err != nil {
		return core.Item{}, fmt.Errorf("error while trying to send event: %w", err)
	}

	return item, nil
}

func (i ItemService) GetAllItems() ([]core.Item, error) {
	var items []core.Item

	res, err := i.Cache.Get(AllItemsCacheKey)
	if err != nil {
		if err != redis.Nil {
			return nil, fmt.Errorf("error while trying to get items: %w", err)
		}
	}

	if res != nil {
		if err := json.Unmarshal(res, &items); err != nil {
			return items, fmt.Errorf("error while trying to unmarshal cache result: %w", err)
		}
		return items, nil
	}

	items, err = i.ItemRepository.GetAllItems()
	if err != nil {
		return nil, fmt.Errorf("error while trying to get all items: %w", err)
	}

	if err := i.Cache.Set(AllItemsCacheKey, items); err != nil {
		return nil, fmt.Errorf("error while trying to set items: %w", err)
	}

	return items, nil
}

func (i ItemService) UpdateItem(update core.Item) (core.Item, error) {
	item, err := i.ItemRepository.GetItemById(update.Id)
	if err != nil {
		return core.Item{}, fmt.Errorf("error while trying to get item by id: %w", err)
	}

	item.Name = update.Name
	if update.Description != "-" {
		item.Description = update.Description
	}

	if err := i.ItemRepository.UpdateItem(item); err != nil {
		return core.Item{}, fmt.Errorf("error while trying to update item: %w", err)
	}

	if err := i.Cache.Delete(AllItemsCacheKey); err != nil {
		return core.Item{}, fmt.Errorf("error while trying to invalidate cache: %w", err)
	}

	if err := i.sendEvent(item); err != nil {
		return core.Item{}, fmt.Errorf("error while trying to send event: %w", err)
	}

	return item, nil
}

func (i ItemService) RemoveItem(id, campaignId int) (core.Item, error) {
	item, err := i.ItemRepository.RemoveItem(id, campaignId)
	if err != nil {
		return core.Item{}, fmt.Errorf("error while trying to remove item: %w", err)
	}

	if err := i.sendEvent(item); err != nil {
		return core.Item{}, fmt.Errorf("error while trying to send event: %w", err)
	}

	return item, nil
}

func (i ItemService) sendEvent(item core.Item) error {
	jsonData, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("error while trying to marshal updated item: %w", err)
	}

	if err := i.MessageBroker.Publish("items_update", jsonData); err != nil {
		return fmt.Errorf("error while trying to ")
	}

	return nil
}
