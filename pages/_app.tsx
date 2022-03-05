import Head from 'next/head';
import { useState } from 'react';
import Navbar from '../components/Navbar'
import Notif from '../components/notif';
import { NotifContext, UserContext } from '../lib/context';
import { useUserData } from '../lib/hooks';
import { INotifState } from '../lib/typeAliases';
import '../styles/globals.css'

function MyApp({ Component, pageProps }) {

  const userData = useUserData();

  // Refactor this into a hook?
  const [notifShow, setNotifShow] = useState<boolean>(false)
  const [notifState, setNotifState] = useState<INotifState | undefined>(undefined)

  return(
    <>
    <Head>
      <meta name="theme-color" content="#fff5e0" />
    </Head>
    <UserContext.Provider value={ userData }>
        <div className="bg-primary dark:bg-primary-dark w-screen min-h-screen">
          <div className="mx-auto max-w-screen-xl">
            <Navbar />
            <NotifContext.Provider value={{ setNotifShow, setNotifState }}>
              <Component {...pageProps} />
            <Notif show={notifShow} state={notifState} />
                  </NotifContext.Provider>
          </div>
        </div>
    </UserContext.Provider>
    </>

  )

}

export default MyApp
