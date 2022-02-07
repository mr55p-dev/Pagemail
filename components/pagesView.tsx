import { collection, deleteDoc, doc, getFirestore, onSnapshot } from "firebase/firestore";
import { useContext, useEffect, useState } from "react"
import { UserContext } from "../lib/context"
import { firestore } from "../lib/firebase";
import PageCard from "./pageCard";

export default function PagesView() {

    const { user } = useContext(UserContext);
    const [ pages, setPages ] = useState([]);

    const deleteCallback = (pageID) => {
        deleteDoc(doc(firestore, "users", user.uid, "pages", pageID))
        .then((stat) => {console.log(stat)})
        .catch((err) => {console.error(err)});
    }

    const unwrapCard = (card) => {
        const data = card.data();
        return(
            <PageCard
                title="Hello, Card!"
                url={data.url}
                documentID={card.id}
                key={card.id}
                deleteCallback={deleteCallback}
                />
        )
    }

    useEffect(() => {
        let unsubscribe;

        if (user) {
            const pagesRef = collection(firestore, "users", user.uid, "pages");
            unsubscribe = onSnapshot(pagesRef, (docs) => {
                if (!docs.empty) {
                    setPages(docs.docs.map(unwrapCard));
                } else {
                    setPages("You have no saved pages")
                }
            })
        }
        return unsubscribe;
    }, [user])

    return(
        <div className="pages-container">
            { pages }
        </div>
    )
}