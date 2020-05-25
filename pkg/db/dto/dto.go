package dto

// DTO is an object that can interact with the DB.
type DTO interface {
	Upload() error
	Download() error
	Delete() error
}
