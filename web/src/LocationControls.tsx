import { useState } from 'react';
import type { KeyboardEvent } from 'react';
import { THEME } from './theme';
import type { LocationSource } from './types';

interface LocationControlsProps {
  latInput: string;
  lonInput: string;
  inputError: string;
  geoButtonLabel: string;
  isGpsActive: boolean;
  locationSource: LocationSource;
  onLatChange: (v: string) => void;
  onLonChange: (v: string) => void;
  onSubmit: () => void;
  onGeoButton: () => void;
}

const LAT_BOUNDS = { min: -90, max: 90 };
const LON_BOUNDS = { min: -180, max: 180 };
const STEP = 0.0001;

function stepValue(current: string, delta: 1 | -1, bounds: { min: number; max: number }): string {
  const parsed = parseFloat(current);
  const base = isNaN(parsed) ? 0 : parsed;
  const next = Math.min(bounds.max, Math.max(bounds.min, base + delta * STEP));
  return next.toFixed(4);
}

// Crosshair/scope icon: filled center dot when GPS tracking is active, hollow otherwise.
function GpsIcon({ color, active, size }: { color: string; active: boolean; size: number }) {
  return (
    <svg viewBox="0 0 24 24" width={size} height={size} fill="none" stroke={color} strokeWidth={2} strokeLinecap="round">
      <circle cx="12" cy="12" r="7" />
      <line x1="12" y1="1" x2="12" y2="4" />
      <line x1="12" y1="20" x2="12" y2="23" />
      <line x1="1" y1="12" x2="4" y2="12" />
      <line x1="20" y1="12" x2="23" y2="12" />
      <circle cx="12" cy="12" r="2" fill={active ? color : 'none'} stroke="none" />
    </svg>
  );
}

function ChevronIcon({ direction }: { direction: 'up' | 'down' }) {
  return (
    <svg viewBox="0 0 10 6" width={8} height={5} fill="none" stroke="currentColor" strokeWidth={1.5} strokeLinecap="round" strokeLinejoin="round">
      <polyline points={direction === 'up' ? '1,5 5,1 9,5' : '1,1 5,5 9,1'} />
    </svg>
  );
}

interface NumberFieldProps {
  id: string;
  label: string;
  value: string;
  bounds: { min: number; max: number };
  focused: boolean;
  onChange: (v: string) => void;
  onKeyDown: (e: KeyboardEvent) => void;
  onFocus: () => void;
  onBlur: () => void;
}

function NumberField({ id, label, value, bounds, focused, onChange, onKeyDown, onFocus, onBlur }: NumberFieldProps) {
  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
      <label htmlFor={id} style={labelStyle}>{label}</label>
      <div style={{ position: 'relative' }}>
        <input id={id} type="number" min={bounds.min} max={bounds.max} step={STEP}
          value={value} onChange={e => onChange(e.target.value)} onKeyDown={onKeyDown}
          onFocus={onFocus} onBlur={onBlur}
          style={{ ...inputStyle, ...(focused ? focusRingStyle : {}) }} />
        <div style={stepperWrapperStyle}>
          <button type="button" tabIndex={-1} aria-label={`Increase ${label}`}
            onClick={() => onChange(stepValue(value, 1, bounds))} style={stepperButtonStyle}>
            <ChevronIcon direction="up" />
          </button>
          <button type="button" tabIndex={-1} aria-label={`Decrease ${label}`}
            onClick={() => onChange(stepValue(value, -1, bounds))} style={{ ...stepperButtonStyle, borderTop: `1px solid ${THEME.border}` }}>
            <ChevronIcon direction="down" />
          </button>
        </div>
      </div>
    </div>
  );
}

