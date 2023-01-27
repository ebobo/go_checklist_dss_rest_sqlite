package model

// Item represents a list iteam metadata.
type Item struct {
	// ID is the unique identifier of the item.
	ID int32 `json:"id" db:"id"`
	// Name is the name of the item.
	Name string `json:"name" db:"name"`
	// Position of item in a list.
	Position int32 `json:"position" db:"position"`
	// Tag is the tag of the item for conmmunicate with the other server.
	Tag string `json:"tag" db:"tag"`
	// Status is the status of the item. on or off.
	Status bool `json:"status" db:"status"`
}
