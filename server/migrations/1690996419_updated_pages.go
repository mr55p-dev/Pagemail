package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("97rubl1eg0xog0y")
		if err != nil {
			return err
		}

		// add
		new_last_crawled := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "4cbcfpay",
			"name": "last_crawled",
			"type": "date",
			"required": false,
			"unique": false,
			"options": {
				"min": "",
				"max": ""
			}
		}`), new_last_crawled)
		collection.Schema.AddField(new_last_crawled)

		// add
		new_title := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "jb2atkwm",
			"name": "title",
			"type": "text",
			"required": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), new_title)
		collection.Schema.AddField(new_title)

		// add
		new_description := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "uj6yvahv",
			"name": "description",
			"type": "text",
			"required": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), new_description)
		collection.Schema.AddField(new_description)

		// add
		new_image_url := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "cnsnntot",
			"name": "image_url",
			"type": "url",
			"required": false,
			"unique": false,
			"options": {
				"exceptDomains": null,
				"onlyDomains": null
			}
		}`), new_image_url)
		collection.Schema.AddField(new_image_url)

		// add
		new_readability_status := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "geanelw2",
			"name": "readability_status",
			"type": "select",
			"required": false,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"values": [
					"UNKNOWN",
					"PROCESSING",
					"FAILED",
					"COMPLETE"
				]
			}
		}`), new_readability_status)
		collection.Schema.AddField(new_readability_status)

		// add
		new_readability_task_data := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "kzkymdv6",
			"name": "readability_task_data",
			"type": "json",
			"required": false,
			"unique": false,
			"options": {}
		}`), new_readability_task_data)
		collection.Schema.AddField(new_readability_task_data)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("97rubl1eg0xog0y")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("4cbcfpay")

		// remove
		collection.Schema.RemoveField("jb2atkwm")

		// remove
		collection.Schema.RemoveField("uj6yvahv")

		// remove
		collection.Schema.RemoveField("cnsnntot")

		// remove
		collection.Schema.RemoveField("geanelw2")

		// remove
		collection.Schema.RemoveField("kzkymdv6")

		return dao.SaveCollection(collection)
	})
}
