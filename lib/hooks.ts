import { useAuthState } from "react-firebase-hooks/auth";
import { useState, useEffect, EffectCallback } from "react";
import { doc, getFirestore, onSnapshot } from "@firebase/firestore";
import { getAuth } from "@firebase/auth";


export function useUserData() {
    const [user] = useAuthState(getAuth());
    const [username, setUsername] = useState(null);

    useEffect(() => {
        let unsubscribe;
        setUsername(user?.displayName);
    //   if (user) {
    //     setUsername(user.displayName)
    //   } else {
    //     setUsername(null);
    //   }
      return unsubscribe;
    }, [user])

    return [user, username];
}