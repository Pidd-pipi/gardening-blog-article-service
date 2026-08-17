export const UserRole = {
  ADMIN: 'admin',
  VISITOR: 'visitor',
} as const;

export const UserRoleLabels: Record<string, string> = {
  [UserRole.ADMIN]: '博主',
  [UserRole.VISITOR]: '游客',
};
