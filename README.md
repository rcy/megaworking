M-x work
========

Quickly jump into focussed work. Set goals outcomes for the
work. Eliminate distractions.  Inspired by Ultraworking.

## Irc bot

Bot lives in a channel, like ##megaworking, and runs the time keeping
for work cycles.

Work cycles are 30 minutes of focussed work followed by a 10 minute
break.  This could be configurable.

The bot announces the start of the work time and the end of the work time.

The bot has op status and will voice and devoice users in the channel
so talking can only happen during the breaks.

Taking a break from work and chatting is encouraged.

Ideas
=====

## Emacs client

When you are ready to work, just hit `M-x work RET` in Emacs.

You will be prompted to fill out:
- how long the session will be
- what you want to accomplish

Then, a timer will be set.

## UTC Workcycle schedule 

interval | activity
---|---
00:00-00:10 | break
00:10-00:40 | work
00:40-00:50 | break
00:50-01:20 | work
01:20-01:30 | break
01:30-02:00 | work
...|...

## Installing Headquarters
https://github.com/NunoSempere/ultraworking-headquarters
