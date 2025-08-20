package internal

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func GetCleanedHTML(email *Email) (string, error) {
	doc, err := html.Parse(strings.NewReader(email.HTML))
	if err != nil {
		return "", err
	}

	walkNode(doc, email)

	if head := findNode(doc, "head"); head != nil {
		injectEncoding(head)
		injectTitle(head, email)
		injectJavaScript(head)
	}

	var buf bytes.Buffer
	err = html.Render(&buf, doc)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func walkNode(n *html.Node, email *Email) {
	if n.Type == html.ElementNode {
		switch strings.ToLower(n.Data) {
		case "img":
			externalImageToPlaceholder(n)
			embedCidSource(n, email)

		case "a":
			safeguardLink(n)

		case "link":
			if hasExternalReference(n, "href") && n.Parent != nil {
				n.Parent.RemoveChild(n)
			}

		case "script", "iframe", "object", "embed", "applet":
			if n.Parent != nil {
				n.Parent.RemoveChild(n)
				return
			}

		default:
			cleanAttributes(n)
		}
	}

	for c := n.FirstChild; c != nil; {
		next := c.NextSibling
		walkNode(c, email)
		c = next
	}
}

func findNode(n *html.Node, name string) *html.Node {
	if n.Type == html.ElementNode && strings.EqualFold(n.Data, name) {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if element := findNode(c, name); element != nil {
			return element
		}
	}

	return nil
}

func injectEncoding(n *html.Node) {
	metaNode := &html.Node{
		Type: html.ElementNode,
		Data: "meta",
		Attr: []html.Attribute{
			{Key: "charset", Val: "UTF-8"},
		},
	}
	n.AppendChild(metaNode)
}

func injectTitle(n *html.Node, email *Email) {
	titleNode := &html.Node{
		Type: html.ElementNode,
		Data: "title",
		FirstChild: &html.Node{
			Type: html.TextNode,
			Data: fmt.Sprintf("%s: %s", email.From, email.Subject),
		},
	}
	n.AppendChild(titleNode)
}

//go:embed runtime.js
var javaScript string

func injectJavaScript(n *html.Node) {
	scriptNode := &html.Node{
		Type: html.ElementNode,
		Data: "script",
		FirstChild: &html.Node{
			Type: html.TextNode,
			Data: javaScript,
		},
	}
	n.AppendChild(scriptNode)
}

// Convert <img> to <div> with placeholder styling
func externalImageToPlaceholder(n *html.Node) {
	if !hasExternalReference(n, "src") {
		return
	}
	n.Data = "div"

	// extract / keep description, size and style
	var newAttrs []html.Attribute
	var style, class string
	placeholder := "Image"
	for _, attr := range n.Attr {
		switch attr.Key {
		case "alt":
			placeholder = attr.Val

		case "width", "height":
			if attr.Val != "" {
				style = fmt.Sprintf("%s %s: %spx;", style, attr.Key, attr.Val)
			}

		case "style":
			style = fmt.Sprintf("%s %s", style, attr.Val)

		case "class":
			class = fmt.Sprintf("%s %s", class, attr.Val)

		case "id":
			newAttrs = append(newAttrs, attr)

		case "src":
			loadImage := fmt.Sprintf("loadImage(this, %q)", attr.Val)
			newAttrs = append(newAttrs, html.Attribute{Key: "onclick", Val: loadImage})
		}
	}

	newAttrs = append(newAttrs, html.Attribute{Key: "style", Val: style})
	newAttrs = append(newAttrs, html.Attribute{Key: "class", Val: class})
	n.Attr = newAttrs

	textNode := &html.Node{
		Type: html.TextNode,
		Data: fmt.Sprintf("[%s]", placeholder),
	}
	n.AppendChild(textNode)
}

func embedCidSource(n *html.Node, email *Email) {
	i, attr := getAttr(n, "src")
	if attr == nil || !strings.HasPrefix(attr.Val, "cid:") {
		return
	}

	attachment := AttachmentByCID(email, attr.Val)
	if attachment != nil {
		encoded := base64.StdEncoding.EncodeToString(attachment.Content)
		data := fmt.Sprintf("data:%s;charset=utf-8;base64,%s", attachment.ContentType, encoded)
		n.Attr[i].Val = data
	}
}

func hasExternalReference(n *html.Node, referencingAttr string) bool {
	if _, attr := getAttr(n, referencingAttr); attr != nil {
		return isExternalHref(strings.ToLower(attr.Val))
	}
	return false
}

func safeguardLink(n *html.Node) {
	if i, attr := getAttr(n, "href"); attr != nil {
		if isExternalHref(strings.ToLower(attr.Val)) {
			loadURL := fmt.Sprintf("loadURL(this, '%s')", attr.Val)
			n.Attr = append(n.Attr, html.Attribute{Key: "onclick", Val: loadURL})
			n.Attr[i].Val = fmt.Sprintf("#%s", attr.Val)
		}
	}
}

func getAttr(n *html.Node, name string) (int, *html.Attribute) {
	for i, attr := range n.Attr {
		if attr.Key == name {
			return i, &n.Attr[i]
		}
	}
	return -1, nil
}

func isExternalHref(href string) bool {
	return strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") || strings.HasPrefix(href, "//")
}

func cleanAttributes(n *html.Node) {
	var cleanAttrs []html.Attribute

	for _, attr := range n.Attr {
		switch strings.ToLower(attr.Key) {
		case "style":
			if css := cleanStyleValue(attr.Val); css != "" {
				cleanAttrs = append(cleanAttrs, html.Attribute{Key: attr.Key, Val: css})
			}

		case "class", "id", "colspan", "rowspan", "cellpadding", "cellspacing", "border", "width", "height", "align", "valign":
			cleanAttrs = append(cleanAttrs, attr)

		case "onclick", "onload", "onerror", "onmouseover":
			// Remove event handlers
			continue

		default:
			if !isDangerousAttribute(attr.Key) {
				cleanAttrs = append(cleanAttrs, attr)
			}
		}
	}

	n.Attr = cleanAttrs
}

// regexes to remove external ressource loading
var urlRegex *regexp.Regexp = regexp.MustCompile(`url\s*\(\s*["']?[^"')]*["']?\s*\)`)
var importRegex *regexp.Regexp = regexp.MustCompile(`@import[^;]*;?`)
var jsRegex *regexp.Regexp = regexp.MustCompile(`javascript\s*:[^;]*;?`)

func cleanStyleValue(css string) string {
	css = urlRegex.ReplaceAllLiteralString(css, "")
	css = importRegex.ReplaceAllLiteralString(css, "")
	css = jsRegex.ReplaceAllLiteralString(css, "")
	return css
}

func isDangerousAttribute(attr string) bool {
	attr = strings.ToLower(attr)
	switch attr {
	case "src", "href", "action", "formaction", "background":
		return true
	}

	// event handlers are named on*
	return strings.HasPrefix(attr, "on")
}
