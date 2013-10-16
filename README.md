# Snowflake

Snowflake is a library for generating unique ID numbers 
on Google App Engine with Golang.

- High Performance
- (Roughly) Time Ordered


## Spec

ID number is composed of:
- time - 41 bits (millisecond precision w/ a custom epoch gives us 69 years)
- random number - 22 bits

A custom epoch of this library is equal to Twitter's Snowflake.


## Random number

Google App Engine Datastore is robust and scalable storage.
For high performace, entity writing should be scattered.

Snowflake uses rundom number in lower 22 bits.
That number will be used to decide Key of Entity for scattering.

Snowflake checks last ID generated time on the entity
to make sure unique.


## Links

- https://github.com/twitter/snowflake
