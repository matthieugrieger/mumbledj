# MumbleDJ
# By Matthieu Grieger
# run_bot.rb

require_relative "mumbledj"
require_relative "config"
require "thread"

bot = MumbleDJ.new(username=BOT_USERNAME, server_address=MUMBLE_SERVER_ADDRESS, port=MUMBLE_SERVER_PORT, 
        default_channel=DEFAULT_CHANNEL, password=MUMBLE_PASSWORD)
bot.connect

begin
  t = Thread.new do
    $stdin.gets
  end
  
  t.join
  rescue Interrupt => e
end
