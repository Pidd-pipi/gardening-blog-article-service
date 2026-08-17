import request from '../utils/request';
import type { User } from '../types';

export interface LoginResult {
  token: string;
  user: User;
}

export function login(email: string, password: string): Promise<LoginResult> {
  return request.post('/auth/login', { email, password });
}

export function register(username: string, email: string, password: string, nickname: string): Promise<LoginResult> {
  return request.post('/auth/register', { username, email, password, nickname });
}

export function getMe(): Promise<User> {
  return request.get('/users/me');
}

export function updateProfile(data: Partial<User>): Promise<User> {
  return request.put('/users/me', data);
}
