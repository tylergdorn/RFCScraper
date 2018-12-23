# RFCScraper

RFCScraper is a tool written in golang that lets you download as many RFC's published on the <a href="https://www.ietf.org/standards/rfcs/">IEFT website</a>
It is a rewrite of a previous tool I wrote to do the same thing but in python, and not multithreaded (go makes this very easy to)

You can also view one printed to the console

Usage: `./RFCScraper STARTNUMBER ENDNUMBER`
The start and end numbers should be real numbers and start should be smaller than end.

To view one you can do the following: `./RFCScraper -view NUMBER`
This will print the RFC out to the console

By default this tool saves to a folder `rfc` it makes in the current working directory