package services

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/danbernal999/kiradoc-backend/pkg/models"
	"github.com/danbernal999/kiradoc-backend/pkg/repository"
)

type ExportService struct {
	repo *repository.Repository
}

func NewExportService(repo *repository.Repository) *ExportService {
	return &ExportService{repo: repo}
}

type ExportOptions struct {
	IncludeWatermark bool
	IncludeTOC       bool
	CompanyLogo      string
	CompanyName      string
	CompanyColors    string
	Draft            bool
}

type ExportResult struct {
	FileName    string
	ContentType string
	Data        []byte
}

func (es *ExportService) ExportPDF(doc *repository.Document, opts ExportOptions) (*ExportResult, error) {
	pdfBuffer := &bytes.Buffer{}

	content := es.buildDocumentContent(doc, opts)

	if err := es.generatePDF(pdfBuffer, doc, content, opts); err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return &ExportResult{
		FileName:    es.sanitizeFileName(doc.Title) + ".pdf",
		ContentType: "application/pdf",
		Data:        pdfBuffer.Bytes(),
	}, nil
}

func (es *ExportService) ExportDOCX(doc *repository.Document, opts ExportOptions) (*ExportResult, error) {
	docxBuffer := &bytes.Buffer{}

	content := es.buildDocumentContent(doc, opts)

	if err := es.generateDOCX(docxBuffer, doc, content, opts); err != nil {
		return nil, fmt.Errorf("failed to generate DOCX: %w", err)
	}

	return &ExportResult{
		FileName:    es.sanitizeFileName(doc.Title) + ".docx",
		ContentType: "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		Data:        docxBuffer.Bytes(),
	}, nil
}

func (es *ExportService) ExportHTML(doc *repository.Document, opts ExportOptions) (*ExportResult, error) {
	htmlBuffer := &bytes.Buffer{}

	content := es.buildDocumentContent(doc, opts)

	html := es.generateHTML(doc, content, opts)

	htmlBuffer.WriteString(html)

	return &ExportResult{
		FileName:    es.sanitizeFileName(doc.Title) + ".html",
		ContentType: "text/html",
		Data:        htmlBuffer.Bytes(),
	}, nil
}

func (es *ExportService) buildDocumentContent(doc *repository.Document, opts ExportOptions) string {
	var content strings.Builder

	if opts.IncludeTOC {
		content.WriteString("TABLE OF CONTENTS\n\n")
		content.WriteString("1. Document Overview\n")
		content.WriteString("2. Main Content\n")
		content.WriteString("3. Signatures\n\n")
	}

	content.WriteString("DOCUMENT: ")
	content.WriteString(doc.Title)
	content.WriteString("\n")
	content.WriteString("Type: ")
	content.WriteString(doc.Type)
	content.WriteString("\n")
	content.WriteString("Created: ")
	createdAtStr := es.formatTime(doc.CreatedAt)
	content.WriteString(createdAtStr)
	content.WriteString("\n\n")

	if opts.Draft {
		content.WriteString("[DRAFT - NOT FOR EXECUTION]\n\n")
	}

	content.WriteString("CONTENT:\n")
	content.WriteString(doc.Content)
	content.WriteString("\n\n")

	if opts.IncludeWatermark && opts.Draft {
		content.WriteString("\n[WATERMARK: DRAFT]\n")
	}

	return content.String()
}

func (es *ExportService) formatTime(t interface{}) string {
	if t == nil {
		return "N/A"
	}
	switch v := t.(type) {
	case time.Time:
		return v.Format("2006-01-02 15:04:05")
	case string:
		return v
	default:
		return "N/A"
	}
}

func (es *ExportService) generatePDF(w io.Writer, doc *repository.Document, content string, opts ExportOptions) error {
	htmlContent := es.generateHTML(doc, content, opts)

	pdf := bytes.NewBufferString(htmlContent)

	data := pdf.Bytes()
	_, err := w.Write(data)
	return err
}

func (es *ExportService) generateDOCX(w io.Writer, doc *repository.Document, content string, opts ExportOptions) error {
	docxContent := es.buildDOCXContent(doc, content, opts)

	_, err := w.Write([]byte(docxContent))
	return err
}

func (es *ExportService) buildDOCXContent(doc *repository.Document, content string, opts ExportOptions) string {
	var result strings.Builder

	result.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	result.WriteString("<document>\n")
	result.WriteString(fmt.Sprintf("  <title>%s</title>\n", es.escapeXML(doc.Title)))
	result.WriteString(fmt.Sprintf("  <type>%s</type>\n", es.escapeXML(doc.Type)))
	result.WriteString(fmt.Sprintf("  <created>%s</created>\n", es.formatTime(doc.CreatedAt)))

	if opts.Draft {
		result.WriteString("  <draft>true</draft>\n")
	}

	result.WriteString(fmt.Sprintf("  <content>%s</content>\n", es.escapeXML(content)))
	result.WriteString("</document>\n")

	return result.String()
}

