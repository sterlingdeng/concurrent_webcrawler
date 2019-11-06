# Concurrent Webcrawler

## The flow
- generate works that will pull url's off of a channel to process them

put first url into the channel

workers pull a url string from the channel to process

workers get all the of urls on the page, creates a struct and fans into a worker that updates 
the adjacency list

worker also pushes the link urls to the queue