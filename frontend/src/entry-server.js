import { createApp } from './app'
import { renderToString } from 'vue/server-renderer'

export async function render(url, manifest) {
  const { app, router } = createApp()

  await router.push(url)
  await router.isReady()

  const ctx = {}
  const html = await renderToString(app, ctx)

  return { html }
}
