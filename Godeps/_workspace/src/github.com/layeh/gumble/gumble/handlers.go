package gumble

import (
	"bytes"
	"crypto/x509"
	"encoding/binary"
	"errors"
	"math"
	"net"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/layeh/gopus"
	"github.com/layeh/gumble/gumble/MumbleProto"
	"github.com/layeh/gumble/gumble/varint"
)

type handlerFunc func(*Client, []byte) error

var (
	errUnimplementedHandler = errors.New("the handler has not been implemented")
	errIncompleteProtobuf   = errors.New("protobuf message is missing a required field")
	errInvalidProtobuf      = errors.New("protobuf message has an invalid field")
	errUnsupportedAudio     = errors.New("unsupported audio codec")
)

var handlers = map[uint16]handlerFunc{
	0:  (*Client).handleVersion,
	1:  (*Client).handleUdpTunnel,
	2:  (*Client).handleAuthenticate,
	3:  (*Client).handlePing,
	4:  (*Client).handleReject,
	5:  (*Client).handleServerSync,
	6:  (*Client).handleChannelRemove,
	7:  (*Client).handleChannelState,
	8:  (*Client).handleUserRemove,
	9:  (*Client).handleUserState,
	10: (*Client).handleBanList,
	11: (*Client).handleTextMessage,
	12: (*Client).handlePermissionDenied,
	13: (*Client).handleACL,
	14: (*Client).handleQueryUsers,
	15: (*Client).handleCryptSetup,
	16: (*Client).handleContextActionModify,
	17: (*Client).handleContextAction,
	18: (*Client).handleUserList,
	19: (*Client).handleVoiceTarget,
	20: (*Client).handlePermissionQuery,
	21: (*Client).handleCodecVersion,
	22: (*Client).handleUserStats,
	23: (*Client).handleRequestBlob,
	24: (*Client).handleServerConfig,
	25: (*Client).handleSuggestConfig,
}

func parseVersion(packet *MumbleProto.Version) Version {
	var version Version
	if packet.Version != nil {
		version.Version = *packet.Version
	}
	if packet.Release != nil {
		version.Release = *packet.Release
	}
	if packet.Os != nil {
		version.OS = *packet.Os
	}
	if packet.OsVersion != nil {
		version.OSVersion = *packet.OsVersion
	}
	return version
}

func (c *Client) handleVersion(buffer []byte) error {
	var packet MumbleProto.Version
	if err := proto.Unmarshal(buffer, &packet); err != nil {
		return err
	}
	return nil
}

func (c *Client) handleUdpTunnel(buffer []byte) error {
	reader := bytes.NewReader(buffer)
	var bytesRead int64

	var audioType byte
	var audioTarget byte
	var user *User
	var audioLength int

	// Header byte
	typeTarget, err := varint.ReadByte(reader)
	if err != nil {
		return err
	}
	audioType = (typeTarget >> 5) & 0x7
	audioTarget = typeTarget & 0x1F
	// Opus only
	if audioType != 4 {
		return errUnsupportedAudio
	}
	bytesRead++

	// Session
	session, n, err := varint.ReadFrom(reader)
	if err != nil {
		return err
	}
	user = c.Users[uint32(session)]
	if user == nil {
		return errInvalidProtobuf
	}
	bytesRead += n

	// Sequence
	sequence, n, err := varint.ReadFrom(reader)
	if err != nil {
		return err
	}
	bytesRead += n

	// Length
	length, n, err := varint.ReadFrom(reader)
	if err != nil {
		return err
	}
	// Opus audio packets set the 13th bit in the size field as the terminator.
	audioLength = int(length) &^ 0x2000
	if audioLength > reader.Len() {
		return errInvalidProtobuf
	}
	audioLength64 := int64(audioLength)
	bytesRead += n

	opus := buffer[bytesRead : bytesRead+audioLength64]
	pcm, err := user.decoder.Decode(opus, AudioMaximumFrameSize, false)
	if err != nil {
		return err
	}
	event := AudioPacketEvent{
		Client: c,
	}
	event.AudioPacket.Sender = user
	event.AudioPacket.Target = int(audioTarget)
	event.AudioPacket.Sequence = int(sequence)
	event.AudioPacket.PositionalAudioBuffer.AudioBuffer = pcm

	reader.Seek(audioLength64, 1)
	binary.Read(reader, binary.LittleEndian, &event.AudioPacket.PositionalAudioBuffer.X)
	binary.Read(reader, binary.LittleEndian, &event.AudioPacket.PositionalAudioBuffer.Y)
	binary.Read(reader, binary.LittleEndian, &event.AudioPacket.PositionalAudioBuffer.Z)

	c.audioListeners.OnAudioPacket(&event)
	return nil
}

