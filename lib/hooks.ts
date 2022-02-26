import { useAuthState } from "react-firebase-hooks/auth";
import { useState, useEffect } from "react";
import { getAuth } from "@firebase/auth";
import { IPage, IPageMetadata, IUserData, IUserDoc } from "./typeAliases";
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


export function usePageMetadata(url: URL, token: string) {

  // process.env.PAGEMAIL_API_ORIGIN ||

  const [pData, setPData] = useState<IPageMetadata>(undefined);
  const emptyMetadata: IPageMetadata = {
    url: "",
    title: "",
    author: "",
    description: "",
    image: ""
  }

  useEffect(() => {
    // If the URL is bad then dont try and request
    if (!url || !token) {
      setPData(undefined);
    } else {
      // Allow the request to be cancelled if the side effect is refreshed
      const controller = new AbortController()
      const { signal } = controller;

      // Get the API address
      const apiAddress = new URL(window.location.origin)

      // Modify the path and query parameters
      apiAddress.pathname = "/api/scrape";
      apiAddress.searchParams.set("url", encodeURIComponent(url.toString()))

      // Get a response
      fetch(apiAddress.toString(), {
        method: "GET",
        mode: "same-origin",
        credentials: "same-origin",
        headers: {
            token: token
        },
        signal: signal
      })
      .then((resp) => {
        if (!resp.ok) {
            return {} as IPageMetadata
        }
        return resp.json()
      })
      .then((body) => {
        setPData({
          url: url.toString(),
          title: body.title,
          author: body.author,
          description: body.description,
          image: body.image
        })
      })
      .catch((err) => {
        if (err.name !== "AbortError") {
          console.error(err);
          setPData(emptyMetadata)
        }
      })
      return () => controller.abort()
    }

  }, [url])

  return pData
}