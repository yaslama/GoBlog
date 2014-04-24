package utils

import (
	"github.com/russross/blackfriday"
	"html/template"
	"regexp"
	"strings"
)

func Html2str(html string) string {
	src := string(html)

	//WillHTMLTags to lowercase
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)

	//RemovalSTYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")

	//RemovalSCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")

	//Remove all angle bracketsHTMLCode，And replace with line breaks
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")

	//Remove consecutive line breaks
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")

	return strings.TrimSpace(src)
}

func Markdown2Html(text string) string {
	return string(blackfriday.MarkdownCommon([]byte(text)))
}

func Markdown2HtmlTemplate(text string) template.HTML {
	return template.HTML(Markdown2Html(text))
}
