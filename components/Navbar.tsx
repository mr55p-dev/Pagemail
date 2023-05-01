import { auth } from "../lib/firebase";
import { signOut } from "@firebase/auth";
import Link from "next/link";
import { useState } from "react";
import { useAuth } from "../lib/context";

export default function Navbar() {
    const { user } = useAuth();
    const [ mobileShow, setMobileShow ] = useState(false);

    const navClickHandler = () => setTimeout(() => {setMobileShow(false)}, 100)

    // Some resources for the navbar
    const photoURL = user?.photoURL || "empty-avatar.png"
    const menuIcon = (<svg className="h-8 w-8 fill-current" viewBox="0 0 24 24">
        <path fillRule="evenodd" d="M4 5h16a1 1 0 0 1 0 2H4a1 1 0 1 1 0-2zm0 6h16a1 1 0 0 1 0 2H4a1 1 0 0 1 0-2zm0 6h16a1 1 0 0 1 0 2H4a1 1 0 0 1 0-2z"/>
    </svg>)

    // Components which are conditional on being signed in
    const signedInDisplay = user ? (
        <>
        <div className="nav-border" />
        <Link href="/account" passHref>
            <a className="inline" onClick={navClickHandler}>
                <div className="btn-shape btn-colour mx-1 my-2 px-3 py-2 hover:cursor-pointer">
                    <img className="rounded-full object-cover inline w-8 h-8 mr-2" src={photoURL} />
                    { user.displayName }
                </div>
            </a>
        </Link>
        <div className="nav-border" />
        <Link href="/upload" passHref>
            <a className="nav-link" onClick={navClickHandler}>
                <div className="btn-shape btn-colour mx-1 my-2 p-3 hover:cursor-pointer">Upload</div>
            </a>
        </Link>
        <div className="nav-border" />
        <Link passHref href="/pages">
            <a className="nav-link" onClick={navClickHandler}>
                <div className="btn-shape btn-colour mx-1 my-2 p-3 hover:cursor-pointer">My Pages</div>
            </a>
        </Link>
        <div className="nav-border" />
        <div className="btn-shape btn-colour mx-1 my-2 p-3 hover:cursor-pointer">
            <button onClick={() => {navClickHandler(); signOut(auth)}}>
                <a>Sign Out</a>
            </button>
        </div>
        <div className="nav-border" />
        </>
    ) : (
        <>
        <div className="nav-border" />
        <Link passHref href="/enter">
            <a className="nav-link" onClick={navClickHandler}>
                <div className="btn-shape btn-colour mx-1 my-2 p-3 hover:cursor-pointer">
                    Sign in
                </div>
            </a>
        </Link>
        <div className="nav-border" />
        </>
    )

    return(
        <nav className="md:flex md:justify-between md:items-center text-secondary dark:text-secondary-dark text-center
                        max-w-screen-xl mx-auto">
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
                <div className="pb-1
                                transition-all md:flex md:items-center md:py-0">
                    { signedInDisplay }
                </div>
            </div>
        </nav>
    )
}
