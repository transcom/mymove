import { CREATE_SIGNED_CERT, signedCertificationReducer } from './ducks';

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
});
