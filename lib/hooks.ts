import { useAuthState } from "react-firebase-hooks/auth";
import { useState, useEffect } from "react";
import { getAuth } from "@firebase/auth";
import { INotifState } from "./typeAliases";


export function useUserData() {
    const [user] = useAuthState(getAuth());
    const [username, setUsername] = useState<string>(undefined);

    useEffect(() => {
        setUsername(user?.displayName);
    }, [user])

    return [user, username];
}

export function useUserToken() {
  const [user] = useAuthState(getAuth());
  const [token, setToken] = useState<string>("")

  useEffect(() => {
    if (user) {
      user.getIdToken()
      .then(t => setToken(t))
      .catch((e) => {
        setToken("")
        console.error(`Failed to retrieve user token: ${e.message}`)
      })
    } else {
      setToken("")
    }
  }, [user])

  return token
}

export function useRendered(): boolean {
  const [isBrowser, setIsBrowser] = useState<boolean>(false);

  useEffect(() => {
      setIsBrowser(true);
  })

  return isBrowser
}

export function useNotif(contents: INotifState) {
  
}