package utils

import (
    "testing"
)

func TestSubParser(t *testing.T) {
    regResubMessage := `@badge-info=subscriber/12;badges=subscriber/12;color=#008000;display-name=vexthorne;emotes=;flags=;id=aa408023-4633-4fa7-bbac-cbf2777c8e34;login=vexthorne;mod=0;msg-id=resub;msg-param-cumulative-months=12;msg-param-months=0;msg-param-multimonth-duration=0;msg-param-multimonth-tenure=0;msg-param-should-share-streak=1;msg-param-streak-months=1;msg-param-sub-plan-name=Channel\sSubscription\s(a_seagull);msg-param-sub-plan=1000;msg-param-was-gifted=false;room-id=19070311;subscriber=1;system-msg=vexthorne\ssubscribed\sat\sTier\s1.\sThey've\ssubscribed\sfor\s12\smonths,\scurrently\son\sa\s1\smonth\sstreak!;tmi-sent-ts=1703446121758;user-id=27979720;user-type=;vip=0 :tmi.twitch.tv USERNOTICE #a_seagull :merry Christmas`
    got := SubParser(regResubMessage)
    want := SubMessage{ DisplayName: "vexthorne", Message: "merry Christmas", Months: "12", Streak: "1" }

    if got != want {
        t.Errorf("got %q, want %q", got, want)
    }
}
