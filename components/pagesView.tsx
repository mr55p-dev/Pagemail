import { collection, deleteDoc, doc, DocumentSnapshot, getFirestore, onSnapshot } from "firebase/firestore";
import { useContext, useEffect, useState } from "react"
import { useAuth } from "../lib/context";
import { firestore } from "../lib/firebase";
import { usePageMetadata, useUserToken } from "../lib/hooks";
import { ICard, IPage, IPageMetadata } from "../lib/typeAliases";
import { AuthCheck } from "./AuthCheck";
import PageCard from "./pageCard";

export default function PagesView() {

    const { user } = useAuth();
    const [ pages, setPages ] = useState<IPage[]>([]);
    const [ metas, setMetas ] = useState<IPageMetadata[]>([]);
    const [ scrapedPages, setScrapedPages ] = useState<IPage[]>([]);
    const token = useUserToken()

    const deleteCallback = (pageID: string): void => {
        deleteDoc(doc(firestore, "users", user.uid, "pages", pageID))
        .then((stat) => {console.log(stat)})
        .catch((err) => {console.error(err)});
    }

    // Page collector
    useEffect(() => {
        if (user) {
            const pagesRef = collection(firestore, "users", user.uid, "pages");
            const unsubscribe = onSnapshot(pagesRef, (docs) => {
                if (!docs.empty) {
                    setPages(docs.docs.map((card) => {
                        const data = card.data() as IPage;
                        return {
                            id: card.id,
                            ...data
                        }
                    }));
                } else {
                    setPages([])
                }
            })
            return unsubscribe;
        }
    }, [user])


    const collectMeta = async (page: IPage): Promise<IPageMetadata> => {
      // Get the API address
      const apiAddress = new URL(window.location.origin)

      // Modify the path and query parameters
      apiAddress.pathname = "/api/scrape";
      apiAddress.searchParams.set("url", encodeURIComponent(page.url.toString()))

      // Get a response
      const response = await fetch(apiAddress.toString(), {
        method: "GET",
        mode: "same-origin",
        credentials: "same-origin",
        headers: {
            token: token
        },
      })
      if (!response.ok) {
          throw new Error("Error collecting the metadata")
      }
      const body = await response.json()
      return {
          url: page.url,
          title: body.title,
          author: body.author,
          description: body.description,
          image: body.image
      }
    }

    // Metadata scraper
    useEffect(() => {
        if (!pages) {
            return
        }

        const metaPromises = pages.map(collectMeta)
        Promise.allSettled(metaPromises)
        .then((results) => {
            return results.map((result) => {
                if (result.status === "fulfilled") {
                    return result.value
                }
            })
        })
        .then((collectedPages) => {
            setMetas(collectedPages)
        })
    }, [pages])


    // Stiches the page and metadata
    useEffect(() => {
        setScrapedPages(pages.map((page) => {
            const matches = metas.filter((m) => m?.url === page.url)
            if (matches) {
                return {
                    ...page,
                    metadata: matches[0]
                }
            }
        }))
    }, [pages, metas])

    return(
        <AuthCheck>
            <div className="grid auto-rows-max pb-4">
                { pages ?
                scrapedPages.sort((p, q) => q.timeAdded - p.timeAdded).map((d: ICard) => <PageCard data={d} deleteCallback={deleteCallback} key={d.id} />) :
                "You have no saved pages." }
            </div>
        </AuthCheck>
    )
}