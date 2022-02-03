import { createContext } from "react";

const defaultContext = { user: null, username: null };

export const UserContext = createContext(defaultContext);
