import { useState, useEffect } from "react";
import { AuthCheck } from "../components/AuthCheck";
import { storeUserURL } from "../lib/firebase";
import { usePageMetadata, useUserToken } from "../lib/hooks";
import Head from "next/head";
import { useAuth } from "../lib/context";

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
    const [borderColour, setBorderColour] = useState<string>("border-tertiary")
    const [canSubmit, setCanSubmit] = useState<boolean>(false);

    const { user } = useAuth();


    const token = useUserToken();
    // const pageMetadata = usePageMetadata(userURL, token);

    const onSubmit: React.FormEventHandler<HTMLFormElement> = (e) => {
        // Prevent the default redirection
        e.preventDefault();

        // Set loading
        setLoading(true);

        // Break if the user is not valid
        if (!token) {
            throw ("User not signed in")
        }

        // Check that the user URL is not undefined
        if (userURL === undefined) {
            console.error("Provide a valid URL")
            return
        }


        // Only do this if everything works fine
        e.currentTarget.reset()

        // Save the URL
        storeUserURL(user.uid, userURL)

        // Done!
        setLoading(false)
        setUserURL(undefined)
    }

    useEffect(() => {
        if (userURL === undefined) {
            setBorderColour("border-tertiary")
            setCanSubmit(false)
        } else {
            setBorderColour("border-green-600")
            setCanSubmit(true)
        }
    }, [userURL])


    const onChange: React.ChangeEventHandler<HTMLInputElement> = (e) => {
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
        <div className="text-secondary dark:text-secondary-dark p-3">
            <Head>
                <title>Save a new page</title>
            </Head>
            <div className="">
                <h1 className="page-heading">Upload</h1>
            </div>
            <AuthCheck>
                <p className="py-2">Use this form to save new pages to your space. Changes will be reflected instantly under your pages!</p>
                <form onSubmit={onSubmit} className="grid grid-rows-3 grid-cols-1
                    md:grid-rows-2 md:grid-cols-12 md:gap-4">
                    <input required name="url" placeholder="URL" onChange={onChange}
                    className={`w-full bg-primary dark:bg-primary-dark border-2 outline-none ${borderColour}
                    md:col-span-10 my-2 p-2 }`} autoComplete="off"/>
                    <div className={`border-2 ${borderColour} md:col-span-10 my-2 p-2`}>
                        <p className="">{userURL !== undefined ? "Valid URL!" : "Invalid URL"}</p>
                    </div>
                    <button type="submit" disabled={!canSubmit} className="border-2 submit-enabled submit-disabled
                    md:col-span-2 md:row-span-2 md:row-start-1 md:col-start-11 my-2">Submit</button>
                </form>

                <div className="flex justify-around max-w-screen-xl">
                </div>
            </AuthCheck>
        </div>
    )
}