# MumbleDJ v2
# By Matthieu Grieger

require "mumble-ruby"

# Class that defines MumbleDJ behavior.
class MumbleDJ

  # Initializes a new instance of MumbleDJ. The parameters are as follows:
  # username: Desired username of the Mumble bot
  # server_address: IP address/web address of Mumble server to connect to
  # server_port: Port number of Mumble server (generally 64738)
  # default_channel: The channel you would like the bot to connect to by
  #   default. If the channel does not exist, the bot will connect to
  #   the root channel of the server instead.
  def initialize(username, server_address, server_port=64738, default_channel="", password="")
    @username = username
    @password = password
    @cert = cert
    @server_address = server_address
    @server_port = server_port
    @default_channel = default_channel
    
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
      # TODO: Call message parser here
    end
  end
  
  # Safely disconnects the bot from the server.
  def disconnect
    @client.disconnect
  end
end
