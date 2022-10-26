import Link from "next/link";
import { useAuth } from "../lib/context";

export function AuthCheck(props) {
    const { authUser } = useAuth();

    return authUser ?
        props.children :
        props.fallback || <Link href="/enter">Sign up or Log in here.</Link>
}

export function AuthCheckQuiet(props) {
    const { authUser } = useAuth()
    return authUser ? props.children : null
}
