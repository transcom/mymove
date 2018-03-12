import {
  CREATE_CERTIFICATION_SUCCESS,
  CREATE_CERTIFICATION_FAILURE,
  signedCertificationReducer,
} from './ducks';

describe('Feedback Reducer', () => {
  it('Should handle CREATE_CERTIFICATION_SUCCESS', () => {
    const initialState = { pendingValue: '', confirmationText: '' };

    const newState = signedCertificationReducer(initialState, {
      type: CREATE_CERTIFICATION_SUCCESS,
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
      type: CREATE_CERTIFICATION_FAILURE,
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
