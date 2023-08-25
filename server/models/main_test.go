package models

import "testing"

func TestPageToMap(t *testing.T) {
	p := Page{
		Id:          "123456",
		Title:       "Hello, world!",
		Description: "Lorem ipsum.",
	}
	m := p.ToMap()
	if m["id"] != p.Id {
		t.Errorf("Map id %s does not match struct field %s", m["id"], p.Id)
	}
	if m["title"] != p.Title {
		t.Errorf("Map title %s does not match struct field %s", m["title"], p.Title)
	}
	if m["description"] != p.Description {
		t.Errorf("Map description %s does not match struct field %s", m["description"], p.Description)
	}
	if m["image_url"] != nil {
		t.Errorf("Map image_url != nil (%s)", m["image_url"])
	}
	
}
