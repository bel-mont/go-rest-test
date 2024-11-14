package html

import "html/template"

// BaseLayoutTemplate parses the provided template files along with the base header and footer templates.
func BaseLayoutTemplate(templateFiles ...string) (*template.Template, error) {
	baseHeader := "web/views/layouts/base-header.gohtml"
	baseFooter := "web/views/layouts/base-footer.gohtml"
	allFiles := append([]string{baseHeader, baseFooter}, templateFiles...)

	return template.ParseFiles(allFiles...)
}

func LoggedLayoutTemplate(templateFiles ...string) (*template.Template, error) {
	loggedHeader := "web/views/layouts/logged-header.gohtml"
	loggedFooter := "web/views/layouts/logged-footer.gohtml"
	allFiles := append([]string{loggedHeader, loggedFooter}, templateFiles...)

	return template.ParseFiles(allFiles...)
}
