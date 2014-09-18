-------------------------
--      MumbleDJ       --
-- By Matthieu Grieger --
------------------------------------------------------------------
-- mumbledj.lua                                                 --
-- The main file which defines most of MumbleDJ's behavior. All --
-- commands are found here, and most of their implementation.   --
------------------------------------------------------------------

local config = require("config")
local song_queue = require("song_queue")

-- Connects to Mumble server.
function piepan.onConnect()
	print(piepan.me.name .. " has connected to the server!")
	local user = piepan.users[piepan.me.name]
	local channel = user.channel("Bot Testing")
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
	
	-- Play command
	if command == config.PLAY_ALIAS then
		local has_permission = check_permissions(config.ADMIN_PLAY, message.user.name)
		
		if has_permission then
			if config.OUTPUT then 
				print(message.user.name .. " has told the bot to start playing music.")
			end
			if song_queue.getLength() == 0 then
				message.user:send(config.NO_SONGS_AVAILABLE)
			else
				if piepan.Audio.isPlaying() then
					message.user:send(config.MUSIC_PLAYING_MSG)
				else
					piepan.me.channel:play("song-converted.ogg", SongQueue.get_next_song)
			end
		end
			
		else
			message.user:send(config.NO_PERMISSION_MSG)
		end
	-- Pause command
	elseif command == config.PAUSE_ALIAS then
		local has_permission = check_permissions(config.ADMIN_PAUSE, message.user.name)
		
		if has_permission then
			if config.OUTPUT then 
				print(message.user.name .. " has told the bot to pause music playback.")
			end
			
			if piepan.Audio.isPlaying() then
				piepan.me.channel:send(string.format(config.SONG_PAUSED_HTML, message.user.name))
				piepan.Audio.stop()
			else
				message.user:send(config.NO_MUSIC_PLAYING_MSG)
			end
		else
			message.user:send(config.NO_PERMISSION_MSG)
		end
	-- Add command
	elseif command == config.ADD_ALIAS then
		local has_permission = check_permissions(config.ADMIN_ADD, message.user.name)
		
		if has_permission then
			if config.OUTPUT then 
				print(message.user.name .. " has told the bot to add the following URL to the queue: " .. argument .. ".")
				if not song_queue.add_song(argument, message.user.name) then
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
					if config.LOWEST_VOLUME < argument < config.HIGHEST_VOLUME then
						config.VOLUME = argument
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
	-- This is just where I put commands for testing. These will most likely be removed
	-- in the "final" version.
	elseif command == "musicplaying" then
		if piepan.Audio.isPlaying() then
			message.user:send("Music is currently playing.")
		else
			message.user:send("Music is not currently playing.")
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
	if song_queue:add_skip(username) then
		local skip_ratio = song_queue:count_skippers() / count_users()
		if skip_ratio > config.SKIP_RATIO then
			piepan.me.channel:send(config.SONG_SKIPPED_HTML)
			next_song()
		else
			piepan.me.channel:send(string.format(config.USER_SKIP_HTML, username))
	else
		message.user:send("You have already voted to skip this song.")
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
	song_queue:reset_skips()
	if song_queue:get_length() ~= 0 then
		local success = song_queue:get_next_song()
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
