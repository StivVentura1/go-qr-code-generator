package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"

	"github.com/skip2/go-qrcode"

	//"image/draw"
	"image/png"
	"os"
	"time"

	"golang.org/x/image/draw"
)

type encodeQR struct {
	Timestamp   time.Time `json:"timestamp"`
	AI          string    `json:"ai"`
	MintingID   string    `json:"minting_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

// Function to resize the logo to fit within the specified width and height
func resizeImage(img image.Image, width, height int) image.Image {
	resized := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(resized, resized.Bounds(), img, img.Bounds(), draw.Over, nil)
	return resized
}

func generateQRCodeWithLogo(data encodeQR, qrPath, logoPath, outputPath string) error {
	// Convert struct to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal struct: %v", err)
	}

	// Generate QR code and save as PNG
	err = qrcode.WriteFile(string(jsonData), qrcode.Highest, 2048, qrPath)
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

	// Open logo image (expecting PNG with transparency)
	logoFile, err := os.Open(logoPath)
	if err != nil {
		return fmt.Errorf("failed to open logo file: %v", err)
	}
	defer logoFile.Close()

	logoImg, err := png.Decode(logoFile)
	if err != nil {
		return fmt.Errorf("failed to decode logo image: %v", err)
	}

	// Define QR code bounds
	qrBounds := qrImg.Bounds()

	// Define the 500x500 white rectangle in the center
	centerSize := 500
	centerX := (qrBounds.Dx() - centerSize) / 2
	centerY := (qrBounds.Dy() - centerSize) / 2

	// Create an RGBA image from the QR code and clear the center
	outputImg := image.NewRGBA(qrBounds)
	draw.Draw(outputImg, qrBounds, qrImg, image.Point{}, draw.Src)

	// Fill the center area with white color (400x400 pixels)
	white := color.RGBA{255, 255, 255, 255}
	centerRect := image.Rect(centerX, centerY, centerX+centerSize, centerY+centerSize)
	draw.Draw(outputImg, centerRect, &image.Uniform{white}, image.Point{}, draw.Src)

	// Resize the logo to fit within the 400x400 space
	resizedLogo := resizeImage(logoImg, centerSize, centerSize)

	// Overlay the resized logo onto the cleared white space
	draw.Draw(outputImg, centerRect, resizedLogo, image.Point{}, draw.Over)

	// Save final image as PNG (to preserve transparency)
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outputFile.Close()

	err = png.Encode(outputFile, outputImg)
	if err != nil {
		return fmt.Errorf("failed to encode final image: %v", err)
	}

	fmt.Println("QR code with white center (400x400) and logo successfully saved to", outputPath)
	return nil
}

func main() {
	sampleData := encodeQR{
		Timestamp:   time.Now(),
		AI:          "123456",
		MintingID:   "MINT123",
		Name:        "TROMBOLOTTO MINTED SUCCESFULLY",
		Description: "Ciao Papino, W il TROMBOLOTTO!",
	}

	err := generateQRCodeWithLogo(
		sampleData,
		"C:/Users/andre/go/src/go-qr-code-generator/qrcode/qrcode.png",                           // Path to save QR code
		"C:/Users/andre/go/src/go-qr-code-generator/logo/400_5e8a4aba11076-removebg-preview.png", // Path to logo file (JPEG format)
		"C:/Users/andre/go/src/go-qr-code-generator/qrcode_with_logo/qrcodeWLogo.jpg",            // Output QR with logo
	)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
