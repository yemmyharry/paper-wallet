package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"
	"log"
	"time"
)

func generatePaperWallet() (*gofpdf.Fpdf, string, string, error) {
	privateKey, err := btcec.NewPrivateKey()
	if err != nil {
		log.Fatalf("failed to generate a private key: %v", err)
	}

	pubKey := privateKey.PubKey()
	address, err := btcutil.NewAddressPubKey(pubKey.SerializeCompressed(), &chaincfg.MainNetParams)
	if err != nil {
		log.Fatalf("Error generating address: %v", err)
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Bitcoin Paper Wallet")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, "Bitcoin Address: "+address.String())
	pdf.Ln(10)
	pdf.Cell(0, 10, "Private Key: "+hex.EncodeToString(privateKey.Serialize()))

	return pdf, address.String(), hex.EncodeToString(privateKey.Serialize()), nil
}

func generateQRCode(data string, filename string) error {
	err := qrcode.WriteFile(data, qrcode.Medium, 256, filename)
	if err != nil {
		return fmt.Errorf("failed to generate QR code: %v", err)
	}
	return nil
}

func main() {
	outputDirFlag := flag.String("outputDir", ".", "Output directory for paper wallet PDF and QR code images")
	flag.Parse()

	pdf, address, privateKey, err := generatePaperWallet()
	if err != nil {
		log.Fatalf("Error generating paper wallet: %v", err)
	}

	pdfFilename := fmt.Sprintf("bitcoin_paper_wallet_%s.pdf", time.Now().Format("2006-01-02_15-04-05"))
	pdfFilePath := fmt.Sprintf("%s/%s", *outputDirFlag, pdfFilename)
	if err := pdf.OutputFileAndClose(pdfFilePath); err != nil {
		log.Fatalf("Error saving paper wallet PDF: %v", err)
	}
	fmt.Printf("Paper wallet PDF saved to: %s\n", pdfFilePath)

	if err := generateQRCode(address, fmt.Sprintf("%s/bitcoin_address.png", *outputDirFlag)); err != nil {
		log.Fatalf("Error generating QR code for Bitcoin address: %v", err)
	}
	fmt.Printf("QR code for Bitcoin address saved to: %s/bitcoin_address.png\n", *outputDirFlag)
	if err := generateQRCode(privateKey, fmt.Sprintf("%s/private_key.png", *outputDirFlag)); err != nil {
		log.Fatalf("Error generating QR code for private key: %v", err)
	}
	fmt.Printf("QR code for private key saved to: %s/private_key.png\n", *outputDirFlag)
}
