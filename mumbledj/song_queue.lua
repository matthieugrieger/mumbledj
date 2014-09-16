-------------------------
--      MumbleDJ       --
-- By Matthieu Grieger --
-------------------------------------------------
-- song_queue.lua                              --
-- Contains the definition of the song queue   --
-- used for queueing up songs.                 --
-------------------------------------------------
local deque = require("deque")

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
			piepan.Thread.new(getYoutubeInfo, youtubeInfoCompleted, {video_id, username})
			return true
		end
	end
	
	return false
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
	
	return {
		id = id,
		title = name,
		duration = string.format("%d:%02d", duration / 60, duration % 60),
		thumbnail = thumbnail,
		username = username
	}
end

function youtubeInfoCompleted(info)
	if info == nil then
		return
	end
	
	song_queue.push_left(info)
	
	if song_queue:length() == 1 then
		os.execute("python download_audio.py")
		piepan.me.channel:play("song.ogg")
	end
	
	if piepan.Audio:isPlaying() then
		local message = string.format(config.NOW_PLAYING_HTML, info.thumbnail, info.id, info.title, info.duration, info.username)
		piepan.me.channel:send(message)
	else
		local message = string.format(config.SONG_ADDED_HTML, info.username, info.title)
		piepan.me.channel:send(message)
	end
end

function SongQueue.getNextSong(url)
	
end

function SongQueue.getLength()
	return song_queue:length()
end

return SongQueue
