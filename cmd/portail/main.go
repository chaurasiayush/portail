package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"github.com/chaurasiayush/portail/internal/config"
	"github.com/chaurasiayush/portail/internal/forwarder"


	"github.com/fsnotify/fsnotify"
)

func startForwarders(ctx context.Context, config *config.Config, wg *sync.WaitGroup) {
	for _, rule := range config.Forwards {
		wg.Add(1)
		switch rule.Protocol {
		case "tcp":
			go forwarder.StartTCPForward(ctx, wg, rule)
		case "udp":
			go forwarder.StartUDPForward(ctx, wg, rule)
		default:
			log.Printf("Unsupported protocol: %s", rule.Protocol)
			wg.Done()
		}
	}
}

func StartTCPForward(ctx context.Context, wg *sync.WaitGroup, rule config.ForwardRule) {
	panic("unimplemented")
}

func main() {
	configPath := flag.String("config", "config.yaml", "Path to YAML config file")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	startForwarders(ctx, cfg, &wg)

	// Setup config watcher
	go watchConfig(*configPath, &wg, &cancel)

	// Handle graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	log.Println("Shutting down...")
	cancel()
	wg.Wait()
	log.Println("Stopped.")
}

func watchConfig(path string, wg *sync.WaitGroup, cancelOld *context.CancelFunc) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Watcher error:", err)
	}
	defer watcher.Close()

	// configDir := "."
	// if stat, err := os.Stat(path); err == nil && !stat.IsDir() {
	// 	configDir = "."
	// }

	err = watcher.Add(path)
	if err != nil {
		log.Fatalf("Watcher add error: %v", err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				log.Println("Config changed, reloading...")
				newCfg, err := config.LoadConfig(path)
				if err != nil {
					log.Println("Reload failed:", err)
					continue
				}

				// Cancel old context
				(*cancelOld)()
				wg.Wait()

				// Start new context
				ctx, newCancel := context.WithCancel(context.Background())
				*cancelOld = newCancel
				startForwarders(ctx, newCfg, wg)
			}
		case err, ok := <-watcher.Errors:
			if ok {
				log.Println("Watcher error:", err)
			}
		}
	}
}
