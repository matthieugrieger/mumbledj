-------------------------
--      MumbleDJ       --
-- By Matthieu Grieger --
------------------------------------------------------------------
-- mumbledj.lua                                                 --
-- The main file which defines most of MumbleDJ's behavior. All --
-- commands are found here, and most of their implementation.   --
------------------------------------------------------------------

local config = require("config")
local deque = require("deque")

-- Connects to Mumble server.
function piepan.onConnect()
	print(piepan.me.name .. " has connected to the server!")
	local user = piepan.users[piepan.me.name]
	local channel = user.channel(config.DEFAULT_CHANNEL)
	if channel == nil then
		print("The channel '" .. config.DEFAULT_CHANNEL .. "' does not exist. Moving bot to root of server...")
		channel = piepan.channels[0]
	end
	piepan.me:moveTo(channel)
end

-- Function that is called when a new message is posted to the channel.
function piepan.onMessage(message)
	if message.user == nil then
		return
	end

	if string.sub(message.text, 0, 1) == config.COMMAND_PREFIX then
		parse_command(message)
	end
end

-- Parses commands and its arguments (if they exist), and calls the appropriate
-- functions for doing the requested task.
function parse_command(message)
	local command = ""
	local argument = ""
	if string.find(message.text, " ") then
		command = string.sub(message.text, 2, string.find(message.text, ' ') - 1)
		argument = string.sub(message.text, string.find(message.text, ' ') + 1)
	else
		command = string.sub(message.text, 2)
	end
	
	-- Add command
	if command == config.ADD_ALIAS then
		local has_permission = check_permissions(config.ADMIN_ADD, message.user.name)
		
		if has_permission then
			if config.OUTPUT then 
				print(message.user.name .. " has added a song to the queue.")
				if not add_song(argument, message.user.name) then
					message.user:send(config.INVALID_URL_MSG)
				end
			end
		else
			message.user:send(config.NO_PERMISSION_MSG)
		end
	-- Skip command
	elseif command == config.SKIP_ALIAS then
		local has_permission = check_permissions(config.ADMIN_SKIP, message.user.name)
		
		if has_permission then
			if config.OUTPUT then 
				print(message.user.name .. " has voted to skip the current song.")
			end
			
			skip(message.user.name)
		else
			message.user:send(config.NO_PERMISSION_MSG)
		end
	-- Volume command
	elseif command == config.VOLUME_ALIAS then
		local has_permission = check_permissions(config.ADMIN_VOLUME, message.user.name)
		
		if has_permission then
			if config.OUTPUT then
				print(message.user.name .. " has changed the volume to the following: " .. argument .. ".")
				if argument ~= nil then
					if config.LOWEST_VOLUME <= tonumber(argument) and tonumber(argument) <= config.HIGHEST_VOLUME then
						config.VOLUME = tonumber(argument)
						message.user:send(string.format(config.VOLUME_SUCCESS, argument))
					else
						message.user:send(config.NOT_IN_VOLUME_RANGE)
					end
				else
					message.user:send(config.NO_ARGUMENT)
				end
			end
		end
	-- Move command
	elseif command == config.MOVE_ALIAS then
		local has_permission = check_permissions(config.ADMIN_MOVE, message.user.name)
		
		if has_permission then
			if config.OUTPUT then 
				print(message.user.name .. " has told the bot to move to the following channel: " .. argument .. ".")
			end
			if not move(argument) then
				message.user:send(config.CHANNEL_DOES_NOT_EXIST_MSG)
			end
		else
			message.user:send(config.NO_PERMISSION_MSG)
		end
	-- Kill command
	elseif command == config.KILL_ALIAS then
		local has_permission = check_permissions(config.ADMIN_KILL, message.user.name)
		
		if has_permission then
			if config.OUTPUT then 
				print(message.user.name .. " has told the bot to kill itself.")
			end
			kill()
		else
			message.user:send(config.NO_PERMISSION_MSG)
		end
	else
		message.user:send("The command you have entered is not valid.")
	end
end

-- Handles a skip request through the use of helper functions found within
-- song_queue.lua. Once done processing, it will compare the skip ratio with
-- the one defined in the settings and decide whether to skip the current song
-- or not.
function skip(username)
	if add_skip(username) then
		local skip_ratio = count_skippers() / count_users()
		piepan.me.channel:send(string.format(config.USER_SKIP_HTML, username))
		if skip_ratio > config.SKIP_RATIO then
			piepan.me.channel:send(config.SONG_SKIPPED_HTML)
			piepan.Audio:stop()
			next_song()
		end
	end
end

-- Moves the bot to the channel specified by the "chan" argument.
-- NOTE: This only supports moving to a sibling channel at the moment.
function move(chan)
	local user = piepan.users[piepan.me.name]
	local channel = user.channel("../" .. chan)
	if channel == nil then
		return false
	else
		piepan.me:moveTo(channel)
		return true
	end
end

-- Performs functions that allow the bot to safely exit.
function kill()
	os.remove("song.ogg")
	os.remove("song-converted.ogg")
	os.remove(".video_fail")
	piepan.disconnect()
	os.exit(0)
end

-- Checks the permissions of a user against the config to see if they are
-- allowed to execute a certain command.
function check_permissions(ADMIN_COMMAND, username)
	if config.ENABLE_ADMINS and ADMIN_COMMAND then
		return is_admin(username)
	end
	
	return true
end

-- Checks if a user is an admin, as specified in config.lua.
function is_admin(username)
	for _,user in pairs(config.ADMINS) do
		if user == username then
			return true
		end
	end
	
	return false
end

-- Switches to the next song.
function next_song()
	reset_skips()
	if get_length() ~= 0 then
		local success = get_next_song()
		if not success then
			piepan.me.channel:send("An error occurred while preparing the next track. Skipping...")
		end
	end
end

-- Checks if a file exists.
function file_exists(file)
	local f=io.open(file,"r")
	if f~=nil then io.close(f) return true else return false end
end

-- Returns the number of users in the Mumble server.
function count_users()
	local user_count = -1 -- Set to -1 to account for the bot
	for name,_ in pairs(piepan.users) do
		user_count = user_count + 1
	end
	return user_count
end

-------------------------------------------------
-- Song Queue Stuff                            --
-- Contains the definition of the song queue   --
-- used for queueing up songs.                 --
-------------------------------------------------

local song_queue = deque.new()
local skippers = {}

-- Begins the process of adding a new song to the song queue.
function add_song(url, username)
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
		return get_next_song()
	end
	
	return true
end

-- Deletes the old song and begins the process of retrieving a new one.
function get_next_song()
	reset_skips()
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
	os.execute("python download_audio.py " .. info.id)
	if not file_exists(".video_fail") then
		while not file_exists("song-converted.ogg") do
			os.execute("sleep " .. tonumber(2))
		end
		if piepan.Audio:isPlaying() then
			piepan.Audio:stop()
		end
		piepan.me.channel:play({filename="song-converted.ogg", volume=config.VOLUME}, get_next_song)
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
function add_skip(username)
	local already_skipped = false
	for _,name in pairs(skippers) do
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
function count_skippers()
	local skipper_count = 0
	for name,_ in pairs(skippers) do
		skipper_count = skipper_count + 1
	end

	return skipper_count
end

-- Resets the list of users who would like to skip a song. Called during a
-- transition between songs.
function reset_skips()
	skippers = {}
end

-- Retrieves the length of the song queue and returns it.
function get_length()
	return song_queue:length()
end
