package db

type EventOutput[T User | Page] chan Event[T]
type EventType string
type Event[T User | Page] struct {
	Event  EventType
	Record *T
}

func (c *DBDriver) AddPageListener(id string, output EventOutput[Page]) {
	c.log.Info("Added listener")
	PageEventMap[id] = output
}

func (c *DBDriver) RemovePageListener(id string) {
	c.log.Info("Removed listener")
	delete(PageEventMap, id)
}
