iso8601duration
===============

ISO8601 Duration Parser for Golang

Inspired by https://github.com/ChannelMeter/iso8601duration, which was in turn
adapted from http://github.com/BrianHicks/finch.

The main difference between this package and ChannelMeter's is that this package
does not define its own `Duration` type -- Golang's `time.Duration` is used
instead. This choice was made from the perspective that `time.Duration` is often
what you want anyway. A consequence is that you cannot "round-trip" from an
ISO8601-formatted duration, through `time.Duration`, and back to ISO8601 with
the guarantee that you'll get the same string. You should however be able to get
an equivalent ISO8601 value.

Also, this package supports decimal fractions in the smallest time value, e.g.
`PT0.25M` is 15 seconds, `PT0.001S` is 1 millisecond, etc.
