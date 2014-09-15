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
	if string.find(message.text, ' ') then
		command = string.sub(message.text, 2, string.find(message.text, ' ') - 1)
		argument = string.sub(message.text, string.find(message.text, ' ') + 1)
	else
		command = string.sub(message.text, 2)
	end
	
	if command == "play" then
		print(message.user.name .. " has told the bot to start playing music.")
	elseif command == "pause" then
		print(message.user.name .. " has told the bot to pause music playback.")
	elseif command == "add" then
		print(message.user.name .. " has told the bot to add the following URL to the queue: " .. argument .. ".")
	elseif command == "skip" then
		print(message.user.name .. " has voted to skip the current song.")
	elseif command == "volumeup" then
		print(message.user.name .. " has told the bot to raise the playback volume.")
	elseif command == "volumedown" then
		print(message.user.name .. " has told the bot to lower the playback volume.")
	elseif command == "move" then
		print(message.user.name .. " has told the bot to move to the following channel: " .. argument .. ".")
	elseif command == "kill" then
		print(message.user.name .. " has told the bot to kill itself.")
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
