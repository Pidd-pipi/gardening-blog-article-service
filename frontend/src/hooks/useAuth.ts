import { useAuth as useAuthContext } from '../stores/authStore';

// 登录态与角色判断
export function useAuth() {
  const { user, token, setAuth, updateUser, logout } = useAuthContext();
  return {
    user,
    token,
    setAuth,
    updateUser,
    logout,
    isAdmin: user?.role === 'admin',
    isLoggedIn: !!token,
  };
}
