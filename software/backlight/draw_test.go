package backlight

import "testing"

func BenchmarkDrawScreen(b *testing.B) {
	worker := Worker{
		Opt: &Options{Width: 4, Height: 4, IsDebug: false},
	}

	rs, err := worker.DrawScreen()
	if err != nil {
		b.Errorf("Benchmark completed with error %s", err)
	}

	if len(rs) != 16 {
		b.Errorf("Unexpected result, expected %d, but was %d", 16, len(rs))
	}
}
