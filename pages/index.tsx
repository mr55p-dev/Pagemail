import Head from "next/head"
import Image from "next/image"

export default function Home() {
  return (
    <>
    <Head>
      <title>PageMail - a simple Read-it-Later!</title>
    </Head>
    <div className="w-screen max-w-screen-xl pt-12 flex flex-col items-start justify-evenly mx-auto flex-wrap ">
      <div className="p-3">
        <h1 className="text-7xl text-sky-700 font-serif">Welcome to PageMail</h1>
        <h2 className="mt-2 text-sky-700">The lightest read-it-later you can imagine!</h2>
      </div>
      <div className="flex ">
        <img src={"/full-browser-img.png"} className="h-96"/>
        <img src={"/pagemail-img.png"} className="h-96"/>
      </div>
    </div>
    </>
  )
}
