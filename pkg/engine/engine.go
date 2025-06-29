package engine

import (
	"github.com/sirowain/notecli/pkg/models"
)

type NoteEngine interface {
	// Initialize the engine with a database connection
	Initialize(dbPath string, options map[string]string) error
	// Create a new note
	CreateNote(content, headline string, tags []string) (models.Note, error)
	// Read a note by ID
	ReadNote(noteId string) (*models.Note, error)
	// Update an existing note
	UpdateNote(noteId, content, headline string, tags []string) error
	// Delete a note by ID
	DeleteNote(noteId string) error
	// List all notes with optional tag filter
	ListNotes(tags []string) ([]*models.Note, error)
	// Search notes by content or headline
	SearchNotes(query string) ([]*models.Note, error)
	// Close the engine and release resources
	Close() error
}
