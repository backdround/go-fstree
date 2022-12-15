// Package config parses user config to abstract fs node structure.
// For exapmle:
/*
configs:
  "config1.txt":
    type: file
    data: "format: txt"
pkg:
  pkg1:
    type: link
    path: "../../pkg1"
*/
// will be parsed in an appropriate node tree structure based with
// entries.DirectoryEntry type.
package config

import (
	"errors"
	"fmt"
	"path"

	"github.com/backdround/go-fstree/entries"
	"github.com/backdround/go-indent"
	"gopkg.in/yaml.v3"
)

type rawEntry = map[string]any

type ParseError struct {
	Message string
	Path    string
}

func (e *ParseError) Error() string {
	indentedMessage := indent.Indent(e.Message, "  ", 1)
	resultMessage := fmt.Sprintf("unable to parse %v:\n%v", e.Path,
		indentedMessage)
	return resultMessage
}

func parseAny(name string, entry rawEntry) (parsedEntry entries.Entry,
	err *ParseError) {
	entryType, ok := entry["type"]
	if !ok {
		return parseDirectory(name, entry)
	}

	switch entryType {
	case "file":
		return parseFile(name, entry)
	case "link":
		return parseLink(name, entry)
	default:
		err := &ParseError{
			Message: fmt.Sprintf(`unknown type: %v`, entryType),
			Path:    name,
		}
		return nil, err
	}
}

func parseDirectory(name string, entry rawEntry) (entries.DirectoryEntry,
	*ParseError) {
	if _, ok := entry["type"]; ok {
		panic(`unexpected "type" property`)
	}

	// A constructed entry
	currentEntry := entries.DirectoryEntry{
		Name:    name,
		Entries: make([]entries.Entry, 0),
	}

	// Parses sub entires
	for subEntryName, subEntryAny := range entry {
		subEntry, ok := subEntryAny.(rawEntry)
		if subEntryAny != nil && !ok {
			parseError := ParseError{
				Message: "unable to convert to dictionary",
				Path:    path.Join(name, subEntryName),
			}
			return entries.DirectoryEntry{}, &parseError
		}

		parsedEntry, err := parseAny(subEntryName, subEntry)
		if err != nil {
			err.Path = path.Join(name, err.Path)
			return entries.DirectoryEntry{}, err
		}

		currentEntry.Entries = append(currentEntry.Entries, parsedEntry)
	}

	return currentEntry, nil
}

func parseFile(name string, entry rawEntry) (entries.FileEntry,
	*ParseError) {
	// Asserts type property
	typeValue, ok := entry["type"]
	if !ok || typeValue != "file" {
		panic(fmt.Sprintf("unexpected type property: %v", typeValue))
	}
	delete(entry, "type")

	// Returns error result
	errorResult := func(errorMessage string) (entries.FileEntry, *ParseError) {
		parseError := ParseError{
			Message: errorMessage,
			Path:    name,
		}
		return entries.FileEntry{}, &parseError
	}

	// A constructed entry
	fileEntry := entries.FileEntry{
		Name: name,
	}

	// Parses file properties
	for name, valueAny := range entry {
		switch name {
		case "data":
			value, ok := valueAny.(string)
			if !ok {
				message := fmt.Sprintf("unable to convert data to string: %v",
					valueAny)
				return errorResult(message)
			}
			fileEntry.Data = []byte(value)
		default:
			return errorResult("unknown property: " + name)
		}
	}

	return fileEntry, nil
}

func parseLink(name string, entry rawEntry) (entries.LinkEntry,
	*ParseError) {
	typeValue, ok := entry["type"]
	if !ok || typeValue != "link" {
		panic(fmt.Sprintf("unexpected type property: %v", typeValue))
	}
	delete(entry, "type")

	// Returns error result
	errorResult := func(errorMessage string) (entries.LinkEntry,
		*ParseError) {
		parseError := ParseError{
			Message: errorMessage,
			Path:    name,
		}
		return entries.LinkEntry{}, &parseError
	}

	// A constructed entry
	linkEntry := entries.LinkEntry{
		Name: name,
	}

	// Gets path property
	pathValueAny, ok := entry["path"]
	if !ok {
		return errorResult("path property must be set for link")
	}
	delete(entry, "path")

	pathValue, ok := pathValueAny.(string)
	if !ok {
		message := fmt.Sprintf("unable to convert path to string: %v",
			pathValueAny)
		return errorResult(message)
	}
	linkEntry.Path = pathValue

	// Parses link properties
	for name := range entry {
		switch name {
		default:
			return errorResult("unknown property: " + name)
		}
	}

	return linkEntry, nil
}

// Parse parses filetree structure from yaml to the entries.
func Parse(yamlData string) (*entries.DirectoryEntry, error) {
	// Unmarshales to a rawTree
	rawTree := make(rawEntry)
	yamlErr := yaml.Unmarshal([]byte(yamlData), rawTree)
	if yamlErr != nil {
		return nil, yamlErr
	}

	// Checks that a root type property doesn't exist
	if _, ok := rawTree["type"]; ok {
		return nil, errors.New(`unexpected "type" property at root`)
	}

	// Parses the root directory
	rootEntry, err := parseDirectory(".", rawTree)
	if err != nil {
		return nil, err
	}

	return &rootEntry, nil
}
