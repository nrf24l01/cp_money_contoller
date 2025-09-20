package dbsync

import (
	"log"
	"time"
)

func (h *Handler) StartSyncing() {
	log.Printf("Starting periodic task synchronization...")
	log.Printf("Full cleanup")
	if err := h.ClearAll(); err != nil {
		log.Printf("Error during full cleanup: %v", err)
		return
	}
	log.Printf("Cleanup done, starting sync")
	
	for {
		h.SyncTasks()
		log.Printf("Go sleep(%v), will wake up at %v", h.Cfg.SyncInterval, time.Now().Add(h.Cfg.SyncInterval).Format("15:04:05 02-01-06"))
		<-time.After(h.Cfg.SyncInterval)
	}
}