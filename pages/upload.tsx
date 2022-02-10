import { useState, useEffect, useContext } from "react";
import { AuthCheck } from "../components/AuthCheck";
import { UserContext } from "../lib/context";
import { storeUserURL } from "../lib/firebase";
import { scrapePageMetadata } from "../lib/scraping";
import { toast } from "react-toastify";
import Modal from "../components/modal";


export default function UploadPage() {
    const [userURL, setUserURL] = useState(null);
    const [userIdToken, setUserIdToken] = useState("");

    const [pageMetadata, setPageMetadata] = useState({});
    const [loading, setLoading] = useState(false);
    const [showModal, setShowModal] = useState(false);

    const { user } = useContext(UserContext);


    useEffect(() => {
        if (user) {
            user.getIdToken()
            .then(token => setUserIdToken(token))
            .catch(err => console.error(err))
        } else {
            setUserIdToken("");
        }
    }, [user])

    useEffect(() => {
        console.log(pageMetadata);
    }, [pageMetadata])

    const onSubmit = (e): void => {
        // Break if the user is not valid
        if (!user) {
            throw ("User not signed in")
        }

        // Set loading
        setLoading(true);

        // Prevent the default redirection
        e.preventDefault();

        // Fetch the page metadata
        scrapePageMetadata(userURL, userIdToken)
            .then(meta => setPageMetadata(meta))
            .catch(() => {
                throw "Failed to scrape the page metadata";
            })

        // Store the URL in firebase
        storeUserURL(user.uid, userURL);

        // Unset loading
        setLoading(false);
        setShowModal(true);
    }

    const onChange = (e): void => {
        try {
            const inputURL = new URL(e.target.value)
            setUserURL(inputURL);
        }
        catch {
            console.error("Invalid URL");
            setUserURL("")
        }
    }

    return(
        <main>
            <AuthCheck>
                <div className="heading sidebar">
                    <h1 className="heading">{loading ? "Loading..." : "Loaded"}</h1>
                </div>
                <div className="form-container">
                    <form onSubmit={onSubmit} className="form">
                        <input name="url" placeholder="URL" onChange={onChange} className="form-input" autoComplete="off"/>
                        <button type="submit" className="form-button">Submit</button>
                    </form>
                </div>
                <Modal show={showModal} onClose={() => setShowModal(false)}>
                    <h4>{pageMetadata?.title}</h4>
                    <p>{pageMetadata?.description}</p>
                    <img src={pageMetadata?.image} className="modal-image" ></img>
                </Modal>
            </AuthCheck>
        </main>
    )
}