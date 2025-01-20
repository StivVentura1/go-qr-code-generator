package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"time"

	"github.com/skip2/go-qrcode"
)

type encodeQR struct {
	Timestamp   time.Time `json:"timestamp"`
	AI          string    `json:"ai"`
	MintingID   string    `json:"minting_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

func generateQRCodeWithLogo(data encodeQR, qrPath, logoPath, outputPath string) error {
	// Convert struct to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal struct: %v", err)
	}

	// Generate QR code and save as PNG
	err = qrcode.WriteFile(string(jsonData), qrcode.Highest, 1024, qrPath)
	if err != nil {
		return fmt.Errorf("failed to generate QR code: %v", err)
	}

	// Open generated QR code (expecting PNG format)
	qrFile, err := os.Open(qrPath)
	if err != nil {
		return fmt.Errorf("failed to open QR code file: %v", err)
	}
	defer qrFile.Close()

	qrImg, err := png.Decode(qrFile)
	if err != nil {
		return fmt.Errorf("failed to decode QR code image: %v", err)
	}

	// Open logo image (expecting JPEG format)
	logoFile, err := os.Open(logoPath)
	if err != nil {
		return fmt.Errorf("failed to open logo file: %v", err)
	}
	defer logoFile.Close()

	logoImg, err := jpeg.Decode(logoFile)
	if err != nil {
		return fmt.Errorf("failed to decode logo image: %v", err)
	}

	// Calculate position to place the logo in the center
	qrBounds := qrImg.Bounds()
	logoBounds := logoImg.Bounds()
	posX := (qrBounds.Dx() - logoBounds.Dx()) / 2
	posY := (qrBounds.Dy() - logoBounds.Dy()) / 2

	// Merge QR code and logo
	outputImg := image.NewRGBA(qrBounds)
	draw.Draw(outputImg, qrBounds, qrImg, image.Point{}, draw.Src)
	draw.Draw(outputImg, logoBounds.Add(image.Point{posX, posY}), logoImg, image.Point{}, draw.Over)

	// Save final image as JPEG
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outputFile.Close()

	options := &jpeg.Options{Quality: 100}
	err = jpeg.Encode(outputFile, outputImg, options)
	if err != nil {
		return fmt.Errorf("failed to encode final image: %v", err)
	}

	fmt.Println("QR code with logo successfully generated and saved to", outputPath)
	return nil
}

func main() {
	sampleData := encodeQR{
		Timestamp:   time.Now(),
		AI:          "123456",
		MintingID:   "MINT123",
		Name:        "Sample QR",
		Description: "This is a sample QR code",
	}

	err := generateQRCodeWithLogo(sampleData, "C:/Users/andre/go/src/go-qr-code-generator/qrcode/qrcode.png", "logo/20161226_163042_HDR.jpg", "C:/Users/andre/go/src/go-qr-code-generator/qrcode_with_logo/qrcodeWLogo.jpg")
	if err != nil {
		fmt.Println("Error:", err)
	}
}
