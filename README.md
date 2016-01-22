## Twitch Stream, get the stream URL easily


Hello there! `twitch-stream` is an app to retrieve Twitch's `m3u8` streaming URLs by performing a few requests to Twitch servers to determine both the stream and the available qualities. This is useful if you don't like the Flash-based player that Twitch has on their desktop website but instead you want to use an app like VLC to watch the stream —you'll still need to solve how you're going to use the chat though—.

### Download

You can **always download the latest app available** for both Macintosh and Windows by going to the **[Releases page](https://github.com/patrickdappollonio/twitch-stream/releases)**. There you'll find always the latest version available.

### Usage

Open your command line app —on Windows, it's called "Command Prompt" while in Macintosh it's called "Terminal"—, then drag the recently downloaded executable and drop it on the command line. That will show the full path to the executable as text in the command line. After that, just add the corresponding extra parameters, like the Twitch Streamer Username you want to see (here we're going to use `patrickdappollonio`) and the stream quality (here we're going to use `best`), like this:

```
twitch-stream patrickdappollonio best
```

That will open the stream with the selected quality (or the best found if none match the requested quality) available for the user `patrickdappollonio` in the VLC app.

If you're using Macintosh and you can't run the app, you might have to give it the proper permissions by writing `chmod +x twitch-stream` at the folder where you downloaded the executable.

### Adding the app to the Operating System Path

If you want to run the app without prepending the full folder path before, then you need to add the app to your **path** set up on your computer. The **path** its just a list of folders where, if an executable is found, then doesn't need to be called with the full path on it, instead, you can just use the executable name (in our case, `twitch-stream`).

**On Windows**, that's an easy task: you could copy your application to `C:\Windows\` or `C:\Windows\System32` and you're good to go. That will make the app globally available and instead of using the full path trick described on [Usage](#usage), you can just call it without the drag-and-drop part and just run `twitch-stream

**On Macintosh**, you probably already figured that out. But if you don't, simply copy the executable to `/usr/local/bin`, by writing `cp twitch-stream /usr/local/bin`. Then simply execute the app without the drag-and-drop step described on [Usage](#usage).

### Internals

The app is written in Go —that's what you see above, those files ending in `.go`— because it's fun and allows the same application to run on both Windows and Macintosh without changing a line of code.

### Pull Requests?

If you think you can make something better, I'll be glad to accept any pull request. Just send them out!

### Bug reports

Please, open an issue at the [Issue Tracker](https://github.com/patrickdappollonio/twitch-stream/issues). I'll try to do my best to solve any kind of problem you might encounter.
