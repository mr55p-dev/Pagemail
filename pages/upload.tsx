import { doc, getFirestore } from "@firebase/firestore";
import { useState, useEffect, useContext } from "react";
import { AuthCheck } from "../components/AuthCheck";
import { UserContext } from "../lib/context";
import { storeUserURL } from "../lib/firebase";

export default function UploadPage() {
    const [userURL, setUserURL] = useState(null);
    const { user } = useContext(UserContext);

    const onSubmit = (e: Event): void => {
        // Prevent the default redirection
        e.preventDefault();

        //
        console.log(`Submit URL: ${userURL}`);

        //
        storeUserURL(user.uid, userURL);
    }

    const onChange = (e: Event): void => {
        setUserURL(e.target.value)
    }

    // useEffect(() => {
    //     console.log(userURL)
    // }, [userURL])

    return(
        <main>
            <AuthCheck>
                <h1 className="heading">You are authorised to upload.</h1>
                <div className="form-container">
                    <form onSubmit={onSubmit} className="form">
                        <input name="url" placeholder="URL" onChange={onChange} className="form-input" autoComplete="off"/>
                        <button type="submit" className="form-button">Submit</button>
                    </form>
                </div>
            </AuthCheck>
        </main>
    )
}