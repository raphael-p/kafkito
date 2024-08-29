package resolvers

var help string = `All commands (except help) require the config file to be properly defined, especially the port that the kafkito server should be running on.

Commands:
	help

	start	- starts a kafkito server

	stop	- stops the server if it is running (asynchronous)

	info	- tells you if the server is running or not

	create	- creates a new queue
		args:
			queueName	- the name given to the new queue
		usage: 'kafkito create <queueName>'

	rename	- renames an existing queue
		args:
			oldQueueName
			newQueueName
		usage: 'kafkito rename <oldQueueName> <newQueueName>'

	delete	- deletes an existing queue
		args:
			queueName
		usage: 'kafkito delete <queueName>'

	list	- lists all queues or messages of a queue (message body may be truncated)
		args:
			queueName (optional)
		usage:
			'kafkito list'
			'kafkito list <queueName>'

	read	- prints the full body of a message
		args:
			messageID	- can be found by listing all messages in a queue
		usage: 'kafkito read <messageID>'

	consume	- removes a message from the queue and prints its full body
		args:
			messageID	- can be found by listing all messages in a queue
		usage: 'kafkito consume <messageID>'

	publish - adds a message to the queue
		args:
			queueName
			messageHeader	- short title or category
			messageBody		- the full message
		usage: 'kafkito publish <queueName> <messageHeader> <messageBody>'
`