func (c *Client) handleAuthenticate(buffer []byte) error {
	return errUnimplementedHandler
}

func (c *Client) handlePing(buffer []byte) error {
	var packet MumbleProto.Ping
	if err := proto.Unmarshal(buffer, &packet); err != nil {
		return err
	}
	c.pingStats.TCPPackets++
	return nil
}

func (c *Client) handleReject(buffer []byte) error {
	var packet MumbleProto.Reject
	if err := proto.Unmarshal(buffer, &packet); err != nil {
		return err
	}

	if packet.Type != nil {
		c.disconnectEvent.Type = DisconnectType(*packet.Type)
	}
	if packet.Reason != nil {
		c.disconnectEvent.String = *packet.Reason
	}
	c.Conn.Close()
	return nil
}

func (c *Client) handleServerSync(buffer []byte) error {
	var packet MumbleProto.ServerSync
	if err := proto.Unmarshal(buffer, &packet); err != nil {
		return err
	}
	event := ConnectEvent{
		Client: c,
	}

	if packet.Session != nil {
		c.Self = c.Users[*packet.Session]
	}
	if packet.WelcomeText != nil {
		event.WelcomeMessage = *packet.WelcomeText
	}
	if packet.MaxBandwidth != nil {
		event.MaximumBitrate = int(*packet.MaxBandwidth)
	}
	c.State = StateSynced

	c.listeners.OnConnect(&event)
	return nil
}

func (c *Client) handleChannelRemove(buffer []byte) error {
	var packet MumbleProto.ChannelRemove
	if err := proto.Unmarshal(buffer, &packet); err != nil {
		return err
	}

	if packet.ChannelId == nil {
		return errIncompleteProtobuf
	}
	var channel *Channel
	{
		channelID := *packet.ChannelId
		channel = c.Channels[channelID]
		if channel == nil {
			return errInvalidProtobuf
		}
		delete(c.Channels, channelID)
		delete(c.permissions, channelID)
		if parent := channel.Parent; parent != nil {
			delete(parent.Children, channel.ID)
		}
	}

	if c.State == StateSynced {
		event := ChannelChangeEvent{
			Client:  c,
			Type:    ChannelChangeRemoved,
			Channel: channel,
		}
		c.listeners.OnChannelChange(&event)
	}
	return nil
}

