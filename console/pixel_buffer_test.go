package console

import "testing"

type testPos struct {
	testValue     pos
	expectedValue pos
}

func TestPixelToChar(t *testing.T) {
	tests := []testPos{
		{testValue: pos{x: 10, y: 10}, expectedValue: pos{x: 40, y: 80}},
		{testValue: pos{x: 0, y: 0}, expectedValue: pos{x: 0, y: 0}},
	}
	for _, test := range tests {
		pixelPos := charToPixel(test.testValue)
		if pixelPos != test.expectedValue {
			t.Error(
				"For", test.testValue,
				"expected", test.expectedValue,
				"got", pixelPos,
			)
		}
	}
}

func TestCharToPixel(t *testing.T) {
	tests := []testPos{
		{testValue: pos{x: 64, y: 64}, expectedValue: pos{x: 16, y: 8}},
		{testValue: pos{x: 0, y: 0}, expectedValue: pos{x: 0, y: 0}},
	}
	for _, test := range tests {
		charPos := pixelToChar(test.testValue)
		if charPos != test.expectedValue {
			t.Error(
				"For", test.testValue,
				"expected", test.expectedValue,
				"got", charPos,
			)
		}
	}
}
