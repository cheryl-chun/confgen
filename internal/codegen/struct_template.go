package codegen

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/cheryl-chun/confgen/internal/analyzer"
)

const structTemplate = `
type {{.Name}} struct {
{{- range .Fields}}
	{{pad .Name $.MaxNameLen}} {{pad .Type $.MaxTypeLen}} ` + "`" + `{{.Tags}}` + "`" + `{{if and $.AddComments .Comment}} // {{.Comment}}{{end}}
{{- end}}
{{- if .IncludeTree}}

	// ConfigTree provides the underlying representation for dynamic queries.
	{{pad "ConfigTree" $.MaxNameLen}} *runtime.Tree
{{- end}}
}
`

// writeStructWithTree renders a Go struct definition using text/template.
// It maintains the aesthetic columnar alignment through a custom FuncMap.
func (g *Generator) writeStructWithTree(def *analyzer.StructDef, includeTree bool) {
	if def == nil {
		return
	}

	// 1. Pre-calculate the maximum dimensions for alignment heuristics.
	maxNameLen, maxTypeLen := g.calculateAlignment(def, includeTree)

	// 2. Define helper functions for the template engine.
	funcMap := template.FuncMap{
		"pad": func(s string, maxLen int) string {
			padding := maxLen - len(s) + 1
			return s + strings.Repeat(" ", padding)
		},
	}

	// 3. Prepare the data context for the template.
	type fieldData struct {
		Name    string
		Type    string
		Tags    string
		Comment string
	}

	fields := make([]fieldData, 0, len(def.Fields))
	for _, f := range def.Fields {
		fields = append(fields, fieldData{
			Name:    f.Name,
			Type:    f.Type,
			Tags:    g.formatTags(f),
			Comment: f.Comment,
		})
	}

	data := struct {
		Name        string
		Fields      []fieldData
		MaxNameLen  int
		MaxTypeLen  int
		IncludeTree bool
		AddComments bool
	}{
		Name:        def.Name,
		Fields:      fields,
		MaxNameLen:  maxNameLen,
		MaxTypeLen:  maxTypeLen,
		IncludeTree: includeTree,
		AddComments: g.opts.AddComments,
	}

	// 4. Execute the template and stream the output to the generator's buffer.
	tmpl, err := template.New("struct").Funcs(funcMap).Parse(structTemplate)
	if err != nil {
		// In a generator, template parsing errors are typically handled as panics
		// if the template is a hardcoded constant.
		panic(fmt.Errorf("failed to parse struct template: %w", err))
	}

	if err := tmpl.Execute(g.buf, data); err != nil {
		panic(fmt.Errorf("failed to execute struct template: %w", err))
	}
}

// calculateAlignment computes the optimal column widths for the struct fields.
func (g *Generator) calculateAlignment(def *analyzer.StructDef, includeTree bool) (int, int) {
	maxNameLen := 0
	maxTypeLen := 0
	for _, field := range def.Fields {
		if len(field.Name) > maxNameLen {
			maxNameLen = len(field.Name)
		}
		if len(field.Type) > maxTypeLen {
			maxTypeLen = len(field.Type)
		}
	}

	if includeTree {
		if len("ConfigTree") > maxNameLen {
			maxNameLen = len("ConfigTree")
		}
		if len("*runtime.Tree") > maxTypeLen {
			maxTypeLen = len("*runtime.Tree")
		}
	}
	return maxNameLen, maxTypeLen
}

// formatTags is a helper to centralize struct tag serialization.
func (g *Generator) formatTags(field *analyzer.FieldDef) string {
	var tags []string
	if field.JSONTag != "" {
		tags = append(tags, fmt.Sprintf(`json:"%s"`, field.JSONTag))
	}
	if field.YAMLTag != "" {
		tags = append(tags, fmt.Sprintf(`yaml:"%s"`, field.YAMLTag))
	}
	if field.MapStructTag != "" {
		tags = append(tags, fmt.Sprintf(`mapstructure:"%s"`, field.MapStructTag))
	}
	return strings.Join(tags, " ")
}
