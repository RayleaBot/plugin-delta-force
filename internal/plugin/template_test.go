package plugin

import (
	"bytes"
	"encoding/json"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderTemplatesParseWithPreviewData(t *testing.T) {
	templates := []string{"delta-loot", "delta-passwords"}
	for _, name := range templates {
		t.Run(name, func(t *testing.T) {
			root := filepath.Join("..", "..", "templates", name)
			markup, err := os.ReadFile(filepath.Join(root, "template.html"))
			if err != nil {
				t.Fatal(err)
			}
			styles, err := os.ReadFile(filepath.Join(root, "styles.css"))
			if err != nil {
				t.Fatal(err)
			}
			preview, err := os.ReadFile(filepath.Join(root, "preview.json"))
			if err != nil {
				t.Fatal(err)
			}
			var data map[string]any
			if err := json.Unmarshal(preview, &data); err != nil {
				t.Fatal(err)
			}
			data["Stylesheet"] = template.CSS(styles)
			parsed, err := template.New(name).Funcs(template.FuncMap{
				"safeHTML": func(value string) template.HTML { return template.HTML(value) },
			}).Parse(string(markup))
			if err != nil {
				t.Fatal(err)
			}
			var output bytes.Buffer
			if err := parsed.Execute(&output, data); err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(output.String(), "seed 2e2cc0c4") || !strings.Contains(output.String(), "三角洲助手") {
				t.Fatalf("rendered template lost direction contract or product name")
			}
			if !strings.Contains(output.String(), "Created By RayleaBot") {
				t.Fatal("rendered template lost the host-injected render footer")
			}
		})
	}
}

func TestRenderTemplateSchemasAcceptHostFooter(t *testing.T) {
	templates := []string{"delta-loot", "delta-passwords"}
	for _, name := range templates {
		t.Run(name, func(t *testing.T) {
			root := filepath.Join("..", "..", "templates", name)
			schemaBytes, err := os.ReadFile(filepath.Join(root, "input.schema.json"))
			if err != nil {
				t.Fatal(err)
			}
			var schema struct {
				Properties map[string]struct {
					Type string `json:"type"`
				} `json:"properties"`
			}
			if err := json.Unmarshal(schemaBytes, &schema); err != nil {
				t.Fatal(err)
			}
			footer, exists := schema.Properties["render_footer"]
			if !exists || footer.Type != "string" {
				t.Fatal("input schema must accept the host-injected render_footer string")
			}
		})
	}
}
