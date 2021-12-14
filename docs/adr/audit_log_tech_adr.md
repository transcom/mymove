# History & Audit Logging

MilMove requires an audit log of actions taken on the application. In this application, several different users have different use cases for the audit logging feature.

- Customer user - veteran user who is requesting the move
- Office users - approves moves, helps veterans with their move questions
- Prime users - does the actual moving process for the customer
- Data warehouse users - consumes log data for analytics

As a part of the user requirements gathering, the audit log data must include the following data

- When the action took place
- What action took place
- Who took the action

When considering how to implement the application, there are two main choices worth considering in the Milmove application. NoSQL vs. SQL Route

## Considered Alternatives

### SQL Route

This approach leverages the existing technology to implement the audit log. This is the path of least resitance because engineers do not need to introduce new technology to the application. Introducing new technology will need new considerations such as CI/CD, security, knowledge, Transcom approval, etc. 

#### Feasability 

The more complex we make the underlying tables for the audit logging, the write speeds will be slower. Therefore, if this route is chosen, a flatter database design will help mitigate performance issues especially as the application scales.  

### NoSQL Route

This approach will introduce new technology. If the amount of data that is expected to be stored into the database is extremely large, a NoSQL route will help maintain write speeds. NoSQL by default has a flat data design so there is no need for extra data modeling considerations because of the flexible nature of NoSQL

#### Feasability

Introducing technology will have a lot of hurdles to jump through. Given the state of the Milmove application, a quicker time to launch may be preferred. If this route was chosen, developers will be introducing new AWS services to the application such as DynamoDB, Lambda, API Gateway, etc. 

NoSQL may slow down consumers in retrieving data to our application. Doing any type of bulk work whether getting or updating data proves more dificult with a NoSQL table. Querying the application may be cumbersome because there is unstructured data. However, AWS includes services such as AWS Glue to help normalize unstructured data to more readabale formats. 

## Comparing the Alternatives

### Performance

NoSQL has quicker write speeds so that the application does not slow down due to a new history log call in the API. NoSQL also has workarounds for its worse/slower readability for potential consumers

### Ease of use

SQL route will leverage existing technology to build out the application. Therefore, there is no additional approval needed or infrastructure resources need to build the feature. 
