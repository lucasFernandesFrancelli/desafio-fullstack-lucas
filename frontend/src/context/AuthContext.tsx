import { createContext, useCallback, useContext, useEffect, useState, type ReactNode } from "react";

import { useLoginMutation, useMe } from "@/api/auth";
import { getStoredToken, setStoredToken } from "@/api/httpClient";
import type { User } from "@/types/api";

interface AuthContextValue {
  user: User | null;
  /** true enquanto tenta reidratar a sessão a partir do token salvo. */
  isLoading: boolean;
  login: (userId: string) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isRehydrating, setIsRehydrating] = useState<boolean>(Boolean(getStoredToken()));
  const loginMutation = useLoginMutation();
  const meQuery = useMe(Boolean(getStoredToken()) && isRehydrating);

  useEffect(() => {
    if (!isRehydrating) {
      return;
    }
    if (meQuery.isSuccess) {
      setUser(meQuery.data);
      setIsRehydrating(false);
    } else if (meQuery.isError) {
      setStoredToken(null);
      setUser(null);
      setIsRehydrating(false);
    }
  }, [isRehydrating, meQuery.isSuccess, meQuery.isError, meQuery.data]);

  const login = useCallback(
    async (userId: string) => {
      const response = await loginMutation.mutateAsync(userId);
      setStoredToken(response.token);
      setUser(response.user);
    },
    [loginMutation],
  );

  const logout = useCallback(() => {
    setStoredToken(null);
    setUser(null);
  }, []);

  const value: AuthContextValue = { user, isLoading: isRehydrating, login, logout };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthContextValue {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth deve ser usado dentro de <AuthProvider>");
  }
  return context;
}
