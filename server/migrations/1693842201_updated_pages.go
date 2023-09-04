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
		new_readability_task_id := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "qrvp1oye",
			"name": "readability_task_id",
			"type": "text",
			"required": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), new_readability_task_id)
		collection.Schema.AddField(new_readability_task_id)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("97rubl1eg0xog0y")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("qrvp1oye")

		return dao.SaveCollection(collection)
	})
}
