interface ClockDisplayProps {
  solarTime: string | null;
}

export default function ClockDisplay({ solarTime }: ClockDisplayProps) {
  return (
    <p
      role="timer"
      aria-live="off"
      aria-atomic="true"
      style={{ textAlign: 'center', fontSize: '3rem', fontVariantNumeric: 'tabular-nums', margin: '0 0 8px', letterSpacing: '0.05em' }}
    >
      {solarTime ?? '--:--:--'}
    </p>
  );
}
