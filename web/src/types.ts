export interface SolarEvent {
  solar_time: string;
  utc: string;
}

export interface Today {
  sunrise: SolarEvent | null;
  sunset: SolarEvent | null;
}

export interface PolarCap {
  reason: string;
}

export interface SolarUpdate {
  solar_time: string | null;
  equation_of_time_minutes: number;
  utc_offset_seconds: number | null;
  altitude_deg: number;
  azimuth_deg: number | null;
  solar_noon: SolarEvent | null;
  today: Today | null;
  previous_sunrise: SolarEvent | null;
  next_sunrise: SolarEvent | null;
  previous_sunset: SolarEvent | null;
  next_sunset: SolarEvent | null;
  polar_cap?: PolarCap;
}
