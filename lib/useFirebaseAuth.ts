import { useEffect, useState } from "react";
import { auth } from "./firebase";


export default function useFirebaseAuth() {
    const [authUser, setAuthUser] = useState(null);
    const [loading, setLoading] = useState(true);

    const authStateChanged = async (user) => {
        if (!user) {
            setAuthUser(null);
            setLoading(false);
            return;
        } else {
            setAuthUser(user);
            setLoading(false);
        }
    }

    useEffect(() => {
        const unsub = auth.onAuthStateChanged(authStateChanged);
        return () => unsub();
    }, [])

    return {authUser, loading};
}