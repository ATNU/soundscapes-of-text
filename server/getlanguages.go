package main

// GetLanguages retrieved supported languages
//
// Returns an array of languages composed of their name and code
func GetLanguages() []struct {
	Name string
	Code string
} {

	return []struct {
		Name string
		Code string
	}{
		{"Danish", "da - DK"},
		{"Dutch", "nl-NL"},
		{"English (Australian)", "en-AU"},
		{"English (British)", "en-GB"},
		{"English (Indian)", "en-IN"},
		{"English (US)", "en-US"},
		{"English (Welsh)", "en-GB-WL"},
		{"French", "fr-FR"},
		{"French (Canadian)", "fr-CA"},
		{"German", "de-DE"},
		{"Icelandic", "is-IS"},
		{"Italian", "it-IT"},
		{"Japanese", "ja-JP"},
		{"Korean", "ko-KR"},
		{"Polish", "pl-PL"},
		{"Portuguese (Brazilian)", "pt-BR"},
		{"Portuguese (European)", "pt-PT"},
		{"Romanian", "ro-RO"},
		{"Russian", "ru-RU"},
		{"Spanish (European)", "es-ES"},
		{"Spanish (Latin American)", "es-US"},
		{"Swedish", "sv-SE"},
		{"Turkish", "tr-TR"},
		{"Welsh", "cy-GB"},
	}
}
