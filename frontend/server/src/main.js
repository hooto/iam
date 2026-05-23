import { mount } from 'svelte'
import './app.css'
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap/dist/js/bootstrap.js'
import App from './App.svelte'

const el = document.getElementById('app')
// clear loading placeholder
if (el) el.innerHTML = ''

const app = mount(App, { target: /** @type {Element} */ (el) })

export default app