-------------------------
--      MumbleDJ       --
-- By Matthieu Grieger --
-------------------------

local config = require("config")
require "song_queue"

function piepan.onConnect()
	print("MumbleDJ has connected to the server!")
	local user = piepan.users["MumbleDJ"]
	local channel = user.channel("Bot Testing")
	piepan.me:moveTo(channel)
end

function piepan.onMessage(message)
	if message.user == nil then
		return
	end

	if string.sub(message.text, 0, 1) == config.COMMAND_PREFIX then
		parseCommand(message)
	end
end

function parseCommand(message)
	local command = ""
	local argument = ""
	if string.find(message.text, " ") then
		command = string.sub(message.text, 2, string.find(message.text, ' ') - 1)
		argument = string.sub(message.text, string.find(message.text, ' ') + 1)
	else
		command = string.sub(message.text, 2)
	end
	
	if command == "play" then
		local has_permission = checkPermissions(config.ADMIN_PLAY, message.user.name)
		
		if has_permission then
			if config.OUTPUT then 
				print(message.user.name .. " has told the bot to start playing music.")
			end
		else
			message.user:send(config.NO_PERMISSION_MSG)
		end
	elseif command == "pause" then
		local has_permission = checkPermissions(config.ADMIN_PAUSE, message.user.name)
		
		if has_permission then
			if config.OUTPUT then 
				print(message.user.name .. " has told the bot to pause music playback.")
			end
		else
			message.user:send(config.NO_PERMISSION_MSG)
		end
	elseif command == "add" then
		local has_permission = checkPermissions(config.ADMIN_ADD, message.user.name)
		
		if has_permission then
			if config.OUTPUT then 
				print(message.user.name .. " has told the bot to add the following URL to the queue: " .. argument .. ".")
			end
		else
			message.user:send(config.NO_PERMISSION_MSG)
		end
	elseif command == "skip" then
		local has_permission = checkPermissions(config.ADMIN_SKIP, message.user.name)
		
		if has_permission then
			if config.OUTPUT then 
				print(message.user.name .. " has voted to skip the current song.")
			end
		else
			message.user:send(config.NO_PERMISSION_MSG)
		end
	elseif command == "volumeup" then
		local has_permission = checkPermissions(config.ADMIN_VOLUMEUP, message.user.name)
		
		if has_permission then
			if config.OUTPUT then 
				print(message.user.name .. " has told the bot to raise the playback volume.")
			end
		else
			message.user:send(config.NO_PERMISSION_MSG)
		end
	elseif command == "volumedown" then
		local has_permission = checkPermissions(config.ADMIN_VOLUMEDOWN, message.user.name)
		
		if has_permission then
			if config.OUTPUT then 
				print(message.user.name .. " has told the bot to lower the playback volume.")
			end
		else
			message.user:send(config.NO_PERMISSION_MSG)
		end
	elseif command == "move" then
		local has_permission = checkPermissions(config.ADMIN_MOVE, message.user.name)
		
		if has_permission then
			if config.OUTPUT then 
				print(message.user.name .. " has told the bot to move to the following channel: " .. argument .. ".")
			end
		else
			message.user:send(config.NO_PERMISSION_MSG)
		end
	elseif command == "kill" then
		local has_permission = checkPermissions(config.ADMIN_KILL, message.user.name)
		
		if has_permission then
			if config.OUTPUT then 
				print(message.user.name .. " has told the bot to kill itself.")
			end
		else
			message.user:send(config.NO_PERMISSION_MSG)
		end
	else
		message.user:send("The command you have entered is not valid.")
	end
end

function play()
	return
end

function pause()
	return
end

function add()
	return
end

function skip()
	return
end

function volumeup()
	return
end

function volumedown()
	return
end

function move()
	return
end

function kill()
	return
end

function checkPermissions(ADMIN_COMMAND, username)
	if config.ENABLE_ADMINS and ADMIN_COMMAND then
		return isAdmin(username)
	end
	
	return true
end

function isAdmin(username)
	for _,user in pairs(config.ADMINS) do
		if user == username then
			return true
		end
	end
	
	return false
end