func (c *Client) handleChannelState(buffer []byte) error {
	var packet MumbleProto.ChannelState
	if err := proto.Unmarshal(buffer, &packet); err != nil {
		return err
	}

	if packet.ChannelId == nil {
		return errIncompleteProtobuf
	}
	event := ChannelChangeEvent{
		Client: c,
	}
	channelID := *packet.ChannelId
	channel := c.Channels[channelID]
	if channel == nil {
		channel = c.Channels.create(channelID)
		channel.client = c

		event.Type |= ChannelChangeCreated
	}
	event.Channel = channel
	if packet.Parent != nil {
		if channel.Parent != nil {
			delete(channel.Parent.Children, channelID)
		}
		newParent := c.Channels[*packet.Parent]
		if newParent != channel.Parent {
			event.Type |= ChannelChangeMoved
		}
		channel.Parent = newParent
		if channel.Parent != nil {
			channel.Parent.Children[channel.ID] = channel
		}
	}
	if packet.Name != nil {
		if *packet.Name != channel.Name {
			event.Type |= ChannelChangeName
		}
		channel.Name = *packet.Name
	}
	if packet.Description != nil {
		if *packet.Description != channel.Description {
			event.Type |= ChannelChangeDescription
		}
		channel.Description = *packet.Description
		channel.DescriptionHash = nil
	}
	if packet.Temporary != nil {
		channel.Temporary = *packet.Temporary
	}
	if packet.Position != nil {
		if *packet.Position != channel.Position {
			event.Type |= ChannelChangePosition
		}
		channel.Position = *packet.Position
	}
	if packet.DescriptionHash != nil {
		event.Type |= ChannelChangeDescription
		channel.DescriptionHash = packet.DescriptionHash
		channel.Description = ""
	}

	if c.State == StateSynced {
		c.listeners.OnChannelChange(&event)
	}
	return nil
}

func (c *Client) handleUserRemove(buffer []byte) error {
	var packet MumbleProto.UserRemove
	if err := proto.Unmarshal(buffer, &packet); err != nil {
		return err
	}

	if packet.Session == nil {
		return errIncompleteProtobuf
	}
	event := UserChangeEvent{
		Client: c,
		Type:   UserChangeDisconnected,
	}
	{
		session := *packet.Session
		event.User = c.Users[session]
		if event.User == nil {
			return errInvalidProtobuf
		}
		if event.User.Channel != nil {
			delete(event.User.Channel.Users, session)
		}
		delete(c.Users, session)
	}
	if packet.Actor != nil {
		event.Actor = c.Users[*packet.Actor]
		if event.Actor == nil {
			return errInvalidProtobuf
		}
		event.Type |= UserChangeKicked
	}
	if packet.Reason != nil {
		event.String = *packet.Reason
	}
	if packet.Ban != nil && *packet.Ban {
		event.Type |= UserChangeBanned
	}
	if event.User == c.Self {
		if packet.Ban != nil && *packet.Ban {
			c.disconnectEvent.Type = DisconnectBanned
		} else {
			c.disconnectEvent.Type = DisconnectKicked
		}
	}

	if c.State == StateSynced {
		c.listeners.OnUserChange(&event)
	}
	return nil
}

