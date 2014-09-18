-------------------------
--      MumbleDJ       --
-- By Matthieu Grieger --
-------------------------

local config = require("config")
local song_queue = require("song_queue")

-- NOTE: Skippers will need to be fixed since all the song-related stuff is being moved to song_queue.lua
local skippers = {}

function piepan.onConnect()
	print(piepan.me.name .. " has connected to the server!")
	local user = piepan.users[piepan.me.name]
	local channel = user.channel("Bot Testing")
	piepan.me:moveTo(channel)
end

function piepan.onMessage(message)
	if message.user == nil then
		return
	end

	if string.sub(message.text, 0, 1) == config.COMMAND_PREFIX then
		parse_command(message)
	end
end

function parse_command(message)
	local command = ""
	local argument = ""
	if string.find(message.text, " ") then
		command = string.sub(message.text, 2, string.find(message.text, ' ') - 1)
		argument = string.sub(message.text, string.find(message.text, ' ') + 1)
	else
		command = string.sub(message.text, 2)
	end
	
	if command == "play" then
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
	elseif command == "pause" then
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
	elseif command == "add" then
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
	elseif command == "skip" then
		local has_permission = check_permissions(config.ADMIN_SKIP, message.user.name)
		
		if has_permission then
			if config.OUTPUT then 
				print(message.user.name .. " has voted to skip the current song.")
			end
			
			skip(message.user.name)
		else
			message.user:send(config.NO_PERMISSION_MSG)
		end
	elseif command == "volume" then
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
	elseif command == "move" then
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
	elseif command == "kill" then
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

function skip(username)
	local user_count = 0
	local skipper_count = 0
	local already_skipped = false
	for name,_ in pairs(piepan.users) do
		user_count = user_count + 1
	end
	
	user_count = user_count - 1 -- So that we do not count the bot.
	
	for name,_ in pairs(skippers) do
		if name == username then
			already_skipped = true
		end
		skipper_count = skipper_count + 1
	end
	
	if not already_skipped then
		table.insert(skippers, username)
		skipper_count = skipper_count + 1
		local skip_ratio = skipper_count / user_count
		if skip_ratio > config.SKIP_RATIO then
			piepan.me.channel:send("The number of votes required for a skip has been met. Skipping song!")
			next_song()
		else
			piepan.me.channel:send("<b>" .. username .. "</b> has voted to skip this song.")
		end
	else
		message.user:send("You have already voted to skip this song.")
	end
end

function move(chan)
	local user = piepan.users["MumbleDJ"]
	local channel = user.channel("../" .. chan)
	if channel == nil then
		return false
	else
		piepan.me:moveTo(channel)
		return true
	end
end

function kill()
	os.remove("song.ogg")
	os.remove("song-converted.ogg")
	os.remove(".video_fail")
	os.exit(0)
end

function check_permissions(ADMIN_COMMAND, username)
	if config.ENABLE_ADMINS and ADMIN_COMMAND then
		return is_admin(username)
	end
	
	return true
end

function is_admin(username)
	for _,user in pairs(config.ADMINS) do
		if user == username then
			return true
		end
	end
	
	return false
end

function next_song()
	skippers = {}
	if song_queue:get_length() ~= 0 then
		local success = song_queue:get_next_song()
		if not success then
			piepan.me.channel:send("An error occurred while preparing the next track. Skipping...")
		end
	end
end

function file_exists(file)
	local f=io.open(file,"r")
	if f~=nil then io.close(f) return true else return false end
end
