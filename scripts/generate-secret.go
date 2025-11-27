package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
)

func main() {
	// Generate cryptographically secure random bytes
	bytes := make([]byte, 64) // 64 bytes = 512 bits
	if _, err := rand.Read(bytes); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating random bytes: %v\n", err)
		os.Exit(1)
	}

	// Encode to base64 untuk mudah di-copy paste
	secret := base64.StdEncoding.EncodeToString(bytes)

	fmt.Println("╔══════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║              JWT SECRET GENERATOR                                    ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("✅ Cryptographically secure JWT_SECRET generated!")
	fmt.Println()
	fmt.Println("🔐 Your JWT_SECRET:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println(secret)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println("📋 Copy dan paste ke .env file:")
	fmt.Printf("   JWT_SECRET=%s\n", secret)
	fmt.Println()
	fmt.Println("⚠️  IMPORTANT:")
	fmt.Println("   - JANGAN commit secret ini ke Git!")
	fmt.Println("   - Simpan secret ini dengan aman")
	fmt.Println("   - Untuk production, gunakan secret manager (AWS Secrets, Vault, dll)")
	fmt.Println("   - Generate secret yang berbeda untuk setiap environment")
	fmt.Println()
	fmt.Printf("ℹ️  Secret length: %d characters (%d bytes)\n", len(secret), len(bytes))
	fmt.Println()
}
