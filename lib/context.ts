import { createContext } from "react";

const defaultUserContext = { user: null, username: null };
const defaultNotifContext = { showNotifCallback };

export const UserContext = createContext(defaultUserContext);
export const NotifContext = createContext(defaultNotifContext);