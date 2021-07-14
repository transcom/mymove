import interceptorReducer, { initialState } from './reducer';

describe('interceptorReducer', () => {
  it('returns the initial state by default', () => {
    expect(interceptorReducer(undefined, undefined)).toEqual(initialState);
  });
});
