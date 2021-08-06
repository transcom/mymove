/*
 * This code exists to allow a requestInterceptor to arbitrarily dispatch events to the
 * active redux store. Redux, by design, doesn't want to do this! But:
 *
 * 1. The most sensible place to put universal response interception for analysis was
 *    attaching it to the SwaggerClient, which by its nature is going to be disconnected
 *    from the redux store scope and the actions that caused it to fire in the first place
 *
 * 2. All of the common methods for redux side effects (sagas and thunks, which we use elsewhere)
 *    need to be connected logically to redux actions that caused them to fire; this isn't
 *    a useful pattern for us unfortunately, as there's a ton of user actions that result in
 *    swagger requests being made. To set up a saga watcher we'd need to exhaustively track them
 *    all, and somehow incorporate adding any future actions which result in XHRs into the array
 *    of events being watched.
 *
 * The conclusion I came to was just closing over a reference to the redux store's dispatch, and
 * then exporting a function that has access to that closure, for use in the requestInterceptor.
 */

/*
 * Start out a closed over reference as a no-op function.
 * There shouldn't be a logical way for this not be defined when a user makes a request,
 * as several actions get fired through the application before the user will have a chance
 * to make any XHR request, but let's just be on the safe side and prevent
 * "cannot execute undefined" errors as a class.
 */
let dispatchReference = () => {};

/*
 * Create some middleware, which is a double curried function. The outermost function gets a
 * reference to the active redux store (in the shape of { getState(), dispatch() }) and we
 * only need the dispatch half of it. This does nothing practical to the application
 * besides getting and closing over a reference to the store's dispatch method.
 */
export const interceptorInjectionMiddleware =
  ({ dispatch }) =>
  (next) =>
  (action) => {
    dispatchReference = dispatch;
    return next(action);
  };

/*
 * Export a function that has access to the closed over dispatch, which can get used by
 * the responseInterceptor attached to the SwaggerClient instance.
 */
export const interceptInjection = (action) => {
  dispatchReference(action);
};
