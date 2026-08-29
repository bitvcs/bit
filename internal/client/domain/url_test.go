package domain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseNipaUrl_Valid(t *testing.T) {
	tests := []struct {
		name      string
		urlString string
		want      *NipaUrl
	}{
		{
			name:      "https with port",
			urlString: "https://example.com:9000/org/project",
			want: &NipaUrl{
				Url:     "https://example.com:9000/org/project",
				Host:    "example.com:9000",
				Org:     "org",
				Project: "project",
			},
		},
		{
			name:      "http without port",
			urlString: "http://example.com/org/project",
			want: &NipaUrl{
				Url:     "http://example.com/org/project",
				Host:    "example.com",
				Org:     "org",
				Project: "project",
			},
		},
		{
			name:      "path with only project defaults org",
			urlString: "https://example.com/project",
			want: &NipaUrl{
				Url:     "https://example.com/project",
				Host:    "example.com",
				Org:     "default",
				Project: "project",
			},
		},
		{
			name:      "path with no segments defaults org and empty project",
			urlString: "https://example.com",
			want: &NipaUrl{
				Url:     "https://example.com",
				Host:    "example.com",
				Org:     "default",
				Project: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseNipaUrl(tt.urlString)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestParseNipaUrl_InvalidScheme(t *testing.T) {
	_, err := ParseNipaUrl("ftp://example.com/org/project")
	require.Error(t, err)
	require.EqualError(t, err, "invalid URL scheme, expected http or https")

	var domainErr *Error
	require.True(t, errors.As(err, &domainErr))
	require.Equal(t, 400, domainErr.Code)
}

func TestParseNipaUrl_InvalidURL(t *testing.T) {
	_, err := ParseNipaUrl("://bad url")
	require.Error(t, err)
	require.Contains(t, err.Error(), "url can not be parsed")
}

func TestParseNipaUrl_InvalidPath(t *testing.T) {
	_, err := ParseNipaUrl("https://example.com/org/project/extra")
	require.Error(t, err)
	require.EqualError(t, err, "invalid URL path, expected /org/project")
}