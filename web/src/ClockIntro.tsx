import { THEME } from './theme';

// First-visit explanation of what apparent solar time is, followed by a break before the solar info.
// Width and padding mirror SolarInfo so the text lines up with the grid below.
export default function ClockIntro() {
  return (
    <div style={{ maxWidth: '560px', margin: '20px auto 0' }}>
      <div style={{ padding: '0 5%' }}>
        <section aria-labelledby="clock-intro-title" style={{ fontSize: '0.875rem', lineHeight: 1.5 }}>
          <h2 id="clock-intro-title" style={{ margin: '0 0 6px', fontSize: '0.9375rem', fontWeight: 600 }}>
            A clock that follows the sun
          </h2>
          <p style={{ margin: 0, color: THEME.textMuted }}>
            This is local apparent solar time: the time a sundial would show where you are. The sun is highest in your
            sky at 12:00 solar time. Clock time is shared across a whole time zone and often shifted for daylight
            saving, so it can differ from solar time by an hour or more.
          </p>
        </section>
        <hr style={{ border: 'none', borderTop: `1px solid ${THEME.border}`, margin: '16px 0 4px' }} />
      </div>
    </div>
  );
}
