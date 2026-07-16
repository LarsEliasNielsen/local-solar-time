import { Component, type ReactNode } from 'react';
import { THEME } from './theme';

interface Props { children: ReactNode; fallback?: ReactNode }
interface State { error: Error | null }

export class ErrorBoundary extends Component<Props, State> {
  state: State = { error: null };

  static getDerivedStateFromError(error: Error): State {
    return { error };
  }

  render() {
    if (this.state.error) {
      return this.props.fallback ?? (
        <p style={{ color: THEME.textError, padding: '16px' }}>Something went wrong.</p>
      );
    }
    return this.props.children;
  }
}
