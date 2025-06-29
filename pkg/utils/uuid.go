package utils

import "github.com/google/uuid"

type NoteId string

func NewNoteId() NoteId {
	return NoteId(generateUUID())
}

func generateUUID() string {
	id := uuid.New()
	return id.String()[:8]
}
