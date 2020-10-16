package characterTracker

// Tracker tracks the number of characters translated
type Tracker interface {
	AddCharacters(string) int
	CountAfterString(string) int
}
