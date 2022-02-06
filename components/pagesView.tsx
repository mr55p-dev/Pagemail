import { collection } from "firebase/firestore";
import { useContext, useEffect, useState } from "react"
import { UserContext } from "../lib/context"
import { firestore } from "../lib/firebase";

export default function PagesView(props) {

    const { user } = useContext(UserContext);
    const [ pages, setPages ] = useState([]);

    useEffect(() => {
        let unsubscribe;

        if (user) {
            const pagesRef = collection(firestore, "users", user.uid, "pages");

        }

    }, [user])

    return(
        <div className="pages-container">
            <h1>Inside the container</h1>
            {/* For loop to print all the pages */}
        </div>
    )
}