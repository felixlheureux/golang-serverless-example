package store

import (
	"github.com/Masterminds/squirrel"
	"github.com/manta-coder/golang-serverless-example/pkg/db"
	"github.com/manta-coder/golang-serverless-example/pkg/domain"
)

// postgres wants $ instead of ? for variable bindings
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

const (
	usersTable         = "users"
	challengesTable    = "challenges"
	clansTable         = "clans"
	charactersTable    = "characters"
	notificationsTable = "notifications"
	SquadsTable        = "squads"
)

// used to facilitate select all fields without using wildcard (*)
var (
	usersColumns         = db.GetDBColumns(domain.User{})
	challengesColumns    = db.GetDBColumns(domain.Challenge{})
	clansColumns         = db.GetDBColumns(domain.Clan{})
	charactersColumns    = db.GetDBColumns(domain.Character{})
	notificationsColumns = db.GetDBColumns(domain.Notification{})
	SquadsColumns        = db.GetDBColumns(domain.Squad{})
)
