import { useState, useContext, useEffect } from "react";
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
    const [loading, setLoading] = useState<boolean>(undefined);
    const [loadingText, setLoadingText] = useState<string>("");

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

    useEffect(() => {
        console.log(loading)
        setLoadingText(loading === undefined ? "" : loading ? "Loading..." : "Loaded")
    }, [loading])


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
        <div>
            <div className="">
                <h1 className="page-heading">Upload</h1>
                <h3 className="heading">{loadingText}</h3>
            </div>
            <AuthCheck>
                <div className="flex justify-around">
                    <div className="border-2 rounded m-2 p-3 bg-sky-50 border-sky-700 max-w-screen-md">
                        <p className="mt-1 mb-3">Use this form to save new pages to your space. Changes will be reflected instantly under your pages!</p>
                        <form onSubmit={onSubmit} className="form flex flex-col md:flex-row">
                            <input name="url" placeholder="URL" onChange={onChange} className="w-full bg-sky-50 border-2 rounded border-sky-700 outline-none btn-shape inline p-1" autoComplete="off"/>
                            <button type="submit" className="btn-shape border-2 btn-colour p-2 rounded mt-2 md:mt-0 md:mx-1 md:px-2">Submit</button>
                        </form>
                        {/* <p className="">{userURL !== undefined ? "Valid URL!" : "Invalid URL :("}</p> */}
                    </div>
                    {/* {pageMetadata !== undefined ?
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
                    } */}
                </div>
            </AuthCheck>
        </div>
    )
}