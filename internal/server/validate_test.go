package server

import "testing"

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     SubscribeRequest
		wantErr bool
	}{
		{"mid-latitude is valid", SubscribeRequest{Lat: 55.6761, Lon: 12.5683}, false},
		{"exact pole is valid", SubscribeRequest{Lat: 90, Lon: 0}, false},
		{"lat too high", SubscribeRequest{Lat: 200, Lon: 0}, true},
		{"lat too low", SubscribeRequest{Lat: -91, Lon: 0}, true},
		{"lon too high", SubscribeRequest{Lat: 0, Lon: 181}, true},
		{"lon too low", SubscribeRequest{Lat: 0, Lon: -200}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validate(%+v) error = %v, wantErr %v", tt.req, err, tt.wantErr)
			}
		})
	}
}
