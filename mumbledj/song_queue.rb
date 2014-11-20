# MumbleDJ v2
# By Matthieu Grieger
# song_queue.rb

require_relative "song"

# A specialized SongQueue class that handles queueing/unqueueing songs
# and other actions.
class SongQueue
  
  # Initializes a new song queue.
  def initialize
    @queue = []
  end
  
  # Checks if song already exists in the queue, and adds it if it doesn't
  # already exist.
  def add_song?(url, submitter)
    # TODO: Determine which kind of URL is given (probably using regex),
    # and instantiate the correct Song object.
    # Example for a YouTube URL given below.
    
    if @queue.empty?
      song = YouTubeSong.new(url, submitter)
      @queue.push(song)
    else
      @queue.each do |song|
        if song.url == url
          return false
        end
      end
      song = YouTubeSong.new(url, submitter)
      @queue.push(song)
    end
  end
  
  # Processes a song delete request. Searches the queue for songs with
  # titles containing the keyword. If found, the song is deleted if the
  # username of the user who requested the deletion matches the
  # username of who originally added the song.
  def delete_song?(keyword, username)
    if not @queue.empty?
      @queue.each do |song|
        if song.song_title.includes?(keyword)
          if song.get_submitter == username
            @queue.delete(song)
            return true
          end
        end
      end
      return false
    else
      return false
    end
    
  end
  
  # Returns a formatted string that contains information about the next
  # song in the queue.
  def peek_next
  
  end
  
  # Returns the current Song object from the queue.
  def get_current_song
    return @queue[0]
  end
end
