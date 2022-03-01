/** @type {import('next-sitemap').IConfig} */

module.exports = {
    siteUrl: process.env.SITE_URL || 'https://www.pagemail.io',
    generateRobotsTxt: true,
    robotsTxtOptions: {
        policies: [
          {
            userAgent: '*',
            allow: '/',
            disallow: '/api'
          }
        ]
      },
  }