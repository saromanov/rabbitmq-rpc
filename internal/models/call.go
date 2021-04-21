package models

// Call defines model for pendign call
type Call struct {
	Done chan bool
	Data []byte
}