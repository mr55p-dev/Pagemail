import { useAuthState } from "react-firebase-hooks/auth";
import { useState, useEffect } from "react";
import { getAuth, User } from "@firebase/auth";
import { INotifState, IPage, IUserContext, IUserData, IUserDoc } from "./typeAliases";
import { collection, CollectionReference, doc, DocumentReference, getFirestore, onSnapshot } from "firebase/firestore";


export function useUserData(): IUserData {
    const [user] = useAuthState(getAuth());
    const [userData, setUserData] = useState<IUserData>(undefined);

    useEffect(() => {
      if(user) {
        const userRef = doc(getFirestore(), "users", user.uid) as DocumentReference<IUserDoc>
        const unsubscribe = onSnapshot(userRef, (userDoc) => {
          // const userDocData = userDoc.data()
          setUserData({
            // email: userDocData.email,
            // photoURL: userDocData.photoURL,
            // anonymous: userDocData.anonymous,
            // newsletter: userDocData.newsletter,
            ...userDoc.data(),
            user: user,
            pages: collection(getFirestore(), "users", user.uid, "pages") as CollectionReference<IPage>
          })
        })
        return () => unsubscribe()
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