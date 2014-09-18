#---------------------#
#      MumbleDJ       #
# By Matthieu Grieger #
#---------------------#--------------------------------------------------#
# download_audio.py                                                      #
# Downloads audio (ogg format) from specified YouTube ID. If no ogg file #
# exists, it creates an empty file called .video_fail that tells the Lua #
# side of the program that the download failed. .video_fail will get     #
# deleted on the next successful download.                               #
#------------------------------------------------------------------------#

import pafy
from sys import argv
from os.path import isfile
from os import remove, system
from time import sleep

url = argv[1]
volume = argv[2]
video = pafy.new(url)

try:
	video.oggstreams[0].download(filepath = 'song.ogg', quiet = True)
	if isfile('.video_fail'):
		remove('.video_fail')
except:
	with open('.video_fail', 'w+') as f:
		f.close()
		
while isfile('song.ogg.temp'):
	sleep(1)
	
system('ffmpeg -i song.ogg -ar 48000 -ac 1 -vol ' + str(volume) + ' song-converted.ogg -y')

while not isfile('song-converted.ogg') and not isfile(".video_fail"):
	sleep(1)

remove('song.ogg')
		

