import React from "react";
import { Navigate } from "react-router";
import { useUser } from "../../lib/pocketbase";
import { AuthState } from "../../lib/data";

export const Protected = ({ children }: { children: React.ReactNode }) => {
  const { authState } = useUser();
  switch (authState) {
    case AuthState.AUTH:
      return children;
    case AuthState.UNVERIFIED:
      return <Navigate to="/verify" replace />;
    case AuthState.UNAUTHORIZED:
    default:
      return <Navigate to="/auth" replace />;
  }
};

export const NotVerified = ({ children }: { children: React.ReactNode }) => {
  const { authState } = useUser();
  switch (authState) {
    case AuthState.AUTH:
      return <Navigate to="/pages" replace />;
    case AuthState.UNAUTHORIZED:
      return <Navigate to="/auth" replace />;
    case AuthState.UNVERIFIED:
    default:
      return children;
  }
};

export const NotProtected = ({ children }: { children: React.ReactNode }) => {
  const { authState } = useUser();
  if (authState === AuthState.AUTH) {
    return <Navigate to="/pages" replace />;
  } else {
    return children;
  }
};
