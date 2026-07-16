import { useEffect, useMemo, useRef, useState } from 'react';
import { THEME } from './theme';
import type { Today } from './types';

interface SolarClockProps {
  solarTime: string | null;
  today: Today | null;
  altitudeDeg: number;
  locationChanged: boolean;
}

// Semicircle gauge geometry (half-circle).
const CX = 200;
const CY = 200;
const R = 180;

function timeToSeconds(hms: string): number {
  const [h, m, s] = hms.split(':').map(Number);
  return (h ?? 0) * 3600 + (m ?? 0) * 60 + (s ?? 0);
}

// Angle convention:
// - 0s (midnight) -> pi (left horizon)
// - 43200s (noon) -> pi/2 (top)
// - 86400s (next midnight) -> 0 (right horizon)
// Angle decreases as the day progresses.
function solarAngle(secs: number): number {
  return (1 - secs / 86400) * Math.PI;
}

// Polar angle -> SVG canvas point on the arc. Y is flipped since SVG y grows downward.
function arcPoint(theta: number): [number, number] {
  return [CX + R * Math.cos(theta), CY - R * Math.sin(theta)];
}

// Filled pie slice from center to arc segment.
function wedgePath(startTheta: number, endTheta: number): string {
  if (Math.abs(startTheta - endTheta) < 0.001) return '';
  const [x1, y1] = arcPoint(startTheta);
  const [x2, y2] = arcPoint(endTheta);
  // SVG large-arc-flag: needed when the wedge spans more than half the circle.
  const large = startTheta - endTheta > Math.PI ? 1 : 0;
  return `M ${CX} ${CY} L ${x1} ${y1} A ${R} ${R} 0 ${large} 1 ${x2} ${y2} Z`;
}

// Used for the initial sweep-in animation (fast start, slow finish).
function easeOut(t: number): number {
  return 1 - (1 - t) * (1 - t);
}

// Used for the location-change tween (slow start, fast middle, slow finish).
function easeInOut(t: number): number {
  return t < 0.5 ? 2 * t * t : -1 + (4 - 2 * t) * t;
}

interface VisualAngles {
  needle: number | null;
  sunrise: number | null;
  sunset: number | null;
}

function deriveTargetAngles(solarTime: string | null, today: Today | null): VisualAngles {
  if (!solarTime && !today) return { needle: null, sunrise: null, sunset: null };
  const needle = solarTime ? solarAngle(timeToSeconds(solarTime)) : null;
  const sunrise = today?.sunrise ? solarAngle(timeToSeconds(today.sunrise.solar_time)) : null;
  const sunset  = today?.sunset  ? solarAngle(timeToSeconds(today.sunset.solar_time))  : null;
  return { needle, sunrise, sunset };
}

function lerp(a: number | null, b: number | null, t: number): number | null {
  if (a === null || b === null) return b;
  return a + (b - a) * t;
}

