import { signOut } from "@firebase/auth";
import Link from "next/link";
import { useContext, useEffect, useState } from "react";
import { UserContext } from "../lib/context";
import { auth } from "../lib/firebase";

export default function Navbar() {
    const { user } = useContext(UserContext);
    const [ mobileShow, setMobileShow ] = useState(false);

    // Some resources for the navbar
    const photoURL = user?.photoURL || "empty-avatar.png"
    const menuIcon = (<svg className="h-8 w-8 fill-current" viewBox="0 0 24 24">
        <path fillRule="evenodd" d="M4 5h16a1 1 0 0 1 0 2H4a1 1 0 1 1 0-2zm0 6h16a1 1 0 0 1 0 2H4a1 1 0 0 1 0-2zm0 6h16a1 1 0 0 1 0 2H4a1 1 0 0 1 0-2z"/>
    </svg>)

    // Components which are conditional on being signed in
    const SignOut = user && (
        <div className="btn-shape btn-colour mx-1 my-2 p-3 hover:cursor-pointer">
            <button onClick={() => signOut(auth)}>
                <a>Sign Out</a>
            </button>
        </div>)

    const signedInDisplay = user ? (
        <>
        <Link href="/account">
            <div className="btn-shape btn-colour mx-1 my-2 px-3 py-2 hover:cursor-pointer">
                <a className="inline">{ user.displayName }</a>
                <img className="rounded inline w-8 h-8 ml-2" src={photoURL} />
            </div>
        </Link>
        <hr className="border-secondary dark:border-secondary-dark"/>
        <Link passHref href="/upload">
            <div className="btn-shape btn-colour mx-1 my-2 p-3 hover:cursor-pointer">
                <a className="nav-link">Upload</a>
            </div>
        </Link>
        <hr className="border-secondary dark:border-secondary-dark"/>
        <Link passHref href="/pages">
            <div className="btn-shape btn-colour mx-1 my-2 p-3 hover:cursor-pointer">
                <a className="nav-link">My Pages</a>
            </div>
        </Link>
        </>
    ) : (
        <Link passHref href="/enter">
            <div className="btn-shape btn-colour mx-1 my-2 p-3 hover:cursor-pointer">
                <a>Sign in</a>
            </div>
        </Link>
    )

    return(
        <nav className="md:flex md:justify-between md:items-center text-secondary dark:text-secondary-dark text-center
                        max-w-screen-xl mx-auto border-b-2 border-secondary md:border-0">
            <div className="flex justify-between items-center px-3 py-6">
                <div className="text-3xl">
                    <Link passHref href="/">
                        <a className="text-tertiary font-semibold">
                            <span className="nav-brand font-serif">PageMail</span>
                        </a>
                    </Link>
                </div>
                <div className="inline px-1">
                    <button className="h-8 w-8 md:hidden" onClick={() => setMobileShow(!mobileShow)}>
                        { menuIcon }
                    </button>
                </div>
            </div>
            <div className={mobileShow ? "block" : "hidden md:block"}>
                <div className="pb-1 border-t-2 border-secondary dark:border-secondary-dark
                                transition-all md:flex md:items-center md:border-0 md:py-0">
                    { signedInDisplay }
                    {/* <hr className="border-secondary"/> */}
                    {/* <Link passHref href="/about">
                        <div className="btn-shape btn-colour mx-1 my-2 p-3 hover:cursor-pointer">
                            <a>About</a>
                        </div>
                    </Link>
                    <hr className="border-secondary"/>
                    <Link passHref href="/contact">
                        <div className="btn-shape btn-colour mx-1 my-2 p-3 hover:cursor-pointer">
                            <a>Contact</a>
                        </div>
                    </Link> */}
                    { user && <hr className="border-secondary dark:border-secondary-dark"/> }
                    { SignOut }
                </div>
            </div>
        </nav>
    )
}