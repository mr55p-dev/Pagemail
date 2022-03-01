import Head from "next/head"

export default function Home() {
  return (
    <>
    <Head>
      <title>PageMail - a simple Read-it-Later!</title>
    </Head>
    <div className="w-screen h-screen flex flex-col items-center justify-around">
      <div className="">
        <h1 className="text-7xl text-sky-700 font-serif">Welcome to PageMail</h1>
        <h2 className="mt-2 text-sky-700">The lightest read-it-later you can imagine!</h2>
      </div>
    </div>
    </>
  )
}
