import Link from "next/link";
import { useContext } from "react";
import { UserContext } from "../lib/context";

export default function Navbar() {
    const { user, username } = useContext(UserContext);

    return(
        <nav>
            <p>{username}</p>
        </nav>
    )
}