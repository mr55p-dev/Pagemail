import { createContext, useContext, useEffect, useState } from "react";
import { User, onIdTokenChanged } from "firebase/auth";
import { getAuth } from "firebase/auth";
import nookies from "nookies";

const AuthContext = createContext<{user: User}>({user: null});

export const useAuth = () => useContext(AuthContext);

export const AuthProvider = ({ children }) => {
    const [user, setUser] = useState<User | null>(null);

    // Listen for user state change (token or otherwise)
    useEffect(() => {
        return onIdTokenChanged(getAuth(), async (user) => {
            if (!user) {
                setUser(user);
                nookies.set(undefined, 'token', '', { path: '/' });
            } else {
                const token = await user.getIdToken();
                setUser(user);
                nookies.set(undefined, 'token', token, { path: '/' });
            }
        })
    }, [user])

      // force refresh the token every 10 minutes
    useEffect(() => {
        const handle = setInterval(async () => {
            const user = getAuth().currentUser;
            if (user) await user.getIdToken(true);
        }, 10 * 60 * 1000)

        // clean up setInterval
        return () => clearInterval(handle);
    }, []);

    return <AuthContext.Provider value={{user: user}}>{children}</AuthContext.Provider>
}
