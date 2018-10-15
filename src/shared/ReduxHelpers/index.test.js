import * as helpers from '.';

describe('reduxHelpers', () => {
  describe('when generateAsyncActionTypes is called', () => {
    describe('error conditions', () => {
      describe('if there is no resource provided', () => {
        it('should return an error', () => {
          expect(() => helpers.generateAsyncActionTypes()).toThrow();
        });
      });
      describe('if resource is empty string', () => {
        it('should return an error', () => {
          expect(() => helpers.generateAsyncActionTypes('')).toThrow();
        });
      });
      describe('if resource is only whitespace', () => {
        it('should return an error', () => {
          expect(() => helpers.generateAsyncActionTypes(' ')).toThrow();
        });
      });
    });
  });
  describe('given a resource', () => {
    const resourceName = 'RESOURCE';
    describe('when generateAsyncActionTypes is called', () => {
      const actions = helpers.generateAsyncActionTypes(resourceName);
      it('should give me a start action', () => {
        expect(actions).toEqual(expect.objectContaining({ start: 'RESOURCE_START' }));
      });
      it('should give me a success action', () => {
        expect(actions).toEqual(expect.objectContaining({ success: 'RESOURCE_SUCCESS' }));
      });
      it('should give me a failure action', () => {
        expect(actions).toEqual(expect.objectContaining({ failure: 'RESOURCE_FAILURE' }));
      });
    });
    describe('when generateAsyncActionCreator is called', () => {
      describe('when the async call is is successful', () => {
        const mockAsyncAction = jest.fn().mockImplementation(() => Promise.resolve('foo'));
        const dispatch = jest.fn();
        const actions = helpers.generateAsyncActions(resourceName);
        const actionCreator = helpers.generateAsyncActionCreator(resourceName, mockAsyncAction);

        const func = actionCreator();
        func(dispatch);
        it('it should call the asyncAction', () => {
          expect(mockAsyncAction.mock.calls.length).toBe(1);
        });
        it('it should dispatch the start', () => {
          expect(dispatch).toBeCalledWith(actions.start());
        });
        it('it should dispatch the success with the payload of foo', () => {
          expect(dispatch).lastCalledWith(actions.success('foo'));
        });
      });
      describe('when async call fails', () => {
        const mockAsyncAction = jest.fn().mockImplementation(() => Promise.reject('something went wrong'));
        const dispatch = jest.fn();
        const actions = helpers.generateAsyncActions(resourceName);
        const actionCreator = helpers.generateAsyncActionCreator(resourceName, mockAsyncAction);
        const func = actionCreator();
        func(dispatch);
        it('it should call the asyncAction', () => {
          expect(mockAsyncAction.mock.calls.length).toBe(1);
        });
        it('it should dispatch the start', () => {
          expect(dispatch).toBeCalledWith(actions.start());
        });
        it('it should dispatch the failure with the payload of foo', () => {
          expect(dispatch).lastCalledWith(actions.error('something went wrong'));
        });
      });
    });
    describe('when generateAsyncReducer is called', () => {
      const actions = helpers.generateAsyncActions(resourceName);
      const onSuccessTransform = jest.fn().mockImplementation(foo => ({ foo }));
      const reducer = helpers.generateAsyncReducer(resourceName, onSuccessTransform);
      it('random', () => {
        const initialState = {};
        Object.freeze(initialState);
        const newstate = reducer(initialState, { type: 'unsupported' });
        expect(newstate).toEqual({});
      });
      it('start', () => {
        const initialState = {};
        Object.freeze(initialState);
        const newstate = reducer(initialState, actions.start());
        expect(newstate).toEqual({
          isLoading: true,
          hasErrored: false,
          hasSucceeded: false,
        });
      });
      it('success', () => {
        const initialState = {};
        Object.freeze(initialState);
        const newstate = reducer(initialState, actions.success('my payload'));
        expect(newstate).toEqual({
          foo: 'my payload',
          isLoading: false,
          hasErrored: false,
          hasSucceeded: true,
        });
      });
      it('failure', () => {
        const initialState = {};
        Object.freeze(initialState);
        const newstate = reducer(initialState, actions.error('something is wrong'));
        expect(newstate).toEqual({
          error: 'something is wrong',
          isLoading: false,
          hasErrored: true,
          hasSucceeded: false,
        });
      });
    });
  });
});
