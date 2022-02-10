// Import the functions you need from the SDKs you need
import { initializeApp, getApp } from "firebase/app";
import { getAnalytics } from "firebase/analytics";

import { getAuth, GoogleAuthProvider, EmailAuthProvider, connectAuthEmulator } from 'firebase/auth';
import { getFirestore, connectFirestoreEmulator, setDoc, addDoc, doc, collection, serverTimestamp } from 'firebase/firestore';
import { useContext } from "react";
import { UserContext } from "./context";

// Your web app's Firebase configuration
// For Firebase JS SDK v7.20.0 and later, measurementId is optional
const firebaseConfig = {
  apiKey: "AIzaSyCt4lpPzhe_UKbvlOcE7g_HSrz4stQbDjQ",
  authDomain: "pagemail-2bc26.firebaseapp.com",
  projectId: "pagemail-2bc26",
  storageBucket: "pagemail-2bc26.appspot.com",
  messagingSenderId: "556909502728",
  appId: "1:556909502728:web:9392f6243b38ceef2c8cbd",
  measurementId: "G-Q62RYYT55K"
};

function createFirebaseApp(cfg) {
    try {
        return getApp();
    } catch {
        return initializeApp(cfg);
    }
}

// Initialise app
const app = createFirebaseApp(firebaseConfig);
const analytics = getAnalytics(app);

export const auth = getAuth(app);
export const googleAuth = new GoogleAuthProvider();
export const emailAuth = new EmailAuthProvider();
connectAuthEmulator(auth, "http://localhost:9099");

export const firestore = getFirestore(app);
connectFirestoreEmulator(firestore, "localhost", 8080);

export function storeUserData(user) {
    const writableValues = {
        username: user.displayName,
        email: user.email,
        photoURL: user.photoURL,
        anonymous: user.isAnonymous,
        // Add in newsletter preferences
        newsletter: false
    }

    // Add the user to the users collection
    setDoc(doc(firestore, "users", user.uid), writableValues)
    .then(() => console.log("Sucseffully written user doc."))
    .catch(() => console.error("Failed to write user doc"))
}

export function storeUserURL(userid: string, url: URL) {
    const writableValues = {
        url: url.toString(),
        timeAdded: serverTimestamp()
    }
    const PageDoc = collection(firestore, "users", userid, "pages")
    addDoc(PageDoc, writableValues)
    .then(() => console.log("Written document!"))
    .catch(() => console.error("Failed to write document!"));
}