func (c *Client) handleUserState(buffer []byte) error {
	var packet MumbleProto.UserState
	if err := proto.Unmarshal(buffer, &packet); err != nil {
		return err
	}

	if packet.Session == nil {
		return errIncompleteProtobuf
	}
	event := UserChangeEvent{
		Client: c,
	}
	var user, actor *User
	{
		session := *packet.Session
		user = c.Users[session]
		if user == nil {
			user = c.Users.create(session)
			user.Channel = c.Channels[0]
			user.client = c

			event.Type |= UserChangeConnected

			decoder, _ := gopus.NewDecoder(AudioSampleRate, 1)
			user.decoder = decoder

			if user.Channel == nil {
				return errInvalidProtobuf
			}
			event.Type |= UserChangeChannel
			user.Channel.Users[session] = user
		}
	}
	event.User = user
	if packet.Actor != nil {
		actor = c.Users[*packet.Actor]
		if actor == nil {
			return errInvalidProtobuf
		}
		event.Actor = actor
	}
	if packet.Name != nil {
		if *packet.Name != user.Name {
			event.Type |= UserChangeName
		}
		user.Name = *packet.Name
	}
	if packet.UserId != nil {
		if *packet.UserId != user.UserID && !event.Type.Has(UserChangeConnected) {
			if *packet.UserId != math.MaxUint32 {
				event.Type |= UserChangeRegistered
				user.UserID = *packet.UserId
			} else {
				event.Type |= UserChangeUnregistered
				user.UserID = 0
			}
		} else {
			user.UserID = *packet.UserId
		}
	}
	if packet.ChannelId != nil {
		if user.Channel != nil {
			delete(user.Channel.Users, user.Session)
		}
		newChannel := c.Channels[*packet.ChannelId]
		if newChannel == nil {
			return errInvalidProtobuf
		}
		if newChannel != user.Channel {
			event.Type |= UserChangeChannel
			user.Channel = newChannel
		}
		user.Channel.Users[user.Session] = user
	}
	if packet.Mute != nil {
		if *packet.Mute != user.Muted {
			event.Type |= UserChangeAudio
		}
		user.Muted = *packet.Mute
	}
	if packet.Deaf != nil {
		if *packet.Deaf != user.Deafened {
			event.Type |= UserChangeAudio
		}
		user.Deafened = *packet.Deaf
	}
	if packet.Suppress != nil {
		if *packet.Suppress != user.Suppressed {
			event.Type |= UserChangeAudio
		}
		user.Suppressed = *packet.Suppress
	}
	if packet.SelfMute != nil {
		if *packet.SelfMute != user.SelfMuted {
			event.Type |= UserChangeAudio
		}
		user.SelfMuted = *packet.SelfMute
	}
	if packet.SelfDeaf != nil {
		if *packet.SelfDeaf != user.SelfDeafened {
			event.Type |= UserChangeAudio
		}
		user.SelfDeafened = *packet.SelfDeaf
	}
	if packet.Texture != nil {
		event.Type |= UserChangeTexture
		user.Texture = packet.Texture
		user.TextureHash = nil
	}
	if packet.Comment != nil {
		if *packet.Comment != user.Comment {
			event.Type |= UserChangeComment
		}
		user.Comment = *packet.Comment
		user.CommentHash = nil
	}
	if packet.Hash != nil {
		user.Hash = *packet.Hash
	}
	if packet.CommentHash != nil {
		event.Type |= UserChangeComment
		user.CommentHash = packet.CommentHash
		user.Comment = ""
	}
	if packet.TextureHash != nil {
		event.Type |= UserChangeTexture
		user.TextureHash = packet.TextureHash
		user.Texture = nil
	}
	if packet.PrioritySpeaker != nil {
		if *packet.PrioritySpeaker != user.PrioritySpeaker {
			event.Type |= UserChangePrioritySpeaker
		}
		user.PrioritySpeaker = *packet.PrioritySpeaker
	}
	if packet.Recording != nil {
		if *packet.Recording != user.Recording {
			event.Type |= UserChangeRecording
		}
		user.Recording = *packet.Recording
	}

	if c.State == StateSynced {
		c.listeners.OnUserChange(&event)
	}
	return nil
}

func (c *Client) handleBanList(buffer []byte) error {
	var packet MumbleProto.BanList
	if err := proto.Unmarshal(buffer, &packet); err != nil {
		return err
	}

	event := BanListEvent{
		Client:  c,
		BanList: make(BanList, 0, len(packet.Bans)),
	}

	for _, banPacket := range packet.Bans {
		ban := &Ban{
			Address: net.IP(banPacket.Address),
		}
		if banPacket.Mask != nil {
			size := net.IPv4len * 8
			if len(ban.Address) == net.IPv6len {
				size = net.IPv6len * 8
			}
			ban.Mask = net.CIDRMask(int(*banPacket.Mask), size)
		}
		if banPacket.Name != nil {
			ban.Name = *banPacket.Name
		}
		if banPacket.Hash != nil {
			ban.Hash = *banPacket.Hash
		}
		if banPacket.Reason != nil {
			ban.Reason = *banPacket.Reason
		}
		if banPacket.Start != nil {
			ban.Start, _ = time.Parse(time.RFC3339, *banPacket.Start)
		}
		if banPacket.Duration != nil {
			ban.Duration = time.Duration(*banPacket.Duration) * time.Second
		}
		event.BanList = append(event.BanList, ban)
	}

	c.listeners.OnBanList(&event)
	return nil
}

