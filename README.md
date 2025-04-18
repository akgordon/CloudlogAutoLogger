# CloudlogAutoLogger
Provide automatic logging to CloudLog for multiple 3rd party applications

This program and source my be freely copied and used as you wish. There is NO warranty and your use is "AS IS". 
Please contact me if you have suggests or additions you would like to make.

Hope this helps your logging to the CloundLog application. 

Alan N7AKG  (N7AKG@ARRL.net)

--------------------------------------
Install and build
--------------------------------------
Install GO on your computer (https://go.dev/)
The GO language is supported on most PC OSs and
this program does not need any UI so should work well in any environment.

Suggested compile steps:
   1 Create a directory "build" 
   2 Compile using the GO BUILD command:
      **go build -o build/cloudlogautologger.exe ./cmd/main**

A Windows executable can be downloaded from "http://www.akgordon.com/ham/cloudlogautologger.zip"

---------------------------------------
Usage
---------------------------------------

1 Copy the executable to any directory of your choice

2 Run the program to set parameters. You will need:
   - A Cloudlog API key
   - A Station profile id - This can be found when editing a station profile on CloudLog. its a number and displayed in the URL string.
   - Your CloudLog URL. For example "https://xxxx.cloudlog.co.uk/"
   - The UDP port number for WSJT-X  (optional)
   - The UDP port number for JS8CALL  (optional)
   - The UDP port number for VARAC   (optional)

3 Select the "S" option.
   - This will run through a list of options for you to respond
   - The result is an xxxxx.ini file
   - The xxxx.ini file is a simple text file, but the API key is encrypted so do not edit that.

4 Now you can either run the program by selecting the "R" option, or have it run automatically by using command line arguments.
   - Command line use example:  "cloundlogautologger run"   or "cloundlogautologger run log"  (if want logging to file)

5 To stop program
   - If in interactive mode then press "ENTER" at anytime in the program's window
   - If using command line mode then ctrl-C or just close the program's window


In you 3rd party programs, set your UDP server to local host address "127.0.0.1" and UDP port to any value you want. 
Since UDP is broadcast and not point-to-point you can share the port ID with any other loggers 
you have on system.

You must use a different port for each 3rd party program in this program but you can try using just one for all programs and see if works.
If you do use just one then only set the port for WSJT-X and set the others to 0 when doing setup.

---------------------------------------
WSJT-X
---------------------------------------
Set address amd port in menu "settings" --> "Reporting"  

---------------------------------------
JS8Call
---------------------------------------
Set address amd port in Menu "File" --> "Settings" --> "Reporting"
Enable the N1MM logger option at bottom page and enter address and port number.

---------------------------------------
VarAC
---------------------------------------
In menu "Settings" --> "Rig control ...." --> "Logging"
- Set "Send log to" --> "N1MM(UDP)"
- Set the IP address and port number
