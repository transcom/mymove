# Use React Query for Office App API interactions

We are currently making heavy use of [Redux](https://redux.js.org/) in both the Office app & customer-facing app to manage API data in the browser. While there are other valid reasons to continue using Redux, it does put a lot of overhead on our code to handle fetching data, providing the correct data to the correct UI components, and deciding when to invalidate or re-fetch data that has become stale. [This article (_Why I Quit Redux_)](https://dev.to/g_abud/why-i-quit-redux-1knl) explains pretty well why this is less than ideal, and also proposes using a library called [React Query](https://react-query.tanstack.com/) as an alternative.

After seeing that React Query follows many of the same conventions as ApolloClient (a client data-layer tool used with GraphQL), I decided to investigate further, since using ApolloClient on previous projects significantly reduced the amount of overhead work when hooking up the API to the frontend. After testing this out with a basic existing query, I believe if we take advantage of this tool, it will both reduce complexity of our frontend codebase and speed up our progress. Please see the [Pros and Cons](#pros-and-cons-of-the-alternatives) section for more technical details and comparison points.

## Considered Alternatives

- Continue using SwaggerRequest thunk as-is (do nothing)
- Use [Redux-Saga](https://redux-saga.js.org/) for managing API interactions
- Use React-Query

## Decision Outcome

- Chosen Alternative: React-Query for Office app (starting with TXO)
- It's important to note that none of these options are incompatible with each other, meaning this is not an all-or-nothing decision. _However_ we want to avoid using too many patterns simultaneously, so I recommend scoping this change to an explicitly defined area of the codebase (in this case, the TXO pages).
- Additionally, there is some risk in React Query co-existing with Redux entities in that it means API data is cached in two locations and could fall out of sync. My recommendation is that we start by using React Query for _all_ TXO pages, and plan to migrate PPM Office pages when they undergo redesigns. Since the customer-facing app has a significantly different use-case and does not involve fetching data simultaneously with the Office pages, migrating that app is less of an immediate concern.
- Assuming this decision is well-communicated and learning resources are provided to the team members responsible for implementation, I believe this will have a significant impact on how quickly we are able to build robust UI for the Office users and meet our deadlines.
- If we begin to use React Query and run into issues or decide that it is not the right choice, this decision can be reversed by changing the code we've written to use one of the alternatives.

## Pros and Cons of the Alternatives

### SwaggerRequest as-is

- `+` no changes needed to what we're already doing (and it's working for us, for now)
- `-` in most places is fetching data on `componentDidMount`, and needs explicit handling to trigger re-fetches if data is changed (such as after submitting a form or `PUT` request), or to re-fetch if props that API calls rely on change (such as the ID of a resource being viewed)
- `-` uses the heavily abstracted SwaggerRequest function, which is currently responsible for all of:
  - dispatching actions to log the request start and success or error
  - making the request itself using SwaggerClient fetch
  - normalizing the data that comes back in the response
  - updating the response data in Redux
  - Because all of this is handled by a single function, it can be difficult to debug, and high-risk to change if, for example, edge-cases need to be handled.
- `-` Requires a significant amount of boilerplate code in order to select data from Redux and handle entity actions & reducer

### Redux-Saga

- `+` gives us the tools to start splitting up the things happening in SwaggerRequest (logging, fetching, normalizing, updating data) iteratively and safely
- `+` makes it easier to chain together fetch flows (load resource A -> load resource B -> etc.)
- `+` makes it easier to test data flows
- `-` continues to rely on Redux for managing API data
- `-` would mean adding even more boilerplate code (for example, to start saga watchers and dispatch actions that trigger sagas)
- `-` has a significant learning curve, for both ES6 generators and the saga effects API
- `-` doesn’t solve the problem of having to explicitly re-fetch or invalidate cached data (but does make it more visible and easier to control)

### React-Query

- `+` removes the need to use Redux for any API data caching, meaning a whole lot of code _could_ be removed
- `+` handles caching and API optimization "for free" (i.e., when you call `useQuery` it will decide whether to re-fetch data or to use cached data) and provides configuration options to fine-tune this
- `+` provides patterns for invalidating cached data, updating data on mutation responses (such as after submitting a form), and chaining API calls. These are all things we do often, and can be difficult to get right when done manually.
- `+` has a few bells and whistles around data-layer management (i.e.: refresh queries on window focus, query retries, pagination/infinite scroll)
- `+` `useQuery` and `useMutation` patterns lend themselves to easily mocking responses from API endpoints that don't exist yet or aren't fully implemented
- `+` not reliant on ES6 generators which need an extra runtime polyfill for some browsers
- `-` would require writing (a) function(s) to use instead of `SwaggerRequest` that have to handle making the API calls and returning normalized data (without Redux)
- `-` dependent on using [hooks](https://reactjs.org/docs/hooks-intro.html) for making API interactions, which can only be used in functional components
- `-` React-Query is a relatively new library (as of Fall 2019) and is also new to us, so there could always be surprises or unknown issues
- `=` documentation is robust and there is a devtools extension for debugging
- `=` API includes [`ReactQueryCacheProvider`](https://react-query.tanstack.com/docs/api#reactquerycacheprovider) which can be used to provide data for unit tests (similar to Redux's Provider)
