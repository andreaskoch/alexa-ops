package main

import (
	"fmt"
	"sort"
)

type localizer interface {
	// Localize the given key and parameters for the specified culture (e.g. "en-US", "de-DE")
	Localize(culture string, key string, p ...interface{}) (string, error)
}

func newLocalization(key, defaultCulture, defaultValue string) localization {
	defaultValues := make(map[string]string)
	defaultValues[defaultCulture] = defaultValue

	return localization{
		key:             key,
		valuesByCulture: defaultValues,
	}
}

type localization struct {
	key             string
	valuesByCulture map[string]string
}

func (localization *localization) Key() string {
	return localization.key
}

func (localization *localization) Cultures() []string {
	var keys []string
	for key := range localization.valuesByCulture {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}

func (localization *localization) Value(culture string, p ...interface{}) (string, error) {
	value, exists := localization.valuesByCulture[culture]
	if !exists {
		return localization.Key(), fmt.Errorf("The key %q has no value for culture %q", localization.Key(), culture)
	}

	return fmt.Sprintf(value, p...), nil
}

func (localization *localization) Add(culture, value string) *localization {
	if localization.valuesByCulture == nil {
		localization.valuesByCulture = make(map[string]string)
	}

	localization.valuesByCulture[culture] = value
	return localization
}

type localizations map[string]localization

type inMemoryLocalizer struct {
	localizations localizations
}

// Localize the given key and parameters for the specified culture (e.g. "en-US", "de-DE")
func (localizer *inMemoryLocalizer) Localize(key, culture string, p ...interface{}) (string, error) {
	localization, keyExists := localizer.localizations[key]
	if !keyExists {
		return key, fmt.Errorf("The localization key %q does not exist", key)
	}

	// first try (by culture code)
	localizedValue, localizationError := localization.Value(culture, p...)
	if localizationError != nil {

		// second try (by locale)
		locale, localeError := getLocaleFromCultureCode(culture)
		if localeError != nil {
			return fmt.Sprintf(key, p...), fmt.Errorf("The localization key %q does not have value for culture %q", localization.Key(), culture)
		}

		localizedValueByLocale, secondLocalizationError := localization.Value(locale, p...)
		if secondLocalizationError != nil {
			return fmt.Sprintf(key, p...), fmt.Errorf("The localization key %q does not have value for culture %q or %q", localization.Key(), culture, locale)
		}

		localizedValue = localizedValueByLocale
	}

	return localizedValue, nil
}

func getLocaleFromCultureCode(cultureCode string) (string, error) {
	if len(cultureCode) > 2 {
		return cultureCode[:2], nil
	}

	return "", fmt.Errorf("Cannot determine the locale from %q", cultureCode)
}