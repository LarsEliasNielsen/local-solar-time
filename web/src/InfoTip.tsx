import { useEffect, useId, useRef, useState } from 'react';
import { THEME } from './theme';

interface InfoTipProps {
  text: string;
  label: string;
  align?: 'left' | 'right';
  placement?: 'top' | 'bottom';
  ariaLabel?: string;
}

// Small "?" badge that shows a plain-language hint on hover, focus, or tap.
export default function InfoTip({ text, label, align = 'left', placement = 'bottom', ariaLabel }: InfoTipProps) {
  const [hovered, setHovered] = useState(false);
  const [pinned, setPinned] = useState(false);
  const id = useId();
  const rootRef = useRef<HTMLSpanElement>(null);
  const open = hovered || pinned;

  useEffect(() => {
    if (!pinned) return;
    const onPointerDown = (e: PointerEvent) => {
      if (!rootRef.current?.contains(e.target as Node)) setPinned(false);
    };
    const onKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') setPinned(false);
    };
    document.addEventListener('pointerdown', onPointerDown);
    document.addEventListener('keydown', onKeyDown);
    return () => {
      document.removeEventListener('pointerdown', onPointerDown);
      document.removeEventListener('keydown', onKeyDown);
    };
  }, [pinned]);

  return (
    <span
      ref={rootRef}
      style={{ position: 'relative', display: 'inline-flex' }}
      onMouseEnter={() => setHovered(true)}
      onMouseLeave={() => setHovered(false)}
    >
      <button
        type="button"
        aria-label={ariaLabel ?? `What is ${label}?`}
        aria-describedby={open ? id : undefined}
        aria-expanded={open}
        onClick={() => setPinned(true)}
        onFocus={() => setPinned(true)}
        onBlur={() => setPinned(false)}
        style={{
          width: '16px',
          height: '16px',
          padding: 0,
          borderRadius: '50%',
          border: `1px solid ${THEME.borderStrong}`,
          background: open ? THEME.surfaceRaised : 'transparent',
          color: THEME.textMuted,
          fontSize: '0.625rem',
          fontWeight: 700,
          lineHeight: '14px',
          cursor: 'help',
        }}
      >
        ?
      </button>
      {open && (
        <span
          id={id}
          role="tooltip"
          style={{
            position: 'absolute',
            [placement === 'bottom' ? 'top' : 'bottom']: 'calc(100% + 6px)',
            [align]: 0,
            zIndex: 10,
            width: 'max-content',
            maxWidth: 'min(240px, 80vw)',
            padding: '8px 10px',
            borderRadius: '6px',
            border: `1px solid ${THEME.border}`,
            background: THEME.surface,
            color: THEME.textPrimary,
            fontSize: '0.75rem',
            lineHeight: 1.4,
            textAlign: 'left',
            boxShadow: '0 4px 12px rgba(0, 0, 0, 0.4)',
          }}
        >
          {text}
        </span>
      )}
    </span>
  );
}
