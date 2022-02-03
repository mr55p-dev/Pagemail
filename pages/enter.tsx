import { GoogleAuthProvider, signInWithPopup, signOut } from "@firebase/auth";
import { auth, googleAuth } from "../lib/firebase"
import { useContext } from "react";
import { UserContext } from "../lib/context";


export default function Enter({ }) {

    const { user, username } = useContext(UserContext);

    function SignInButton() {
        // const signInWithGoogle = async () => {
        //     const result = await signInWithPopup(auth, googleAuth)

        // };
        const signInWithGoogle = () => {
            signInWithPopup(auth, googleAuth)
            .then((result) => {
                const credential = GoogleAuthProvider.credentialFromResult(result);
                const token = credential.accessToken;

                const user = result.user
                const username = user.displayName
            })
            .catch((err) => {
                const errorCode = err.code;
                const errorMessage = err.message;
    
                const email = err.email;
                const credential = GoogleAuthProvider.credentialFromError(err);
            })
        }

        return(
            <button className="signin-google-btn" onClick={signInWithGoogle}>
                <img className="signin-google-img" src="/google-signin@1x.png" />
            </button>
        )
    };

    function SignOutButton() {
        const SignOut = () => {
            console.log("Called signout")
            signOut(auth)
        };

        return(
            <button className="signout-btn" onClick={SignOut}>
                Sign out
            </button>
        )
    };

    return(
        <main>
            <h1>Registration</h1>
            {user ? <SignOutButton /> : <SignInButton />}
        </main>
    )
}