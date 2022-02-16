import { collection, deleteDoc, doc, DocumentSnapshot, getFirestore, onSnapshot } from "firebase/firestore";
import { useContext, useEffect, useState } from "react"
import { UserContext } from "../lib/context"
import { firestore } from "../lib/firebase";
import { useUserToken } from "../lib/hooks";
import { IPage } from "../lib/typeAliases";
import { AuthCheck } from "./AuthCheck";
import PageCard from "./pageCard";

export default function PagesView() {

    const { user } = useContext(UserContext);
    const [ pages, setPages ] = useState<JSX.Element[]>([]);
    const token = useUserToken()

    const deleteCallback = (pageID) => {
        deleteDoc(doc(firestore, "users", user.uid, "pages", pageID))
        .then((stat) => {console.log(stat)})
        .catch((err) => {console.error(err)});
    }

    const unwrapCard = (card: DocumentSnapshot) => {
        const data = card.data() as IPage;
        const dateCreated = new Date(1000 * data.timeAdded.seconds)

        return(
            <PageCard
                title="Hello, Card!"
                url={data.url}
                documentID={card.id}
                deleteCallback={deleteCallback}
                dateCreated={dateCreated.toISOString()}
                token={token}
                key={card.id}
                />
        )
    }

    useEffect(() => {
        if (user) {
            const pagesRef = collection(firestore, "users", user.uid, "pages");
            const unsubscribe = onSnapshot(pagesRef, (docs) => {
                if (!docs.empty) {
                    setPages(docs.docs.map(unwrapCard));
                } else {
                    setPages([])
                }
            })
            return unsubscribe;
        }
    })

    return(
        <AuthCheck>
            <div className="pages-container">
                { pages ? pages : "You have no saved pages." }
            </div>
        </AuthCheck>
    )
}