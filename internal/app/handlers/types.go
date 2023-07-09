package handlers

//go:generate mockgen -destination=mocks/get_shortened_url.go -package=mocks . RequestContext
type RequestContext interface {
	GetRawData() ([]byte, error)
	JSON(code int, obj any)
}

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}
