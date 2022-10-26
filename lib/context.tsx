import { User } from "firebase/auth";
import { createContext, useContext } from "react";
import useFirebaseAuth from "./useFirebaseAuth";

interface TAuthCtx {
    authUser: User,
    loading: boolean
};

const baseCtx: TAuthCtx = { authUser: null, loading: true }
const AuthUserContext = createContext(baseCtx);

export function AuthUserProvider({ children }) {
    const auth = useFirebaseAuth();
    return <AuthUserContext.Provider value={auth}>{children}</AuthUserContext.Provider>
}

export const useAuth = () => useContext(AuthUserContext);
