```
                        __
                       / _|
 __  __             _ | |_       _     _ 
|  \/  | __ _ _ __ (_)|  _| ___ | | __| |
| |\/| |/ _` | '_ \| || |_ / _ \| |/ _` |
| |  | | (_| | | | | ||  _| (_) | | (_| |
|_|  |_|\__,_|_| |_|_||_|  \___/|_|\__,_|
```
# Instalation
```
$ git clone https://github.com/DomesticMoth/manifold.git
$ cd manifold
$ make build
# make install
```
# Description
```
+--------+
| User A |-->|
+--------+   |   +---------+
             |-->| VK chat |
+--------+   |   +----+----+
| User B |-->|        |          Manifold allows you
+--------+            |           to create bridges between chats in
                      |           various messengers and social media
                +-----+------+    to synchronize their contents.
                |            |
                |  Manifold  |
                |            |
                +-----+------+
                      |
+--------+            |
| User C |-->|        |
+--------+   |   +----+----------+
             |-->| Telegram chat |
+--------+   |   +---------------+
| User D |-->|
+--------+
```
There are several ways to connect a chat to the Manifold.
## One bot
The easiest way is to place one bot connected to the  in each chat, through which all messages from other chats will be forwarded.  
```
+-------------+
|Chat 1       |
|+-----------+|                     
|| Bridgebot +---+  +-------------+  
|+-----------+|  |  |Chat 3       |  
+-------------+  |  |+-----------+|
                 +---+ Bridgebot ||
+-------------+  |  |+-----------+|
|Chat 2       |  |  +-------------+
|+-----------+|  |
|| Bridgebot +---+
|+-----------+|
+-------------+
```
But in this case, messages from all users will be forwarded through one account (bot).
```
         +-------------------+
         |<By User A>        |
         |Some text          |
         +-------------------+
         +-------------------+
+------+ |<By User B>        |
|Bot   | |Some other text    |
+------+ +-------------------+

+------+ +-------------------+
|User C| |Yet another text   |
+------+ +-------------------+
```

## Puppets
On the other hand, you can set up an individual puppet bot in one chat for each user of another chat.  
```
+------------+ +------------+
| Chat 1     | | Chat 2     |
|+----------+| |+----------+|
|| Puppet A +<--+ User   A ||
|+----------+| |+----------+|
|+----------+| |+----------+|
|| Puppet B +<--+ User   B ||
|+----------+| |+----------+|
|+----------+| |+----------+|
|| User   C +-->+ Puppet C ||
|+----------+| |+----------+|
|+----------+| |+----------+|
|| User   D +-->+ Puppet D ||
|+----------+| |+----------+|
+------------+ +------------+
```

# Reasons
One of the problems of large (as well as small) Internet communities is the choice of the platform on which members of this community gather.  
Different users prefer different platforms for convenience or security reasons.   
Also, many users are often dissatisfied with the need to have separate clients for some platform just for the sake of a single community.  
  
To solve this problem, two tools have already been created, each with its own strengths and weaknesses. This is a [universal client "pidgin"](https://www.pidgin.im) and [matrix protocol](https://matrix.org).  
Each of them has to make some compromises.  
Manifold is another tool offering an alternative solution to the problem.  

## Comparison
[Pidgin](https://www.pidgin.im) is a chat program which lets you login to account on multiple chat platforms. This saves you from having a lot of clients on your device. However, you still need to have accounts on these platforms. In addition, many services (for example, discord) prohibit alternative clients.  
  
[Matrix](https://matrix.org) is an open protocol for decentralized, secure real-time communications. It aims to create a variety of independent client and server implementations, as well as to create bridges to other platforms. If you need to consolidate a disparate community on different platforms, you can create a central matrix chat and connect it with the rest via bridges.
```
             +--------+
     +-------+ Matrix +-------+
     |       +----+---+       |
     |            |           |
  +--+---+    +---+--+     +--+---+
  |Bridge|    |Bridge|     |Bridge|
  +--+---+    +---+--+     +--+---+
     |            |           |
+----+-----+  +---+-----+  +--+---+
| Telegram |  | Discord |  | XMPP |
+----------+  +---------+  +------+
```
However, in some cases, this approach is not optimal, since messages from different chats require a long chain of many intermediate elements.
```
+---------------+
|    Telegram   |
+-------+-------+
        |
+-------+-------+
|    Bridge 1   |
+-------+-------+
        |
+-------+-------+
| Matrix server |
+-------+-------+
        |
+-------+-------+
|    Bridge 2   |
+-------+-------+
        |
+-------+-------+
|    Discord    |
+---------------+
```
This can create a noticeable delay, especially if you use a third-party matrix server. Also, message metadata may be lost in this pipeline.  
Among other things, this approach is redundant if you need to link two chat rooms that already exist on different platforms and you don't need a third chat in matrix. Moreover, this is redundant if you need to link two chats on the same platform.  
And in the end, the configuration of many individual bridges can be inconvenient.  
  
**Manifold** is in many ways similar to several matrix bridges combined into a single whole. However, unlike the ecosystem of bridges in matrix, **Manifold** exists by itself and does not require raising your own chat server or using someone else's for its work.  
In addition, **Manifold** provides the integration of chats on the principle of p2p, instead of connecting them through the main chat as in matrix.  
Also, **Manifold** has access to all metadata of relayed messages, which allows them to be displayed more correctly in connected chats.  

# Architecture
The main structural element of the Manifold is the "unit". Each unit is a separate entity that interacts with messages passing through the Manifold.  
A unit can be, for example, a connection to a chat on an external platform, a chat bot, or a message logger.  
All units inside a running instance of Manifold are connected to each other according to the p2p principle.  
Each unit can publish events (for example, a chat message event) and receives events from other units.  
  
To mark entities transmitted in events (users, messages, chats, etc.) there is a global system of uint64 identifiers.  
For new entities, a random id is generated using a high-quality RNG and stored. Due to the number of possible variations of the uint64 number, the probability of a collision is extremely small.  
For individual users, IDs can be set explicitly in the configuration.  
Manifold provides a simple ID-based event filtering system for implementing ban lists, etc.  

# Usage
The program accepts a single optional command line argument pointing to the path to the configuration file.  
If the path is not explicitly passed, the program looks for the config in "/etc/manifold/config.toml".  

# Configuration
Manifold configuration file is described in toml format and is located along the path "/etc/manifest/config.toml" or along an arbitrary path passed through command line arguments.  

## Global Parameters
The following global parameters can be specified in the configuration file:  
  
**LogLevel** - Optional string parameter that defines the logging level. Acceptable values (in any case): trace, debug, info, warn, error, fatal, panic. The default value is "Info".  
**Db** - Optional string parameter defining the path to the sqlite database. If the database does not exist, it will be created. The default value is "/etc/manifold/manifold.db".  
**BlockList** - An optional array of numeric IDs globally blocked in the entire instance of the Manifold. Events containing them will be discarded. Empty by default.  
## Units
In addition to global parameters, the configuration contains an array of units named "Unit".  
For each unit, a string field "Name" must be specified that is unique within the config (it is used for log entries, database work, etc.).  
Also, for each unit, you can specify unique lists of blocked identifiers for outgoing (from this unit to the rest) and incoming (from the rest to this unit) events.  
``` toml
[[ Unit ]]
    Name = "Unit 1"

[[ Unit ]]
    Name = "Unit 2"
    BlockListInternal = [0, 1, 2]
    BlockListExternal = [3, 4, 5]
```
In addition, the unit configuration must contain one nested configuration of a specific unit type:
### Log unit
This unit logs all events from all other units to stdout.  
The configuration of units does not contain parameters.  
``` toml
[[ Unit ]]
    Name = "Logger"
    [ Unit.Log ]
```

### Ping unit
This unit responds to any message with the text "ping" or "pong" with the opposite message.  
The configuration of units does not contain parameters.  
``` toml
[[ Unit ]]
    Name = "Ping"
    [ Unit.Ping ]
```

### Vk unit
This unit connects the [Vk](https://vk.com) chat to the Manifold.  
In order to use this unit, you must have a Vk bot added to the chat you are interested in with access to all messages.  
You can also grant admin rights to this bot if you want to manage the chat via Manifold.  
You can read more about Vk chatbots here: [getting-started](https://dev.vk.com/api/bots/getting-started), [bots_docs](https://vk.com/dev/bots_docs).  
There are two mandatory parameters in the configuration of this unit:  
**Token** - Bot token.  
**PeerId** - Chat ID. Note that these identifiers are unique for each bot.  
``` toml
[[ Unit ]]
    Name = "Vk"
    [ Unit.Vk ]
        Token = "<your token here>"
        PeerId = 0 # your chat id here
```
You can also specify a mapping table between Vk user IDs and local IDs used for incoming messages.  
``` toml
[[ Unit ]]
    Name = "Vk"
    [ Unit.Vk ]
        Token = "<your token here>"
        PeerId = 0 # your chat id here
        UsersInc = [{Vk=0, Local=0},
                    {Vk=1, Local=1},
                    {Vk=2, Local=2}]
```
In order to configure individual puppet bots for specific users, you will also need the following optional parameters:  
**Puppet** - This is an array, each element of which is a configuration of the puppet bot.  
``` toml
[[ Unit ]]
    Name = "Vk"
    [ Unit.Vk ]
        Token = "<your token here>"
        PeerId = 0 # your chat id here
        UsersInc = [{Vk=0, Local=0},
                    {Vk=1, Local=1},
                    {Vk=2, Local=2}]
        [[ Unit.Vk.Puppet ]] # Puppet 0
            Token = "<your token here>"
            PeerId = 0 # your chat id here
        [[ Unit.Vk.Puppet ]] # Puppet 1
            Token = "<your token here>"
            PeerId = 0 # your chat id here
        [[ Unit.Vk.Puppet ]] # Puppet 2
            Token = "<your token here>"
            PeerId = 0 # your chat id here
```
Then you have to specify the mapping of local user IDs to these puppets using the **UsersOutg** parameter.  
``` toml
[[ Unit ]]
    Name = "Vk"
    [ Unit.Vk ]
        Token = "<your token here>"
        PeerId = 0 # your chat id here
        UsersInc = [{Vk=0, Local=0},
                    {Vk=1, Local=1},
                    {Vk=2, Local=2}]
        UsersOutg = [{Local=0, Puppet=0},
                     {Local=1, Puppet=1},
                     {Local=2, Puppet=2}]
        [[ Unit.Vk.Puppet ]] # Puppet 0
            Token = "<your token here>"
            PeerId = 0 # your chat id here
        [[ Unit.Vk.Puppet ]] # Puppet 1
            Token = "<your token here>"
            PeerId = 0 # your chat id here
        [[ Unit.Vk.Puppet ]] # Puppet 2
            Token = "<your token here>"
            PeerId = 0 # your chat id here
```
