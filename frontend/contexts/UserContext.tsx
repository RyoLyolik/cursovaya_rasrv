"use client";

import { createContext, useContext, useEffect, useMemo, useState } from "react";
import axios from "axios";
import { useRouter } from "next/navigation";

interface Role {
  title: string;
}

interface UserData {
  email: string;
  creation_date: string;
  role: Role;
}

interface UserContextType {
  user: UserData | null;
  loading: boolean;
  error: string | null;
  logout: () => Promise<void>;
  refreshUser: () => Promise<void>;
}

const UserContext = createContext<UserContextType | undefined>(undefined);

export function UserProvider({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const [user, setUser] = useState<UserData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchUser = async () => {
    try {
      const response = await axios.get("http://localhost:3377/api/account", {
        withCredentials: true,
      });
      
      if (response.data.status === "Error") {
        throw new Error(response.data.error?.detail || "Failed to fetch user data");
      }
      
      setUser(response.data?.data);
      setError(null);
    } catch (err) {
      setUser(null);
      setError(
        err instanceof Error ? err.message : "Failed to fetch user data"
      );
      if (axios.isAxiosError(err) && err.response?.status === 401) {
        router.push("/login");
      }
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = async () => {
    try {
      await axios.post(
        "http://localhost:3377/api/logout",
        {},
        { withCredentials: true }
      );
      setUser(null);
      router.push("/login");
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Logout failed"
      );
    }
  };

  useEffect(() => {
    fetchUser();
  }, []);

  const contextValue = useMemo(() => ({
    user,
    loading,
    error,
    logout: handleLogout,
    refreshUser: fetchUser
  }), [user, loading, error])

  return (
    <UserContext.Provider 
      value={contextValue}
    >
      {children}
    </UserContext.Provider>
  );
}

export function useUser() {
  const context = useContext(UserContext);
  if (context === undefined) {
    throw new Error("useUser must be used within a UserProvider");
  }
  return context;
}