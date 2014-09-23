-------------------------
--      MumbleDJ       --
-- By Matthieu Grieger --
-------------------------------------------------
-- song_queue.lua                              --
-- Contains the definition of the song queue   --
-- used for queueing up songs.                 --
-------------------------------------------------
local deque = require("deque")
local config = require("config")

local song_queue = deque.new()
local skippers = {}

SongQueue = {}

-- Begins the process of adding a new song to the song queue.
function SongQueue.add_song(url, username)
	local patterns = {
		"https?://www%.youtube%.com/watch%?v=([%d%a_%-]+)",
		"https?://youtube%.com/watch%?v=([%d%a_%-]+)",
		"https?://youtu.be/([%d%a_%-]+)",
		"https?://youtube.com/v/([%d%a_%-]+)",
		"https?://www.youtube.com/v/([%d%a_%-]+)"
	}
	
	for _,pattern in ipairs(patterns) do
		local video_id = string.match(url, pattern)
		if video_id ~= nil and string.len(video_id) < 20 then
			return get_youtube_info(video_id, username)
		else
			return false
		end
	end
end

-- Retrieves the metadata for the specified YouTube video via the gdata API.
function get_youtube_info(id, username)
	if id == nil then
		return false
	end
	local cmd = [[
		wget -q -O - 'http://gdata.youtube.com/feeds/api/videos/%s?v=2&alt=jsonc' |
		jshon -Q -e data -e title -u -p -e duration -u -p -e thumbnail -e hqDefault -u
	]]
	local jshon = io.popen(string.format(cmd, id))
	local name = jshon:read()
	local duration = jshon:read()
	local thumbnail = jshon:read()
	if name == nil or duration == nil then
		return false
	end
	
	return youtube_info_completed({
		id = id,
		title = name,
		duration = string.format("%d:%02d", duration / 60, duration % 60),
		thumbnail = thumbnail,
		username = username
	})
end

-- Notifies the channel that a song has been added to the queue, and plays the
-- song if it is the first one in the queue.
function youtube_info_completed(info)
	if info == nil then
		return false
	end
	
	song_queue:push_right(info)
	
	local message = string.format(config.SONG_ADDED_HTML, info.username, info.title)
	piepan.me.channel:send(message)
	
	if not piepan.Audio.isPlaying() then
		return SongQueue.get_next_song()
	end
	
	return true
end

-- Deletes the old song and begins the process of retrieving a new one.
function SongQueue.get_next_song()
	SongQueue.reset_skips()
	if file_exists("song-converted.ogg") then
		os.remove("song-converted.ogg")
	end
	if song_queue:length() ~= 0 then
		local next_song = song_queue:pop_left()
		return start_song(next_song)
	end
end

-- Downloads/encodes the audio file and then begins to play it.
function start_song(info)
	os.execute("python download_audio.py " .. info.id .. " " .. config.VOLUME)
	while not file_exists("song-converted.ogg") do
		os.execute("sleep " .. tonumber(2))
	end
	if not file_exists(".video_fail") then
		if piepan.Audio:isPlaying() then
			piepan.Audio:stop()
		end
		piepan.me.channel:play("song-converted.ogg", SongQueue.get_next_song)
	else
		return false
	end
	
	if piepan.Audio:isPlaying() then
		local message = string.format(config.NOW_PLAYING_HTML, info.thumbnail, info.id, info.title, info.duration, info.username)
		piepan.me.channel:send(message)
		return true
	end
	
	return false
end

-- Adds the username of a user who requested a skip. If their name is
-- already in the list nothing will happen.
function SongQueue.add_skip(username)
	local already_skipped = false
	for name,_ in pairs(skippers) do
		if name == username then
			already_skipped = true
		end
	end
	if not already_skipped then
		table.insert(skippers, username)
		return true
	end
	return false
end

-- Counts the number of users who would like to skip the current song and
-- returns it.
function SongQueue.count_skippers()
	local skipper_count = 0
	for name,_ in pairs(skippers) do
		skipper_count = skipper_count + 1
	end
	return skipper_count
end

-- Resets the list of users who would like to skip a song. Called during a
-- transition between songs.
function SongQueue.reset_skips()
	skippers = {}
end

-- Retrieves the length of the song queue and returns it.
function SongQueue.get_length()
	return song_queue:length()
end

-- Checks if a file exists.
function file_exists(file)
	local f=io.open(file,"r")
	if f~=nil then io.close(f) return true else return false end
end

return SongQueue
