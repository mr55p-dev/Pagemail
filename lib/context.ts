import { createContext } from "react";

const defaultUserContext = { user: null, username: null };

export const UserContext = createContext(defaultUserContext);
