package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/signintech/gopdf"
)

// ProseMirrorNode represents a node in ProseMirror JSON
type ProseMirrorNode struct {
	Type    string            `json:"type"`
	Content []ProseMirrorNode `json:"content,omitempty"`
	Text    string            `json:"text,omitempty"`
	Marks   []Mark            `json:"marks,omitempty"`
	Attrs   Attributes        `json:"attrs,omitempty"`
}

// Mark represents text styles (bold, italic, underline, etc.)
type Mark struct {
	Type  string     `json:"type"`
	Attrs Attributes `json:"attrs,omitempty"`
}

// Attributes hold additional styling info
type Attributes struct {
	Level int    `json:"level,omitempty"` // for headings
	Color string `json:"color,omitempty"` // for text color
}

// NewPdfDocument creates a PDF with gopdf, adds basic styles, and returns as byte slice
func NewPdfDocument(content []ProseMirrorNode) ([]byte, error) {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()
	pdf.SetMargins(10, 10, 10, 10)

	// Set up fonts
	if err := pdf.AddTTFFont("regular", "./fonts/arial/ARIAL.TTF"); err != nil {
		return nil, err
	}
	if err := pdf.AddTTFFont("bold", "./fonts/arial/ARIALBD.TTF"); err != nil {
		return nil, err
	}
	if err := pdf.AddTTFFont("italic", "./fonts/arial/ARIALI.TTF"); err != nil {
		return nil, err
	}

	pdf.SetFont("regular", "", 12)
	for _, node := range content {
		err := renderNode(&pdf, node)
		if err != nil {
			return nil, err
		}
	}

	// Write PDF to a byte slice
	var pdfBytes []byte
	pdfBytes, err := pdf.GetBytesPdfReturnErr()
	if err != nil {
		return nil, err
	}
	return pdfBytes, nil
}

// renderNode recursively renders ProseMirror nodes to the PDF
func renderNode(pdf *gopdf.GoPdf, node ProseMirrorNode) error {
	switch node.Type {
	case "paragraph":
		return renderParagraph(pdf, node)
	case "heading":
		return renderHeading(pdf, node)
	case "bullet_list":
		return renderList(pdf, node, "bullet")
	case "ordered_list":
		return renderList(pdf, node, "ordered")
	default:
		return nil
	}
}

// renderText applies text styling and writes to PDF
func renderText(pdf *gopdf.GoPdf, text string, marks []Mark) error {
	font := "regular"
	size := 12.0
	color := "black"

	for _, mark := range marks {
		switch mark.Type {
		case "bold":
			font = "bold"
		case "italic":
			font = "italic"
		case "underline":
			// Custom handling if needed
		case "textColor":
			color = mark.Attrs.Color
		}
	}

	pdf.SetFont(font, "", size)
	if color == "blue" {
		pdf.SetTextColor(0, 0, 255)
	} else {
		pdf.SetTextColor(0, 0, 0) // default black
	}

	pdf.Cell(nil, text)
	return nil
}

// renderParagraph formats paragraphs
func renderParagraph(pdf *gopdf.GoPdf, node ProseMirrorNode) error {
	for _, child := range node.Content {
		if child.Text != "" {
			err := renderText(pdf, child.Text, child.Marks)
			if err != nil {
				return err
			}
		}
	}
	pdf.Br(20) // line break for next paragraph
	return nil
}

// renderHeading applies heading styles based on level
func renderHeading(pdf *gopdf.GoPdf, node ProseMirrorNode) error {
	level := node.Attrs.Level
	font := "bold"
	size := 18 - float64(level*2) // Adjust size based on level

	pdf.SetFont(font, "", size)
	for _, child := range node.Content {
		if child.Text != "" {
			renderText(pdf, child.Text, child.Marks)
		}
	}
	pdf.Br(20)
	return nil
}

// renderList formats bullet or ordered lists
func renderList(pdf *gopdf.GoPdf, node ProseMirrorNode, listType string) error {
	for i, child := range node.Content {
		if child.Type == "list_item" {
			pdf.SetX(20)
			if listType == "ordered" {
				pdf.Cell(nil, fmt.Sprintf("%d. ", i+1))
			} else {
				pdf.Cell(nil, "• ")
			}
			renderParagraph(pdf, child.Content[0]) // render list item text
		}
	}
	return nil
}

// pdfHandler handles incoming JSON data and returns the generated PDF
func pdfHandler(w http.ResponseWriter, r *http.Request) {
	// Read JSON from the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Parse JSON into ProseMirrorNode
	var content []ProseMirrorNode
	err = json.Unmarshal(body, &content)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	// Generate PDF from parsed content
	pdfBytes, err := NewPdfDocument(content)
	if err != nil {
		http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
		return
	}

	// Set response headers and send PDF as response
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=document.pdf")
	w.Write(pdfBytes)
}

func main() {
	http.HandleFunc("/generate-pdf", pdfHandler)
	fmt.Println("Server started on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
