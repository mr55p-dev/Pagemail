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
		new_is_readable := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "1tatfjqk",
			"name": "is_readable",
			"type": "bool",
			"required": false,
			"unique": false,
			"options": {}
		}`), new_is_readable)
		collection.Schema.AddField(new_is_readable)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("97rubl1eg0xog0y")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("1tatfjqk")

		return dao.SaveCollection(collection)
	})
}
