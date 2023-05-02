module.exports = {
  reactStrictMode: true,
  env: {
	  ENV: process.env.NODE_ENV,
	  USE_EMULATOR: process.env.USE_EMULATOR,
	  VERCEL_ENV: process.env.VERCEL_ENV
  }
}
