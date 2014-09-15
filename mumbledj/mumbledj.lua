-------------------------
--      MumbleDJ       --
-- By Matthieu Grieger --
-------------------------

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
	
	local patterns = {
		"https?://www%.youtube%.com/watch%?v=([%d%a_%-]+)",
        "https?://youtube%.com/watch%?v=([%d%a_%-]+)",
        "https?://youtu.be/([%d%a_%-]+)",
        "https?://youtube.com/v/([%d%a_%-]+)",
        "https?://www.youtube.com/v/([%d%a_%-]+)"
	}
end
