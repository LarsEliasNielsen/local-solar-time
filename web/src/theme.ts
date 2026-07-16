// Design tokens for the frontend.

const SOLAR_CLOCK_NEEDLE_STROKE_WIDTH = 5;

export const THEME = {
  // --- shared, cross-component tokens ---
  textPrimary: '#F9FAFB',
  textMuted: '#9CA3AF',
  textDim: '#6B7280',
  textError: '#F87171',
  textWarning: '#FBBF24',
  surface: '#1F2937',
  surfaceRaised: '#374151',
  border: '#374151',
  borderStrong: '#4B5563',
  accent: '#1D4ED8',
  accentBorder: '#2563EB',
  pageBackground: '#111827',

  // --- component-local tokens ---
  solarClock: {
    night: '#4e5355',
    day: '#ffb703',
    needle: 'white',
    pivot: 'white',
    needleStrokeWidth: SOLAR_CLOCK_NEEDLE_STROKE_WIDTH,
    needleTipRadius: 15,
    pivotRadius: SOLAR_CLOCK_NEEDLE_STROKE_WIDTH / 2,
  },
} as const;