export default function LocationControls({
  latInput, lonInput, inputError,
  geoButtonLabel, isGpsActive, locationSource,
  onLatChange, onLonChange, onSubmit, onGeoButton,
}: LocationControlsProps) {
  const [focusedField, setFocusedField] = useState<'lat' | 'lon' | 'apply' | null>(null);
  const isGeoError = geoButtonLabel === 'Location denied';

  function handleKeyDown(e: KeyboardEvent) {
    if (e.key === 'Enter') onSubmit();
  }

  return (
    <div style={{ maxWidth: '560px', margin: '12px auto 0' }}>
      {/* Hide native number-input spinners; replaced by the custom stepper buttons above. */}
      <style>{`
        input[type=number]::-webkit-outer-spin-button,
        input[type=number]::-webkit-inner-spin-button { -webkit-appearance: none; margin: 0; }
        input[type=number] { -moz-appearance: textfield; }
      `}</style>
      <div style={{ display: 'flex', gap: '10px', alignItems: 'flex-end', justifyContent: 'center', flexWrap: 'wrap', padding: '0 5%' }}>
        <NumberField id="lat-input" label="Lat" value={latInput} bounds={LAT_BOUNDS}
          focused={focusedField === 'lat'} onChange={onLatChange} onKeyDown={handleKeyDown}
          onFocus={() => setFocusedField('lat')} onBlur={() => setFocusedField(null)} />
        <NumberField id="lon-input" label="Lon" value={lonInput} bounds={LON_BOUNDS}
          focused={focusedField === 'lon'} onChange={onLonChange} onKeyDown={handleKeyDown}
          onFocus={() => setFocusedField('lon')} onBlur={() => setFocusedField(null)} />
        <button onClick={onSubmit}
          onFocus={() => setFocusedField('apply')} onBlur={() => setFocusedField(null)}
          style={{ ...applyButtonStyle, ...(focusedField === 'apply' ? focusRingStyle : {}) }}>
          Apply
        </button>
        <button onClick={onGeoButton} disabled={isGpsActive}
          aria-label={isGpsActive ? 'GPS active' : geoButtonLabel}
          title={isGpsActive ? 'GPS active' : geoButtonLabel}
          style={gpsButtonStyle(isGpsActive, isGeoError)}>
          <GpsIcon color={gpsIconColor(isGpsActive, isGeoError)} active={isGpsActive} size={16} />
        </button>
      </div>
      {locationSource === 'default' && !inputError && (
        <p style={{ fontSize: '0.75rem', color: THEME.textDim, margin: '4px 0 0', padding: '0 5%' }}>Default location (Copenhagen)</p>
      )}
      {inputError && (
        <p style={{ fontSize: '0.75rem', color: THEME.textError, margin: '4px 0 0', padding: '0 5%' }}>{inputError}</p>
      )}
    </div>
  );
}

const labelStyle: React.CSSProperties = {
  fontSize: '0.75rem',
  fontWeight: 500,
  color: THEME.textPrimary,
};

const focusRingStyle: React.CSSProperties = {
  boxShadow: `0 0 0 3px ${THEME.accentBorder}66`,
};

const inputStyle: React.CSSProperties = {
  display: 'block',
  background: THEME.surface,
  border: `1.5px solid ${THEME.accentBorder}`,
  borderRadius: '6px',
  color: THEME.textPrimary,
  fontSize: '0.9rem',
  padding: '6px 24px 6px 10px',
  width: '110px',
};

const stepperWrapperStyle: React.CSSProperties = {
  position: 'absolute',
  top: '1.5px',
  right: '1.5px',
  bottom: '1.5px',
  width: '20px',
  display: 'flex',
  flexDirection: 'column',
  borderLeft: `1px solid ${THEME.border}`,
  borderTopRightRadius: '4px',
  borderBottomRightRadius: '4px',
  overflow: 'hidden',
};

const stepperButtonStyle: React.CSSProperties = {
  flex: 1,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  padding: 0,
  border: 'none',
  background: 'transparent',
  color: THEME.textMuted,
  cursor: 'pointer',
};

const baseButtonStyle: React.CSSProperties = {
  borderRadius: '6px',
  color: THEME.textPrimary,
  fontWeight: 500,
  cursor: 'pointer',
};

const applyButtonStyle: React.CSSProperties = {
  ...baseButtonStyle,
  background: THEME.surface,
  border: `1.5px solid ${THEME.accentBorder}`,
  fontSize: '0.9rem',
  padding: '6px 16px',
};

function gpsButtonStyle(active: boolean, error: boolean): React.CSSProperties {
  return {
    ...baseButtonStyle,
    width: '32px',
    height: '32px',
    padding: 0,
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    borderRadius: '50%',
    background: active ? `${THEME.accent}33` : THEME.surface,
    border: `1.5px solid ${error ? THEME.textError : active ? THEME.accent : THEME.border}`,
    cursor: active ? 'default' : 'pointer',
  };
}

function gpsIconColor(active: boolean, error: boolean): string {
  if (error) return THEME.textError;
  if (active) return THEME.accent;
  return THEME.textMuted;
}
