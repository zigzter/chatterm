package utils

import "github.com/zigzter/chatterm/types"

var (
	JoinMessage            = ":saruman!saruman@saruman.tmi.twitch.tv JOIN #gandalf"
	PartMessage            = ":saruman!saruman@saruman.tmi.twitch.tv PART #gandalf"
	RawChatMessage         = "@badge-info=subscriber/77;badges=moderator/1,subscriber/36,bits/100;color=#00F8FF;display-name=gandalf;emotes=;first-msg=0;flags=;id=e923f2fb-bfd6-415b-b796-91f5a4748357;mod=1;returning-chatter=0;room-id=19070311;subscriber=1;tmi-sent-ts=1703527538237;turbo=0;user-id=20816785;user-type=mod :gandalf!gandalf@gandalf.tmi.twitch.tv PRIVMSG #gandalf :All we have to decide is what to do with the time that is given to us."
	RegResubMessage        = `@badge-info=subscriber/12;badges=subscriber/12;color=#008000;display-name=gimli;emotes=;flags=;id=aa408023-4633-4fa7-bbac-cbf2777c8e34;login=gimli;mod=0;msg-id=resub;msg-param-cumulative-months=12;msg-param-months=0;msg-param-multimonth-duration=0;msg-param-multimonth-tenure=0;msg-param-should-share-streak=1;msg-param-streak-months=1;msg-param-sub-plan-name=Channel\sSubscription\s(gandalf);msg-param-sub-plan=1000;msg-param-was-gifted=false;room-id=19070311;subscriber=1;system-msg=gimli\ssubscribed\sat\sTier\s1.\sThey've\ssubscribed\sfor\s12\smonths,\scurrently\son\sa\s1\smonth\sstreak!;tmi-sent-ts=1703446121758;user-id=27979720;user-type=;vip=0 :tmi.twitch.tv USERNOTICE #gandalf :Nobody tosses a dwarf`
	GlobalUserStateMessage = `:tmi.twitch.tv 001 frodo :Welcome, GLHF!
    :tmi.twitch.tv 002 frodo :Your host is tmi.twitch.tv
    :tmi.twitch.tv 003 frodo :This server is rather new
    :tmi.twitch.tv 004 frodo :-
    :tmi.twitch.tv 375 frodo :-
    :tmi.twitch.tv 372 frodo :You are in a maze of twisty passages, all alike.
    :tmi.twitch.tv 376 frodo :>
    @badge-info=;badges=twitch-recap-2023/1;color=#00F8FF;display-name=frodo;emote-sets=0,15067,19194,791602,1512303,300206297,300374282,301690833,301800850,302210471,303528731,340690726,367062000,369287198,410301729,459130136,472873131,477339272,485767368,488737509,537206155,564265402,592920959,610186276,1709863403,45f285b1-d1bf-40c5-9e81-eed39de4f4d1,7d0b9d11-479c-465a-b4fa-f732a5790599;user-id=20816785;user-type= :tmi.twitch.tv GLOBALUSERSTATE`
	ParsedSubMessage = types.SubMessage{
		DisplayName: "gimli",
		Message:     "Nobody tosses a dwarf",
		Months:      "12",
		Streak:      "1",
	}
	ParsedChatMessage = types.ChatMessage{
		Color:          "#00F8FF",
		DisplayName:    "gandalf",
		IsFirstMessage: false,
		IsMod:          true,
		IsVIP:          false,
		Message:        "All we have to decide is what to do with the time that is given to us.",
		Timestamp:      "10:05",
	}
)
