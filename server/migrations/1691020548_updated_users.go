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

		collection, err := dao.FindCollectionByNameOrId("_pb_users_auth_")
		if err != nil {
			return err
		}

		// add
		new_readability_enabled := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "qxwb4qt3",
			"name": "readability_enabled",
			"type": "bool",
			"required": false,
			"unique": false,
			"options": {}
		}`), new_readability_enabled)
		collection.Schema.AddField(new_readability_enabled)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("_pb_users_auth_")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("qxwb4qt3")

		return dao.SaveCollection(collection)
	})
}
