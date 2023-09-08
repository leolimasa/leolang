# Implements a listener that fills up a distributed queue

struct Message
    payload String
    received DateTime


# Actors start listening when they are created and stop when they
# go out of scope.
actor Queue
    generic I

    :: (mut (List I))
    queue = []

    :: (mut (List (Chan I)))
    requests = []

    push = fn (self msg)
        req = pop requests
	if req == None
            queue.push msg
	    (send req msg)


    :: Queue (Chan I)
    request = fn (self)
        :: Chan I
    	c = new-channel
	msg = pop queue
    	if (msg == None)
	    then
		push requests c
		c
	    else
	    	
		
actor Worker
    start = fn (queue)
        # Block on waiting on the channel returned by queue.request
        msg = wait-on (queue.request)
	if (msg == None)
	    start queue
	    else
   	    	process msg
		start queue


start = fn ()
    # Create 4 workers
    workers = map
        fn (i) Worker
	range 4

    # Start all workers with the specified queue
    :: (Queue Message)
    queue = (Queue)
    map 
    	fn (w) (w.start queue)
	workers

    # Push stuff to the queue so it gets processed



