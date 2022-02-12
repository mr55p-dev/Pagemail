import { useAuthState } from "react-firebase-hooks/auth";
import { useState, useEffect } from "react";
import { getAuth } from "@firebase/auth";
import { INotifState, IPage, IUserData, IUserDoc } from "./typeAliases";
import { collection, CollectionReference, doc, DocumentReference, getFirestore, onSnapshot } from "firebase/firestore";


export function useUserData(): IUserData {
  const emptyUser: IUserData = {
    user: null,
    username: null,
    photoURL: null,
    email: null,
    newsletter: null,
    anonymous: null,
    pages: null
  }

  const [user] = useAuthState(getAuth());
  const [userData, setUserData] = useState<IUserData>(emptyUser);

  useEffect(() => {
    if(user) {
      const userRef = doc(getFirestore(), "users", user.uid) as DocumentReference<IUserDoc>
      const unsubscribe = onSnapshot(userRef, (userDoc) => {
        setUserData({
          ...userDoc.data(),
          user: user,
          pages: collection(
            getFirestore(),
            "users", user.uid, "pages"
          ) as CollectionReference<IPage>
        })
      })
      return () => unsubscribe()
    } else {
      setUserData(emptyUser);
    }
  }, [user])
  return userData;
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
  }, [])

  return isBrowser
}

export function useNotif(contents: INotifState) {
  
}