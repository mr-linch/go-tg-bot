package locales

import (
	"embed"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

var (
	//go:embed active.*.yaml
	fs embed.FS

	// Default is the default language code.
	Default = language.English

	Meta = map[string]struct {
		Label string
		Emoji string
	}{
		"en": {
			Label: "English",
			Emoji: "🇺🇸",
		},
		"uk": {
			Label: "Українська",
			Emoji: "🇺🇦",
		},
	}
)

func NewBundle() (*i18n.Bundle, error) {
	bundle := i18n.NewBundle(Default)

	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	for _, path := range []string{
		"active.uk.yaml",
	} {
		if _, err := bundle.LoadMessageFileFS(fs, path); err != nil {
			return nil, errors.Wrapf(err, "load locale %s", path)
		}
	}

	return bundle, nil
}
