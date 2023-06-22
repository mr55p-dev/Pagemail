import React from "react";
import { Navigate } from "react-router";
import { useUser } from "../../lib/pocketbase";
import { AuthState } from "../../lib/data";

export const Protected = ({ children }: { children: React.ReactNode }) => {
  const { user } = useUser();
  return user !== null ? children : <Navigate to="/auth" replace/>;
};

export const NotProtected = ({ children }: { children: React.ReactNode}) => {
  const { authState } = useUser()
  return authState !== AuthState.AUTH ? children : <Navigate to="/pages" replace/>
}
