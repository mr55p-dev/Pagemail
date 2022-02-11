// Following instructions from https://devrecipes.net/modal-component-with-next-js/
import Document, { Html, Head, Main, NextScript } from "next/document";

class MainDocument extends Document {
  static async getInitialProps(ctx) {
    const initialProps = await Document.getInitialProps(ctx);
    return { ...initialProps };
  }

  render() {
    return (
      <Html>
        <Head />
        <body>
          <Main />
          <NextScript />
          <div id="modal-root"></div>
          <div id="notif-root"></div>
        </body>
      </Html>
    );
  }
}

export default MainDocument;