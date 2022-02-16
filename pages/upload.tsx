import { useState, useContext } from "react";
import { AuthCheck } from "../components/AuthCheck";
import { UserContext } from "../lib/context";
import { storeUserURL } from "../lib/firebase";
import Modal from "../components/modal";
import { usePageMetadata, useUserToken } from "../lib/hooks";

function validateURL(inputString: string): URL {
    // First coerce the input into a URL
    if (!(inputString.startsWith("https://") || inputString.startsWith("http://"))) {
        console.log("Prepending https://")
        inputString = `https://${inputString}`;
    }
    const inputURL = new URL(inputString)

    // Allow only http or https
    if (!["http:", "https:"].includes(inputURL.protocol)) {
        throw "Invalid URL protocol"
    }
    if (!(/^\S+\.\S+$/.test(inputURL.hostname))) {
        throw "Invalid URL hostname"
    }
    return inputURL
}

export default function UploadPage() {
    const [userURL, setUserURL] = useState<URL>(undefined);
    const [loading, setLoading] = useState<boolean>(false);

    const { user } = useContext(UserContext);

    const [showModal, setShowModal] = useState<boolean>(false);

    const token = useUserToken();
    const pageMetadata = usePageMetadata(userURL, token);

    const onSubmit: React.FormEventHandler<HTMLFormElement> = (e) => {
        // Prevent the default redirection
        e.preventDefault();

        // Break if the user is not valid
        if (!token) {
            throw ("User not signed in")
        }

        // Check that the user URL is not undefined
        if (userURL === undefined) {
            console.error("Provide a valid URL")
            return
        }

        // Set loading
        setLoading(true);

        // Only do this if everything works fine
        e.currentTarget.reset()

        // Save the URL
        storeUserURL(user.uid, userURL)

        // Enable the modal
        setShowModal(true);

        // Done!
        setLoading(false)
    }


    const onChange: React.ChangeEventHandler<HTMLInputElement> = (e) => {
        setShowModal(false);
        try {
            const valid = validateURL(e.target.value);
            setUserURL(valid);
        }
        catch (error) {
            console.error("invalid url")
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
                        <p>{userURL !== undefined ? "Valid URL!" : "Invalid URL :("}</p>
                        <button type="submit" className="form-button">Submit</button>
                    </form>
                </div>
                {pageMetadata !== undefined ?
                    pageMetadata.title ?
                        <p>Page Metadata: {pageMetadata.title}</p>
                        :
                        null
                    : null
                }
                {pageMetadata !== undefined ?
                    pageMetadata.title ?
                        <Modal show={showModal} onClose={() => setShowModal(false)}>
                            <h4>{pageMetadata?.title}</h4>
                            <p>{pageMetadata?.description}</p>
                            <img alt="" src={pageMetadata?.image} className="modal-image" />
                        </Modal>
                        : <p>No metadata</p>
                    :
                    <p>Loading...</p>
                }
            </AuthCheck>
        </main>
    )
}