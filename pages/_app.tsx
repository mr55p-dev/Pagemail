import Head from 'next/head';
import Navbar from '../components/Navbar'
import { AuthProvider } from '../lib/context';
import '../styles/globals.css'

function MyApp({ Component, pageProps }) {
  return(
    <>
    <Head>
      <meta name="theme-color" content="#fff5e0" />
    </Head>
    <AuthProvider>
        <div className="bg-primary dark:bg-primary-dark w-screen min-h-screen">
          <div className="mx-auto max-w-screen-xl">
            <Navbar />
            <Component {...pageProps} />
          </div>
        </div>
    </AuthProvider>
    </>

  )

}

export default MyApp
