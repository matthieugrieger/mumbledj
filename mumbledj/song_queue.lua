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

SongQueue = {}

function SongQueue.addSong(url, username)
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
			print("YouTube URL is valid!")
			getYoutubeInfo(video_id, username)
		end
	end
end

function getYoutubeInfo(id, username)
	if id == nil then
		return
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
		return
	end
	
	print("Finished getting info.")
	youtubeInfoCompleted({
		id = id,
		title = name,
		duration = string.format("%d:%02d", duration / 60, duration % 60),
		thumbnail = thumbnail,
		username = username
	})
end

function youtubeInfoCompleted(info)
	if info == nil then
		return false
	end
	
	song_queue:push_right(info)
	
	if song_queue:length() == 1 then
		os.execute("python download_audio.py " .. info.id .. " " .. config.VOLUME)
		while not file_exists("song-converted.ogg") do
			os.execute("sleep " .. tonumber(2))
		end
		piepan.me.channel:play("song-converted.ogg", nextSong)
	end
	
	if piepan.Audio:isPlaying() then
		local message = string.format(config.NOW_PLAYING_HTML, info.thumbnail, info.id, info.title, info.duration, info.username)
		piepan.me.channel:send(message)
	else
		local message = string.format(config.SONG_ADDED_HTML, info.username, info.title)
		piepan.me.channel:send(message)
	end
	
	return true
end

function SongQueue.getNextSong(url)
	
end

function SongQueue.getLength()
	return song_queue:length()
end

function file_exists(file)
	local f=io.open(file,"r")
	if f~=nil then io.close(f) return true else return false end
end

return SongQueue
