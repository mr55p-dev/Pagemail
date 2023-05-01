// Import the functions you need from the SDKs you need
import { initializeApp, getApp } from "firebase/app";
import { getAuth, connectAuthEmulator, User } from 'firebase/auth';
import { getFirestore, connectFirestoreEmulator, setDoc, addDoc, doc, collection, serverTimestamp } from 'firebase/firestore';
import { IPage, IUserDoc } from "./typeAliases";


// Your web app's Firebase configuration
// For Firebase JS SDK v7.20.0 and later, measurementId is optional
interface IFirebaseConfig {
	apiKey: string;
	authDomain: string;
	projectId: string;
	storageBucket: string;
	messagingSenderId: string;
	appId: string;
	measurementId?: string;
};

function getFirebaseconfig(): IFirebaseConfig {
	if (process.env.NEXT_PUBLIC_ENVIRONMENT === "production") {
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

// Create the application
function createFirebaseApp(cfg: IFirebaseConfig) {
    try {
        return getApp();
    } catch {
        return initializeApp(cfg);
    }
}

const firebaseConfig = getFirebaseconfig();
export const app = createFirebaseApp(firebaseConfig);
export const auth = getAuth(app);
export const firestore = getFirestore(app);

// Register auth and firebase, configure emulators if necessary
if (process.env.NEXT_PUBLIC_USE_EMULATOR === '1') {
	console.log('Connecting auth, firestore emulator');
	const auth = getAuth(app);
	const firestore = getFirestore(app);
	connectAuthEmulator(auth, "http://localhost:9099");
	connectFirestoreEmulator(firestore, "localhost", 8080);
}

export function storeUserData(user: User) {
    const writableValues: IUserDoc = {
        username: user.displayName,
        email: user.email,
        photoURL: user.photoURL,
        anonymous: user.isAnonymous
    }
	const app = getApp();
	const firestore = getFirestore(app);

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
	const app = getApp();
	const firestore = getFirestore(app);

    const PageDoc = collection(firestore, "users", userid, "pages")
    addDoc(PageDoc, writableValues)
    .then(() => console.log("Written document!"))
    .catch(() => console.error("Failed to write document!"));
}

