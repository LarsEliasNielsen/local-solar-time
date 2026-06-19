package server

import "fmt"

// validate checks that req contains geometrically valid coordinates.
func validate(req SubscribeRequest) error {
	if req.Lat < -90 || req.Lat > 90 {
		return fmt.Errorf("lat must be between -90 and 90, got %v", req.Lat)
	}
	if req.Lon < -180 || req.Lon > 180 {
		return fmt.Errorf("lon must be between -180 and 180, got %v", req.Lon)
	}
	return nil
}
