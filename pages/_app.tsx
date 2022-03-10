import Head from 'next/head';
import Navbar from '../components/Navbar'
import { AuthProvider } from '../lib/context';
import '../styles/globals.css'

function MyApp({ Component, pageProps }) {
  return(
    <>
    <Head>
      <meta charSet="utf-8"/>
      <meta httpEquiv="X-UA-Compatible" content="IE=edge"/>
      <meta name="viewport" content="width=device-width,initial-scale=1,minimum-scale=1,maximum-scale=1,user-scalable=no"/>
      <meta name="description" content="A simplistic, easy to use and free link-saving service." key="description"/>
      <meta name="keywords" content="Keywords" key="keywords"/>
      <title key="title">PageMail</title>


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
