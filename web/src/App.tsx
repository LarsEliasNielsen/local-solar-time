import { useCallback, useEffect, useRef, useState } from 'react';
import SolarClock from './SolarClock';
import ClockDisplay from './ClockDisplay';
import SolarInfo from './SolarInfo';
import LocationControls from './LocationControls';
import { getLocation } from './geo';
import { openSocket } from './ws';
import type { SolarUpdate, LocationSource } from './types';

const DEFAULT_LAT = 55.6761;
const DEFAULT_LON = 12.5683;

const RECONNECT_INITIAL_MS = 1000;
const RECONNECT_MAX_MS = 30000;

interface Coords {
  lat: number;
  lon: number;
}

type AppState =
  | { phase: 'connecting'; coords: Coords; locationSource: LocationSource }
  | { phase: 'connected'; update: SolarUpdate; coords: Coords; locationSource: LocationSource }
  | { phase: 'reconnecting'; delay: number; coords: Coords; locationSource: LocationSource };

export default function App() {
  const [appState, setAppState] = useState<AppState>({
    phase: 'connecting',
    coords: { lat: DEFAULT_LAT, lon: DEFAULT_LON },
    locationSource: 'default',
  });

  const [latInput, setLatInput] = useState(String(DEFAULT_LAT));
  const [lonInput, setLonInput] = useState(String(DEFAULT_LON));
  const [inputError, setInputError] = useState('');
  const [geoButtonLabel, setGeoButtonLabel] = useState('Use my location');
  const [locationChanged, setLocationChanged] = useState(false);

  const lastUpdateRef = useRef<SolarUpdate | null>(null);
  const cleanupRef = useRef<(() => void) | null>(null);
  const reconnectTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const reconnectDelayRef = useRef(RECONNECT_INITIAL_MS);
  // Session counter: each connect() call gets a unique ID. Stale onClose/onMessage callbacks
  // from a socket that was intentionally closed check this and bail out.
  const sessionRef = useRef(0);

  const connect = useCallback((coords: Coords, source: LocationSource) => {
    // Cancel any pending reconnect and close the current socket before opening a new one.
    if (reconnectTimerRef.current) clearTimeout(reconnectTimerRef.current);
    if (cleanupRef.current) cleanupRef.current();

    const session = ++sessionRef.current;
    reconnectDelayRef.current = RECONNECT_INITIAL_MS;

    setAppState(prev => {
      if (prev.phase === 'connected') setLocationChanged(v => !v);
      return { phase: 'connecting', coords, locationSource: source };
    });
    setLatInput(coords.lat.toFixed(4));
    setLonInput(coords.lon.toFixed(4));

    const cleanup = openSocket(
      coords.lat,
      coords.lon,
      (update) => {
        if (session !== sessionRef.current) return;
        lastUpdateRef.current = update;
        setAppState({ phase: 'connected', update, coords, locationSource: source });
      },
      () => {
        // Ignore close events from sockets that were intentionally replaced.
        if (session !== sessionRef.current) return;

        const delay = reconnectDelayRef.current;
        reconnectDelayRef.current = Math.min(delay * 2, RECONNECT_MAX_MS);

        setAppState(prev => ({
          ...prev,
          phase: 'reconnecting',
          delay: Math.round(delay / 1000),
        }));

        reconnectTimerRef.current = setTimeout(() => {
          if (session !== sessionRef.current) return;
          if (source === 'gps') {
            getLocation().then(
              gc => connect({ lat: gc.latitude, lon: gc.longitude }, 'gps'),
              () => connect({ lat: DEFAULT_LAT, lon: DEFAULT_LON }, 'default'),
            );
          } else {
            connect(coords, source);
          }
        }, delay);
      },
    );
    cleanupRef.current = cleanup;
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  // Initial connect + background geolocation
  useEffect(() => {
    connect({ lat: DEFAULT_LAT, lon: DEFAULT_LON }, 'default');

    getLocation().then(
      gc => connect({ lat: gc.latitude, lon: gc.longitude }, 'gps'),
      () => { /* stay on default */ },
    );

    return () => {
      if (cleanupRef.current) cleanupRef.current();
      if (reconnectTimerRef.current) clearTimeout(reconnectTimerRef.current);
    };
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  function handleManualSubmit() {
    const lat = parseFloat(latInput);
    const lon = parseFloat(lonInput);
    if (isNaN(lat) || lat < -90 || lat > 90) {
      setInputError('Latitude must be between -90 and 90.');
      return;
    }
    if (isNaN(lon) || lon < -180 || lon > 180) {
      setInputError('Longitude must be between -180 and 180.');
      return;
    }
    setInputError('');
    connect({ lat, lon }, 'manual');
  }

  function handleGeoButton() {
    if (appState.locationSource === 'gps') return;
    getLocation().then(
      gc => connect({ lat: gc.latitude, lon: gc.longitude }, 'gps'),
      () => {
        setGeoButtonLabel('Location denied');
        setTimeout(() => setGeoButtonLabel('Use my location'), 2000);
      },
    );
  }

  const update      = appState.phase === 'connected' ? appState.update : lastUpdateRef.current;
  const solarTime   = update?.solar_time ?? null;
  const today       = update?.today ?? null;
  const altitudeDeg = Math.round((update?.altitude_deg ?? 0) * 10) / 10;
  const solarNoon   = update?.solar_noon ?? null;
  const polarCap    = update?.polar_cap;
  const isGpsActive = appState.locationSource === 'gps';

  return (
    <div style={{ maxWidth: '640px', margin: '0 auto', padding: '16px', fontFamily: 'system-ui, sans-serif', color: '#F9FAFB' }}>
      <ClockDisplay solarTime={solarTime} />
      <SolarClock solarTime={solarTime} today={today} altitudeDeg={altitudeDeg} locationChanged={locationChanged} />
      <SolarInfo today={today} solarNoon={solarNoon} altitudeDeg={altitudeDeg} polarCap={polarCap} />
      {appState.phase === 'reconnecting' && (
        <p role="status" style={{ fontSize: '0.8rem', color: '#FBBF24', margin: '4px 0' }}>
          Reconnecting in {appState.delay}s&hellip;
        </p>
      )}
      <LocationControls
        latInput={latInput} lonInput={lonInput} inputError={inputError}
        geoButtonLabel={geoButtonLabel} isGpsActive={isGpsActive} locationSource={appState.locationSource}
        onLatChange={setLatInput} onLonChange={setLonInput}
        onSubmit={handleManualSubmit} onGeoButton={handleGeoButton}
      />
    </div>
  );
}
