import { useEffect, useMemo, useRef, useState } from 'react';
import type { Today } from './types';

interface SolarClockProps {
  solarTime: string | null;
  today: Today | null;
  altitudeDeg: number;
  locationChanged: boolean;
}

const CX = 200;
const CY = 200;
const R = 180;

function timeToSeconds(hms: string): number {
  const [h, m, s] = hms.split(':').map(Number);
  return (h ?? 0) * 3600 + (m ?? 0) * 60 + (s ?? 0);
}

function solarAngle(secs: number): number {
  return (1 - secs / 86400) * Math.PI;
}

function arcPoint(theta: number): [number, number] {
  return [CX + R * Math.cos(theta), CY - R * Math.sin(theta)];
}

// Filled pie slice from center to arc segment
function wedgePath(startTheta: number, endTheta: number): string {
  if (Math.abs(startTheta - endTheta) < 0.001) return '';
  const [x1, y1] = arcPoint(startTheta);
  const [x2, y2] = arcPoint(endTheta);
  const large = startTheta - endTheta > Math.PI ? 1 : 0;
  return `M ${CX} ${CY} L ${x1} ${y1} A ${R} ${R} 0 ${large} 1 ${x2} ${y2} Z`;
}

function easeOut(t: number): number {
  return 1 - (1 - t) * (1 - t);
}

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

  // Start animation: sweeps needle from midnight to current time, arc boundaries grow from edges
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

  // When location changes, capture the current rendered angles as the tween start point.
  // The new update hasn't arrived yet, so we can't start the tween here - just set a flag.
  useEffect(() => {
    if (!hasAnimated.current) return;
    tweenFromRef.current = { ...visualRef.current };
    tweenPendingRef.current = true;
  }, [locationChanged]);

  // When a new target arrives after a location change, play the tween from the captured state.
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

  // Real-time updates: sync visual to target on each tick; bail out if unchanged to avoid
  // re-renders, and skip while an animation is running.
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
  let baseColor = '#3D4A5C';
  if (today === null) baseColor = '#6B7280';

  // Wedge paths only change when sunrise/sunset angles change - memoize to skip per-tick recomputation
  const arcPaths = useMemo(() => {
    if (visual.sunrise === null || visual.sunset === null) return null;
    return {
      left:  wedgePath(Math.PI, visual.sunrise),
      day:   wedgePath(visual.sunrise, visual.sunset),
      right: wedgePath(visual.sunset, 0),
    };
  }, [visual.sunrise, visual.sunset]);

  // Polar day/night: today exists but no sunrise or sunset
  let polarBaseColor: string | null = null;
  if (today !== null && visual.sunrise === null && visual.sunset === null) {
    polarBaseColor = altitudeDeg > 0 ? '#F97316' : '#3D4A5C';
  }

  return (
    <svg
      viewBox="0 0 400 210"
      style={{ width: '100%', maxWidth: '560px', display: 'block', margin: '0 auto' }}
      aria-label="Solar time clock"
    >
      {/* Filled half-circle: full wedge as base, then day/night wedges on top */}
      {!arcPaths && <path d={fullWedge} fill={polarBaseColor ?? baseColor} />}
      {arcPaths && (
        <>
          <path d={arcPaths.left}  fill="#3D4A5C" />
          <path d={arcPaths.day}   fill="#F97316" />
          <path d={arcPaths.right} fill="#3D4A5C" />
        </>
      )}

      <circle cx={CX} cy={CY - R} r="4" fill="rgba(255,255,255,0.5)" />

      {showNeedle && (
        <>
          <line
            x1={CX} y1={CY} x2={needleEnd[0]} y2={needleEnd[1]}
            stroke="white" strokeWidth="2" strokeLinecap="round"
          />
          <circle cx={needleEnd[0]} cy={needleEnd[1]} r="10" fill="#F97316" stroke="white" strokeWidth="2" />
        </>
      )}

      <circle cx={CX} cy={CY} r="5" fill="#9CA3AF" />
    </svg>
  );
}
