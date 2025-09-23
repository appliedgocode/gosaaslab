# The Go SaaS Lab

A minimalist SaaS platform, developed for the [Applied Go Weekly Newsletter](https://newsletter.appliedgo.net/archive). I'm building a platform with *just enough* features to be a viable basis for building SaaS apps or maybe just a web presence for a small business. 

## Why? 

For small business, sophisticated application stacks are an overkill and probably are more of a maintenance burden than a quick start. This project aims at demonstrating that it doesn't need a cluster of Kubernetes clusters with microservices and distributed databases and Redis caches everywhere to serve web apps to clients. If you're not aiming to become the next Netflix serving movies to hundreds of millions of subscribers, you can (and should) dare to start smaller than that. 

This repo is a proof of concept. I focus on understanding the fundamental ideas and concepts behind each part of the platform, with code as clear as possible. Take the code and expand it to meet your needs, or use the acquired insights to select the optimal 3rd-party solution for your business.

## The series

I plan to publish a weekly "Spotlight" section in my newsletter (link above) that adds new functionality to the project, discussing the ins and outs as well as alternatives. This repository will grow alongside these weekly updates, with Git tags marking each step of this journey. 

To follow the journey, check out each tag from oldest to newest, and see what got added and changed. Or if you just want to study or use the code, check out each component's documentation.

## The project's chronology

### 1. A not-so-basic Web server

...serving a static page. Go famously lets you whip together a web server in a few lines of code. But the typical minimal example lacks two crucial ingredients: proper timeouts and graceful shutdown.

[01 - A Core Web Server](./doc/01-web-server.md)

*(to be continued)*