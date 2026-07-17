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

export default function LocationControls({
  latInput, lonInput, inputError,
  geoButtonLabel, isGpsActive, locationSource,
  onLatChange, onLonChange, onSubmit, onGeoButton,
}: LocationControlsProps) {
  function handleKeyDown(e: KeyboardEvent) {
    if (e.key === 'Enter') onSubmit();
  }

  return (
    <div style={{ maxWidth: '560px', margin: '12px auto 0' }}>
      <div style={{ display: 'flex', gap: '8px', alignItems: 'center', flexWrap: 'wrap', padding: '0 5%' }}>
        <label style={{ fontSize: '0.8rem', color: THEME.textMuted }}>
          Lat
          <input type="number" min={-90} max={90} step={0.0001}
            value={latInput} onChange={e => onLatChange(e.target.value)} onKeyDown={handleKeyDown}
            style={inputStyle} />
        </label>
        <label style={{ fontSize: '0.8rem', color: THEME.textMuted }}>
          Lon
          <input type="number" min={-180} max={180} step={0.0001}
            value={lonInput} onChange={e => onLonChange(e.target.value)} onKeyDown={handleKeyDown}
            style={inputStyle} />
        </label>
        <button onClick={onSubmit} style={applyButtonStyle}>Apply</button>
        <button onClick={onGeoButton} disabled={isGpsActive}
          style={isGpsActive ? mutedButtonStyle : activeButtonStyle}>
          {isGpsActive ? 'GPS active' : geoButtonLabel}
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

const inputStyle: React.CSSProperties = {
  display: 'block',
  background: THEME.surface,
  border: `1px solid ${THEME.border}`,
  borderRadius: '4px',
  color: THEME.textPrimary,
  padding: '4px 8px',
  width: '100px',
  marginTop: '2px',
};

const baseButtonStyle: React.CSSProperties = {
  borderRadius: '4px',
  color: THEME.textPrimary,
  padding: '4px 12px',
  cursor: 'pointer',
  alignSelf: 'flex-end',
  marginBottom: '1px',
};

const applyButtonStyle:  React.CSSProperties = { ...baseButtonStyle, background: THEME.surfaceRaised, border: `1px solid ${THEME.borderStrong}` };
const activeButtonStyle: React.CSSProperties = { ...baseButtonStyle, background: THEME.accent, border: `1px solid ${THEME.accentBorder}` };
const mutedButtonStyle:  React.CSSProperties = { ...baseButtonStyle, background: THEME.surface, border: `1px solid ${THEME.border}`, color: THEME.textDim, cursor: 'default' };
