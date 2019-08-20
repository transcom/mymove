package unit

import "golang.org/x/text/language"
import "golang.org/x/text/message"

// Miless represents mile value in int
type Miles int

// String gives a string representation of Miles
func (miles Miles) String() string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d", int(miles))
}
