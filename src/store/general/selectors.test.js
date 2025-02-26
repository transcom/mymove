import { selectCurrentMoveId, selectRefetchQueue } from './selectors';

describe('selectCurrentMoveId', () => {
  it('returns the moveId value', () => {
    const testState = {
      generalState: {
        moveId: 'test',
        refetchQueue: false,
      },
    };

    expect(selectCurrentMoveId(testState)).toEqual(testState.generalState.moveId);
    expect(selectRefetchQueue(testState)).toEqual(testState.generalState.refetchQueue);
  });
});
