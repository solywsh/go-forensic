package pathx

import "testing"

func TestFilterSubDir(t *testing.T) {
	type args struct {
		subDirs []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "TestFilterSubDir",
			args: args{
				subDirs: []string{
					"/data/user/0/",
					"/data/user/12/",
					"/data/user/",
					"/data/data/com.android.xxx/",
					"/data/data/com.android.xxx/database",
					"/data/data/com.android.xxx/file",
					"/data/data/com.android.xxx/cache",
				},
			},
			want: []string{
				"/data/user/",
				"/data/data/com.android.xxx/",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterSubDir(tt.args.subDirs); len(got) != len(tt.want) {
				t.Errorf("FilterSubDir() = %v, want %v", got, tt.want)
			}
		})
	}
}
