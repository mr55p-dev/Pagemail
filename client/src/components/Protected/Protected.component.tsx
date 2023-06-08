import React from "react";
import { Navigate } from "react-router";
import { useUser } from "../../lib/pocketbase";

const Protected = ({ children }: { children: React.ReactNode }) => {
  const { user } = useUser();
  return user !== null ? children : <Navigate to="/auth" replace/>;
};

export default Protected;
