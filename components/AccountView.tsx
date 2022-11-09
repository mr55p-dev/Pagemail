import { doc, setDoc, onSnapshot } from "@firebase/firestore";
import Head from "next/head";
import { useEffect, useState } from "react";
import { AuthCheck } from "../components/AuthCheck"
import { useAuth } from "../lib/context";
import { firestore } from "../lib/firebase";
import { IUserData } from "../lib/typeAliases";

export function AccountView ({ }): JSX.Element {
    const { authUser } = useAuth();
    useEffect(() => {

    })
    const email = authUser?.email
    const username = authUser?.displayName
    const [newsletterPref, setNewsletterPref] = useState<boolean>(undefined);

    // Subscribe to the user document and listen for changes
    useEffect(() => {
        const userDoc = doc(firestore, "users", authUser?.uid);
        const unsubscribe = onSnapshot(userDoc, (userData) => {
            setNewsletterPref((userData.data() as IUserData).newsletter)
        })
        return () => unsubscribe()
    }, [])

    // Handle preference updation
    const handleNewsletterStateChange = (e) => {
        e.preventDefault()
        console.log("HandlingNewsletter")
        const userDoc = doc(firestore, "users", authUser.uid)
        setDoc(userDoc, {newsletter: !newsletterPref}, {merge: true})
        .then(() => console.log("Updated preferences."))
    }

    return (
        <AuthCheck>
            <main className="p-3">
                <Head>
                    <title>Your account</title>
                </Head>
                <h1 className="page-heading">Account information</h1>
                <div>
                    <form className="grid grid-rows-5 grid-cols-12 gap-2">
                        <p className="col-span-12 md:col-span-4">Username: </p>
                            <input className="col-span-12 md:col-span-8" defaultValue={username} readOnly={true} />
                        <p className="col-span-12 md:col-span-4">Email address: </p>
                            <input className="col-span-12 md:col-span-8" defaultValue={email} readOnly={true} />
                        <p className="col-span-10 md:col-span-4">Subscribe to emails: </p>
                            {newsletterPref === undefined ?
                            "Loading newsletter preference..." :
                            <input id="newsletter" className=" col-span-2 md:col-span-8" type="checkbox" checked={newsletterPref} readOnly={true}/>}
                        <button onClick={handleNewsletterStateChange}>{newsletterPref !== true ? "Subscribe to" : "Unsubscribe from"} the newsletter</button>
                    </form>
                </div>
            </main>
        </AuthCheck>
     );
}
