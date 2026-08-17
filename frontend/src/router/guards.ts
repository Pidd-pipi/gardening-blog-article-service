import { createElement, type ReactNode } from 'react';
import { Navigate } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';

export function RequireAuth({ children }: { children: ReactNode }) {
  const { isLoggedIn } = useAuth();
  if (!isLoggedIn) return createElement(Navigate, { to: '/login', replace: true });
  return createElement(ReactFragment, null, children);
}

export function RequireAdmin({ children }: { children: ReactNode }) {
  const { isAdmin } = useAuth();
  if (!isAdmin) return createElement(Navigate, { to: '/', replace: true });
  return createElement(ReactFragment, null, children);
}

function ReactFragment({ children }: { children: ReactNode }) {
  return createElement('div', { style: { display: 'contents' } }, children);
}
