package sqlitestore

import (
	"github.com/ebobo/dss_checklist/pkg/model"
)

func (s *SqliteStore) AddItem(item model.Item) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.NamedExec(
		`INSERT INTO items (
			id,
			name,
			position,
			tag,
			status)
		 VALUES(
			:id,
			:name,
			:position,
			:tag,
			:status)`, item)
	return err
}

func (s *SqliteStore) GetItem(id string) (model.Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var item model.Item
	return item, s.db.QueryRowx("SELECT * FROM keys WHERE id = ?", id).StructScan(&item)
}

func (s *SqliteStore) SetItemStatus(id string, status bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return CheckForZeroRowsAffected(s.db.Exec("UPDATE items SET status = ? WHERE id = ?", status, id))
}

func (s *SqliteStore) ListItems() ([]model.Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var items []model.Item
	return items, s.db.Select(&items, "SELECT * FROM items")
}

func (s *SqliteStore) DeleteItem(id string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, err := s.db.Exec("DELETE FROM items WHERE id = ?", id)
	return err
}

// func (s *SqliteStore) PrintItems() error {
// 	s.mu.RLock()
// 	defer s.mu.RUnlock()

// 	log.Println("Printing items")

// 	row, err := s.db.Query("SELECT * FROM items ORDER BY name")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer row.Close()
// 	for row.Next() { // Iterate and fetch the records from result cursor
// 		var id string
// 		var name string
// 		var position int
// 		var tag string
// 		var status bool
// 		row.Scan(&id, &name, &position, &tag, &status)
// 		log.Println("Item: ", id, " ", name)
// 	}
// 	return nil
// }
