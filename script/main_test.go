package main

import (
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestCalculateDay(t *testing.T) {
	type args struct {
		FirstRunDate time.Time
		CurrentDate  time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "休み",
			args: args{
				FirstRunDate: parseTime("2024-02-14"),
				CurrentDate:  parseTime("2024-03-06"),
			},
			want: true,
		},
		{
			name: "飲む",
			args: args{
				FirstRunDate: parseTime("2024-02-14"),
				CurrentDate:  parseTime("2024-03-05"),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := calculateDay(tt.args.FirstRunDate, tt.args.CurrentDate); got != tt.want {
				t.Errorf("calculateDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func parseTime(dateString string) time.Time {
	date, _ := time.Parse("2006-01-02", dateString)
	return date
}
