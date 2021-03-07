Requirements:

	Each trader can only attempt to load trade data based on the following conditions
	
    Daily data:
        -   Upto 3 load attempts with a total USD value of 5K. Other currencies will follow soon!

    Weekly data:
        - No limit on load attempts with a total USD of 20K.

Solution design:

	A Rest API in Golang. This API will perform a data load based on the input file as per the requirement.

	There are also two simple clients - one using HTML and another using TypeScript.

	Below are the components in the system:
	
	Rest API: Data Models -> Controller -> Service -> Redis Layer -> Redis (on cloud)
		Data Models:
			- We have input and output data models. 
			- Input data model has fields used for loading data and an additional field to keep track of number of transactions for a day.
			- Output data model has fields used for demonstrating the output. This doesnt have any extra fields than the template.

		Controller:
			- API endpoints are defined here that would map to a function. This function passes on the request to service which will operate on the input data.

		Service:
			- The entire business logic for data processing and loading resides here.

		Redis Layer:
			- The persistence / key-value store chosen for this solution is Redis. Service layer passes the information on to this which has a template for the required Redis operations.

		Redis (on cloud):
			- Redis instance chosen for this solution is on cloud.

		Further, there are two supporting packages to hold constants and utilities for this solution.

	Client:
		HTML Client 
			-   A simple html page that can be used to upload the input file. This will call the data load API and will fetch the results of the load status
			-   To load the client, just navigate to the project root and launch the file in the browser: index.html
