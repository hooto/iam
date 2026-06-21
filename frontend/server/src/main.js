import { mount } from 'svelte'
import './app.css'
import 'bootstrap/dist/css/bootstrap.css'
import App from './App.svelte'

// Keep the index.html loading placeholder visible for at least 300ms
// to avoid a flash of content on fast initial loads.
const LOADING_DELAY_MS = 300

const app = await new Promise((resolve) => {
  setTimeout(() => {
    const el = document.getElementById('app')
    // clear loading placeholder
    if (el) el.innerHTML = ''
    resolve(mount(App, { target: /** @type {Element} */ (el) }))
  }, LOADING_DELAY_MS)
})

export default app
