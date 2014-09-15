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
	print(string.sub(message.text, 0, 1))
	if string.sub(message.text, 0, 1) == config.COMMAND_PREFIX then
		parseCommand(message)
		print("Command has been found!")
	end
end

function parseCommand(message)
	return
end
