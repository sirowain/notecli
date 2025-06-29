package utils

import "errors"

// ErrEmptyHeadlineAndContent is returned when both headline and content are empty
var ErrEmptyHeadlineAndContent = errors.New("headline and content cannot be both empty")

// ErrNoteIdRequired is returned when a note ID is required but not provided
var ErrNoteIdRequired = errors.New("note ID is required")

// ErrNoteNotFound is returned when a note with the specified ID does not exist
var ErrNoteNotFound = errors.New("note not found")

// ErrNoteCreationFailed is returned when note creation fails
var ErrNoteCreationFailed = errors.New("note creation failed")

// ErrNoteUpdateFailed is returned when note update fails
var ErrNoteUpdateFailed = errors.New("note update failed")

// ErrNoteDeletionFailed is returned when note deletion fails
var ErrNoteDeletionFailed = errors.New("note deletion failed")
