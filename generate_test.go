package substate_test

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/nickcorin/substate"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "Updates golden files")

func assertGolden(t *testing.T, goldenFile string, actual []byte) {
	t.Helper()

	if *update {
		err := ioutil.WriteFile(goldenFile, actual, 0644)
		require.NoError(t, err)
		return
	}

	_, err := os.Stat(goldenFile)
	require.NoError(t, err)

	expected, err := ioutil.ReadFile(goldenFile)
	require.NoError(t, err)
	require.NotNil(t, actual)

	require.True(t, bytes.Equal(expected, actual))
}

func TestGenerate(t *testing.T) {

	type test struct {
		input  string
		golden string
		err    error
	}

	tests := map[string]test{
		"simple valid interface": {
			input:  "simple.go",
			golden: "simple.golden",
			err:    nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dir, err := os.MkdirTemp("/tmp", "gensubstate")
			require.NoError(t, err)

			tmp, err := os.CreateTemp(dir, "generate-*.go")
			require.NoError(t, err)

			t.Cleanup(func() {
				err := os.Remove(tmp.Name())
				require.NoError(t, err)
			})

			path := filepath.Join("testdata", test.input)
			err = substate.Generate(path, tmp.Name(), "Substate")
			require.Equal(t, test.err, err)

			actual, err := ioutil.ReadFile(tmp.Name())
			require.NoError(t, err)

			goldenPath := filepath.Join("testdata", test.golden)
			assertGolden(t, goldenPath, actual)
		})
	}
}
