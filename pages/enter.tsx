import { GoogleAuthProvider, signInWithPopup, signOut } from "@firebase/auth";
import { auth, googleAuth } from "../lib/firebase"
import { useContext } from "react";
import { UserContext } from "../lib/context";


export default function Enter({ }) {

    const { user, username } = useContext(UserContext);

    function SignInButton() {
        const signInWithGoogle = () => {
            signInWithPopup(auth, googleAuth)
            .then((result) => {
                const credential = GoogleAuthProvider.credentialFromResult(result);
                // const token = credential.accessToken;

                // const user = result.user
                // const username = user.displayName
            })
            .catch((err) => {
                const errorCode = err.code;
                // const errorMessage = err.message;

                // const email = err.email;
                // const credential = GoogleAuthProvider.credentialFromError(err);
            })
        }

        return(
            <button className="signin-google-btn" onClick={signInWithGoogle}>
                <img className="signin-google-img" src="/google-signin@1x.png" />
            </button>
        )
    };

    function SignOutButton() {
        const SignOut = () => signOut(auth);

        return(
            <button className="signout-btn" onClick={SignOut}>
                Sign out
            </button>
        )
    };

    return(
        <main>
            {
            user
            ?
            <>
                <h1 className="text-3xl font-bold underline">Signed in as {username}</h1>
                <SignOutButton />
            </>
            :
            <>
                <h1 className="text-3xl font-bold underline">Sign in or register</h1>
                <SignInButton />
            </>
            }
        </main>
    )
}