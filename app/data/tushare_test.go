package data

import (
	"log"
	"testing"

	"github.com/99MyCql/duffett/pkg"
)

func TestGetDailyData(t *testing.T) {
	pkg.InitConfig("../../conf.yaml")
	type args struct {
		tsCode    string
		tradeDate string
	}
	tests := []struct {
		name    string
		args    args
		want    *DailyData
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				tsCode:    "000001.SZ",
				tradeDate: "20201228",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDailyData(tt.args.tsCode, tt.args.tradeDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDailyData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			log.Print(got)
		})
	}
}
