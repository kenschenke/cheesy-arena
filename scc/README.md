SCC Installation and Configuration
============
The field uses three Raspberry Pi boxes to provide some of the functionality of a FIRST field.
FIRST calls these SCC (Station Control Cabinets), so that is the name adopted here.

These boxes monitor and report on the emergency stop buttons at each driver station, control
field lights using FadeCandy light controllers, monitor the emergency stop button on the
scoring table, and provide the scorekeeper / FTA with buttons to control the field lights.

## Installation

1. Prepare a Raspberry Pi 3b+ or 4 with a default installation of Raspian.
1. Connect the Raspberry Pi to the field Ethernet.  WiFi is not recommended for reliability.  The Pi must be connected to VLAN 100.
1. (Optional but recommended) configure the Pi with a static IP.
1. (Optional but recommended) enable SSH using rasp-config
1. Add these two lines to the /boot/config.txt file:

`# Enable the power/activity LED`

`enable_uart=1`

6. Install Node.js using the instructions here: https://www.w3schools.com/nodejs/nodejs_raspberrypi.asp
7. Create a directory for the scc JavaScript files:

`mkdir -p ~/scc/logs`

8. Copy scc.js and launcher.sh to the ~/scc directory
9. Install the required node modules using these commands:

`cd ~/scc`

`npm install ws onoff`

10. Create the file in the /boot/scc that tells scc.js which box it is running on.  The file needs to contain one line with the word red, blue, or scoring.

`sudo nano /boot/scc`

11. Set the SCC JavaScript file to run automatically on boot by adding the following line to the CronTab file.  First, type:

`sudo crontab -e`

Then add this line to the bottom of the file

`@reboot sh /home/pi/scc/launcher.sh >/home/pi/scc/logs/crontab 2>&1`

12. If running field lights using FadeCandy, install it with instructions at https://learn.adafruit.com/1500-neopixel-led-curtain-with-raspberry-pi-fadecandy/fadecandy-server-setup

13. (Optional) The fcserver.json is the configuration we use on our FadeCandy servers.

## Field Hardware

The Raspberry Pi's are placed inside custom cases that provide power and connection status LEDs as well as connection points for the emergency stop buttons.  CAD files for the cases are found here: https://cad.onshape.com/documents/43c157a9e200950f05fc2766/w/f1a71a60f29cc2de381a2c21/e/21e196e834ff44fcab3826cb

### Emerency Stop Buttons

We use these emergency stop buttons from Amazon:  https://amz.run/4Ob0

We used Ethernet cables to attach the stop buttons to the Raspberry Pi.  Ethernet cables are designed for low signal attenuation over very long distances, are widely available, and have robust connectors.

You can find panel mount Ethernet jacks here: https://amz.run/4Ob1
These fit perfectly in the holes on the emergency stop button boxes.