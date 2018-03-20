import {
  CREATE_DOCUMENT_SUCCESS,
  CREATE_DOCUMENT_FAILURE,
  documentReducer,
} from './ducks';

describe('Document Reducer', () => {
  it('Should handle CREATE_DOCUMENT_SUCCESS', () => {
    const initialState = {
      hasSubmitError: false,
      hasSubmitSuccess: false,
      confirmationText: '',
    };

    const newState = documentReducer(initialState, {
      type: CREATE_DOCUMENT_SUCCESS,
      item: 'Successful item!',
    });

    expect(newState).toEqual({
      confirmationText: 'Document uploaded!',
      hasSubmitError: false,
      hasSubmitSuccess: true,
    });
  });

  it('Should handle CREATE_DOCUMENT_FAILURE', () => {
    const initialState = {
      hasSubmitError: false,
      hasSubmitSuccess: false,
      confirmationText: '',
    };

    const newState = documentReducer(initialState, {
      type: CREATE_DOCUMENT_FAILURE,
      error: 'No bueno.',
    });

    expect(newState).toEqual({
      confirmationText: 'Upload error.',
      hasSubmitError: true,
      hasSubmitSuccess: false,
    });
  });
});
