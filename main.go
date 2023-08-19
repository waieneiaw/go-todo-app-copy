package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/waieneiaw/go-todo-app-copy/config"
	"golang.org/x/sync/errgroup"
)

func run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.New()
	if err != nil {
		return err
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", cfg.Port, err)
	}

	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with: %v", url)

	s := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// time.Sleep(5 * time.Second) // コマンドラインで実験するため
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}

	eg, ctx := errgroup.WithContext(ctx)

	// 別ゴルーチンでHTTPサーバーを起動する
	eg.Go(func() error {
		if err := s.Serve(l); err != nil &&
			// Http.ErrServerClosedはシャットダウンが正常終了したことを示すため除外する
			err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}

		return nil
	})

	<-ctx.Done() // チャネルからの通知を待機する

	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}

	// 別ゴルーチンの終了を待つ
	return eg.Wait()
}

func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}
