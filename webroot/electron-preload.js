import { contextBridge } from 'electron'
import { spawn } from 'child_process'
import { fileURLToPath } from 'url'
import path from 'path'

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)

const proc = spawn(path.join(__dirname, '..', 'fractgen'), ['serve', '--listen', '127.0.0.1:12345'])
proc.stdout.on('data', (data) => {
  console.log(`${data.toString()}`)
})
proc.on('close', (code) => {
  console.log(`child process exited with code ${code}`)
})

// contextBridge.exposeInMainWorld('versions', {
//   node: () => process.versions.node,
//   chrome: () => process.versions.chrome,
//   electron: () => process.versions.electron,
//   // we can also expose variables, not just functions
// })
