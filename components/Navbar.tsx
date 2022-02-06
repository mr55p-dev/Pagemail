import { signOut } from "@firebase/auth";
import Link from "next/link";
import { useContext } from "react";
import { UserContext } from "../lib/context";
import { auth } from "../lib/firebase";
import { AuthCheckQuiet } from "./AuthCheck";

export default function Navbar() {
    const { user, username } = useContext(UserContext);

    const SignOut = () => signOut(auth);

    const signedOutDisplay = () => {
        return(
            <>
                <Link href="/enter">
                    <button className="nav-btn nav-user-signin"><a>Sign in</a></button>
                </Link>
            </>
        )
    }

    const signedInDisplay = () => {
        const photoURL = user?.photoURL ? user.photoURL : "/empty-avatar.png"
        return(
            <>
                <img className="nav-user-profile"
                    src={photoURL} alt="Profile image" />
                <button className="nav-btn nav-user-signout" onClick={SignOut}>
                        Sign out
                </button>
            </>
        )
    }

    return(
        <nav className="">
            <div className="nav-left">
                <Link href="/">
                    <a className="nav-link"><span className="nav-brand">PageMail</span></a>
                </Link>
                <Link href="/blog"><a className="nav-link">Blog</a></Link>
                <AuthCheckQuiet>
                    <Link href="/upload"><a className="nav-link">Upload</a></Link>
                    <Link href="/pages"><a className="nav-link">My Pages</a></Link>
                </AuthCheckQuiet>
            </div>
            <div className="nav-right">
                { user ? signedInDisplay() : signedOutDisplay() }
            </div>
        </nav>
    )
}