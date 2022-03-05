import Link from "next/link";
import { useContext } from "react";
import { useAuth } from "../lib/context";

export function AuthCheck(props) {
    const { user } = useAuth();

    return user ?
        props.children :
        props.fallback || <Link href="/enter">Sign up or Log in here.</Link>
}

export function AuthCheckQuiet(props) {
    const { user } = useAuth()
    return user ? props.children : null
}