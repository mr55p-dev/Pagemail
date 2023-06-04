import React from "react";
import { DataState } from "./data";

interface UserContextType {
  user: UserRecord | null;
  checkUser: () => void;
  authStatus: DataState;
  setAuthStatus: (status: DataState) => void;
}

const UserContext = React.createContext<UserContextType>({
  user: null,
  // eslint-disable-next-line @typescript-eslint/no-empty-function
  checkUser: () => {},
  authStatus: DataState.UNKNOWN,
  // eslint-disable-next-line @typescript-eslint/no-empty-function
  setAuthStatus: () => {},
});

export default UserContext;
