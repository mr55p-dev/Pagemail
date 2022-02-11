import { createContext } from "react";
import { IUserContext, INotifContext } from "./typeAliases";


const defaultUserContext: IUserContext = undefined;
const defaultNotifContext: INotifContext = undefined;

export const UserContext = createContext(defaultUserContext);
export const NotifContext = createContext(defaultNotifContext);