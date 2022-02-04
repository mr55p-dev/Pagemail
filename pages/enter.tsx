import { GoogleAuthProvider, signInWithPopup, signOut } from "@firebase/auth";
import { auth, googleAuth } from "../lib/firebase"
import { useContext } from "react";
import { UserContext } from "../lib/context";


export default function Enter({ }) {

    const { user, username } = useContext(UserContext);

    function SignInForm() {
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
            <div className="grid overflow-hidden grid-cols-3 grid-rows-2 gap-0">
                <div className="col-start-2">
                    <div className="bg-white border rounded-lg shadow-xl">
                        <form className="p-3">
                            <input type="email" id="email" className="block w-full h-12 px-3 mt-2 text-sm
                                            border border-gray-200 rounded-md ring-gray-300 ring-opacity-50
                                            focus:outline-none focus:bg-white focus:border-blue-500" placeholder="Email or phone number"/>
                            <input type="password" id="password-1" className="block w-full h-12 px-3 mt-3 text-sm
                                            border border-gray-200 rounded-md ring-gray-300 ring-opacity-50
                                            focus:outline-none focus:bg-white focus:border-blue-500" placeholder="Password" autoComplete="current-password"/>
                            <button className="w-full h-12 px-20 py-2 mt-3 text-base font-bold text-white bg-blue-500 rounded-md hover:bg-blue-600">Log In</button>
                            <div className="mt-4 text-sm font-bold text-center">
                                <a className="text-blue-500 hover:underline hover:cursor-pointer">Forgot password?</a>
                            </div>
                            <div className="mt-5">
                                <hr/>
                            </div>
                            <div className="my-4 text-center">
                                <button className="w-5/6 h-12 px-1 py-2 mt-3 text-base font-bold
                                        text-white bg-green-500 rounded-md hover:bg-green-600">Create new account</button>
                            </div>
                        </form>
                        <button className="" onClick={signInWithGoogle}>
                            <img className="signin-google-img" src="/google-signin@1x.png" />
                        </button>
                    </div>
                </div>
            </div>
        )
    };

    return(
        <main>
            {
            user
            ?
            <>
                <h1 className="text-3xl font-bold underline">Hello, {username}!</h1>
            </>
            :
            <>
                <h1 className="text-3xl font-bold underline">Sign in or register</h1>
                <SignInForm />
            </>
            }


        </main>
    )
}