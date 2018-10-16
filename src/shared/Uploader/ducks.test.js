import { CREATE_DOCUMENT_SUCCESS, CREATE_DOCUMENT_FAILURE, documentReducer } from './ducks';

describe('Document Reducer', () => {
  it('Should handle CREATE_DOCUMENT_SUCCESS', () => {
    const initialState = {
      hasErrored: false,
      hasSucceeded: false,
      upload: null,
    };

    const newState = documentReducer(initialState, {
      type: CREATE_DOCUMENT_SUCCESS,
      item: { url: 'nino.com', filename: 'nino', type: 'image' },
    });

    expect(newState).toEqual({
      upload: { url: 'nino.com', filename: 'nino', type: 'image' },
      hasErrored: false,
      hasSucceeded: true,
    });
  });

  it('Should handle CREATE_DOCUMENT_FAILURE', () => {
    const initialState = {
      hasErrored: false,
      hasSucceeded: false,
      upload: null,
    };

    const newState = documentReducer(initialState, {
      type: CREATE_DOCUMENT_FAILURE,
      error: 'No bueno.',
    });

    expect(newState).toEqual({
      upload: {},
      hasErrored: true,
      hasSucceeded: false,
    });
  });
});
