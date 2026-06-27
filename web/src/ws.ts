import type { SolarUpdate } from './types';

export function openSocket(
  lat: number,
  lon: number,
  onMessage: (update: SolarUpdate) => void,
  onClose: () => void,
): () => void {
  const proto = window.location.protocol === 'https:' ? 'wss' : 'ws';
  const ws = new WebSocket(`${proto}://${window.location.host}/ws`);
  ws.addEventListener('open', () => ws.send(JSON.stringify({ lat, lon })));
  ws.addEventListener('message', e => onMessage(JSON.parse(e.data) as SolarUpdate));
  ws.addEventListener('close', onClose);
  ws.addEventListener('error', () => ws.close());
  return () => ws.close();
}
