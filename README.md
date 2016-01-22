## Twitch Stream, get the stream URL easily

Hello there! `twitch-stream` is a replacement app, similar to `livestreamer`, to watch Twitch's live streams on VLC, on Windows, Mac and Linux. This is always useful since Twitch's Player currently works by wrapping the stream with a flash-player one. Flash is well known to be a battery-hungry application and it will definitely make your battery unhappy if you're watching streams on the run. Also, `twitch-streams` can be used in conjunction with Chatty: you can chat and watch streams at the same time. Great!

### Download

You can **always download the latest version** for Windows, Macintosh and Linux by going to the **[Releases page](https://github.com/patrickdappollonio/twitch-stream/releases)**. For easy and convenience, I offer a Windows Installer so you don't need to worry about where to place the application.

### Usage

Open your command line app —on Windows, it's called "Command Prompt" while in Macintosh it's called "Terminal"—, and depending on your platform and how you installed the app, just execute it by calling it with the defined parameters. If you downloaded the Windows Installer, just execute the command below. If you downloaded the standalone executables, you maybe added it to your system path or you have them in a folder. In any case, call the command with one required parameter —the stream name you want to watch— and one optional, the stream quality, like this:

```
twitch-stream patrickdappollonio [quality]
```

That will open my stream —if I'm streaming— with the selected quality (or the "best" found if there was no other quality available) for the user `patrickdappollonio` in the VLC app.

If you're using Macintosh and you can't run the app, you might have to give it the proper permissions by writing `chmod +x twitch-stream` at the folder where you downloaded the executable.

### Adding the app to the Operating System Path

If you want to run the app without prepending the full folder path before, then you need to add the app to your **path** set up on your computer. The **path** its just a list of folders where, if an executable is found, then doesn't need to be called with the full path on it, instead, you can just use the executable name (in our case, `twitch-stream`).

**On Windows**, you can download the installer and have it automatically set up for you. Alternatively, if you want to add the executable all by yourself, then that's an easy task: you could copy your application to `C:\Windows\` or `C:\Windows\System32` and you're good to go. That will make the app globally available and instead of using the full path trick described on [Usage](#usage), you can just call it without the drag-and-drop part and just run `twitch-stream`.

**On Macintosh and Linux**, you probably already figured that out. But if you don't, simply copy the executable to `/usr/local/bin`, by writing `cp twitch-stream /usr/local/bin`. Then simply execute the app like any other command line app. Alternatively you can run the app by dragging-and-dropping the executable to the Terminal Window with the needed parameters too.

### Internals

The app is written in Go —that's what you see above, those files ending in `.go`— because it's fun and allows the same application to run on both Windows and Macintosh without changing a line of code.

### Do you accept Pull Requests?

Absolutely! If you think you can make something better, I'll be glad to accept any pull request. Just send them out!

### Bug reports / Suggestions / Improvements?

Please, open an issue at the [Issue Tracker](https://github.com/patrickdappollonio/twitch-stream/issues). I'll try to do my best to solve any kind of problem you might encounter. If it's a suggestion or improvement, then be my guest and file an issue there as well. Any suggestions are always welcome.