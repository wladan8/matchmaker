package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.com/matchmaker/config"
	"gitlab.com/matchmaker/internal/handler"
	"gitlab.com/matchmaker/internal/matchmaker"
)

func Start(cfg *config.ServerConfig) {
	m := matchmaker.New(&matchmaker.Config{
		GroupSize:       cfg.GroupSize,
		DiffSkill:       cfg.DiffSkill,
		DiffLatency:     cfg.DiffLatency,
		TickerFrequency: time.Duration(cfg.TickerFrequency),
	})
	ctx, matchmakerCancel := context.WithCancel(context.Background())
	m.GatherGroupsProcessing(ctx)
	h := handler.New(m)
	h.ReceiveGroupsFromMatchmaker(ctx)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /users/", h.UserToPool)

	server := &http.Server{
		Handler:           mux,
		Addr:              fmt.Sprintf(":%d", cfg.ServerPort),
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
	}
	shutdownChan := make(chan struct{}, 1)
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}

		time.Sleep(1 * time.Second)
		slog.Info("stopped handle new connections.")
		shutdownChan <- struct{}{}
	}()

	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM)
	<-exitChan
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shutdown error: %v", err)
	}
	matchmakerCancel()
	<-shutdownChan
	slog.Info("graceful shutdown complete")
}
