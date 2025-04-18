# CloudlogAutoLogger
Provide automatic logging to CloudLog for multiple 3rd party applications

--------------------------------------
Install and build
--------------------------------------
Install GO on your computer (https://go.dev/)

Create a directory "build" 

Compile using the GO BUILD command:
**go build -o build/cloudlogautologger.exe ./cmd/main**

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

---------------------------------------
WSJT-X
---------------------------------------


