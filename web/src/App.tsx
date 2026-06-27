import { useCallback, useEffect, useRef, useState } from 'react';
import SolarClock from './SolarClock';
import { getLocation } from './geo';
import { openSocket } from './ws';
import type { SolarUpdate } from './types';

const DEFAULT_LAT = 55.6761;
const DEFAULT_LON = 12.5683;

const RECONNECT_INITIAL_MS = 1000;
const RECONNECT_MAX_MS = 30000;

type LocationSource = 'default' | 'gps' | 'manual';

interface Coords {
  lat: number;
  lon: number;
}

type AppState =
  | { phase: 'connecting'; coords: Coords; locationSource: LocationSource }
  | { phase: 'connected'; update: SolarUpdate; coords: Coords; locationSource: LocationSource }
  | { phase: 'reconnecting'; delay: number; coords: Coords; locationSource: LocationSource };

function formatNoonUtc(utc: string | null): string {
  if (!utc) return '--';
  const d = new Date(utc);
  return `${String(d.getUTCHours()).padStart(2, '0')}:${String(d.getUTCMinutes()).padStart(2, '0')} UTC`;
}

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

  function handleKeyDown(e: React.KeyboardEvent) {
    if (e.key === 'Enter') handleManualSubmit();
  }

  const update = appState.phase === 'connected' ? appState.update : lastUpdateRef.current;
  const locationSource = appState.locationSource;
  const isGpsActive = locationSource === 'gps';

  return (
    <div style={{ maxWidth: '640px', margin: '0 auto', padding: '16px', fontFamily: 'system-ui, sans-serif', color: '#F9FAFB' }}>
      <p style={{ textAlign: 'center', fontSize: '3rem', fontVariantNumeric: 'tabular-nums', margin: '0 0 8px', letterSpacing: '0.05em' }}>
        {update?.solar_time ?? '--:--:--'}
      </p>

      <SolarClock update={update} locationChanged={locationChanged} />

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '8px', margin: '12px 0', fontSize: '0.875rem' }}>
        <div>
          <span style={{ color: '#9CA3AF' }}>Sunrise (ST)</span><br />
          {update?.today?.sunrise?.solar_time ?? (update?.today === null ? '--' : 'No sunrise today')}
        </div>
        <div>
          <span style={{ color: '#9CA3AF' }}>Sunset (ST)</span><br />
          {update?.today?.sunset?.solar_time ?? (update?.today === null ? '--' : 'No sunset today')}
        </div>
        <div>
          <span style={{ color: '#9CA3AF' }}>Solar noon</span><br />
          {update ? formatNoonUtc(update.solar_noon?.utc ?? null) : '--'}
        </div>
        <div>
          <span style={{ color: '#9CA3AF' }}>Altitude</span><br />
          {update ? `${update.altitude_deg.toFixed(1)}°` : '--'}
        </div>
      </div>

      {update?.polar_cap && (
        <p style={{ fontSize: '0.75rem', color: '#9CA3AF', margin: '4px 0 8px' }}>{update.polar_cap.reason}</p>
      )}

      {appState.phase === 'reconnecting' && (
        <p style={{ fontSize: '0.8rem', color: '#FBBF24', margin: '4px 0' }}>
          Reconnecting in {appState.delay}s&hellip;
        </p>
      )}

      <div style={{ marginTop: '12px' }}>
        <div style={{ display: 'flex', gap: '8px', alignItems: 'center', flexWrap: 'wrap' }}>
          <label style={{ fontSize: '0.8rem', color: '#9CA3AF' }}>
            Lat
            <input
              type="number" min={-90} max={90} step={0.0001}
              value={latInput}
              onChange={e => setLatInput(e.target.value)}
              onKeyDown={handleKeyDown}
              style={inputStyle}
            />
          </label>
          <label style={{ fontSize: '0.8rem', color: '#9CA3AF' }}>
            Lon
            <input
              type="number" min={-180} max={180} step={0.0001}
              value={lonInput}
              onChange={e => setLonInput(e.target.value)}
              onKeyDown={handleKeyDown}
              style={inputStyle}
            />
          </label>
          <button onClick={handleManualSubmit} style={applyButtonStyle}>Apply</button>
          <button
            onClick={handleGeoButton}
            disabled={isGpsActive}
            style={isGpsActive ? mutedButtonStyle : activeButtonStyle}
          >
            {isGpsActive ? 'GPS active' : geoButtonLabel}
          </button>
        </div>
        {locationSource === 'default' && !inputError && (
          <p style={{ fontSize: '0.75rem', color: '#6B7280', margin: '4px 0 0' }}>Default location (Copenhagen)</p>
        )}
        {inputError && (
          <p style={{ fontSize: '0.75rem', color: '#F87171', margin: '4px 0 0' }}>{inputError}</p>
        )}
      </div>
    </div>
  );
}

const inputStyle: React.CSSProperties = {
  display: 'block',
  background: '#1F2937',
  border: '1px solid #374151',
  borderRadius: '4px',
  color: '#F9FAFB',
  padding: '4px 8px',
  width: '100px',
  marginTop: '2px',
};

const applyButtonStyle: React.CSSProperties = {
  background: '#374151',
  border: '1px solid #4B5563',
  borderRadius: '4px',
  color: '#F9FAFB',
  padding: '4px 12px',
  cursor: 'pointer',
  alignSelf: 'flex-end',
  marginBottom: '1px',
};

const activeButtonStyle: React.CSSProperties = {
  background: '#1D4ED8',
  border: '1px solid #2563EB',
  borderRadius: '4px',
  color: '#F9FAFB',
  padding: '4px 12px',
  cursor: 'pointer',
  alignSelf: 'flex-end',
  marginBottom: '1px',
};

const mutedButtonStyle: React.CSSProperties = {
  background: '#1F2937',
  border: '1px solid #374151',
  borderRadius: '4px',
  color: '#6B7280',
  padding: '4px 12px',
  cursor: 'default',
  alignSelf: 'flex-end',
  marginBottom: '1px',
};
