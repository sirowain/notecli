package localdb

import (
	"fmt"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/sirowain/notecli/pkg/models"
	"github.com/sirowain/notecli/pkg/utils"
)

type LocalDBEngine struct {
	// db is the database connection
	db *bolt.DB
}

// Initialize the engine with a database connection
func (e *LocalDBEngine) Initialize(dbPath string, options map[string]string) error {
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return err
	}
	e.db = db
	// Create the "notes" bucket if it doesn't exist
	checkAndCreateBucket(e.db, "notes")
	return nil
}

// Close the engine and release resources
func (e *LocalDBEngine) Close() error {
	return e.db.Close()
}

// CreateNote creates a new note in the local database
func (e *LocalDBEngine) CreateNote(content, headline string, tags []string) (models.Note, error) {
	if headline == "" && content == "" {
		return models.Note{}, utils.ErrEmptyHeadlineAndContent
	}
	note := models.NewNote(headline, content, tags)
	if note == nil {
		return models.Note{}, utils.ErrNoteCreationFailed
	}
	err := e.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("notes"))
		fmt.Println("Bucket:", bucket)
		if bucket == nil {
			return ErrBucketNotFound
		}
		fmt.Println("Marshalling note to JSON")
		data, err := note.ToJson()
		fmt.Println("Marshalled data:", string(data))
		if err != nil {
			return err
		}
		if err := bucket.Put([]byte(note.Id), data); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return models.Note{}, err
	}
	// Return the created note
	return *note, nil
}

// ReadNote reads a note by ID from the local database
func (e *LocalDBEngine) ReadNote(noteId string) (*models.Note, error) {
	if noteId == "" {
		return nil, utils.ErrNoteIdRequired
	}
	var note *models.Note
	err := e.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("notes"))
		if bucket == nil {
			return ErrBucketNotFound
		}
		data := bucket.Get([]byte(noteId))
		if data == nil {
			return utils.ErrNoteNotFound
		}
		note = &models.Note{}
		if err := note.FromJson(data); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if note == nil {
		return nil, utils.ErrNoteNotFound
	}
	return note, nil
}

// UpdateNote updates an existing note in the local database
func (e *LocalDBEngine) UpdateNote(noteId, content, headline string, tags []string) error {
	if noteId == "" {
		return utils.ErrNoteIdRequired
	}

	return e.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("notes"))
		if bucket == nil {
			return ErrBucketNotFound
		}
		data := bucket.Get([]byte(noteId))
		if data == nil {
			return utils.ErrNoteNotFound
		}
		note := &models.Note{}
		if err := note.FromJson(data); err != nil {
			return err
		}
		if headline != "" {
			note.SetHeadline(headline)
		}
		if content != "" {
			note.SetContent(content)
		}
		if len(tags) > 0 {
			note.SetTags(tags)
		}
		note.UpdatedAt = utils.GetCurrentTimestamp()

		if note.GetContent() == "" && note.GetHeadline() == "" {
			return utils.ErrEmptyHeadlineAndContent
		}

		newData, err := note.ToJson()
		if err != nil {
			return err
		}

		if err := bucket.Put([]byte(noteId), newData); err != nil {
			return utils.ErrNoteUpdateFailed
		}

		return nil
	})
}

// DeleteNote deletes a note by ID from the local database
func (e *LocalDBEngine) DeleteNote(noteId string) error {
	if noteId == "" {
		return utils.ErrNoteIdRequired
	}
	return e.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("notes"))
		if bucket == nil {
			return ErrBucketNotFound
		}
		if err := bucket.Delete([]byte(noteId)); err != nil {
			return utils.ErrNoteDeletionFailed
		}
		return nil
	})
}

// ListNotes lists all notes in the local database with optional filters
func (e *LocalDBEngine) ListNotes(tags []string) ([]*models.Note, error) {
	var notes []*models.Note
	err := e.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("notes"))
		if bucket == nil {
			return ErrBucketNotFound
		}
		return bucket.ForEach(func(k, v []byte) error {
			note := &models.Note{}
			if err := note.FromJson(v); err != nil {
				return err
			}
			// If tags are provided, filter notes by tags
			if len(tags) > 0 {
				matched := false
				for _, tag := range tags {
					if utils.StringInSlice(tag, note.Tags) {
						matched = true
						break
					}
				}
				if !matched {
					return nil // Skip this note if it doesn't match the tags
				}
			}
			// Append the note to the list
			notes = append(notes, note)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return notes, nil
}

// SearchNotes searches notes by content or headline in the local database
func (e *LocalDBEngine) SearchNotes(query string) ([]*models.Note, error) {
	var notes []*models.Note
	err := e.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("notes"))
		if bucket == nil {
			return ErrBucketNotFound
		}
		return bucket.ForEach(func(k, v []byte) error {
			note := &models.Note{}
			if err := note.FromJson(v); err != nil {
				return err
			}
			if strings.Contains(note.Headline, query) || strings.Contains(note.Content, query) {
				notes = append(notes, note)
			}
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return notes, nil
}

// NewLocalDBEngine creates a new instance of LocalDBEngine
func NewLocalDBEngine() *LocalDBEngine {
	return &LocalDBEngine{}
}

func checkAndCreateBucket(db *bolt.DB, bucketName string) error {
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	})
}
