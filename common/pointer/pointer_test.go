package pointer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// useValue as a global prevents most compiler optimizations on benchmarks
var useValue any

func BenchmarkInt_helper_function(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		ptr := Int(1)
		useValue = ptr
	}
}

func BenchmarkInt_declared(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		value := 1
		useValue = &value
	}
}

func BenchmarkString_helper_function(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		ptr := String("words")
		useValue = ptr
	}
}

func BenchmarkString_declared(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		value := "words"
		useValue = &value
	}
}

func Test_TimeOrNil(t *testing.T) {
	now := time.Now()

	ptrNow := TimeOrNil(now)
	require.NotNil(t, ptrNow)
	assert.Equal(t, now, *ptrNow)

	ptrNil := TimeOrNil(time.Time{})
	assert.Nil(t, ptrNil)
}

func Test_StringOrNil(t *testing.T) {
	a := "a"
	s := StringOrNil(a)
	assert.NotNil(t, s)
	assert.Equal(t, &a, s)

	assert.Nil(t, StringOrNil(""))
}

func Benchmark_Create(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		v := To("test")
		useValue = v
	}
}

func Benchmark_Legacy(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		v := String("test")
		useValue = v
	}
}

func TestFrom(t *testing.T) {
	type args struct {
		input interface{} // used/forwarded with pointer
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "int",
			args: args{
				input: 1,
			},
			want: 1,
		},
		{
			name: "int64",
			args: args{
				input: int64(1),
			},
			want: int64(1),
		},
		{
			name: "string",
			args: args{
				input: "test",
			},
			want: "test",
		},
		{
			name: "bool",
			args: args{
				input: true,
			},
			want: true,
		},
		{
			name: "struct",
			args: args{
				input: struct{}{},
			},
			want: struct{}{},
		},
		{
			name: "struct ptr++",
			args: args{
				input: &struct{}{},
			},
			want: &struct{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, From(&tt.args.input), "From(%v)", tt.args.input)
		})
	}
}

func TestFrom_passing_null_int(t *testing.T) {
	tests := []struct {
		name string
		want interface{}
	}{
		{
			name: "int",
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, From[int](nil), "From(int)")
		})
	}
}

func TestFrom_passing_null_string(t *testing.T) {
	tests := []struct {
		name string
		want interface{}
	}{
		{
			name: "string",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, From[string](nil), "From(string)")
		})
	}
}

func TestFrom_passing_null_bool(t *testing.T) {
	tests := []struct {
		name string
		want interface{}
	}{
		{
			name: "bool",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, From[bool](nil), "From(bool)")
		})
	}
}

func TestFrom_passing_null_struct(t *testing.T) {
	tests := []struct {
		name string
		want interface{}
	}{
		{
			name: "string",
			want: struct{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, From[struct{}](nil), "From(struct{})")
		})
	}
}

func TestCastTo(t *testing.T) {
	t.Run("it returns nil given nil", func(t *testing.T) {
		expectedType := (*uint32)(nil)

		actual := Cast[uint, uint32]((*uint)(nil))

		require.Nil(t, actual)
		require.IsType(t, expectedType, actual)
	})

	t.Run("it casts pointer to A to pointer to B", func(t *testing.T) {
		v := uint(123)
		expectedType := uint32(123)

		actual := Cast[uint, uint32](&v)

		require.NotNil(t, actual)
		require.IsType(t, &expectedType, actual)
		require.Equal(t, expectedType, *actual)
	})
}
