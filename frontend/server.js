import fs from 'node:fs'
import path from 'node:path'
import { fileURLToPath } from 'node:url'
import express from 'express'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const isProd = process.env.NODE_ENV === 'production'
const port = process.env.PORT || 3000
const apiUrl = process.env.API_URL || 'http://localhost:8080'

async function createServer() {
  const app = express()

  let vite
  if (!isProd) {
    const { createServer: createViteServer } = await import('vite')
    vite = await createViteServer({
      server: { middlewareMode: true },
      appType: 'custom',
    })
    app.use(vite.middlewares)
  } else {
    app.use(express.static(path.resolve(__dirname, 'dist/client'), { index: false }))
  }

  // Proxy /api requests to Go backend
  app.use('/api', express.raw({ type: '*/*', limit: '10mb' }), async (req, res) => {
    const url = `${apiUrl}${req.originalUrl}`
    try {
      const hasBody = !['GET', 'HEAD'].includes(req.method) && req.body?.length > 0
      const response = await fetch(url, {
        method: req.method,
        headers: {
          'Content-Type': 'application/json',
          ...(req.headers.authorization ? { authorization: req.headers.authorization } : {}),
        },
        body: hasBody ? req.body : undefined,
      })
      res.status(response.status)
      const text = await response.text()
      res.set('Content-Type', 'application/json').end(text)
    } catch (e) {
      console.error('API proxy error:', e)
      res.status(502).json({ error: 'API unavailable' })
    }
  })

  app.use('*', async (req, res) => {
    const url = req.originalUrl

    try {
      let template, render

      if (!isProd) {
        template = fs.readFileSync(path.resolve(__dirname, 'index.html'), 'utf-8')
        template = await vite.transformIndexHtml(url, template)
        render = (await vite.ssrLoadModule('/src/entry-server.js')).render
      } else {
        template = fs.readFileSync(path.resolve(__dirname, 'dist/client/index.html'), 'utf-8')
        render = (await import('./dist/server/entry-server.js')).render
      }

      const { html: appHtml } = await render(url)
      const html = template.replace('<!--app-html-->', appHtml)

      res.status(200).set({ 'Content-Type': 'text/html' }).end(html)
    } catch (e) {
      vite?.ssrFixStacktrace(e)
      console.error(e)
      res.status(500).end(e.message)
    }
  })

  app.listen(port, () => {
    console.log(`Frontend SSR server: http://localhost:${port}`)
  })
}

createServer()
