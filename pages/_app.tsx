import Navbar from '../components/Navbar'
import { UserContext } from '../lib/context'
import { useUserData } from '../lib/hooks';
import '../styles/globals.css'

function MyApp({ Component, pageProps }) {

  const [user, username] = useUserData();

  return(
    <UserContext.Provider value={{ user, username }}>
      <Navbar />
      <Component {...pageProps} />
    </UserContext.Provider>

  )

}

export default MyApp
