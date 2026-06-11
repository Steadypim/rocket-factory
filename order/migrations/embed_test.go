package migrations

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMigrationsAreEmbedded(t *testing.T) {
	entries, err := FS.ReadDir(".")

	require.NoError(t, err)
	require.NotEmpty(t, entries)
	require.Equal(t, "00001_order_table.sql", entries[0].Name())
}
