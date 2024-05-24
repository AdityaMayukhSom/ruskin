import { defineConfig } from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  base: '/ruskin-docs/',
  title: "Ruskin Docs",
  description: "Ruskin is a modern distributed open source message queuing system built with Go. This website hosts the documentation for Ruskin.",
  lastUpdated: true,
  cleanUrls: true,
  metaChunk: true,
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: 'Home', link: '/' },
      // { text: 'Examples', link: '/markdown-examples' }
    ],
    sidebar: [
      {
        // text: 'Ruskin',
        items: [
          { text: 'Markdown Examples', link: '/markdown-examples' },
          { text: 'Runtime API Examples', link: '/api-examples' },
          { text: 'Producer', link: '/producer' },
          { text: 'Message Queue', link: '/message-queue' },
          { text: 'Queue Identifier Map', link: '/queue-identifier-map' },
        ]
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/AdityaMayukhSom/ruskin' }
    ]
  }, srcDir: "./docs"
})
