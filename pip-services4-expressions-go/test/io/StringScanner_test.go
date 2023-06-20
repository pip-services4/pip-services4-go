package test_io

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/io"
)

func TestStringScannerRead(t *testing.T) {
	content := "Test String\nLine2\rLine3\r\n\r\nLine5"
	scanner := io.NewStringScanner(content)
	fixture := NewScannerFixture(scanner, content)

	fixture.TestRead(t)
}

func TestStringScannerUnread(t *testing.T) {
	content := "Test String\nLine2\rLine3\r\n\r\nLine5"
	scanner := io.NewStringScanner(content)
	fixture := NewScannerFixture(scanner, content)

	fixture.TestUnread(t)
}

func TestStringScannerLineColumn(t *testing.T) {
	content := "Test String\nLine2\rLine3\r\n\r\nLine5"
	scanner := io.NewStringScanner(content)
	fixture := NewScannerFixture(scanner, content)

	fixture.TestLineColumn(t, 3, 's', 1, 3)
	fixture.TestLineColumn(t, 12, '\n', 2, 0)
	fixture.TestLineColumn(t, 15, 'n', 2, 3)
	fixture.TestLineColumn(t, 21, 'n', 3, 3)
	fixture.TestLineColumn(t, 26, '\r', 4, 0)
	fixture.TestLineColumn(t, 30, 'n', 5, 3)
}
