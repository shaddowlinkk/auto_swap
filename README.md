This a a windoes only program to detect a new version of a plugin then stop the server and gupdate the plugin for the new version then restart the server.
Example config json

{
  "watch": "test.txt", <- file that you are looking for
  "dest": "help.txt", <- where you want to move the file
  "start_server": { <- how to start the server ( this example use a bat file)
    "type": "CMD",
    "command": "run.bat",
    "args": []
  },
  "stop_server": { <- how to stio the server (this is a command to input into the server)
    "type": "STDIN",
    "command": "stop",
    "args": []
  }
}
