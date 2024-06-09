package interfaces

type MusicSearchResult struct {
	Name string
	URL  string
}

type MusicSearchProvider interface {
	SearchPlay(string) (MusicSearchResult, error)
}
