package dbsync

func (h *Handler) ClearAll() error {
	if err := h.RMQ.Purge(); err != nil {
		return err
	}
	if err := h.Redis.Purge(); err != nil {
		return err
	}
	return nil
}