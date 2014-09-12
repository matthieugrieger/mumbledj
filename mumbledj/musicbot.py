#----------------------#
#   Mumble Music Bot   #
#  By Matthieu Grieger #
#----------------------#---------------------------#
# musicbot.py                                      #
# Contains definitions of musicbot class & methods #
#--------------------------------------------------#
#import pymumble here
from config import *

class MusicBot:
	# Since all the configuration is set in config.py, we don't really need to do anything here.
	def __init__(self):
		print('Starting up ' + SERVER_USERNAME + '...') 
	
	# Connects to the Mumble server with the credentials specified upon object creation.	
	def connect_to_server(self):
		self.mumble = pymumble.Mumble(SERVER_ADDRESS, SERVER_PORT, SERVER_USERNAME, SERVER_PASSWORD, debug = DEBUG)
		self.mumble.start()
		self.mumble.is_ready()
		self.mumble.channels.find_by_name(DEFAULT_CHANNEL).move_in()
		self.mumble.users.myself.mute()
	
	# Starts to play the first song in the queue when called. If no songs exist in the queue, it will wait until
	# a song is added.	
	def start_music(self):
		pass
		
	# Resumes music when called.
	def play_music(self):
		pass
		
	# Pauses music until told to resume.
	def pause_music(self):
		pass
		
	# Adds a YouTube link to the queue along with its metadata.
	def add_to_queue(self):
		pass
		
	# Sends a message to the chat when a new song starts playing.	
	def announce_new_song(self):
		pass
		
	# Raises the volume by the increment decided by the user. 
	def raise_volume(self, increment):
		pass
	
	# Lowers the volume by the increment decided by the user.	
	def lower_volume(self, increment):
		pass
		
	# Moves bot from the current channel to the specified channel.
	def move_bot(self, channel):
		pass
		
	# Completely stops the bot if a person on the approved list of admins
	# issues the stop command.
	def stop_bot(self):
		pass
		
	
	
	
