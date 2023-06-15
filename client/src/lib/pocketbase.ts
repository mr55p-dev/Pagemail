import PocketBase, { Record } from "pocketbase";
import { UserRecord } from "./datamodels";
import React from "react";
import { AuthState } from "./data";
import { useNavigate } from "react-router-dom";

const pb_url =
  process.env.NODE_ENV === "development"
    ? "http://127.0.0.1:8090"
    : "https://v2.pagemail.io";
export const pb = new PocketBase(pb_url);
pb.autoCancellation(false);

export const getCurrentUser = (): UserRecord | null => {
  if (pb.authStore.model instanceof Record) {
    const mdl = pb.authStore.model;
    return {
      id: mdl.id,
      email: mdl.email,
      created: mdl.created,
      updated: mdl.updated,
      verified: mdl.verified ? true : false,
      subscribed: mdl.subscribed,
    };
  }
  return null;
};

export const useUser = () => {
  const [user, setUser] = React.useState<UserRecord | null>(getCurrentUser());
  const [authErr, setAuthErr] = React.useState<Error | null>(null);
  const [authState, setAuthState] = React.useState<AuthState>(
    user ? AuthState.AUTH : AuthState.NOT_AUTH
  );
  const nav = useNavigate();

  React.useEffect(() => {
    const unsub = pb.authStore.onChange(() => {
      setUser(getCurrentUser());
    });

	return () => unsub()
  }, []);

  React.useEffect(() => {
    setAuthState(user ? AuthState.AUTH : AuthState.NOT_AUTH);
  }, [user]);

  const login = async <T>(callback: () => Promise<T>): Promise<T> => {
    setAuthState(AuthState.PENDING);
    try {
      const rval = await callback();
      setAuthState(AuthState.AUTH);
      setAuthErr(null);
      nav("/pages");
      return rval;
    } catch (err) {
      setAuthState(AuthState.NOT_AUTH);
      setAuthErr(err as Error);
      return Promise.reject(err);
    }
  };

  const logout = () => {
    pb.authStore.clear();
    setAuthState(AuthState.NOT_AUTH);
  };

  return { user, authState, login, logout, authErr };
};
