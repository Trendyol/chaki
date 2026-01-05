package logger

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseLocation(t *testing.T) {
	tests := []struct {
		name     string
		timezone string
		want     *time.Location
		wantErr  bool
	}{
		{
			name:     "Local",
			timezone: "Local",
			want:     time.Local,
			wantErr:  false,
		},
		{
			name:     "local lowercase",
			timezone: "local",
			want:     time.Local,
			wantErr:  false,
		},
		{
			name:     "UTC",
			timezone: "UTC",
			want:     time.UTC,
			wantErr:  false,
		},
		{
			name:     "utc lowercase",
			timezone: "utc",
			want:     time.UTC,
			wantErr:  false,
		},
		{
			name:     "UTC with offset positive",
			timezone: "UTC+3",
			want:     time.FixedZone("UTC+3", 3*60*60),
			wantErr:  false,
		},
		{
			name:     "UTC with offset negative",
			timezone: "UTC-5",
			want:     time.FixedZone("UTC-5", -5*60*60),
			wantErr:  false,
		},
		{
			name:     "IANA location",
			timezone: "Europe/Istanbul",
			want:     func() *time.Location { loc, _ := time.LoadLocation("Europe/Istanbul"); return loc }(),
			wantErr:  false,
		},
		{
			name:     "empty timezone",
			timezone: "",
			want:     time.Local,
			wantErr:  false,
		},
		{
			name:     "invalid offset",
			timezone: "UTC+abc",
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "invalid IANA location",
			timezone: "Invalid/Location",
			want:     nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseLocation(tt.timezone)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.String(), got.String())
			}
		})
	}
}
