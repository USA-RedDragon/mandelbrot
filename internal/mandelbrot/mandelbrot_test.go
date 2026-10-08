package mandelbrot

import "testing"

func TestScaleByCap(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		start, factor, want float64
	}{
		{start: 4, factor: 1.1, want: 4},
		{start: 4, factor: 0.5, want: 2},
		{start: 0.5, factor: 4, want: 1},
		{start: 1, factor: 0.5, want: 0.5},
	} {
		m := NewMandelbrot(4, 4, Settings{MaxIterations: 10, Scale: tc.start, Exponent: 2})
		m.ScaleBy(tc.factor)
		if m.scale != tc.want {
			t.Errorf("start %v, ScaleBy(%v): scale %v, want %v", tc.start, tc.factor, m.scale, tc.want)
		}
	}
}
