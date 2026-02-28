package teautils

import (
	"strings"
	"testing"
)

func TestCalculateCenter(t *testing.T) {
	t.Run("CentersModalInScreen", func(t *testing.T) {
		row, col := CalculateCenter(80, 24, 40, 10)
		if row != 7 {
			t.Errorf("expected row=7, got %d", row)
		}
		if col != 20 {
			t.Errorf("expected col=20, got %d", col)
		}
	})

	t.Run("ClampsNegativeToZero", func(t *testing.T) {
		// Modal larger than screen
		row, col := CalculateCenter(10, 5, 20, 10)
		if row != 0 {
			t.Errorf("expected row=0 when modal taller than screen, got %d", row)
		}
		if col != 0 {
			t.Errorf("expected col=0 when modal wider than screen, got %d", col)
		}
	})

	t.Run("ModalLargerThanScreen", func(t *testing.T) {
		row, col := CalculateCenter(40, 12, 80, 24)
		if row != 0 {
			t.Errorf("expected row=0, got %d", row)
		}
		if col != 0 {
			t.Errorf("expected col=0, got %d", col)
		}
	})
}

func TestCalculateCenter_OddDimensions(t *testing.T) {
	// screenW=81, screenH=25, modalW=41, modalH=11
	// row = (25-11)/2 = 7, col = (81-41)/2 = 20
	row, col := CalculateCenter(81, 25, 41, 11)
	if row != 7 {
		t.Errorf("expected row=7, got %d", row)
	}
	if col != 20 {
		t.Errorf("expected col=20, got %d", col)
	}

	// screenW=79, screenH=23, modalW=40, modalH=10
	// row = (23-10)/2 = 6, col = (79-40)/2 = 19
	row, col = CalculateCenter(79, 23, 40, 10)
	if row != 6 {
		t.Errorf("expected row=6, got %d", row)
	}
	if col != 19 {
		t.Errorf("expected col=19, got %d", col)
	}
}

func TestMeasureRenderedView(t *testing.T) {
	view := "Hello World\nSecond Line\nThird"
	width, height := MeasureRenderedView(view)
	if width != 11 {
		t.Errorf("expected width=11, got %d", width)
	}
	if height != 3 {
		t.Errorf("expected height=3, got %d", height)
	}
}

func TestMeasureRenderedView_ANSI(t *testing.T) {
	// ANSI escape sequences should not inflate width
	view := "\x1b[31mHello\x1b[0m"
	width, height := MeasureRenderedView(view)
	if width != 5 {
		t.Errorf("expected width=5 (ANSI-aware), got %d", width)
	}
	if height != 1 {
		t.Errorf("expected height=1, got %d", height)
	}
}

func TestMeasureRenderedView_Empty(t *testing.T) {
	width, height := MeasureRenderedView("")
	if width != 0 {
		t.Errorf("expected width=0, got %d", width)
	}
	// strings.Split("", "\n") returns [""], so height=1
	if height != 1 {
		t.Errorf("expected height=1 for empty string (single empty line), got %d", height)
	}
}

func TestMeasureRenderedView_SingleLine(t *testing.T) {
	width, height := MeasureRenderedView("Hello")
	if width != 5 {
		t.Errorf("expected width=5, got %d", width)
	}
	if height != 1 {
		t.Errorf("expected height=1, got %d", height)
	}
}

func TestCenterModal(t *testing.T) {
	// Build a 3-line, 10-char-wide view
	view := strings.Join([]string{
		"1234567890",
		"1234567890",
		"1234567890",
	}, "\n")

	width, height, row, col := CenterModal(view, 80, 24)
	if width != 10 {
		t.Errorf("expected width=10, got %d", width)
	}
	if height != 3 {
		t.Errorf("expected height=3, got %d", height)
	}

	// CalculateCenter(80,24,10,3) => row=(24-3)/2=10, col=(80-10)/2=35
	// CenterModal shifts up by 1 => row=9
	if row != 9 {
		t.Errorf("expected row=9 (centered minus shift-up-by-1), got %d", row)
	}
	if col != 35 {
		t.Errorf("expected col=35, got %d", col)
	}
}
