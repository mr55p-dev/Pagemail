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
            <>
                <a className="inline-block text-sm px-4 py-2 h-12 w-auto justify-center leading-none border rounded text-white border-white
                    hover:border-transparent hover:text-sky-600 hover:bg-white
                    mt-4 lg:mt-0" href="/enter">
                        Sign In
                </a>
            </>
        )
    }

    const signedInDisplay = () => {
        return(
            <>
                <img className="inline object-cover w-12 h-12 mr-2 rounded-full"
                    src={user?.photoURL} alt="Profile image" />
                <button className="inline-block text-sm px-4 py-2 h-12 w-auto justify-center leading-none border rounded text-white border-white
                    hover:border-transparent hover:text-sky-600 hover:bg-white
                    mt-4 lg:mt-0" onClick={SignOut}>
                        Sign out
                </button>
            </>
        )
    }

    return(
    <>
        <nav className="flex items-center justify-between flex-wrap bg-sky-600 p-6">
            <div className="flex-shrink-0 text-black mr-5">
                <a href="/">
                    <span className="font-bold text-3xl tracking-tight text-white">PageMail</span>
                </a>
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
        {/* <nav className="flex items-center justify-between flex-wrap bg-teal-500 p-6">
            <div className="flex items-center flex-shrink-0 text-white mr-6">
                <img className="fill-current h-8 w-8 mr-2" width="54" height="54" src="/icon.png" />
                <span className="font-semibold text-xl tracking-tight">PageMail</span>
            </div>
            <div className="block lg:hidden">
                <button className="flex items-center px-3 py-2 border rounded text-teal-200 border-teal-400 hover:text-white hover:border-white">
                    <svg className="fill-current h-3 w-3" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><title>Menu</title><path d="M0 3h20v2H0V3zm0 6h20v2H0V9zm0 6h20v2H0v-2z"/></svg>
                </button>
            </div>
            <div className="w-full block flex-grow lg:flex lg:items-center lg:w-auto">
                <div className="text-sm lg:flex-grow">
                <a href="#responsive-header" className="block mt-4 lg:inline-block lg:mt-0 text-teal-200 hover:text-white mr-4">
                    Docs
                </a>
                <a href="#responsive-header" className="block mt-4 lg:inline-block lg:mt-0 text-teal-200 hover:text-white mr-4">
                    Examples
                </a>
                <a href="#responsive-header" className="block mt-4 lg:inline-block lg:mt-0 text-teal-200 hover:text-white">
                    Blog
                </a>
                </div>
                <div>
                <a href="#" className="inline-block text-sm px-4 py-2 leading-none border rounded text-white border-white hover:border-transparent hover:text-teal-500 hover:bg-white mt-4 lg:mt-0">Download</a>
                </div>
            </div>


            <p>{username}</p>
            <img className="inline object-cover w-16 h-16 mr-2 rounded-full" src={user?.photoURL}></img>
        </nav> */}
        </>

    )
}