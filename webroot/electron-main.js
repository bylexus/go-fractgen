import { app, BrowserWindow } from 'electron'
import { fileURLToPath } from 'url'
import path from 'path'
import { spawn } from 'child_process'

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)
const fractgenBin = path.join(__dirname, '..', 'fractgen')
const fractgetnDir = path.join(__dirname, '..')
const listen = '127.0.0.1:12345'

console.log(`Start local fractgen server on ${listen}`)
const proc = spawn(fractgenBin, ['serve', '--listen', '127.0.0.1:12345'], {
  cwd: fractgetnDir,
})
proc.stdout.on('data', (data) => {
  console.log(`${data.toString()}`)
})
proc.on('close', (code) => {
  console.log(`child process exited with code ${code}`)
})

const createWindow = () => {
  const win = new BrowserWindow({
    width: 800,
    height: 600,
    title: 'FractGen - alexi.ch',
    webPreferences: {
      preload: path.join(__dirname, 'electron-preload.js'),
    },
  })

  // win.loadFile('dist/index.html')
  win.loadURL(`http://${listen}/dist/index.html`)
  win.maximize()
}

app.whenReady().then(() => {
  // Start window with delay, to let the server start
  setTimeout(() => {
    createWindow()
  }, 1000)
})
