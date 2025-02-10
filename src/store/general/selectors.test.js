import { selectCurrentMoveId, selectShouldRefetchQueue } from './selectors';

describe('selectCurrentMoveId', () => {
  it('returns the moveId value', () => {
    const testState = {
      generalState: {
        moveId: 'test',
        shouldRefetchQueue: false,
      },
    };

    expect(selectCurrentMoveId(testState)).toEqual(testState.generalState.moveId);
    expect(selectShouldRefetchQueue(testState)).toEqual(testState.generalState.shouldRefetchQueue);
  });
});
