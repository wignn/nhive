package events

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/novelhive/novel-service/internal/domain"
)

type NATSPublisher struct {
	js nats.JetStreamContext
}

func NewNATSPublisher(nc *nats.Conn) (*NATSPublisher, error) {
	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}
	// Create stream if not exists
	js.AddStream(&nats.StreamConfig{
		Name:     "NOVELS",
		Subjects: []string{"novel.>", "chapter.>"},
	})
	return &NATSPublisher{js: js}, nil
}

type NovelEvent struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Slug          string   `json:"slug"`
	Synopsis      string   `json:"synopsis"`
	Author        string   `json:"author"`
	CoverURL      string   `json:"cover_url"`
	Genres        []string `json:"genres"`
	Status        string   `json:"status"`
	TotalChapters int      `json:"total_chapters"`
}

type ChapterEvent struct {
	NovelID    string `json:"novel_id"`
	NovelTitle string `json:"novel_title"`
	NovelSlug  string `json:"novel_slug"`
	ChapterID  string `json:"chapter_id"`
	Number     int    `json:"number"`
	Title      string `json:"title"`
}

func (p *NATSPublisher) PublishNovelCreated(novel *domain.Novel) {
	evt := NovelEvent{
		ID: novel.ID, Title: novel.Title, Slug: novel.Slug,
		Synopsis: novel.Synopsis, Author: novel.Author, CoverURL: novel.CoverURL,
		Status: novel.Status, TotalChapters: novel.TotalChapters,
	}
	for _, g := range novel.Genres {
		evt.Genres = append(evt.Genres, g.Name)
	}
	data, _ := json.Marshal(evt)
	if _, err := p.js.Publish("novel.created", data); err != nil {
		log.Printf("Failed to publish novel.created: %v", err)
	}
}

func (p *NATSPublisher) PublishNovelUpdated(novel *domain.Novel) {
	evt := NovelEvent{
		ID: novel.ID, Title: novel.Title, Slug: novel.Slug,
		Synopsis: novel.Synopsis, Author: novel.Author, CoverURL: novel.CoverURL,
		Status: novel.Status, TotalChapters: novel.TotalChapters,
	}
	for _, g := range novel.Genres {
		evt.Genres = append(evt.Genres, g.Name)
	}
	data, _ := json.Marshal(evt)
	if _, err := p.js.Publish("novel.updated", data); err != nil {
		log.Printf("Failed to publish novel.updated: %v", err)
	}
}

func (p *NATSPublisher) PublishChapterPublished(novel *domain.Novel, ch *domain.Chapter) {
	evt := ChapterEvent{
		NovelID: novel.ID, NovelTitle: novel.Title, NovelSlug: novel.Slug,
		ChapterID: ch.ID, Number: ch.Number, Title: ch.Title,
	}
	data, _ := json.Marshal(evt)
	if _, err := p.js.Publish("novel.chapter.published", data); err != nil {
		log.Printf("Failed to publish novel.chapter.published: %v", err)
	}
}
