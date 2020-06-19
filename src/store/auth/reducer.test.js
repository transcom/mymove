import authReducer, { initialState } from './reducer';
import { setActiveRole } from './actions';

describe('authReducer', () => {
  it('returns the initial state by default', () => {
    expect(authReducer(undefined, undefined)).toEqual(initialState);
  });

  it('returns the existing state if activeRole is set for unhandled actions', () => {
    const currentState = {
      ...initialState,
      activeRole: 'myRole',
    };

    expect(authReducer(currentState, undefined)).toEqual(currentState);
  });

  it('handles the setActiveRole action', () => {
    expect(authReducer(initialState, setActiveRole('myRole'))).toEqual({
      ...initialState,
      activeRole: 'myRole',
    });
  });
});
