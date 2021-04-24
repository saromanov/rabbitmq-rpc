package publish

// Unmarshaller defines unmarshaler for client 
type Unmarshaller interface {
	Do([]byte, interface{}) error 
}