func (es *ExportService) generateHTML(doc *repository.Document, content string, opts ExportOptions) string {
	var html strings.Builder

	html.WriteString("<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n")
	html.WriteString("<meta charset=\"UTF-8\">\n")
	html.WriteString("<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n")
	html.WriteString(fmt.Sprintf("<title>%s</title>\n", doc.Title))

	html.WriteString("<style>\n")
	html.WriteString("body { font-family: Arial, sans-serif; margin: 20px; line-height: 1.6; }\n")
	html.WriteString("h1 { color: #333; border-bottom: 2px solid #007bff; padding-bottom: 10px; }\n")
	html.WriteString("h2 { color: #555; margin-top: 20px; }\n")
	html.WriteString(".draft-header { background-color: #fff3cd; border: 1px solid #ffc107; padding: 10px; margin: 10px 0; color: #856404; }\n")
	html.WriteString(".watermark { opacity: 0.1; font-size: 100px; position: fixed; transform: rotate(-45deg); pointer-events: none; }\n")
	html.WriteString(".metadata { background-color: #f8f9fa; padding: 10px; border-radius: 5px; margin: 10px 0; font-size: 12px; }\n")
	html.WriteString("table { width: 100%; border-collapse: collapse; margin: 10px 0; }\n")
	html.WriteString("th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }\n")
	html.WriteString("th { background-color: #f2f2f2; }\n")
	html.WriteString("</style>\n")
	html.WriteString("</head>\n<body>\n")

	if opts.IncludeWatermark && opts.Draft {
		html.WriteString("<div class=\"watermark\">DRAFT</div>\n")
	}

	if opts.Draft {
		html.WriteString("<div class=\"draft-header\">⚠️ DRAFT DOCUMENT - NOT FOR EXECUTION</div>\n")
	}

	html.WriteString(fmt.Sprintf("<h1>%s</h1>\n", doc.Title))
	html.WriteString("<div class=\"metadata\">\n")
	html.WriteString(fmt.Sprintf("<strong>Type:</strong> %s<br>\n", doc.Type))
	html.WriteString(fmt.Sprintf("<strong>Created:</strong> %s<br>\n", es.formatTime(doc.CreatedAt)))
	html.WriteString(fmt.Sprintf("<strong>Status:</strong> %s\n", doc.Status))
	html.WriteString("</div>\n")

	if opts.IncludeTOC {
		html.WriteString("<h2>Table of Contents</h2>\n")
		html.WriteString("<ol>\n")
		html.WriteString("<li><a href=\"#overview\">Document Overview</a></li>\n")
		html.WriteString("<li><a href=\"#content\">Main Content</a></li>\n")
		html.WriteString("</ol>\n")
	}

	html.WriteString("<h2 id=\"content\">Content</h2>\n")
	html.WriteString(fmt.Sprintf("<pre>%s</pre>\n", content))

	html.WriteString("<hr>\n")
	html.WriteString(fmt.Sprintf("<p><small>Generated on %s</small></p>\n", time.Now().Format("2006-01-02 15:04:05")))

	html.WriteString("</body>\n</html>")

	return html.String()
}

func (es *ExportService) sanitizeFileName(title string) string {
	title = strings.TrimSpace(title)
	title = strings.ReplaceAll(title, "/", "-")
	title = strings.ReplaceAll(title, "\\", "-")
	title = strings.ReplaceAll(title, ":", "-")
	title = strings.ReplaceAll(title, "*", "-")
	title = strings.ReplaceAll(title, "?", "-")
	title = strings.ReplaceAll(title, "\"", "-")
	title = strings.ReplaceAll(title, "|", "-")
	title = strings.ReplaceAll(title, "<", "-")
	title = strings.ReplaceAll(title, ">", "-")
	if len(title) > 200 {
		title = title[:200]
	}
	return title
}

func (es *ExportService) escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}

func (es *ExportService) ExportBatch(documents []*repository.Document, format string, opts ExportOptions) (*ExportResult, error) {
	if len(documents) == 0 {
		return nil, fmt.Errorf("no documents to export")
	}

	allData := make(map[string][]byte)

	for _, doc := range documents {
		var result *ExportResult
		var err error

		switch format {
		case "pdf":
			result, err = es.ExportPDF(doc, opts)
		case "docx":
			result, err = es.ExportDOCX(doc, opts)
		case "html":
			result, err = es.ExportHTML(doc, opts)
		default:
			result, err = es.ExportHTML(doc, opts)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to export document %s: %w", doc.ID, err)
		}

		allData[result.FileName] = result.Data
	}

	zipData, err := es.createZIP(allData)
	if err != nil {
		return nil, fmt.Errorf("failed to create ZIP: %w", err)
	}

	return &ExportResult{
		FileName:    fmt.Sprintf("documents-%d.zip", time.Now().Unix()),
		ContentType: "application/zip",
		Data:        zipData,
	}, nil
}

func (es *ExportService) createZIP(files map[string][]byte) ([]byte, error) {
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	defer w.Close()

	for filename, data := range files {
		f, err := w.Create(filename)
		if err != nil {
			return nil, err
		}
		if _, err := f.Write(data); err != nil {
			return nil, err
		}
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (es *ExportService) BuildExportOptions(req *models.ExportRequest) ExportOptions {
	return ExportOptions{
		IncludeWatermark: req.IncludeWatermark,
		IncludeTOC:       req.IncludeTOC,
		CompanyLogo:      req.CompanyLogo,
		CompanyName:      req.CompanyName,
		CompanyColors:    req.CompanyColors,
		Draft:            req.Draft,
	}
}
