package strgen

import (
	"testing"
)

func TestStringGenerator_generateHexValues(t *testing.T) {
	tests := []struct {
		name     string
		valueLen int
		wantErr  bool
	}{
		{"standard size", 5, false},
		{"16 bytes", 16, false},
		{"0 bytes", 0, false},
		{"large size", 1024, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(make(chan<- string))
			got, err := s.generateHexValues(tt.valueLen)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateHexValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			expectedLength := tt.valueLen * 2
			if len(got) != expectedLength {
				t.Errorf("GenerateHex() = %v, length = %d, want length = %d", got, len(got), expectedLength)
			}
		})
	}
}

func BenchmarkGenerateHexValues(b *testing.B) {
	s := New(make(chan<- string))
	sizes := []int{0, 5, 8, 16, 32, 64, 128, 256, 512}
	for _, size := range sizes {
		b.Run("size="+string(rune(size)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = s.generateHexValues(size)
			}
		})
	}
}