export default function SolarClock({ solarTime, today, altitudeDeg, locationChanged }: SolarClockProps) {
  const [visual, setVisual] = useState<VisualAngles>({ needle: null, sunrise: null, sunset: null });

  const hasAnimated = useRef(false);
  const animFrameRef = useRef<number | null>(null);
  const tweenPendingRef = useRef(false);
  const tweenFromRef = useRef<VisualAngles>({ needle: null, sunrise: null, sunset: null });

  const target = deriveTargetAngles(solarTime, today);
  const hasFirstUpdate = target.needle !== null;

  const visualRef = useRef<VisualAngles>({ needle: null, sunrise: null, sunset: null });
  visualRef.current = visual;

  const targetRef = useRef<VisualAngles>({ needle: null, sunrise: null, sunset: null });
  targetRef.current = target;

  // Animation 1/3: Initial sweep-in.
  // Needle sweeps from midnight to current time.
  useEffect(() => {
    if (hasAnimated.current || !hasFirstUpdate) return;
    hasAnimated.current = true;

    const to = { ...targetRef.current };
    const start = performance.now();
    const duration = 800;

    function frame(now: number) {
      const t = Math.min((now - start) / duration, 1);
      const e = easeOut(t);

      setVisual({
        needle: Math.PI - (Math.PI - (to.needle ?? Math.PI)) * e,
        sunrise: to.sunrise !== null ? Math.PI - (Math.PI - to.sunrise) * e : null,
        sunset: to.sunset !== null ? to.sunset * e : null,
      });

      if (t < 1) {
        animFrameRef.current = requestAnimationFrame(frame);
      } else {
        animFrameRef.current = null;
      }
    }

    animFrameRef.current = requestAnimationFrame(frame);
    return () => {
      if (animFrameRef.current !== null) cancelAnimationFrame(animFrameRef.current);
    };
  }, [hasFirstUpdate]);

  // Animation 2/3a: Location-change capture.
  // When location changes, snapshot the current rendered angles as the tween start point.
  useEffect(() => {
    if (!hasAnimated.current) return;
    tweenFromRef.current = { ...visualRef.current };
    tweenPendingRef.current = true;
  }, [locationChanged]);

  // Animation 2/3b: Location-change tween playback.
  // When new target arrives after a location change, play the tween from the captured tween start point.
  useEffect(() => {
    if (!tweenPendingRef.current || target.needle === null) return;
    tweenPendingRef.current = false;

    if (animFrameRef.current !== null) cancelAnimationFrame(animFrameRef.current);

    const from = { ...tweenFromRef.current };
    const to = { ...targetRef.current };
    const start = performance.now();
    const duration = 400;

    function frame(now: number) {
      const raw = Math.min((now - start) / duration, 1);
      const e = easeInOut(raw);
      setVisual({
        needle: lerp(from.needle ?? Math.PI, to.needle ?? Math.PI, e),
        sunrise: lerp(from.sunrise, to.sunrise, e),
        sunset: lerp(from.sunset, to.sunset, e),
      });
      if (raw < 1) {
        animFrameRef.current = requestAnimationFrame(frame);
      } else {
        animFrameRef.current = null;
      }
    }

    animFrameRef.current = requestAnimationFrame(frame);
    return () => {
      if (animFrameRef.current !== null) cancelAnimationFrame(animFrameRef.current);
    };
  }, [target.needle, target.sunrise, target.sunset]);

  // Animation 3/3: Real-time sync.
  // Snap visual straight to target on each server tick.
  // Bail out if unchanged to avoid re-renders, and skip while an animation is running.
  useEffect(() => {
    if (!hasAnimated.current || animFrameRef.current !== null) return;
    const t = target;
    setVisual(prev =>
      prev.needle === t.needle && prev.sunrise === t.sunrise && prev.sunset === t.sunset
        ? prev : t
    );
  }, [target.needle, target.sunrise, target.sunset]);

  const needleEnd = visual.needle !== null ? arcPoint(visual.needle) : null;
  const showNeedle = needleEnd !== null && solarTime != null;

  const fullWedge = wedgePath(Math.PI, 0);
  const baseColor: string = today === null ? THEME.textDim : THEME.solarClock.night;

  // Wedge paths only change when sunrise/sunset angles change, memoize to skip per-tick recomputation.
  const arcPaths = useMemo(() => {
    if (visual.sunrise === null || visual.sunset === null) return null;
    return {
      left:  wedgePath(Math.PI, visual.sunrise),
      day:   wedgePath(visual.sunrise, visual.sunset),
      right: wedgePath(visual.sunset, 0),
    };
  }, [visual.sunrise, visual.sunset]);

  // Polar day/night: today exists but no sunrise or sunset.
  let polarBaseColor: string | null = null;
  if (today !== null && visual.sunrise === null && visual.sunset === null) {
    polarBaseColor = altitudeDeg > 0 ? THEME.solarClock.day : THEME.solarClock.night;
  }

  return (
    <svg
      viewBox="0 0 400 210"
      style={{ width: '100%', maxWidth: '560px', display: 'block', margin: '0 auto' }}
      aria-label="Solar time clock"
    >
      {!arcPaths && <path d={fullWedge} fill={polarBaseColor ?? baseColor} />}
      {arcPaths && (
        <>
          <path d={arcPaths.left}  fill={THEME.solarClock.night} />
          <path d={arcPaths.day}   fill={THEME.solarClock.day} />
          <path d={arcPaths.right} fill={THEME.solarClock.night} />
        </>
      )}

      {showNeedle && (
        <>
          <line
            x1={CX} y1={CY} x2={needleEnd[0]} y2={needleEnd[1]}
            stroke={THEME.solarClock.needle} strokeWidth={THEME.solarClock.needleStrokeWidth} strokeLinecap="round"
          />
          <circle
            cx={needleEnd[0]} cy={needleEnd[1]} r={THEME.solarClock.needleTipRadius}
            fill={THEME.solarClock.day} stroke={THEME.solarClock.needle} strokeWidth={THEME.solarClock.needleStrokeWidth}
          />
        </>
      )}

      <circle cx={CX} cy={CY} r={THEME.solarClock.pivotRadius} fill={THEME.solarClock.pivot} />
    </svg>
  );
}
