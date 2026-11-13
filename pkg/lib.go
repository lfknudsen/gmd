package pkg

import (
	html2 "html"
	"regexp"
	"strconv"
	"strings"

	"gitlab.com/golang-commonmark/markdown"
)

func removeNewlines(s string) string {
	s = strings.ReplaceAll(s, "\n\n", "\\")
	s = strings.ReplaceAll(s, "\r\n\r\n", "\\")
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\\", "\n")
	return s
}

func RenderMarkdown(buffer []byte) string {
	md := markdown.New()
	html := md.RenderToString(buffer)

	rx := regexp.MustCompile("(?s)<p>.*?</p>")
	html = rx.ReplaceAllStringFunc(html, removeNewlines)

	html = strings.Replace(html, "<p>", "", -1)
	html = strings.Replace(html, "</p>", "", -1)
	html = strings.Replace(html, "<br/>", "\n", -1)
	html = strings.Replace(html, "<br>", "\n", -1)
	html = strings.Replace(html, "<i>", "\033[003m", -1)
	html = strings.Replace(html, "</i>", "", -1)

	// Code block rendering
	parts := strings.Split(html, "\n")
	codeBlockStartAt := -1
	codeBlock := false
	maxLength := 0
	prefix := "<pre><code>"
	postfix := "</code>"
	for i, line := range parts {
		if strings.HasPrefix(line, prefix) {
			if strings.Contains(line, postfix) {
				continue
			}
			codeBlock = true
			line = line[len(prefix):]
			maxLength = len(line)
			codeBlockStartAt = i
			continue
		}
		if codeBlock {
			stopAt := strings.Index(line, postfix)
			length := 0
			if stopAt == -1 {
				length = len(line)
			} else {
				length = len(line[:stopAt])
			}
			if length > maxLength {
				maxLength = length
			}
			if stopAt != -1 {
				for j := range parts[codeBlockStartAt : i+1] {
					length = len(parts[codeBlockStartAt+j])
					repetitions := maxLength - length
					if repetitions > 0 {
						padding := strings.Repeat(" ", repetitions)
						parts[codeBlockStartAt+j] = parts[codeBlockStartAt+j] + padding
					}
				}
				codeBlock = false
				codeBlockStartAt = -1
				maxLength = 0
			}
		}
	}

	html = strings.Join(parts, "\n")

	html = strings.Replace(html, "<pre>", "\n", -1)
	html = strings.Replace(html, "</pre>", "", -1)
	html = strings.Replace(html, "<code>", "\033[007m", -1)
	html = strings.Replace(html, "</code>", "\033[0m", -1)

	// Indent
	parts = strings.Split(html, "\n")
	colour := "0"
	indent := ""
	for i, line := range parts {
		colourIndex := strings.LastIndex(line, "\033[")
		if colourIndex != -1 {
			colourIndex += len("\033[")
			colourEndsAt := strings.IndexByte(line[colourIndex:], 'm')
			if colourEndsAt != -1 {
				colour = line[colourIndex : colourIndex+colourEndsAt]
			}
		}
		if !strings.HasPrefix(line, "<h") {
			parts[i] = "\033[0m" + indent + "\033[" + colour + "m" + parts[i]
		} else if strings.HasPrefix(line, "<h") {
			d, err := strconv.Atoi(line[2:3])
			if err != nil {
				continue
			}
			indent = strings.Repeat("  ", d)
		}
	}

	html = strings.Join(parts, "\n")
	html = strings.Replace(html, "<h1>", "\033[94m", -1)
	html = strings.Replace(html, "</h1>", "\033[0m\n", -1)
	html = strings.Replace(html, "<h2>", "  \033[95m", -1)
	html = strings.Replace(html, "</h2>", "\033[0m", -1)
	html = strings.Replace(html, "<h3>", "    \033[96m", -1)
	html = strings.Replace(html, "</h3>", "\033[0m", -1)
	html = strings.Replace(html, "<h4>", "      \033[001m", -1)
	html = strings.Replace(html, "</h4>", "\033[0m", -1)
	html = strings.Replace(html, "<h5>", "        \033[002m", -1)
	html = strings.Replace(html, "</h5>", "\033[0m", -1)

	html = html2.UnescapeString(html)

	return html
}
