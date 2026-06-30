import type { SolarEvent, Today, PolarCap } from './types';

interface SolarInfoProps {
  today: Today | null;
  solarNoon: SolarEvent | null;
  altitudeDeg: number;
  polarCap?: PolarCap;
}

function formatNoonUtc(utc: string | null): string {
  if (!utc) return '--';
  const d = new Date(utc);
  return `${String(d.getUTCHours()).padStart(2, '0')}:${String(d.getUTCMinutes()).padStart(2, '0')} UTC`;
}

export default function SolarInfo({ today, solarNoon, altitudeDeg, polarCap }: SolarInfoProps) {
  return (
    <>
      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '8px', margin: '12px 0', fontSize: '0.875rem' }}>
        <div><span style={{ color: '#9CA3AF' }}>Sunrise (ST)</span><br />{today?.sunrise?.solar_time ?? (today === null ? '--' : 'No sunrise today')}</div>
        <div><span style={{ color: '#9CA3AF' }}>Sunset (ST)</span><br />{today?.sunset?.solar_time ?? (today === null ? '--' : 'No sunset today')}</div>
        <div><span style={{ color: '#9CA3AF' }}>Solar noon</span><br />{formatNoonUtc(solarNoon?.utc ?? null)}</div>
        <div><span style={{ color: '#9CA3AF' }}>Altitude</span><br />{`${altitudeDeg.toFixed(1)}°`}</div>
      </div>
      {polarCap && <p style={{ fontSize: '0.75rem', color: '#9CA3AF', margin: '4px 0 8px' }}>{polarCap.reason}</p>}
    </>
  );
}
