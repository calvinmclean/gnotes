package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"time"
)

const noteTemplate = `
{{- $dayOne := (index .Dates 0) -}}
# Week of the {{ DayWithSuffix $dayOne.Day }} ({{ $dayOne.Month }} {{ $dayOne.Year }})
{{ if .TODOs }}
### Last Week's Remaining TODOs
{{ range $index, $todo := .TODOs }}
{{- $todo }}
{{ end }}
{{ end -}}

{{ range $index, $date := .Dates }}
### {{ $date.Weekday }} {{ printf "%02d" $date.Day }}{{ if ne $date.Month $dayOne.Month }} ({{ $date.Month }} {{ $date.Year }}){{ end }}

{{ end }}`

func main() {
	now := time.Now()
	dates := GetDates(now)

	file, err := CreateDirectoryAndFile(dates[0])
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer file.Close()

	if file == nil {
		return
	}

	err = GenerateNote(dates, file)
	if err != nil {
		log.Fatalf("error creating notes from template: %v", err)
	}
}

func GenerateNote(dates []time.Time, wr io.Writer) error {
	todos, err := GetLastWeekTODOs(dates[0])
	if err != nil {
		fmt.Println("error getting TODOs from last week: %w", err)
	}

	return template.Must(
		template.
			New("note").
			Funcs(template.FuncMap{
				"DayWithSuffix": DayWithSuffix,
			}).
			Parse(noteTemplate),
	).Execute(wr, map[string]interface{}{
		"Dates": dates,
		"TODOs": todos,
	})
}

func GetDates(now time.Time) []time.Time {
	dates := make([]time.Time, 5)
	for i := 0; i < 5; i++ {
		dates[i] = now.AddDate(0, 0, i-int(now.Weekday())+1)
	}
	return dates
}

func DayWithSuffix(day int) string {
	return fmt.Sprintf("%02d%s", day, GetSuffix(day))
}

func GetSuffix(day int) string {
	switch day {
	case 1, 21, 31:
		return "st"
	case 2, 22:
		return "nd"
	case 3, 23:
		return "rd"
	default:
		return "th"
	}
}

func CreateDirectoryAndFile(dayOne time.Time) (*os.File, error) {
	dir, filename := NotesPath(dayOne)

	contents, err := os.ReadDir(dir)
	if err != nil && !DirNotExists(err) {
		return nil, fmt.Errorf("error reading directory: %w", err)
	} else if DirNotExists(err) {
		err = os.MkdirAll(dir, 0750)
		if err != nil {
			return nil, fmt.Errorf("error creating directory: %w", err)
		}
	}

	fileExists := false
	for _, entry := range contents {
		if entry.Name() == filename {
			fileExists = true
			break
		}
	}

	if fileExists {
		return nil, nil
	}

	return os.Create(path.Join(dir, filename))
}

func IsPathError(err error) bool {
	var pathErr *fs.PathError
	return errors.As(err, &(pathErr))
}

func DirNotExists(err error) bool {
	switch v := err.(type) {
	case *fs.PathError:
		return v.Unwrap().Error() == "no such file or directory"
	default:
		return false
	}
}

func NotesPath(dayOne time.Time) (string, string) {
	dir := path.Join("./", fmt.Sprintf("%d", dayOne.Year()), fmt.Sprintf("%02d%s", dayOne.Month(), dayOne.Month().String()[0:3]))
	filename := fmt.Sprintf("Week_of_the_%s.md", DayWithSuffix(dayOne.Day()))
	return dir, filename
}

func GetLastWeekTODOs(dayOne time.Time) ([]string, error) {
	lastWeekDayOne := dayOne.AddDate(0, 0, -7)

	data, err := os.ReadFile(path.Join(NotesPath(lastWeekDayOne)))
	if err != nil {
		return nil, fmt.Errorf("error reading last week's note: %w", err)
	}

	todos := []string{}
	for _, line := range bytes.Split(data, []byte{'\n'}) {
		line = bytes.TrimSpace(line)

		if len(line) < 6 {
			continue
		}

		if bytes.Equal(line[:6], []byte{'-', ' ', '[', ' ', ']', ' '}) {
			todos = append(todos, string(line))
		}
	}

	return todos, nil
}
