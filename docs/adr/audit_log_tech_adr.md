# Activity Log

**User Story:** *[ticket/issue-number]*

MilMove requires an audit log of actions taken on the application. In this application, several different users have different use cases for the audit logging feature.

- Customer user - veteran user who is requesting the move
- Office users - approves moves, helps veterans with their move questions
- Prime users - does the actual moving process for the customer
- Data warehouse users - consumes log data for analytics

As a part of the user requirements gathering, the audit log data must include the following data

- When the action took place
- What action took place
- Who took the action

## Considered Alternatives

- SQL as-is
- SQL with new messaging tech
- NoSQL

## Decision Outcome

Between SQL as-is and SQL with messaging tech, if we start with SQL as-is, we can later iterate and switch to the messaging system if performance becomes an issue

- Chosen Alternative: TBD

## Pros and Cons of the Alternatives

### *SQL as-is*

- `+` This approach leverages the existing technology to implement the audit log. 
- `+` This is the path of least resitance because engineers do not need to introduce new technology to the application. Introducing new technology will need new considerations such as CI/CD, security, knowledge, Transcom approval, etc. 
- `-` Write speeds may slow down the application especially now that the application is introducing many new writes to the DB.
- `-` The more complex we make the underlying tables for the audit logging, the write speeds will be slower. Therefore, if this route is chosen, a flatter database design will help mitigate performance issues especially as the application scales.  

### SQL with messaging system

This option uses the same API and the same DB currently used, but adds a queue-based messaging system such as AWS Eventbridge and AWS 

- `+` Leverages exisitng systems
- `+` Alleviates pressure on DB writes 
- `+` Since the user does not care whether an audit record was created, we can decouple the audit creating process from the user experience
- `+` No significant amount of security considerations needed
- `-` Introduces some amount of new technology

### NoSQL

This includes options such as DynamoDB, MongoDB, etc.

- `+` NoSQL will help maintain write speeds
- `+` No need for extensive data modeling considerations
- `+` More scalable due to lower storage costs. If using DynamoDB
- `-` Querying data is not as straightforward. Will need additional steps to make it digestable for consumers such as the IGC datawarehouse team and Milmove UI
- `-` NoSQL may slow down consumers in retrieving data to our application. Doing any type of bulk work whether getting or updating data proves more dificult with a NoSQL table. 
