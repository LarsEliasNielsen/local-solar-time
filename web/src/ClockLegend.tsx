import type { ReactNode } from 'react';
import { THEME } from './theme';

const C = THEME.solarClock;

// One legend entry: a 12px marker plus its label.
function Item({ marker, label }: { marker: ReactNode; label: string }) {
  return (
    <span style={{ display: 'inline-flex', alignItems: 'center', gap: '6px' }}>
      <svg width="12" height="12" viewBox="0 0 12 12" aria-hidden="true">{marker}</svg>
      {label}
    </span>
  );
}

// Key for the solar clock colors and the sun marker.
export default function ClockLegend() {
  return (
    <div
      role="list"
      aria-label="Clock legend"
      style={{ display: 'flex', justifyContent: 'center', flexWrap: 'wrap', gap: '16px', margin: '10px 0 0', fontSize: '0.75rem', color: THEME.textMuted }}
    >
      <span role="listitem"><Item label="Day" marker={<rect width="12" height="12" rx="2" fill={C.day} />} /></span>
      <span role="listitem"><Item label="Night" marker={<rect width="12" height="12" rx="2" fill={C.night} />} /></span>
      <span role="listitem"><Item label="Sun" marker={<circle cx="6" cy="6" r="4.5" fill={C.day} stroke={C.needle} strokeWidth="2" />} /></span>
    </div>
  );
}
