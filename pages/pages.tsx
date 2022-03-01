import { AuthCheck } from "../components/AuthCheck"
import PagesView from "../components/pagesView"

export default function PagesRoute() {
    return(
        <main className="bg-sky-50 md:bg-white">
            <AuthCheck>
                <h1 className="page-heading">Your Pages</h1>
                <PagesView />
            </AuthCheck>
        </main>

    )
}