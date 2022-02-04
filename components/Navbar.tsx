import { signOut } from "@firebase/auth";
import Link from "next/link";
import { useContext } from "react";
import { UserContext } from "../lib/context";
import { auth } from "../lib/firebase";

export default function Navbar() {
    const { user, username } = useContext(UserContext);

    const SignOut = () => signOut(auth);

    const scale = 12;

    const signedOutDisplay = () => {
        return(
            <><Link href="/enter">
                <button className="inline-block text-sm px-4 py-2 h-12 w-auto justify-center leading-none border rounded text-white border-white
                    hover:border-transparent hover:text-sky-600 hover:bg-white
                    mt-4 lg:mt-0">
                        <a>Sign in</a>
                </button>
            </Link></>
        )
    }

    const signedInDisplay = () => {
        const photoURL = user?.photoURL ? user.photoURL : "/empty-avatar.png"
        return(
            <>
                <img className="inline object-cover w-12 h-12 mr-2 rounded-full"
                    src={photoURL} alt="Profile image" />
                <button className="inline-block text-sm px-4 py-2 h-12 w-auto justify-center leading-none border rounded text-white border-white
                    hover:border-transparent hover:text-sky-600 hover:bg-white
                    mt-4 lg:mt-0" onClick={SignOut}>
                        Sign out
                </button>
            </>
        )
    }

    return(
        <nav className="flex items-center justify-between flex-wrap bg-sky-600 p-6 shadow-xl">
            <div className="flex-shrink-0 text-black mr-5">

                <Link href="/">
                    <a><span className="font-bold text-3xl tracking-tight text-white">PageMail</span></a>
                </Link>
            </div>
            <div className="w-full block flex-grow flex w-auto">
                <div className="lg:flex-grow">
                    <a href="/blog" className="block lg:inline-block text-neutral-300 hover:text-white">Blog</a>
                </div>
            </div>
            <div>
                { user ? signedInDisplay() : signedOutDisplay() }
            </div>
        </nav>
    )
}