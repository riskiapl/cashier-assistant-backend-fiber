package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var (
	bundle *i18n.Bundle
	once   sync.Once
)

// I18n represents the translation utility
type I18n struct {
	bundle       *i18n.Bundle
	defaultLang  string
	localizer    *i18n.Localizer
	supportLangs []language.Tag
}

// NewI18n creates a new instance of I18n
func NewI18n() (*I18n, error) {
	once.Do(func() {
		bundle = i18n.NewBundle(language.English)
		bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

		// Load all translation files from locales directory
		localesDir := "locales"
		files, err := os.ReadDir(localesDir)
		if err != nil {
			fmt.Printf("Error reading locales directory: %v\n", err)
			return
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			if filepath.Ext(file.Name()) == ".json" {
				// Load the translation file
				bundle.MustLoadMessageFile(filepath.Join(localesDir, file.Name()))
			}
		}
	})

	// Default supported languages
	supportedLangs := []language.Tag{
		language.English,
		language.Indonesian,
	}

	// Create a default localizer with English
	defaultLocalizer := i18n.NewLocalizer(bundle, language.English.String())

	return &I18n{
		bundle:       bundle,
		defaultLang:  language.English.String(),
		localizer:    defaultLocalizer,
		supportLangs: supportedLangs,
	}, nil
}

// SetLanguage sets the language for the current request
func (i *I18n) SetLanguage(lang string) {
	if lang == "" {
		lang = i.defaultLang
	}

	// If lang is a BCP 47 language tag (e.g., "en-US"), extract the base language
	tag, err := language.Parse(lang)
	if err == nil {
		base, _ := tag.Base()
		lang = base.String()
	}

	i.localizer = i18n.NewLocalizer(i.bundle, lang, i.defaultLang)
}

// Translate translates a message ID with optional template data
func (i *I18n) Translate(messageID string, templateData map[string]interface{}) string {
	msg, err := i.localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
	})

	if err != nil {
		// Fallback to the message ID if translation fails
		return messageID
	}

	return msg
}

// Global instance for direct access
var GlobalI18n *I18n

// InitGlobalI18n initializes the global I18n instance
func InitGlobalI18n() error {
	i18n, err := NewI18n()
	if err != nil {
		return err
	}
	GlobalI18n = i18n
	return nil
}

// T is a shorthand function for translation
func T(messageID string, templateData ...map[string]interface{}) string {
	if GlobalI18n == nil {
		return messageID
	}

	// Handle optional template data
	data := map[string]interface{}{}
	if len(templateData) > 0 {
		data = templateData[0]
	}

	return GlobalI18n.Translate(messageID, data)
}
