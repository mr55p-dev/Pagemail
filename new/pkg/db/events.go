package db

type EventOutput[T User | Page] chan Event[T]
type EventType string
type Event[T User | Page] struct {
	Event  EventType
	Record *T
}

func (c *Client) AddPageListener(id string, output EventOutput[Page]) {
	c.log.Info().Msg("Added listener")
	PageEventMap[id] = output
}

func (c *Client) RemovePageListener(id string) {
	c.log.Info().Msg("Removed listener")
	delete(PageEventMap, id)
}
