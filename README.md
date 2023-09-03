# Electronic program guidereadme

I chose to implement this in Go, thinking Go being so simple, everyone can read it.
I could probably have done a lot of things smarter by choosing a language that supports some functional paradigms, but I think the code came out pretty simple in the end.

I decided not to add abstractions or layers that wasn't needed, simply because this is a simple input/output problem, thus we don't need any abstractions to test the code properly.

I realise not everyone has Go installed though, so I provided a Dockerfile so it's easier to run this:
```
docker build -t epg-service .
docker run -it -p 8080:8080 epg-service
```

We can then simply hit the service using for example curl:
```
curl -vX POST http://localhost:8080 -d @payload.json
```

## Assumptions

I found a few inconsistencies in the spec and the example/input data provided, I'll try to cover these here as best I can:

#### Multiple begins for the same program in a row.
In the example payload provided, we can see the example input for monday, has "Dybvaaaaad" with 2 begins and no ends:
```
 { 
 "title": "Dybvaaaaad", 
 "state": "begin", 
 "time": 36000 
 }, 
 { 
 "title": "Dybvaaaaad", 
 "state": "begin", 
 "time": 38100 
 } 
```

I scanned through the document multiple times, but could not find anything mentioning this case, and as special cases are already highlighted, I assumed this to be an error in the input.

## Rethinking the Datamodel

#### Locale
The current datamodel comes with a few challenges, most importantly when dealing with time. The data can be laid out in a multitude of ways, however the way time is currently being managed would lead to curiosities when mixing locales, for example:
We send in a model with time 3600 on a monday, but I am in a Datamodel-6 timezone. This entry should no longer belong to a monday but drop back to a sunday.

The problem here is that normally Unix timestamp are absolute from a specific point in time, but here we keep resetting the time per day, thus losing the absolute trait, and become dependant on being within the same day.

#### Identity
Say we need the structure to stay as is, one beneficial trait could be to add an identity to each, and not rely on names. Names are subject to changing and could lead to bugs if renaming is not done properly.
While in the scope of this task it wouldn't matter, it's always good practice to rely on unique identities rather than data carrying fields.


#### Ordering
Let's assume that this payload is used entirely for the purpose of being run through this program, and nothing else, thus we can model the structure of this data entirely to fit our usecase here, a beneficial trait would be to add ordering.
Maps are not ordered, thus a map consisting of weekdays as string needs to be ordered. (For the way I decided to implement things :) )

We could easily achieve this by using an array and use the index as the index of the day of the week, 0 -> monday, 1 -> tuesday on so on, as I have done in the code.

#### Structure
As above, let's assume the payload is entirely just for our little program here, we can almost solve the entire thing, just by modelling data a little differently.
Putting together all the points from above into this, and fiddeling a bit with the layout, we can come up with something like this:
```
{
    "program": [
        { 
            "id": "cd2ee67b-e32c-4eff-8e34-eedb0b06bb50",
            "title": "Nyhederne",
            "slots": [
                { "start": 1693800000, "end": 1693814400 },
                { "start": 1693854000, "end": 1693855800 }
            ]
        },
        {
            "id": "967d1afa-a620-46c1-86f9-b2b34d1ada1b",
            "title": "Dybvaaaaad",
            "slots": [
                { "start": 1693857600, "end": 1693859700 }
            ]
        },
        ...
        {
            "id": "75277a5d-793c-4307-8ee6-505a8dcb66a3",
            "title": "Nybyggerne",
            "slots": [
                { "start": 1694292000, "end": 1694296800 }
            ]
        },
        {
            "id": "4369b2c9-c136-4b06-b6a1-a4ab91375008",
            "title": "Dybvaaaaad",
            "slots": [
                { "start": 1694296800, "end": 1694296800 }
            ]
        },


    ]
}
```

This layout of data reads quite simply in the context of this program: Print what you have.
This simplifies the state by maintaining data that belongs together in the same structure, and uses proper unix timestamps, meaning we make no assumptions about days.


