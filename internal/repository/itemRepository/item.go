package itemRepository

import (
	"database/sql"
	"errors"
	"fmt"
	"hezzl/internal/core"
)

var (
	ItemNotFound = errors.New("item not found")
)

type ItemPostgres struct {
	db *sql.DB
}

func NewItemPostgres(db *sql.DB) *ItemPostgres {
	return &ItemPostgres{
		db: db,
	}
}

func (i ItemPostgres) CreateItem(name string) (core.Item, error) {
	var id int
	if err := i.db.QueryRow("INSERT INTO items(name) VALUES ($1) RETURNING id", name).Scan(&id); err != nil {
		return core.Item{}, fmt.Errorf("error while trying to insert into items: %w", err)
	}

	return i.GetItemById(id)
}

func (i ItemPostgres) GetItemById(id int) (core.Item, error) {
	var item core.Item

	res, err := i.db.Query("SELECT * FROM items WHERE id=$1", id)
	if err != nil {
		return item, fmt.Errorf("error while trying to query item: %w", err)
	}
	defer res.Close()

	for res.Next() {
		if err := res.Scan(&item.Id, &item.CampaignId, &item.Name, &item.Description, &item.Priority, &item.Removed, &item.CreatedAt); err != nil {
			return item, fmt.Errorf("error while trying to scan item query result: %w", err)
		}
	}

	if item.Id == 0 {
		return core.Item{}, ItemNotFound
	}

	return item, nil
}

func (i ItemPostgres) GetAllItems() ([]core.Item, error) {
	var allItems []core.Item

	res, err := i.db.Query("SELECT * FROM items")
	if err != nil {
		return allItems, fmt.Errorf("error while trying to query all items: %w", err)
	}
	defer res.Close()

	for res.Next() {
		var item core.Item
		if err := res.Scan(&item.Id, &item.CampaignId, &item.Name, &item.Description, &item.Priority, &item.Removed, &item.CreatedAt); err != nil {
			return allItems, fmt.Errorf("error while trying to scan all items query result: %w", err)
		}
		allItems = append(allItems, item)
	}

	return allItems, nil
}

func (i ItemPostgres) UpdateItem(update core.Item) error {
	tx, err := i.db.Begin()
	if err != nil {
		return fmt.Errorf("error while trying to start transaction: %w", err)
	}

	if _, err := tx.Exec("UPDATE items SET name=$1, description=$2 WHERE id=$3", update.Name, update.Description, update.Id); err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("error while trying to rollback transaction")
		}
		return fmt.Errorf("error while trying to update items: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error while trying to commit transaction")
	}

	return nil
}

func (i ItemPostgres) RemoveItem(id, campaignId int) (core.Item, error) {
	if _, err := i.db.Exec("UPDATE items SET removed=TRUE WHERE id=$1 AND campaign_id=$2", id, campaignId); err != nil {
		return core.Item{}, fmt.Errorf("error while trying to remove item: %w", err)
	}

	return i.GetItemById(id)
}
