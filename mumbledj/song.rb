# MumbleDJ v2
# By Matthieu Grieger
# song.rb

# Base Song class that defines default behavior for any kind of song.
class Song

  # Starts the song.
  def start
  
  end
  
  # Gets the name of the user who submitted the song.
  def get_submitter
    return @submitter
  end
  
  # Adds a skipper to the skips array for the current song.
  def add_skip?(username)
    if not @skips.include?(username)
      @skips << username
      return true
    else
      return false
    end
  end
  
  # Determines if a skip should occur. Returns true if a skip is needed,
  # false otherwise.
  def skip_now?(total_users)
    return (total_users / @skips.count) >= SKIP_RATIO
  end
end

class YouTubeSong < Song
  
  # Initializes the YouTubeSong object and retrieves the song title,
  # duration, and thumbnail URL from the YouTube API.
  def initialize(url, submitter)
    @url = url
    @submitter = submitter
    @skips = []
    # TODO: Retrieve YouTube information
    @song_title = ""
    @song_duration = ""
    @song_thumbnail_url = ""
  end
end
