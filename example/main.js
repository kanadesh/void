// Interaction with the backend takes place in the `application` channel of BroadcastChannel.
const channel = new BroadcastChannel('application')

const fps = 60
const resolution = 1000

// All metadata should be set up without overextending or under-extending any elements.
const option = {
    fps,
    frames: 300,
    width: resolution,
    height: resolution,
    cacheDir: '.cache',
    resultFile: 'video.mp4',
    number: 8, // Number of pages to run in parallel.
    audios: [
        {
            link: '/example/warp.mp3',
            start: 0
        }
    ]
}

const application = document.getElementById('application')
let frame = 0

channel.addEventListener('message', ({ data }) => {
    switch (
        data.action // The `data.action` sent by the backend contains the string 'request' or 'load'. The build proceeds by returning 'response' and 'ok' respectively.
    ) {
        // It is executed at the very first rendering of the application.
        // Time-consuming actions such as preparing data for the application and fetching from external sources are performed here to speed up the build.
        case 'request':
            application.width = resolution
            application.height = resolution

            // When a 'response' is returned in response to a 'request', the build is initiated.
            // Add an associative array of settings to 'body'.
            channel.postMessage(
                {
                    action: 'response',
                    body: option
                },
                '*'
            )
            break

        // A 'load' is performed for every frame screenshot.
        // The rendering proceeds by simply returning 'ok' to this.
        case 'load':
            frame = data.body.frame
            application.textContent = `Now frame: ${frame}`
            channel.postMessage(
                {
                    action: 'ok',
                    body: null
                },
                '*'
            )
            break

        default:
            break
    }
})
