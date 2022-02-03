import { useAuthState } from "react-firebase-hooks/auth";
import { useState, useEffect, EffectCallback } from "react";
import { doc, getFirestore, onSnapshot } from "@firebase/firestore";
import { getAuth } from "@firebase/auth";


export function useUserData() {
    const [user] = useAuthState(getAuth());
    const [username, setUsername] = useState(null);

    useEffect(() => {
      let unsubscribe;

      if (user) {
        // const ref = doc(getFirestore(), "users", "uid")
        // unsubscribe = onSnapshot(ref, (doc) => {
        //   setUsername(doc.data()?.username);
        // })

        setUsername(user.displayName)
      } else {
        setUsername(null);
      }
      console.log(user)
      return unsubscribe;
    }, [user])

    return [user, username];
}