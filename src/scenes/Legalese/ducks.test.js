import { CREATE_SIGNED_CERT, dateToTimestamp, signedCertificationReducer } from './ducks';

describe('Feedback Reducer', () => {
  it('Should handle CREATE_CERTIFICATION_SUCCESS', () => {
    const initialState = { pendingValue: '', confirmationText: '' };

    const newState = signedCertificationReducer(initialState, {
      type: CREATE_SIGNED_CERT.success,
      item: 'Successful item!',
    });

    expect(newState).toEqual({
      pendingValue: '',
      confirmationText: 'Feedback submitted!',
      hasSubmitError: false,
      hasSubmitSuccess: true,
    });
  });

  it('Should handle CREATE_CERTIFICATION_FAILURE', () => {
    const initialState = { pendingValue: '', confirmationText: '' };

    const newState = signedCertificationReducer(initialState, {
      type: CREATE_SIGNED_CERT.failure,
      error: 'No bueno.',
    });

    expect(newState).toEqual({
      pendingValue: '',
      confirmationText: 'Submission error.',
      hasSubmitError: true,
      hasSubmitSuccess: false,
    });
  });
  describe('format to swagger date-time in users timezone', () => {
    it('adds timestamp to date', () => {
      const dt = dateToTimestamp('2017-07-21');
      expect(dt).toBe('2017-07-21T00:00:00+00:00');
    });
    it('does nothing if passed a timestamp', () => {
      const ts = '2017-07-21T00:00:00+00:00';
      const dt = dateToTimestamp(ts);
      expect(dt).toBe('2017-07-21T00:00:00+00:00');
    });
  });
});
