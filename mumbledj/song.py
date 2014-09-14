#----------------------#
#       MumbleDJ       #
#  By Matthieu Grieger #
#----------------------#---------------------------#
# song.py                                          #
# Contains definitions of Song class & methods     #
#--------------------------------------------------#

import pafy
from config import *

class Song:
	# Constructs a new Song object with all of the necessary attributes
	# for the given YouTube url. No downloading of audio or thumbnails
	# is done at this stage.
	def __init__(self, youtube_url):
		song = pafy.new(youtube_url)
		self.song_title = song.title
		self.song_duration = song.duration
		self.song_thumbnail = song.thumb
		self.song_audio = song.getbestaudio(preftype = 'ogg')
		self.audio_ready = False
		self.song_skips = 0
	
	# Downloads the audio file that was found in __init__. Download progress
	# callbacks are passed to check_download_status.
	def download_song(self):
		self.song_audio.download(quiet = True, callback = self._check_download_status)
		
	# A callback function that checks the download status of an audio file. Simply
	# sets audio_ready to True when the ratio is 1 (download is completed).
	def _check_download_status(total, recvd, ratio, rate, eta):
		if ratio == 1:
			self.audio_ready = True
	
	# Returns the status of a song download.
	def audio_ready(self):
		return self.audio_ready
	
	# Called after a song is done playing. Audio files have no use to us after the song
	# has finished, so they should just be deleted.
	def delete_song(self):
		if DELETE_AUDIO:
			pass # Delete audio
	
	# Downloads thumnail from location specified by song_thumbnail.
	def download_thumbnail(self):
		if SHOW_CHAT_NOTIFICATIONS and SHOW_YT_THUMBNAILS:
			pass # Download thumbnail
	
	# Called after a song is done playing. Much like audio files, thumbnails are of no
	# use to us after the song has finished.
	def delete_thumbnail(self):
		if DELETE_THUMBNAILS and SHOW_CHAT_NOTIFICATIONS and SHOW_YT_THUMBNAILS:
			pass # Delete thumbnail
	
	# Increments the skip count when a user uses the skip command.
	def increment_skip_count(self):
		self.song_skips = self.song_skips + 1
