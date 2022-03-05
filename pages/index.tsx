import Head from "next/head"
import Image from "next/image"

export default function Home() {
  return (
    <>
    <Head>
      <title>PageMail - a simple Read-it-Later!</title>
      <meta name="description" content="A simplistic, easy to use and free link-saving service." />
    </Head>
    <div className="w-screen max-w-screen-xl md:pt-12 flex flex-col items-start justify-evenly mx-auto flex-wrap">
      <div className="px-3 md:py-12">
        <h1 className="text-7xl text-secondary dark:text-secondary-dark font-serif font-bold my-4">A simple Read-It-Later</h1>
        <h2 className="mt-8 text-tertiary text-lg">PageMail is a <b>simplistic</b>, <b>easy to use</b> and <b>free</b> link-saving service</h2>
      </div>
    </div>
    </>
  )
}
