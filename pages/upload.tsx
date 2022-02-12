import { useState, useEffect, useContext } from "react";
import { AuthCheck } from "../components/AuthCheck";
import { UserContext } from "../lib/context";
import { storeUserURL } from "../lib/firebase";
import { scrapePageMetadata } from "../lib/scraping";
import Modal from "../components/modal";
import { INotifState, IPageMetadata } from "../lib/typeAliases";
import { useUserToken } from "../lib/hooks";
import Notif from "../components/notif";


export default function UploadPage() {
    const [userURL, setUserURL] = useState<URL>(undefined);
    const [pageMetadata, setPageMetadata] = useState<IPageMetadata>(undefined);
    const [loading, setLoading] = useState<boolean>(false);
    const [showModal, setShowModal] = useState<boolean>(false);

    const { user } = useContext(UserContext);
    const token = useUserToken()

    const [showNotif, setShowNotif] = useState<boolean>(false);
    const [stateNotif, setStateNotif] = useState<INotifState>(undefined);

    const [error, setError] = useState<string>("")


    const onSubmit: React.FormEventHandler<HTMLFormElement> = (e) => {
        // Prevent the default redirection
        e.preventDefault();

        // Break if the user is not valid
        if (!token) {
            throw ("User not signed in")
        }

        // Set loading
        setLoading(true);

        // Fetch the page metadata
        scrapePageMetadata(userURL, token)
            .then(meta => {
                if (meta) {
                    setPageMetadata(meta)
                    setError("")
                } else {
                    setPageMetadata(undefined)
                }
            })
            .then(() => {
                storeUserURL(user.uid, userURL);
                setShowModal(true);
            })
            .catch(() => {
                setPageMetadata(undefined)
                setError("Failed to scrape the page metadata!")
            })
            .then(() => {
                setLoading(false);
            })

    }

    // Side effect to render the notification for 5 seconds
    useEffect(() => {
        if (!loading) {
            setShowNotif(true);
            setStateNotif({
                title:  error ? error   : "Success!",
                text:   error ? ""      : pageMetadata?.title,
                style:  error ? "error" : "success"
            })
            const timer = setTimeout(() => {
                setShowNotif(false);
                setStateNotif(undefined);
            }, 5000)
            return () => clearTimeout(timer);
        }
    }, [loading, pageMetadata, error])

    const onChange: React.ChangeEventHandler<HTMLInputElement> = (e) => {
        setShowModal(false);
        try {
            const inputURL = new URL(e.target.value)
            setUserURL(inputURL);
        }
        catch {
            console.error("Invalid URL");
            setUserURL(undefined)
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
                    <img alt="" src={pageMetadata?.image} className="modal-image" />
                </Modal>
                <Notif show={showNotif} state={stateNotif}/>
            </AuthCheck>
        </main>
    )
}