func (c *Client) handleTextMessage(buffer []byte) error {
	var packet MumbleProto.TextMessage
	if err := proto.Unmarshal(buffer, &packet); err != nil {
		return err
	}

	event := TextMessageEvent{
		Client: c,
	}
	if packet.Actor != nil {
		event.Sender = c.Users[*packet.Actor]
	}
	if packet.Session != nil {
		event.Users = make([]*User, 0, len(packet.Session))
		for _, session := range packet.Session {
			if user := c.Users[session]; user != nil {
				event.Users = append(event.Users, user)
			}
		}
	}
	if packet.ChannelId != nil {
		event.Channels = make([]*Channel, 0, len(packet.ChannelId))
		for _, id := range packet.ChannelId {
			if channel := c.Channels[id]; channel != nil {
				event.Channels = append(event.Channels, channel)
			}
		}
	}
	if packet.TreeId != nil {
		event.Trees = make([]*Channel, 0, len(packet.TreeId))
		for _, id := range packet.TreeId {
			if channel := c.Channels[id]; channel != nil {
				event.Trees = append(event.Trees, channel)
			}
		}
	}
	if packet.Message != nil {
		event.Message = *packet.Message
	}

	c.listeners.OnTextMessage(&event)
	return nil
}

func (c *Client) handlePermissionDenied(buffer []byte) error {
	var packet MumbleProto.PermissionDenied
	if err := proto.Unmarshal(buffer, &packet); err != nil {
		return err
	}

	if packet.Type == nil || *packet.Type == MumbleProto.PermissionDenied_H9K {
		return errInvalidProtobuf
	}

	event := PermissionDeniedEvent{
		Client: c,
		Type:   PermissionDeniedType(*packet.Type),
	}
	if packet.Reason != nil {
		event.String = *packet.Reason
	}
	if packet.Name != nil {
		event.String = *packet.Name
	}
	if packet.Session != nil {
		event.User = c.Users[*packet.Session]
		if event.User == nil {
			return errInvalidProtobuf
		}
	}
	if packet.ChannelId != nil {
		event.Channel = c.Channels[*packet.ChannelId]
		if event.Channel == nil {
			return errInvalidProtobuf
		}
	}
	if packet.Permission != nil {
		event.Permission = Permission(*packet.Permission)
	}

	c.listeners.OnPermissionDenied(&event)
	return nil
}

func (c *Client) handleACL(buffer []byte) error {
	var packet MumbleProto.ACL
	if err := proto.Unmarshal(buffer, &packet); err != nil {
		return err
	}

	acl := &ACL{
		Inherits: packet.GetInheritAcls(),
	}
	if packet.ChannelId == nil {
		return errInvalidProtobuf
	}
	acl.Channel = c.Channels[*packet.ChannelId]
	if acl.Channel == nil {
		return errInvalidProtobuf
	}

	if packet.Groups != nil {
		acl.Groups = make([]*ACLGroup, 0, len(packet.Groups))
		for _, group := range packet.Groups {
			aclGroup := &ACLGroup{
				Name:         *group.Name,
				Inherited:    group.GetInherited(),
				InheritUsers: group.GetInherit(),
				Inheritable:  group.GetInheritable(),
			}
			if group.Add != nil {
				aclGroup.usersAdd = make(map[uint32]*ACLUser)
				for _, userID := range group.Add {
					aclGroup.usersAdd[userID] = &ACLUser{
						UserID: userID,
					}
				}
			}
			if group.Remove != nil {
				aclGroup.usersRemove = make(map[uint32]*ACLUser)
				for _, userID := range group.Remove {
					aclGroup.usersRemove[userID] = &ACLUser{
						UserID: userID,
					}
				}
			}
			if group.InheritedMembers != nil {
				aclGroup.usersInherited = make(map[uint32]*ACLUser)
				for _, userID := range group.InheritedMembers {
					aclGroup.usersInherited[userID] = &ACLUser{
						UserID: userID,
					}
				}
			}
			acl.Groups = append(acl.Groups, aclGroup)
		}
	}
	if packet.Acls != nil {
		acl.Rules = make([]*ACLRule, 0, len(packet.Acls))
		for _, rule := range packet.Acls {
			aclRule := &ACLRule{
				AppliesCurrent:  rule.GetApplyHere(),
				AppliesChildren: rule.GetApplySubs(),
				Inherited:       rule.GetInherited(),
				Granted:         Permission(rule.GetGrant()),
				Denied:          Permission(rule.GetDeny()),
			}
			if rule.UserId != nil {
				aclRule.User = &ACLUser{
					UserID: *rule.UserId,
				}
			} else if rule.Group != nil {
				var group *ACLGroup
				for _, g := range acl.Groups {
					if g.Name == *rule.Group {
						group = g
						break
					}
				}
				if group == nil {
					group = &ACLGroup{
						Name: *rule.Group,
					}
				}
				aclRule.Group = group
			}
			acl.Rules = append(acl.Rules, aclRule)
		}
	}
	c.tmpACL = acl
	return nil
}

