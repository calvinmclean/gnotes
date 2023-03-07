package main

import (
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
{{- $dayOne := (index . 0) -}}
# Week of the {{DayWithSuffix $dayOne.Day}} ({{$dayOne.Month}} {{$dayOne.Year}})
{{range $index, $date := .}}
### {{$date.Weekday}} {{printf "%02d" $date.Day}}{{if ne $date.Month $dayOne.Month}} ({{$date.Month}} {{$date.Year}}){{end}}

{{end}}`

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
	return template.Must(
		template.
			New("note").
			Funcs(template.FuncMap{
				"DayWithSuffix": DayWithSuffix,
			}).
			Parse(noteTemplate),
	).Execute(wr, dates)
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
		// append
		// return os.OpenFile(path.Join(dir, filename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

// Ideas
// - Maybe I should have it only add headers for a new day when I run it on that day which will help me remember to switch sections and will reduce clutter...
// - Always copy TODOs to a central location and remove them once the original is checked off
