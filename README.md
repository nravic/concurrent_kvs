# Goncurrent KVS 

## Description
A simple goroutine (and thread)-safe key value store using generics that supports concurrent reads and writes. Primarily built to support heavily read-skewed applications.

## Why 
Wanted to give generics a try and learn more about dealing with concurrency in go. 

## Usage 
Probably a good idea to use a real library instead, this is more of an experiment. Check out https://github.com/orcaman/concurrent-map which employs sharding and is generally more high performance. I do plan on extending this though, which I detail below. 


## Path to production-ready KVS/TODO 
- generic `hashFunction`s
- add value for embedded applications by reducing memory footprint
- add tests!
- memory optimizations? 
- better collision avoidance  