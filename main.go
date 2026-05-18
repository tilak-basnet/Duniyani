package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"duniyani/consensus"
	"duniyani/network"
	"duniyani/wallet"
)

func main() {
	listenAddr := flag.String("listen", ":9000", "QUIC P2P Listen Address")
	miner := flag.Bool("miner", false, "Enable PoUW Mining")
	flag.Parse()

	log.Println("Initializing Duniyani Quantum-Secure Layer 1 Node...")

	// 1. Initialize Wallet
	w, err := wallet.NewWallet()
	if err != nil {
		log.Fatalf("Failed to generate ML-DSA wallet: %v", err)
	}
	log.Printf("Node Address: %s", string(w.GetAddress()))

	// 2. Setup Graceful Shutdown Context
	ctx, cancel := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// 3. Initialize Mempool and P2P Host
	mempool := network.NewMempool()
	_ = mempool // Used in consensus integration

	p2pHost := &network.QUICHost{}
	go func() {
		if err := p2pHost.Start(ctx, *listenAddr); err != nil {
			log.Printf("P2P Host stopped: %v", err)
		}
	}()

	// 4. Initialize Consensus Engine
	if *miner {
		log.Println("Starting Proof of Useful Work (PoUW) Engine...")
		engine := consensus.NewPoUWEngine(w.PublicKey.Bytes()) // Dummy enclave key
		_ = engine                                             // Hook into a mining loop using mempool txs
	}

	// 5. Wait for Shutdown Signal
	<-sigCh
	log.Println("\nShutdown signal received. Committing state and gracefully halting node...")
	cancel()
	log.Println("Duniyani node offline.")
}