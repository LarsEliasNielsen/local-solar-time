import type { KeyboardEvent } from 'react';
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
    <div style={{ marginTop: '12px' }}>
      <div style={{ display: 'flex', gap: '8px', alignItems: 'center', flexWrap: 'wrap' }}>
        <label style={{ fontSize: '0.8rem', color: '#9CA3AF' }}>
          Lat
          <input type="number" min={-90} max={90} step={0.0001}
            value={latInput} onChange={e => onLatChange(e.target.value)} onKeyDown={handleKeyDown}
            style={inputStyle} />
        </label>
        <label style={{ fontSize: '0.8rem', color: '#9CA3AF' }}>
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
        <p style={{ fontSize: '0.75rem', color: '#6B7280', margin: '4px 0 0' }}>Default location (Copenhagen)</p>
      )}
      {inputError && (
        <p style={{ fontSize: '0.75rem', color: '#F87171', margin: '4px 0 0' }}>{inputError}</p>
      )}
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

const baseButtonStyle: React.CSSProperties = {
  borderRadius: '4px',
  color: '#F9FAFB',
  padding: '4px 12px',
  cursor: 'pointer',
  alignSelf: 'flex-end',
  marginBottom: '1px',
};

const applyButtonStyle:  React.CSSProperties = { ...baseButtonStyle, background: '#374151', border: '1px solid #4B5563' };
const activeButtonStyle: React.CSSProperties = { ...baseButtonStyle, background: '#1D4ED8', border: '1px solid #2563EB' };
const mutedButtonStyle:  React.CSSProperties = { ...baseButtonStyle, background: '#1F2937', border: '1px solid #374151', color: '#6B7280', cursor: 'default' };
