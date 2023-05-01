import { auth, storeUserData } from "../lib/firebase";
import { GoogleAuthProvider, signInWithPopup } from "@firebase/auth";
import Head from "next/head";
import { useAuth } from "../lib/context";


export default function Enter({ }) {

    const { user } = useAuth();
	const googleAuth = new GoogleAuthProvider();

    function SignInForm() {
        const signInWithGoogle = async () => {
            signInWithPopup(auth, googleAuth)
            .then((result) => { storeUserData(result.user) })
            .catch((err) => {console.error(err.message)})
        }

        return(
            <div className="form-container">
                <form className="form">
                        <button className="form-button button-google-signup" type="button" onClick={signInWithGoogle}>
                            <img className="signin-google-img" src="/google-signin@2x.png" />
                        </button>
                </form>
            </div>
        )
    };

    return(
        <main>
            <Head>
                <title>Sign in - PageMail</title>
            </Head>
                { user ?
                    <h1 className="page-heading">Hello, {user.displayName}!</h1>
                :
                <>
                    <h1 className="page-heading">Sign in or register</h1>
                    <SignInForm />
                </>
                }
        </main>
    )
}
