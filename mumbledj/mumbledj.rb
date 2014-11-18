# MumbleDJ v2
# By Matthieu Grieger
# mumbledj.rb

require "mumble-ruby"
require_relative "config"
require_relative "song_queue"

# Class that defines MumbleDJ behavior.
class MumbleDJ

  # Initializes a new instance of MumbleDJ. The parameters are as follows:
  # username: Desired username of the Mumble bot
  # server_address: IP address/web address of Mumble server to connect to
  # server_port: Port number of Mumble server (generally 64738)
  # default_channel: The channel you would like the bot to connect to by
  #   default. If the channel does not exist, the bot will connect to
  #   the root channel of the server instead.
  def initialize(username, server_address, server_port, default_channel, password)
    @username = username
    @password = password
    @server_address = server_address
    @server_port = server_port
    @default_channel = default_channel
    @song_queue = SongQueue.new
    
    Mumble.configure do |conf|
      conf.sample_rate = 48000
      conf.bitrate = 32000
      conf.ssl_cert_opts[:cert_dir] = File.expand_path("certs")
    end
  end
  
  # Connects to the Mumble server with the credentials specified in
  # initialize.
  def connect
    @client = Mumble::Client.new(@server_address) do |conf|
      conf.username = @username
      if @password != ""
        conf.password = @password
      end
    end
    
    self.set_callbacks
    
    @client.connect
    @client.on_connected do
      if @default_channel != ""
        @client.join_channel(@default_channel)
      end
    end
  end
  
  # Sets various callbacks that can be triggered during the connection.
  def set_callbacks
    @client.on_text_message do |message|
      parse_message(message)
    end
  end
  
  # Parses messages looking for commands, and calls the appropriate
  # methods to complete each requested command.
  def parse_message(message)
	@sender = @client.users[message.actor].name
    if message.message[0] == COMMAND_PREFIX
      if message.message.count(" ") != 0
        @command = message.message[1..(message.message.index(" ") - 1)]
        @argument = message.message[(message.message.index(" ") + 1)..-1]
      else
        @command = message.message[1..-1]
      end
      
      case @command
        when ADD_ALIAS
          if has_permission?(ADMIN_ADD, @sender)
            if OUTPUT_ENABLED
              puts("#{@sender} has added a song to the queue.")
            end
            if song_add_successful?(@argument, @sender)
              @client.text_channel("#(@sender} has added a song to the queue.")
            else
              @client.text_user(@sender, "The URL you provided was not valid.")
            end
          else
            @client.text_user(@sender, NO_PERMISSION_MSG)
          end
        when SKIP_ALIAS
          if has_permission?(ADMIN_SKIP, @sender)
            if OUTPUT_ENABLED
              puts("#{@sender} has voted to skip the current song.")
            end
            @song_queue.get_current_song.add_skip(@sender)
          else
            @client.text_user(@sender, NO_PERMISSION_MSG)
          end
        when VOLUME_ALIAS
          puts("Volume command request.")
        when MOVE_ALIAS
          if has_permission?(ADMIN_MOVE, @sender)
            begin
              @client.join_channel(@argument)
            rescue Mumble::ChannelNotFound
              @client.text_user(@sender, "The channel you provided does not exist.")
            end
          else
            @client.text_user(@sender, NO_PERMISSION_MSG)
          end
        # This one doesn't work for some reason. Gotta do some testing.
        when KILL_ALIAS
          if has_permission?(ADMIN_KILL, @sender)
            disconnect
          else
            @client.text_user(@sender, NO_PERMISSION_MSG)
          end
        else
          @client.text_user(@sender, INVALID_COMMAND_MSG)
      end
    end
  end
  
  # Checks message sender against ADMINS array to verify if they have
  # permission to use a specific command.
  def has_permission?(admin_command, sender)
    if ENABLE_ADMINS and admin_command
      return ADMINS.include?(sender)
    end
  end
  
  # Safely disconnects the bot from the server.
  def disconnect
    @client.disconnect
  end
end
