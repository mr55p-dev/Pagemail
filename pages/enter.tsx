import { signInWithPopup } from "@firebase/auth";
import { auth, googleAuth} from "../lib/firebase"
import { useContext } from "react";
import { UserContext } from "../lib/context";
import { storeUserData } from "../lib/firebase";
import Head from "next/head";


export default function Enter({ }) {

    const { user } = useContext(UserContext);

    function SignInForm() {
        const signInWithGoogle = async () => {
            signInWithPopup(auth, googleAuth)
            .then((result) => { storeUserData(result.user) })
            .catch((err) => {console.log(err.message)})
        }

        return(
                    <div className="form-container">
                        <form className="form">
                            {/* <input type="email" id="email" className="form-input" placeholder="Email or phone number"/> */}
                            {/* <input type="password" id="password-1" className="form-input" placeholder="Password" autoComplete="current-password"/> */}
                            {/* <div className="enter-form-button-group"> */}
                                {/* <button className="form-button button-login">Log In</button>
                                <button className="form-button button-signup">Create new account</button> */}
                                <button className="form-button button-google-signup" type="button" onClick={signInWithGoogle}>
                                    <img className="signin-google-img" src="/google-signin@2x.png" />
                                </button>
                            {/* </div> */}
                            {/* <a className="enter-form-passreset-link">Forgot password?</a> */}
                        </form>
                    </div>
        )
    };

    return(
        <main>
            <Head>
                <title>Sign in</title>
            </Head>
                { user ?
                    <h1 className="heading">Hello, {user.displayName}!</h1>
                :
                <>
                    <h1 className="heading">Sign in or register</h1>
                    <SignInForm />
                </>
                }
        </main>
    )
}