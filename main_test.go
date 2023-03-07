package main

import (
	"fmt"
	"testing"
	"time"
)

func TestGetDates(t *testing.T) {
	tests := []struct {
		name         string
		now          time.Time
		expectedDays []int
	}{
		{
			"RegularMonday",
			time.Date(2023, time.February, 20, 0, 0, 0, 0, time.Local),
			[]int{20, 21, 22, 23, 24},
		},
		{
			"StartOnTuesday",
			time.Date(2023, time.February, 21, 0, 0, 0, 0, time.Local),
			[]int{20, 21, 22, 23, 24},
		},
		{
			"StartOnWednesday",
			time.Date(2023, time.February, 22, 0, 0, 0, 0, time.Local),
			[]int{20, 21, 22, 23, 24},
		},
		{
			"StartOnThursday",
			time.Date(2023, time.February, 23, 0, 0, 0, 0, time.Local),
			[]int{20, 21, 22, 23, 24},
		},
		{
			"StartOnFriday",
			time.Date(2023, time.February, 24, 0, 0, 0, 0, time.Local),
			[]int{20, 21, 22, 23, 24},
		},
		{
			"StartOnMondayWhenWeekCrossesIntoNewMonth",
			time.Date(2023, time.February, 27, 0, 0, 0, 0, time.Local),
			[]int{27, 28, 1, 2, 3},
		},
		{
			"StartOnTuesdayWhenWeekCrossesIntoNewMonth",
			time.Date(2023, time.February, 28, 0, 0, 0, 0, time.Local),
			[]int{27, 28, 1, 2, 3},
		},
		{
			"StartOnWednesdayWhenWeekCrossesIntoNewMonth",
			time.Date(2023, time.March, 1, 0, 0, 0, 0, time.Local),
			[]int{27, 28, 1, 2, 3},
		},
		{
			"StartOnThursdayWhenWeekCrossesIntoNewMonth",
			time.Date(2023, time.March, 2, 0, 0, 0, 0, time.Local),
			[]int{27, 28, 1, 2, 3},
		},
		{
			"StartOnFridayWhenWeekCrossesIntoNewMonth",
			time.Date(2023, time.March, 3, 0, 0, 0, 0, time.Local),
			[]int{27, 28, 1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetDates(tt.now)
			for i, d := range result {
				if tt.expectedDays[i] != d.Day() {
					t.Errorf("expected %d but got %d", tt.expectedDays[i], d.Day())
				}
			}
		})
	}
}

func TestGetDayWithSuffix(t *testing.T) {
	tests := []struct {
		day      int
		expected string
	}{
		{1, "01st"},
		{2, "02nd"},
		{3, "03rd"},
		{4, "04th"},
		{5, "05th"},
		{6, "06th"},
		{7, "07th"},
		{8, "08th"},
		{9, "09th"},
		{10, "10th"},
		{11, "11th"},
		{12, "12th"},
		{13, "13th"},
		{14, "14th"},
		{15, "15th"},
		{16, "16th"},
		{17, "17th"},
		{18, "18th"},
		{19, "19th"},
		{20, "20th"},
		{21, "21st"},
		{22, "22nd"},
		{23, "23rd"},
		{24, "24th"},
		{25, "25th"},
		{26, "26th"},
		{27, "27th"},
		{28, "28th"},
		{29, "29th"},
		{30, "30th"},
		{31, "31st"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d", tt.day), func(t *testing.T) {
			result := DayWithSuffix(tt.day)
			if tt.expected != result {
				t.Errorf("expected %s but got %s", tt.expected, result)
			}
		})
	}
}
