import { IncomingMessage, createServer } from 'http'

createServer(async (req, res) => {
    const body = await readBody(req)
    console.log(body.length)
    res.end(body)
}).listen(9090, () => {
    console.log(`http://127.0.0.1:9090`)
})

/**
 * @param {IncomingMessage} req 
 * @returns {Promise<string>}
 */
const readBody = req => new Promise(resolve => {
    let body = ''
    req.on('data',
        /** @param {Buffer} stream */
        stream => {
            body += stream
        })
    req.on('end', () => {
        resolve(body)
    })
})