package offline_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/wakatime/wakatime-cli/pkg/offline"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueueFilepathLegacy(t *testing.T) {
	home, err := os.UserHomeDir()
	require.NoError(t, err)

	tests := map[string]struct {
		ViperValue string
		EnvVar     string
		Expected   string
	}{
		"default": {
			Expected: filepath.Join(home, ".wakatime.bdb"),
		},
		"env_trailing_slash": {
			EnvVar:   "~/path2/",
			Expected: filepath.Join(home, "path2", ".wakatime.bdb"),
		},
		"env_without_trailing_slash": {
			EnvVar:   "~/path2",
			Expected: filepath.Join(home, "path2", ".wakatime.bdb"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := os.Setenv("WAKATIME_HOME", test.EnvVar)
			require.NoError(t, err)

			defer os.Unsetenv("WAKATIME_HOME")

			queueFilepath, err := offline.QueueFilepathLegacy()
			require.NoError(t, err)

			assert.Equal(t, test.Expected, queueFilepath)
		})
	}
}