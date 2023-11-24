package main

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

type OffsetMapping struct {
	PlainTextOffset int
	HTMLTextOffset  int
}

func main() {
	// Replace with your HTML text
	htmlText := `<p><strong>bold text</strong></p>
<p><i>italic text</i></p>
<p><s>strikethrough</s></p>
<ol>
<li>item one</li></ol>
<ul>
<li>item one</li></ul>
<p><a href="https://www.microsoft.com/" title="https://www.microsoft.com/">Here's</a> a link</p>
<h1>Header text</h1>
<blockquote>
<p>this is a quote</p>
</blockquote>`

	// Create an HTML tokenizer
	tokenizer := html.NewTokenizer(strings.NewReader(htmlText))

	var offsetMappings []OffsetMapping
	var plainText strings.Builder
	var offset, htmlOffset int

	for {
		tokenType := tokenizer.Next()

		switch tokenType {
		case html.ErrorToken:
			// End of the HTML document
			fmt.Println("HTML Text:")
			fmt.Println(htmlText)
			fmt.Println("\nPlain Text:")
			fmt.Println(plainText.String())

			// Print all offset mappings
			for i, mapping := range offsetMappings {
				fmt.Printf("Entry %d: PlainTextOffset=%d %c, HTMLTextOffset=%d %c\n", i, mapping.PlainTextOffset, plainText.String()[mapping.PlainTextOffset], mapping.HTMLTextOffset, htmlText[mapping.HTMLTextOffset])
			}
			return

		case html.TextToken:
			token := tokenizer.Token()
			text := token.Data
			plainText.WriteString(text)
			offsetMappings = append(offsetMappings, OffsetMapping{
				PlainTextOffset: offset,
				HTMLTextOffset:  htmlOffset,
			})
			offset += len(text)
			htmlOffset += len(text)

		case html.StartTagToken, html.EndTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			tagName := token.Data
			// Add the length of the tag, its angle brackets, and attributes to the HTML offset
			htmlOffset += len(tagName) + 2 // +2 accounts for "<" and ">"
			for _, attr := range token.Attr {
				// Add the length of the attribute name, attribute value, and '=' character to the HTML offset
				htmlOffset += len(attr.Key) + len(attr.Val) + 4 // +3 accounts for attribute name, '=', and attribute value
			}
			if tokenType == html.EndTagToken || tokenType == html.SelfClosingTagToken {
				htmlOffset++ // Add 1 for the slash character in end tags
			}
		}
	}
}
