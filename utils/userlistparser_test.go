package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zigzter/chatterm/types"
)

func TestUserListParser(t *testing.T) {
	got := UserListParser(UsersListMessage)
	want := types.UserListMessage{Users: []string{
		"legolas",
		"gimli",
		"gandalf",
		"aragorn",
		"saruman",
		"sauron",
		"bilbo",
		"samwise",
		"merry",
		"pippin",
		"boromir",
		"elrond",
		"galadriel",
		"celeborn",
		"eowyn",
		"eomer",
		"theoden",
		"grima",
		"wormtongue",
		"faramir",
		"denethor",
		"radagast",
		"glorfindel",
		"arwen",
		"haldir",
		"tom",
		"bombadil",
		"goldberry",
		"treebeard",
		"barliman",
		"butterbur",
		"rosie",
		"cotton",
		"lobelia",
		"gaffer",
		"gamgee",
		"belladonna",
		"took",
		"bullroarer",
		"took",
		"isildur",
		"elendil",
		"gil-galad",
		"beren",
		"luthien",
		"turgon",
		"finrod",
		"feanor",
		"melian",
		"thingol",
		"gollum",
		"smeagol",
		"bilbo",
		"frodo",
	}}

	assert.Equal(t, got, want)
}
