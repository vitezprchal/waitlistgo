package models

import "html/template"

type Blog struct {
	Slug    string        `json:"slug"`
	Content template.HTML `json:"content"`
	SEO     *SEO          `json:"seo"`
}
