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
