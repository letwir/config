package epub

import (
	"archive/zip"
	"bytes"
	"html"
	"io"
	"strings"
)

type ParsedFile struct {
	Name string
	Text string
}

// StripTags removes HTML tags from string
func StripTags(htmlStr string) string {
	var builder strings.Builder
	inTag := false
	for _, r := range htmlStr {
		if r == '<' {
			inTag = true
			continue
		}
		if r == '>' {
			inTag = false
			continue
		}
		if !inTag {
			builder.WriteRune(r)
		}
	}
	return html.UnescapeString(builder.String())
}

// ParseEPUBBytes extracts plain text from EPUB ZIP bytes
func ParseEPUBBytes(epubBytes []byte) ([]ParsedFile, error) {
	reader, err := zip.NewReader(bytes.NewReader(epubBytes), int64(len(epubBytes)))
	if err != nil {
		return nil, err
	}

	var parsedFiles []ParsedFile

	for _, file := range reader.File {
		filename := strings.ToLower(file.Name)
		if strings.HasSuffix(filename, ".html") || strings.HasSuffix(filename, ".xhtml") || strings.HasSuffix(filename, ".htm") {
			rc, err := file.Open()
			if err != nil {
				continue
			}
			
			var buf bytes.Buffer
			_, err = io.Copy(&buf, rc)
			rc.Close()
			if err != nil {
				continue
			}

			content := buf.String()
			plainText := StripTags(content)

			parsedFiles = append(parsedFiles, ParsedFile{
				Name: file.Name,
				Text: plainText,
			})
		}
	}

	return parsedFiles, nil
}
