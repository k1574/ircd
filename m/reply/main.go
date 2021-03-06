package reply

import(
	"fmt"
)

const(
	RPL_WELCOME = 1
	RPL_YOURHOST = 2
	RPL_CREATED = 3
	RPL_MYINFO = 4
	RPL_BOUNCE = 5
	RPL_TRACELINK = 200
	RPL_TRACECONNECTING = 201
	RPL_TRACEHANDSHAKE = 202
	RPL_TRACEUNKNOWN = 203
	RPL_TRACEOPERATOR = 204
	RPL_NONE = 300
	RPL_AWAY = 301
	RPL_USERHOST = 302
	RPL_ISON = 303
	RPL_UNAWAY = 305
	RPL_NOAWAY = 306
	RPL_WHOISUSER = 311
	RPL_WHOISSERVER = 312
	RPL_WHOISOPERATOR = 313
	RPL_WHOWASUSER = 314
	RPL_ENDOFWHO = 315
	RPL_WHOISIDLE = 317
	RPL_ENDOFWHOIS = 318
	RPL_WHOISCHANNELS = 319
	RPL_LIST = 322
	RPL_LISTEND = 323
	RPL_CHANNELMODEIS = 324
	RPL_UNIQOPIS = 325
	RPL_NOTOPIC = 331
	RPL_TOPIC = 332
	RPL_INVITING = 341
	RPL_SUMMONING = 342
	RPL_INVITELIST = 346
	RPL_ENDOFINVITELIST = 347
	RPL_EXCEPTLIST = 348
	RPL_ENDOFEXCEPTLIST = 349
	RPL_VERSION = 351
	RPL_WHOREPLY = 352
	RPL_NAMREPLY = 353
	RPL_LINKS = 364
	RPL_ENDOFLINKS = 365
	RPL_ENDOF_NAMES = 366
	RPL_BANLIST = 367
	RPL_ENDOFBANLIST = 368
	RPL_ENDOFWHOWAS = 369
	RPL_INFO = 371
	RPL_MOTD = 372
	RPL_ENDOFINFO = 374
	RPL_MOTDSTART = 375
	ERR_NORECIPIENT = 411
	ERR_NOADMININFO = 423
	ERR_FILEERROR = 424
	ERR_NO_NICKNAMEGIVEN = 431
	ERR_ERRONEUSNICKNAME = 432
	ERR_NICKNAMEINUSE = 433
	ERR_NICKCOLLISION = 436
	ERR_UNAVAILSOURCE = 437
	ERR_USERNOTINCHANNEL = 441
	ERR_NOTONCHANNEL = 442
	ERR_USERONCHANNEL = 443
	ERR_NEEDMOREPARAMS = 461
	ERR_ALREADYREGISTERED = 462
)

var(
	Replies = map[int]string{
		ERR_NEEDMOREPARAMS: "Not enough parameters",
		ERR_ALREADYREGISTERED: "You may not reregister",
	}
)

func
FormatNum(num int) string {
	return fmt.Sprintf("%03d", num)
}

