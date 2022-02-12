import { AuthCheck } from "../components/AuthCheck"
import PagesView from "../components/pagesView"

export default function PagesRoute() {
    return(
        <main>
            <AuthCheck>
                <h1 className="heading">Your Pages</h1>
                <PagesView />
            </AuthCheck>
        </main>

    )
}