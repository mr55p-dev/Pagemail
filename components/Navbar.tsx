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
    const menuIcon = (<svg className="h-8 w-8 fill-current" viewBox="0 0 24 24"><path fill-rule="evenodd" d="M4 5h16a1 1 0 0 1 0 2H4a1 1 0 1 1 0-2zm0 6h16a1 1 0 0 1 0 2H4a1 1 0 0 1 0-2zm0 6h16a1 1 0 0 1 0 2H4a1 1 0 0 1 0-2z"/></svg>)

    // Components which are conditional on being signed in
    const SignOut = user && (
        <div className="btn-shape btn-colour py-3 hover:cursor-pointer">
            <button onClick={() => signOut(auth)}>
                <a>Sign Out</a>
            </button>
        </div>)

    const signedInDisplay = user ? (
        <>
        <Link href="/account">
            <div className="btn-shape btn-colour py-2 hover:cursor-pointer">
                <a className="inline">{ user.displayName }</a>
                <img className="ml-2 h-8 w-8 inline" src={photoURL} />
            </div>
        </Link>
        <hr />
        <Link href="/upload">
            <div className="btn-shape btn-colour py-3 hover:cursor-pointer">
                <a className="nav-link">Upload</a>
            </div>
        </Link>
        <hr />
        <Link href="/pages">
            <div className="btn-shape btn-colour py-3 hover:cursor-pointer">
                <a className="nav-link">My Pages</a>
            </div>
        </Link>
        </>
    ) : (
        <Link href="/enter">
            <div className="btn-shape btn-colour py-3 hover:cursor-pointer">
                <a>Sign in</a>
            </div>
        </Link>
    )

    return(
        <nav className="bg-sky-50 md:flex md:justify-between md:items-center text-sky-800 text-center">
            <div className="flex justify-between items-center px-3 py-6">
                <div className="text-3xl">
                    <Link href="/">
                        <a className="">
                            <span className="nav-brand">PageMail</span>
                        </a>
                    </Link>
                </div>
                <div className="inline px-1">
                    <button className="h-8 w-8 sm:hidden" onClick={() => setMobileShow(!mobileShow)}>
                        { menuIcon }
                    </button>
                </div>
            </div>
            <div className={mobileShow ? "hidden" : "block"}>
                <div className="pb-1 border-t-2 border-sky-100 transition-all md:flex md:items-center md:border-0 md:py-0">
                    { signedInDisplay }
                    <hr className="border-sky-100"/>
                    <Link href="/about">
                        <div className="btn-shape btn-colour py-3 hover:cursor-pointer">
                            <a>About</a>
                        </div>
                    </Link>
                    <hr className="border-sky-100"/>
                    <Link href="/contact">
                        <div className="btn-shape btn-colour py-3 hover:cursor-pointer">
                            <a>Contact</a>
                        </div>
                    </Link>
                    { user && <hr className="border-sky-100"/> }
                    { SignOut }
                </div>
            </div>
        </nav>
                // <AuthCheckQuiet>
                //     <Link href="/upload"><a className="nav-link">Upload</a></Link>
                //     <Link href="/pages"><a className="nav-link">My Pages</a></Link>
                // </AuthCheckQuiet>
    )
}