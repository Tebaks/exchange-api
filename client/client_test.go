package client

import (
	"testing"
)

func Test_exchangeRateClient_GetExchangeRate(t *testing.T) {
	client := NewExchangeRateClient()
	type args struct {
		name string
		url  string
	}
	tests := []struct {
		name    string
		c       *exchangeRateClient
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "Test mock url 1",
			args: args{
				name: "USD",
				url:  "https://run.mocky.io/v3/f05a9177-5719-47b3-94d9-60ab75e538a7",
			},
			want:    9.96,
			wantErr: false,
		},
		{
			name: "Test mock url 2",
			args: args{
				name: "EUR",
				url:  "https://run.mocky.io/v3/7dcc2ac7-85f0-4032-a421-f8d947e20824",
			},
			want:    12.4230,
			wantErr: false,
		},
		{
			name: "Test mock url 3",
			args: args{
				name: "GBP",
				url:  "https://run.mocky.io/v3/377b087e-56bf-4a3f-a2f8-d662b4782705",
			},
			want:    12.0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetExchangeRate(tt.args.name, tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("exchangeRateClient.GetExchangeRate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("exchangeRateClient.GetExchangeRate() = %v, want %v", got, tt.want)
			}
		})
	}
}
