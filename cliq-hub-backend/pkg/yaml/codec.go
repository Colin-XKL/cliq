package yamlcodec

import (
    "encoding/base64"
    "strings"

    "gopkg.in/yaml.v3"

    tpl "cliq-hub-backend/internal/template"
)

func StripFences(s string) string {
    out := strings.TrimSpace(s)
    // remove Markdown fenced code blocks if present
    if strings.HasPrefix(out, "```") {
        // find first newline after fence
        // naive approach: remove triple backticks anywhere
        out = strings.ReplaceAll(out, "```yaml", "")
        out = strings.ReplaceAll(out, "```", "")
        out = strings.TrimSpace(out)
    }
    return out
}

func UnmarshalTemplate(s string) (*tpl.Template, error) {
    var t tpl.Template
    if err := yaml.Unmarshal([]byte(s), &t); err != nil {
        return nil, err
    }
    return &t, nil
}

func MarshalTemplate(t *tpl.Template) (string, error) {
    b, err := yaml.Marshal(t)
    if err != nil {
        return "", err
    }
    return string(b), nil
}

func Base64Encode(s string) string {
    return base64.StdEncoding.EncodeToString([]byte(s))
}