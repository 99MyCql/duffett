package data

import (
	"testing"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/99MyCql/duffett/pkg"
)

func TestGetDailyData(t *testing.T) {
	pkg.InitConfig("../../conf.yaml")
	type args struct {
		tsCode    string
		tradeDate time.Time
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
				tradeDate: time.Now(),
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
			log.Debug(got)
		})
	}
}
