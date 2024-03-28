import { selectCurrentMoveId } from './selectors';

describe('selectCurrentMoveId', () => {
  it('returns the moveId value', () => {
    const testState = {
      generalState: {
        moveId: 'test',
      },
    };

    expect(selectCurrentMoveId(testState)).toEqual(testState.generalState.moveId);
  });
});
