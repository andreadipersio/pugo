package cache

import (
	"log"
)

type CacheRequest struct {
	Filename, Content string
}

type Cache struct {
	Data map[string]string

	GetChan chan string
	PutChan chan CacheRequest
}

func NewCache() *Cache {
	return &Cache{
		Data:    make(map[string]string, 10),
		GetChan: make(chan string),
		PutChan: make(chan CacheRequest),
	}
}

func (c *Cache) Run() {
	log.Println("Starting cache")

	for {
		select {

		case filename := <-c.GetChan:
			log.Printf("cache::get %v", filename)
			content := c.Data[filename]

			if content == "" {
				log.Print("cache::miss")
			}

			c.GetChan <- content

		case put := <-c.PutChan:
			log.Printf("cache::put %v", put.Filename)
			c.Data[put.Filename] = put.Content
		}
	}
}
