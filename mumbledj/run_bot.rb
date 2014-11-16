# MumbleDJ
# By Matthieu Grieger
# run_bot.rb

require_relative "mumbledj"
require 'thread'

bot = MumbleDJ.new(username="MumbleDJTest", server_address="matthieugrieger.com") 
bot.connect

begin
  t = Thread.new do
    $stdin.gets
  end
  
  t.join
  rescue Interrupt => e
end
