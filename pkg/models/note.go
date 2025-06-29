package models

import (
	"encoding/json"

	"github.com/sirowain/notecli/pkg/utils"
)

type Note struct {
	Id        utils.NoteId `json:"id"`
	Headline  string       `json:"headline"`
	Content   string       `json:"content"`
	Tags      []string     `json:"tags"`
	CreatedAt string       `json:"created_at"`
	UpdatedAt string       `json:"updated_at"`
}

func (n *Note) FromJson(data []byte) error {
	return json.Unmarshal(data, &n)
}

// MarshalJSON serializes the Note instance to JSON format.
func (n *Note) ToJson() ([]byte, error) {
	return json.Marshal(n)
}

// NewNote creates a new Note instance with the current timestamp and a generated ID.
func NewNote(headline, content string, tags []string) *Note {
	timestamp := utils.GetCurrentTimestamp()
	noteId := utils.NewNoteId()
	return &Note{
		Id:        noteId,
		Headline:  headline,
		Content:   content,
		Tags:      tags,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}
}

func (n *Note) GetId() utils.NoteId {
	return n.Id
}
func (n *Note) GetHeadline() string {
	return n.Headline
}
func (n *Note) SetHeadline(headline string) {
	n.Headline = headline
}
func (n *Note) GetContent() string {
	return n.Content
}
func (n *Note) SetContent(content string) {
	n.Content = content
}
func (n *Note) GetTags() []string {
	return n.Tags
}
func (n *Note) SetTags(tags []string) {
	n.Tags = tags
}
func (n *Note) GetCreatedAt() string {
	return n.CreatedAt
}
func (n *Note) GetUpdatedAt() string {
	return n.UpdatedAt
}
func (n *Note) SetUpdatedAt(updatedAt string) {
	n.UpdatedAt = updatedAt
}