func (c *Client) handleQueryUsers(buffer []byte) error {
	var packet MumbleProto.QueryUsers
	if err := proto.Unmarshal(buffer, &packet); err != nil {
		return err
	}

	acl := c.tmpACL
	if acl == nil {
		return errIncompleteProtobuf
	}
	c.tmpACL = nil

	userMap := make(map[uint32]string)
	for i := 0; i < len(packet.Ids) && i < len(packet.Names); i++ {
		userMap[packet.Ids[i]] = packet.Names[i]
	}

	for _, group := range acl.Groups {
		for _, user := range group.usersAdd {
			user.Name = userMap[user.UserID]
		}
		for _, user := range group.usersRemove {
			user.Name = userMap[user.UserID]
		}
		for _, user := range group.usersInherited {
			user.Name = userMap[user.UserID]
		}
	}
	for _, rule := range acl.Rules {
		if rule.User != nil {
			rule.User.Name = userMap[rule.User.UserID]
		}
	}

	event := ACLEvent{
		Client: c,
		ACL:    acl,
	}
	c.listeners.OnACL(&event)
	return nil
}

func (c *Client) handleCryptSetup(buffer []byte) error {
	return errUnimplementedHandler
}

func (c *Client) handleContextActionModify(buffer []byte) error {
	var packet MumbleProto.ContextActionModify
	if err := proto.Unmarshal(buffer, &packet); err != nil {
		return err
	}

	if packet.Action == nil || packet.Operation == nil {
		return errInvalidProtobuf
	}

	event := ContextActionChangeEvent{
		Client: c,
	}

	switch *packet.Operation {
	case MumbleProto.ContextActionModify_Add:
		if ca := c.ContextActions[*packet.Action]; ca != nil {
			return nil
		}
		event.Type = ContextActionAdd
		contextAction := c.ContextActions.create(*packet.Action)
		if packet.Text != nil {
			contextAction.Label = *packet.Text
		}
		if packet.Context != nil {
			contextAction.Type = ContextActionType(*packet.Context)
		}
		event.ContextAction = contextAction
	case MumbleProto.ContextActionModify_Remove:
		contextAction := c.ContextActions[*packet.Action]
		if contextAction == nil {
			return nil
		}
		event.Type = ContextActionRemove
		delete(c.ContextActions, *packet.Action)
		event.ContextAction = contextAction
	default:
		return errInvalidProtobuf
	}

	c.listeners.OnContextActionChange(&event)
	return nil
}

func (c *Client) handleContextAction(buffer []byte) error {
	return errUnimplementedHandler
}

