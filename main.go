package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/sirowain/notecli/pkg/engine"
	"github.com/sirowain/notecli/pkg/engine/localdb"
	"github.com/sirowain/notecli/pkg/utils"
	"github.com/urfave/cli/v3"
)

func main() {
	noteEngine, err := setupDatabase()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error initializing note engine:", err)
		panic(err)
	}

	// Close the engine when done
	defer func() {
		if closeErr := noteEngine.Close(); closeErr != nil {
			fmt.Fprintln(os.Stderr, "Error closing note engine:", closeErr)
		}
	}()

	cmd := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "add a note to the list",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "headline",
						Aliases: []string{"h"},
						Usage:   "headline of the note",
					},
					&cli.StringSliceFlag{
						Name:    "tags",
						Aliases: []string{"t"},
						Usage:   "tags of the note (comma-separated)",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					headline := cmd.String("headline")
					tags := cmd.StringSlice("tags")
					return addNote(noteEngine, cmd.Args().First(), headline, tags)
				},
			},
			{
				Name:    "update",
				Aliases: []string{"u"},
				Usage:   "update a note by id",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "content",
						Aliases: []string{"c"},
						Usage:   "content of the note",
					},
					&cli.StringFlag{
						Name:    "headline",
						Aliases: []string{"h"},
						Usage:   "headline of the note",
					},
					&cli.StringSliceFlag{
						Name:    "tags",
						Aliases: []string{"t"},
						Usage:   "tags of the note (comma-separated)",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					content := cmd.String("content")
					headline := cmd.String("headline")
					tags := cmd.StringSlice("tags")
					return updateNote(noteEngine, cmd.Args().First(), content, headline, tags)
				},
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "list all notes on the list",
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:    "tags",
						Aliases: []string{"t"},
						Usage:   "filter notes by tags",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					tags := cmd.StringSlice("tags")
					return listNotes(noteEngine, tags)
				},
			},
			{
				Name:    "show",
				Aliases: []string{"s"},
				Usage:   "show a note by id",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return showNote(noteEngine, cmd.Args().First())
				},
			},
			{
				Name:    "edit",
				Aliases: []string{"e"},
				Usage:   "edit a note in the default editor",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return editNote(noteEngine, cmd.Args().First())
				},
			},
			{
				Name:    "delete",
				Aliases: []string{"d"},
				Usage:   "delete a note by id",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "all",
						Aliases: []string{"a"},
						Usage:   "delete all notes",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Bool("all") {
						// return deleteAllNotes() //TODO
						fmt.Fprintln(os.Stderr, "Delete all notes functionality is not implemented yet.")
						return nil
					}
					return deleteNote(noteEngine, cmd.Args().First())
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func deleteNote(noteEngine engine.NoteEngine, noteId string) error {
	if err := noteEngine.DeleteNote(noteId); err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}
	fmt.Println("Note deleted successfully.")
	return nil
}

func showNote(noteEngine engine.NoteEngine, noteId string) error {
	note, err := noteEngine.ReadNote(noteId)
	if err != nil {
		return fmt.Errorf("failed to read note: %w", err)
	}
	if note == nil {
		fmt.Println("Note not found.")
		return nil
	}
	fmt.Printf("Id: %s\nHeadline: %s\nContent: %s\nTags: %v\nCreated At: %s\nUpdated At: %s\n",
		note.GetId(), note.GetHeadline(), note.GetContent(), note.GetTags(), note.CreatedAt, note.UpdatedAt)
	return nil
}

func addNote(noteEngine engine.NoteEngine, content, headline string, tags []string) error {
	note, err := noteEngine.CreateNote(content, headline, tags)
	if err != nil {
		return fmt.Errorf("failed to create note: %w", err)
	}
	fmt.Printf("Note created successfully with ID: %s\n", note.GetId())
	return nil
}

func updateNote(noteEngine engine.NoteEngine, noteId, content, headline string, tags []string) error {
	err := noteEngine.UpdateNote(noteId, content, headline, tags)
	if err != nil {
		return fmt.Errorf("failed to update note: %w", err)
	}
	fmt.Println("Note updated successfully.")
	return nil
}

func listNotes(noteEngine engine.NoteEngine, tags []string) error {
	notes, err := noteEngine.ListNotes(tags)
	if err != nil {
		return fmt.Errorf("failed to list notes: %w", err)
	}
	if len(notes) == 0 {
		fmt.Println("No notes found.")
		return nil
	}
	fmt.Println("Notes found:")
	for _, note := range notes {
		description := note.GetHeadline()
		if note.Headline == "" {
			description = note.GetContent()
		}
		fmt.Printf("[%s] %s\n",
			note.GetId(), description)
	}
	return nil
}

func editNote(noteEngine engine.NoteEngine, noteId string) error {
	note, err := noteEngine.ReadNote(noteId)
	if err != nil {
		return utils.ErrNoteNotFound
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	// Temp file
	tmpfile, err := os.CreateTemp("", "notecli-buffer-*.txt")
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(note.GetContent())); err != nil {
		return err
	}
	tmpfile.Close()

	// Open the editor synchronously
	cmd := exec.Command(editor, tmpfile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	// Read the edited content
	editedContent, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		return err
	}

	// Update the note with the edited content
	if err := noteEngine.UpdateNote(noteId, string(editedContent), note.GetHeadline(),
		note.GetTags()); err != nil {
		return fmt.Errorf("failed to update note: %w", err)
	}
	fmt.Println("Note updated successfully.")
	return nil
}

func setupDatabase() (engine.NoteEngine, error) {
	confDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	appConfDir := confDir + "/notecli"
	if _, err := os.Stat(appConfDir); os.IsNotExist(err) {
		if err := os.MkdirAll(appConfDir, 0755); err != nil {
			return nil, err
		}
	}
	// Create the database file path
	dbPath := appConfDir + "/my.db"
	dbEngine := localdb.NewLocalDBEngine()
	if err := dbEngine.Initialize(dbPath, nil); err != nil {
		return nil, err
	}
	return dbEngine, nil
}
