![Manifold](https://raw.githubusercontent.com/DomesticMoth/manifold/main/media/Manifold.png) 
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
The easiest way is to place one bot connected to the Manifold in each chat, through which all messages from other chats will be forwarded.  
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

