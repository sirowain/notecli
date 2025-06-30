# notecli

A simple command-line tool for managing notes, with tagging and clipboard integration.

## Features

- Add, update, list, show, edit, and delete notes
- Tagging support for organizing notes
- Clipboard integration for quick note creation and copying
- Edit notes using your preferred editor
- Local database storage (BoltDB)

## Installation

Clone the repository and build the binary:

```sh
$ git clone https://github.com/sirowain/notecli.git
$ cd notecli
$ go build -o notecli
```

## Usage

Run `notecli --help` to see the available commands:

```sh
$ ./notecli --help
```

### Global Usage

```sh
$ notecli [command] [flags]
```

### Commands and Flags

#### add (a)

Add a note to the list.

**Flags:**
- `-h`, `--headline <headline>`: Headline of the note
- `-t`, `--tags <tag1,tag2,...>`: Tags for the note (comma-separated)
- `-e`, `--editor`: Use your `$EDITOR` to add note content
- `-c`, `--clipboard`: Add note from clipboard content

**Examples:**
```sh
$ notecli add -h "Meeting Notes" -t work,meeting "Discussed project timeline"
$ notecli add -e -h "Journal" -t personal
$ notecli add -c -h "Copied from clipboard"
```

#### update (u)

Update a note by ID.

**Flags:**
- `-c`, `--content <content>`: New content for the note
- `-h`, `--headline <headline>`: New headline
- `-t`, `--tags <tag1,tag2,...>`: New tags (comma-separated)

**Example:**
```sh
$ notecli update 1 -c "Updated content" -h "Updated headline" -t updated,tag
```

#### list (l)

List all notes.

**Flags:**
- `-t`, `--tags <tag1,tag2,...>`: Filter notes by tags (comma-separated)

**Example:**
```sh
$ notecli list -t work
Notes found (2):
[001] Meeting Notes
```

#### show (s)

Show a note by ID.

**Flags:**
- `-c`, `--clipboard`: Copy note content to clipboard

**Example:**
```sh
$ notecli show 1

001 | Meeting Notes
--------------------------------------------------
Discussed project timeline

--------------------------------------------------
Tags: work, meeting
Created: 2025-06-30T10:21:46+02:00 | Updated: 2025-06-30T10:22:09+02:00
```

#### edit (e)

Edit a note in your default editor.

**Example:**
```sh
$ notecli edit 1
```

#### delete (d)

Delete a note by ID.

**Example:**
```sh
$ notecli delete 1
```

## Configuration

The database is stored in your user config directory under `notecli/my.db` (e.g., `~/.config/notecli/my.db` on Linux).

## License

MIT
