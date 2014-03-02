CalcGraph
=========

Streaming calculation graph, written in [Go][go-lang].

It is currently a toy project to discover Go's features.

Features
--------

Here are the currently avaiable node types:

* Constant
* Observer
* Sink
* Add/Sub/Mul (variadic, type promotion)
* FanIn
* _(FanOut is automatic and implicit, by branching a node's output to multiple other)_

A node can either be run in a step by step manner by calling its `Run()` method, or it can be looped over.
An `Executor` facility permits running or looping a node concurrently, while being easily interruptible.

An `Assembler` facility permits collecting exactly one value from each input `Node` and returning them as an array.
This helps constructing variadic nodes, as well as synchronizing nodes' work.
It can optionaly verify that each value is of a given type.

License
-------

This project is licensed under the 2-clause BSD license. See `LICENSE` file.

[go-lang]: http://www.golang.org
