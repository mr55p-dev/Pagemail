// Import the functions you need from the SDKs you need
import { initializeApp, getApp, FirebaseApp } from "firebase/app";

import { getAuth, GoogleAuthProvider, EmailAuthProvider, connectAuthEmulator, User, Auth } from 'firebase/auth';
import { getFirestore, connectFirestoreEmulator, setDoc, addDoc, doc, collection, serverTimestamp, Firestore } from 'firebase/firestore';
import { IPage, IUserDoc } from "./typeAliases";


interface IFirebaseConfig {
	apiKey: string;
	authDomain: string;
	projectId: string;
	storageBucket: string;
	messagingSenderId: string;
	appId: string;
	measurementId?: string;
};

// Your web app's Firebase configuration
// For Firebase JS SDK v7.20.0 and later, measurementId is optional
function getFirebaseconfig(): IFirebaseConfig {
	if (process.env.vercel_env === "production") {
		console.log('Using production firebase environment')
		return {
			apiKey: "AIzaSyCt4lpPzhe_UKbvlOcE7g_HSrz4stQbDjQ",
			authDomain: "pagemail-2bc26.firebaseapp.com",
			projectId: "pagemail-2bc26",
			storageBucket: "pagemail-2bc26.appspot.com",
			messagingSenderId: "556909502728",
			appId: "1:556909502728:web:9392f6243b38ceef2c8cbd",
			measurementId: "G-Q62RYYT55K"
		}
	} else {
		console.log('Using preprod firebase environment');
		return {
			apiKey: "AIzaSyCuO0tYBhtwsG0zEYRTFE605aVbLVPNPDs",
			authDomain: "pagemail-preprod.firebaseapp.com",
			projectId: "pagemail-preprod",
			storageBucket: "pagemail-preprod.appspot.com",
			messagingSenderId: "538304052902",
			appId: "1:538304052902:web:2a1def87e86772c9cd7e67"
		}
	}
}


function createFirebaseApp(cfg: IFirebaseConfig) {
    try {
        return getApp();
    } catch {
        return initializeApp(cfg);
    }
}

// Initialise app
const firebaseConfig = getFirebaseconfig();
const app = createFirebaseApp(firebaseConfig);

const getAuthLocal = (app: FirebaseApp): Auth => {
	console.log('Fetching auth');
	const auth = getAuth(app);
	console.log(process.env)
	if (process.env.EMULATE === '1') {
		console.log('Connecting auth emulator');
		connectAuthEmulator(auth, "http://localhost:9099");
	}
	else { console.log('nothing') }
	return auth

}
const getFirestoreLocal = (app: FirebaseApp): Firestore => {
	console.log('Fetching firestore');
    const firestore = getFirestore(app);
	if (process.env.EMULATE == '1') {
		console.log('Connecting firestore emulator');
		connectFirestoreEmulator(firestore, "localhost", 8080);
	}
	return firestore
}

export const auth = getAuthLocal(app);
export const googleAuth = new GoogleAuthProvider();
export const emailAuth = new EmailAuthProvider();
export const firestore = getFirestoreLocal(app);

export function storeUserData(user: User) {
    const writableValues: IUserDoc = {
        username: user.displayName,
        email: user.email,
        photoURL: user.photoURL,
        anonymous: user.isAnonymous
    }

    // Add the user to the users collection
    setDoc(doc(firestore, "users", user.uid), writableValues, { merge: true })
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
