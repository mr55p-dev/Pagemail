import { doc, getFirestore } from "@firebase/firestore";
import { useState, useEffect, useContext } from "react";
import { AuthCheck } from "../components/AuthCheck";
import { UserContext } from "../lib/context";
import { storeUserURL } from "../lib/firebase";
import { scrapePageMetadata } from "../lib/scraping";

export default function UploadPage() {
    const [userURL, setUserURL] = useState(null);
    const [pageMetadata, setPageMetadata] = useState({});
    const [loading, setLoading] = useState(false);

    const { user } = useContext(UserContext);

    const onSubmit = (e): void => {
        // Break if the user is not valid
        if (!user) { return }

        // Set loading
        setLoading(true);

        // Prevent the default redirection
        e.preventDefault();

        let userToken;
        user.getIdToken()
        .then((token) => { userToken = token; })

        // Fetch the page metadata
        scrapePageMetadata(userURL, user?.getIdToken())
            .then((meta) => { setPageMetadata(meta) })
            .catch(() => {})

        // Store the URL in firebase
        storeUserURL(user.uid, userURL);

        // Unset loading
        setLoading(false);
    }

    const onChange = (e): void => {
        setLoading(true);
        // Coerce the entered value to a URL (if possible)
        let inputURL: URL;

        // Try to create the URL
        try { inputURL = new URL(e.target.value) }
        catch { console.log("Invalid URL"); return }

        // Set the local URL state
        setUserURL(inputURL);

        if (!user) {
            console.error("User not signed in?")
            return
        }

        // Save the token
        user.getIdToken()
        .then((token: string) => {
            return scrapePageMetadata(userURL, token)
        })
        .then((meta) => {
            setPageMetadata(meta)
        })
        .then(() => setLoading(false))
        .catch((err) => {
            console.log("hiya")
            // console.error(err)
            setPageMetadata({})
            setLoading(false)
        })
    }

    return(
        <main>
            <AuthCheck>
                <div className="heading sidebar">
                    <h1 className="heading">{loading ? "Loading..." : "Loaded"}</h1>
                    <h4>{pageMetadata?.title}</h4>
                    <p>{pageMetadata?.description}</p>
                </div>
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