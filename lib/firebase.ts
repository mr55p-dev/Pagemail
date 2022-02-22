// Import the functions you need from the SDKs you need
import { initializeApp, getApp } from "firebase/app";

import { getAuth, GoogleAuthProvider, EmailAuthProvider, connectAuthEmulator, User } from 'firebase/auth';
import { getFirestore, connectFirestoreEmulator, setDoc, addDoc, doc, collection, serverTimestamp } from 'firebase/firestore';
import { IPage, IUserDoc } from "./typeAliases";

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

export const auth = getAuth(app);
export const googleAuth = new GoogleAuthProvider();
export const emailAuth = new EmailAuthProvider();
export const firestore = getFirestore(app);

if (process.env.USE_EMULATOR == "1") {
    connectAuthEmulator(auth, "http://localhost:9099");
    connectFirestoreEmulator(firestore, "localhost", 8080);
}

export function storeUserData(user: User) {
    const writableValues: IUserDoc = {
        username: user.displayName,
        email: user.email,
        photoURL: user.photoURL,
        anonymous: user.isAnonymous,
        newsletter: false
    }

    // Add the user to the users collection
    setDoc(doc(firestore, "users", user.uid), writableValues)
        .then(() => console.log("Sucseffully written user doc."))
        .catch(() => console.error("Failed to write user doc"))
}

export function storeUserURL(userid: string, url: URL) {
    const writableValues: IPage = {
        url: url.toString(),
        timeAdded: serverTimestamp()
    }

    const PageDoc = collection(firestore, "users", userid, "pages")
    addDoc(PageDoc, writableValues)
    .then(() => console.log("Written document!"))
    .catch(() => console.error("Failed to write document!"));
}