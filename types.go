package main

type directory struct {
	Channels []channel `json:"channels"`
}

type channel struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Versions    []version `json:"versions"`
}

type version struct {
	Version   string `json:"version"`
	Changelog string `json:"changelog"`
	Timestamp Time   `json:"timestamp"`
	Files     []file `json:"files"`
}

type file struct {
	URL    string `json:"url"`
	Target string `json:"target"`
	Type   string `json:"type"`
	Sha256 string `json:"sha256"`
	Sha512 string `json:"sha512"`
}
