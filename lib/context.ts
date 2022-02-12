import { createContext } from "react";
import { IUserContext, INotifContext } from "./typeAliases";


const defaultUserContext: IUserContext = { user: null };
const defaultNotifContext: INotifContext = { setNotifShow: null, setNotifState: null };

export const UserContext = createContext(defaultUserContext);
export const NotifContext = createContext(defaultNotifContext);