func (c *Client) handleUserList(buffer []byte) error {
	var packet MumbleProto.UserList
	if err := proto.Unmarshal(buffer, &packet); err != nil {
		return err
	}

	event := UserListEvent{
		Client:   c,
		UserList: make(RegisteredUsers, 0, len(packet.Users)),
	}

	for _, user := range packet.Users {
		registeredUser := &RegisteredUser{
			UserID: *user.UserId,
		}
		if user.Name != nil {
			registeredUser.Name = *user.Name
		}
		event.UserList = append(event.UserList, registeredUser)
	}

	c.listeners.OnUserList(&event)
	return nil
}

func (c *Client) handleVoiceTarget(buffer []byte) error {
	return errUnimplementedHandler
}

func (c *Client) handlePermissionQuery(buffer []byte) error {
	var packet MumbleProto.PermissionQuery
	if err := proto.Unmarshal(buffer, &packet); err != nil {
		return err
	}

	if packet.Flush != nil && *packet.Flush {
		oldPermissions := c.permissions
		c.permissions = make(map[uint32]*Permission)
		for channelID := range oldPermissions {
			channel := c.Channels[channelID]
			event := ChannelChangeEvent{
				Client:  c,
				Type:    ChannelChangePermission,
				Channel: channel,
			}
			c.listeners.OnChannelChange(&event)
		}
	}
	if packet.ChannelId != nil {
		channel := c.Channels[*packet.ChannelId]
		if packet.Permissions != nil {
			p := Permission(*packet.Permissions)
			c.permissions[channel.ID] = &p
			event := ChannelChangeEvent{
				Client:  c,
				Type:    ChannelChangePermission,
				Channel: channel,
			}
			c.listeners.OnChannelChange(&event)
		}
	}
	return nil
}

func (c *Client) handleCodecVersion(buffer []byte) error {
	return errUnimplementedHandler
}

func (c *Client) handleUserStats(buffer []byte) error {
	var packet MumbleProto.UserStats
	if err := proto.Unmarshal(buffer, &packet); err != nil {
		return err
	}

	if packet.Session == nil {
		return errIncompleteProtobuf
	}
	user := c.Users[*packet.Session]
	if user == nil {
		return errInvalidProtobuf
	}

	if user.Stats == nil {
		user.Stats = &UserStats{}
	}
	*user.Stats = UserStats{
		User: user,
	}
	stats := user.Stats

	if packet.Version != nil {
		stats.Version = parseVersion(packet.Version)
	}
	if packet.Onlinesecs != nil {
		stats.Connected = time.Now().Add(time.Duration(*packet.Onlinesecs) * -time.Second)
	}
	if packet.Idlesecs != nil {
		stats.Idle = time.Duration(*packet.Idlesecs) * time.Second
	}
	if packet.Bandwidth != nil {
		stats.Bandwidth = int(*packet.Bandwidth)
	}
	if packet.Address != nil {
		stats.IP = net.IP(packet.Address)
	}
	if packet.Certificates != nil {
		stats.Certificates = make([]*x509.Certificate, 0, len(packet.Certificates))
		for _, data := range packet.Certificates {
			if data != nil {
				if cert, err := x509.ParseCertificate(data); err == nil {
					stats.Certificates = append(stats.Certificates, cert)
				}
			}
		}
	}
	stats.StrongCertificate = packet.GetStrongCertificate()
	stats.CELTVersions = packet.GetCeltVersions()
	if packet.Opus != nil {
		stats.Opus = *packet.Opus
	}

	event := UserChangeEvent{
		Client: c,
		Type:   UserChangeStats,
		User:   user,
	}

	c.listeners.OnUserChange(&event)
	return nil
}

func (c *Client) handleRequestBlob(buffer []byte) error {
	return errUnimplementedHandler
}

func (c *Client) handleServerConfig(buffer []byte) error {
	return errUnimplementedHandler
}

func (c *Client) handleSuggestConfig(buffer []byte) error {
	return errUnimplementedHandler
}
