package booking

// AdminHandler handles admin booking operations.
// The time-grid endpoint is handled by the venue module's AdminGetTimeGrid handler.
type AdminHandler struct {
	svc *Service
}

func NewAdminHandler(svc *Service) *AdminHandler {
	return &AdminHandler{svc: svc}
}
