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

url = argv[1]
video = pafy.new(url)

def encode_file(stream, downloaded, ratio, rate, eta):
	if ratio == 1:
		print('Encoding!')
		system("ffmpeg -i song.ogg -ar 48000 -ac 1 song-converted.ogg -y")

try:
	video.oggstreams[0].download(filepath = "song.ogg", quiet = True, callback = encode_file)
	if isfile(".video_fail"):
		remove(".video_fail")
except:
	with open(".video_fail", "w+") as f:
		f.close()
		

