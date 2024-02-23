package main

import (
	"strconv"
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
		{
			name: "飲む",
			args: args{
				FirstRunDate: parseTime("2024-02-14"),
				CurrentDate:  parseTime("2024-02-23"),
			},
			want: false,
		},
		{
			name: "飲む",
			args: args{
				FirstRunDate: parseTime("2024-02-14"),
				CurrentDate:  parseTime("2024-02-14"),
			},
			want: false,
		},
		{
			name: "休み",
			args: args{
				FirstRunDate: parseTime("2024-02-14"),
				CurrentDate:  parseTime("2024-03-10"),
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := calculateDay(tt.args.FirstRunDate, tt.args.CurrentDate); err != nil || got != tt.want {
				if err != nil {
					t.Errorf("calculateDay() error = %v", err)
					return
				}
				t.Errorf("calculateDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func parseTime(dateString string) time.Time {
	date, _ := time.Parse("2006-01-02", dateString)
	return date
}

func TestGetRandomImage(t *testing.T) {
	type args struct {
		imageURLs []string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "URLが3つの場合",
			args: args{
				imageURLs: generateImageURLs(3),
			},
		},
		{
			name: "URLが10の場合",
			args: args{
				imageURLs: generateImageURLs(10),
			},
		},
		{
			name: "URLが100の場合",
			args: args{
				imageURLs: generateImageURLs(100),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// imageが取得できるか
			if got, err := getRandomImage(tt.args.imageURLs); err != nil || len(got) == 0 {
				if err != nil {
					t.Errorf("getRandomImage() error = %v", err)
					return
				}
				t.Errorf("getRandomImage() = %v, want %v", got, tt.want)

				// ランダムに取得されているか
				frequency := make(map[string]float64)
				numIterations := 10000
				for i := 0; i < numIterations; i++ {
					got, _ := getRandomImage(tt.args.imageURLs)
					frequency[got]++
				}
				expectedFrequency := float64(numIterations) / float64(len(tt.args.imageURLs))
				for imageURL, freq := range frequency {
					if freq < expectedFrequency*0.8 || freq > expectedFrequency*1.2 {
						t.Errorf("Image %s の頻度が期待される頻度 %f に対して多すぎるか少なすぎます。", imageURL, expectedFrequency)
					}
				}
			}
		})
	}
}

func generateImageURLs(numImages int) []string {
	imageURLs := make([]string, numImages)
	for i := 0; i < numImages; i++ {
		imageURLs[i] = "image" + strconv.Itoa(i+1) + ".jpg"
	}
	return imageURLs
}
