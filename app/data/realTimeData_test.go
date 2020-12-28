package data

import (
	"testing"
)

func TestGetRealTimeData(t *testing.T) {
	type args struct {
		tsCode string
	}
	tests := []struct {
		name    string
		args    args
		want    *RealTimeData
		wantErr bool
	}{
		{
			name:    "1",
			args:    args{tsCode: "000001.SZ"},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRealTimeData(tt.args.tsCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRealTimeData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got)
		})
	}
}
