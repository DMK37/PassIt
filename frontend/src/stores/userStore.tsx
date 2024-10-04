import { create } from "zustand";
import { persist } from "zustand/middleware";
import { User } from "../types/user";
import { jwtDecode } from "jwt-decode";

interface UserStore {
  user: User | null;
  token: string | null;
  setUser: (user: User, token: string) => void;
  logout: () => void;
}

const useUserStore = create<UserStore>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      setUser: (user: User, token: string) => set({ user, token }),
      logout: () => set({ user: null, token: null }),
    }),
    {
      name: "user-storage",
    }
  )
);

interface TokenPayload {
  exp: number;
}

export const isTokenExpired = (token: string | null): boolean => {
  if (!token) return true;

  try {
    const { exp } = jwtDecode<TokenPayload>(token);
    if (Date.now() >= exp * 1000) {
      return true;
    }
  } catch (error) {
    return true;
  }

  return false;
};

export default useUserStore;
