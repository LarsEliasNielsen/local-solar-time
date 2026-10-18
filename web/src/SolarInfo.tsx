import InfoTip from './InfoTip';
import { THEME } from './theme';
import type { SolarEvent, Today, PolarCap } from './types';

interface SolarInfoProps {
  today: Today | null;
  solarNoon: SolarEvent | null;
  altitudeDeg: number;
  polarCap?: PolarCap;
}

const HINTS = {
  sunrise: 'When the sun peeks over the horizon. Shown in solar time: what a sundial at your spot would read, not your phone clock.',
  sunset: 'When the sun drops below the horizon. Shown in solar time: what a sundial at your spot would read, not your phone clock.',
  solarNoon: 'The moment the sun is highest in the sky today. In solar time that is always 12:00, so it is shown in UTC, the world reference clock.',
  altitude: 'How high the sun is above the horizon right now. 0° is on the horizon, 90° is straight overhead, and below 0° the sun has set.',
} as const;

// Muted item label with its info badge; the badge sits on the outer side so the text lines up with the value.
function Label({ text, hint, align = 'left' }: { text: string; hint: string; align?: 'left' | 'right' }) {
  return (
    <span style={{ display: 'inline-flex', verticalAlign: 'top', alignItems: 'center', gap: '6px', flexDirection: align === 'right' ? 'row-reverse' : 'row' }}>
      <span style={{ color: THEME.textMuted }}>{text}</span>
      <InfoTip label={text} text={hint} align={align} />
    </span>
  );
}

function formatNoonUtc(utc: string | null): string {
  if (!utc) return '--';
  const d = new Date(utc);
  return `${String(d.getUTCHours()).padStart(2, '0')}:${String(d.getUTCMinutes()).padStart(2, '0')} UTC`;
}

export default function SolarInfo({ today, solarNoon, altitudeDeg, polarCap }: SolarInfoProps) {
  return (
    <div style={{ maxWidth: '560px', margin: '0 auto' }}>
      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '8px', margin: '12px 0', padding: '0 5%', fontSize: '0.875rem' }}>
        <div><Label text="Sunrise (ST)" hint={HINTS.sunrise} /><br />{today?.sunrise?.solar_time ?? (today === null ? '--' : 'No sunrise today')}</div>
        <div style={{ textAlign: 'right' }}><Label text="Sunset (ST)" hint={HINTS.sunset} align="right" /><br />{today?.sunset?.solar_time ?? (today === null ? '--' : 'No sunset today')}</div>
        <div><Label text="Solar noon" hint={HINTS.solarNoon} /><br />{formatNoonUtc(solarNoon?.utc ?? null)}</div>
        <div style={{ textAlign: 'right' }}><Label text="Altitude" hint={HINTS.altitude} align="right" /><br />{`${altitudeDeg.toFixed(1)}°`}</div>
      </div>
      {polarCap && <p style={{ fontSize: '0.75rem', color: THEME.textMuted, margin: '4px 0 8px', padding: '0 5%' }}>{polarCap.reason}</p>}
    </div>
  );
}
