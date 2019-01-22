const timer =
  typeof performance !== 'undefined' && performance !== null && typeof performance.now === 'function'
    ? performance
    : Date;
export default function logger({ getState }) {
  return next => action => {
    const logEntry = {};
    let returnedValue;
    logEntry.started = timer.now();
    logEntry.startedTime = new Date();
    logEntry.prevState = getState();
    logEntry.action = action;
    try {
      returnedValue = next(action);
    } catch (e) {
      logEntry.error = e;
    }
    logEntry.took = timer.now() - logEntry.started;
    logEntry.nextState = getState();
    console.log(logEntry.action.type, ' will dispatch ', logEntry); // eslint-disable-line
    return returnedValue;
  };
}
