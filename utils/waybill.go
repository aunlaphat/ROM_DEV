package utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"
)

type WaybillData struct {
	RefNo    string
	Queue    string
	Products []Products
}

type Products struct {
	SKU  string
	Name string
	QTY  int
}

func GenerateWaybill(input WaybillData) (string, error) {
	queue := input.Queue
	refNo := input.RefNo
	fontBold := "TahomaBold"
	font := "Tahoma"

	var buf bytes.Buffer

	png, err := qrcode.Encode(queue, qrcode.Medium, 256)
	if err != nil {
		log.Fatal(err)
	}

	_, err = buf.Write(png)
	if err != nil {
		log.Fatal(err)
	}

	tmpFile, err := os.CreateTemp("", "qrcode-*.png")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	tmpFile.Close()

	pdf := gofpdf.New("P", "mm", "A6", "")
	pdf.AddUTF8Font(font, "", "./assets/tahoma.ttf")
	pdf.AddUTF8Font(fontBold, "", "./assets/tahomabd.ttf")
	pdf.SetFont(font, "", 12)

	pdf.SetHeaderFunc(func() {
		pdf.Image(tmpFile.Name(), 78, 1, 25, 25, false, "", 0, "")
		pdf.Ln(1)

		pdf.SetXY(3, 3)
		pdf.SetFont(fontBold, "", 10)
		pdf.Cell(0, 5, "Order No: "+refNo)

		pdf.SetXY(3, pdf.GetY()+5)
		pdf.SetFont(font, "", 10)
		pdf.Cell(0, 5, "Queue: "+queue)

		pdf.SetXY(3, pdf.GetY()+10)
		pdf.Cell(0, 5, "WH From: MMT_BEWELL")

		pdf.SetXY(3, pdf.GetY()+5)
		pdf.Cell(0, 5, "WH From: RBN")

		pdf.SetFont(fontBold, "", 8)
		pdf.SetXY(5, pdf.GetY()+10)
		pdf.CellFormat(5, 5, "#", "1", 0, "C", false, 0, "")
		pdf.CellFormat(35, 5, "SKU", "1", 0, "C", false, 0, "")
		pdf.CellFormat(47, 5, "ชื่อสินค้า", "1", 0, "C", false, 0, "")
		pdf.CellFormat(7, 5, "QTY", "1", 0, "C", false, 0, "")
		y := pdf.GetY() + 5
		for i, product := range input.Products {
			pdf.SetXY(5, y)
			pdf.SetFont(font, "", 8)
			pdf.CellFormat(5, 5, strconv.Itoa(i+1), "1", 0, "C", false, 0, "")
			pdf.CellFormat(35, 5, product.SKU, "1", 0, "L", false, 0, "")
			pdf.CellFormat(47, 5, product.Name, "1", 0, "L", false, 0, "")
			pdf.CellFormat(7, 5, strconv.Itoa(product.QTY), "1", 0, "C", false, 0, "")
			y += 5.00
		}
	})

	pdf.AddPage()

	/*** FOR PROD ONLY!!! ***/
	var pdfBuffer bytes.Buffer
	err = pdf.Output(&pdfBuffer)
	if err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	// Upload the PDF file
	err = UploadPDF(pdfBuffer, queue)
	if err != nil {
		return "", fmt.Errorf("failed to upload PDF: %w", err)
	}

	/*** FOR TEST ONLY!!! ***/
	// err = pdf.OutputFileAndClose(queue + ".pdf")
	// if err != nil {
	// 	return "", fmt.Errorf("error to save file %v", err)
	// }
	return "success", nil

}

func UploadPDF(pdfBuffer bytes.Buffer, claimID string) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create the form file
	part, err := writer.CreateFormFile("file", claimID+".pdf")
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}
	_, err = io.Copy(part, &pdfBuffer)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	// Write the additional fields
	err = writer.WriteField("order_no", claimID)
	if err != nil {
		return fmt.Errorf("failed to write order_no field: %w", err)
	}
	err = writer.WriteField("tracking_no", "-")
	if err != nil {
		return fmt.Errorf("failed to write tracking_no field: %w", err)
	}

	// Close the writer
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}
	url := os.Getenv("CHECKER_API")
	// Create the HTTP request
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	defer resp.Body.Close()

	// Print response body for debugging
	responseBody, _ := io.ReadAll(resp.Body)
	fmt.Println("Response Body:", string(responseBody))

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to upload file: status code %d", resp.StatusCode)
	}

	return nil
}
