import PocketBase, { Record, RecordAuthResponse } from "pocketbase";
import { UserRecord } from "./datamodels";
import React from "react";
import { AuthState, DataState } from "./data";
import { useNavigate } from "react-router-dom";
import { NotificationCtx } from "./notif";

const pb_url = import.meta.env.VITE_PAGEMAIL_API_HOST;
console.log(pb_url);
export const pb = new PocketBase(pb_url);
console.log(pb.baseUrl);
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

function getAuthState(user: UserRecord | null | undefined): AuthState {
  if (!user) {
    return AuthState.UNAUTHORIZED;
  } else if (!user.verified) {
    return AuthState.UNVERIFIED;
  } else {
    return AuthState.AUTH;
  }
}

export const useUser = () => {
  const [user, setUser] = React.useState<UserRecord | null>(getCurrentUser());
  const [authErr, setAuthErr] = React.useState<Error | null>(null);
  const [authState, setAuthState] = React.useState<AuthState>(
    getAuthState(user)
  );
  const [reqState, setReqState] = React.useState<DataState>(DataState.UNKNOWN);
  const nav = useNavigate();
  const { notifInfo, notifErr } = React.useContext(NotificationCtx);

  React.useEffect(() => {
    const unsub = pb.authStore.onChange(() => {
      setUser(getCurrentUser());
    });

    return () => unsub();
  }, []);

  React.useEffect(() => {
    setAuthState(getAuthState(user));
  }, [user]);

  async function login(
    callback: () => Promise<RecordAuthResponse<UserRecord>>
  ): Promise<UserRecord> {
    setReqState(DataState.PENDING);
    try {
      const rval = await callback();
      setAuthState(getAuthState(user));
      setAuthErr(null);
      setReqState(DataState.SUCCESS);
      nav("/pages");
      return rval.record;
    } catch (err) {
      setAuthState(AuthState.UNAUTHORIZED);
      setAuthErr(err as Error);
      setReqState(DataState.FAILED);
      notifErr((err as Error).message);
      return Promise.reject(err);
    }
  }

  async function refresh() {
    return await pb.collection("users").authRefresh();
  }

  function logout() {
    pb.authStore.clear();
    setAuthState(AuthState.UNAUTHORIZED);
    notifInfo("Signed out");
  }

  return { user, authState, login, refresh, logout, authErr, reqState };
};
