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
  
To solve this problem, two tools have already been created, each with its own strengths and weaknesses. This is a [universal client "pidgin"](https://www.pidgin.im/) and [matrix protocol](https://matrix.org/).  
Each of them has to make some compromises.  
Manifold is another tool offering an alternative solution to the problem.  
