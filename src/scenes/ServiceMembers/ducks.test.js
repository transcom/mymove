import {
  UPDATE_SERVICE_MEMBER,
  GET_SERVICE_MEMBER,
  serviceMemberReducer,
} from './ducks';

describe('Service Member Reducer', () => {
  const sampleServiceMember = { id: 'UUID', first_name: 'bob' };
  describe('UPDATE_SERVICE_MEMBER', () => {
    it('Should handle UPDATE_SERVICE_MEMBER_SUCCESS', () => {
      const initialState = { currentServiceMember: null };

      const newState = serviceMemberReducer(initialState, {
        type: UPDATE_SERVICE_MEMBER.success,
        payload: sampleServiceMember,
      });

      expect(newState).toEqual({
        currentServiceMember: sampleServiceMember,
        hasSubmitError: false,
        hasSubmitSuccess: true,
      });
    });

    it('Should handle UPDATE_SERVICE_MEMBER_FAILURE', () => {
      const initialState = { currentServiceMember: { id: 'bad' } };

      const newState = serviceMemberReducer(initialState, {
        type: UPDATE_SERVICE_MEMBER.failure,
        error: 'No bueno.',
      });

      expect(newState).toEqual({
        currentServiceMember: { id: 'bad' },
        hasSubmitError: true,
        hasSubmitSuccess: false,
        error: 'No bueno.',
      });
    });
  });
  describe('GET_SERVICE_MEMBER', () => {
    it('Should handle GET_SERVICE_MEMBER_SUCCESS', () => {
      const initialState = { currentServiceMember: null };
      const newState = serviceMemberReducer(initialState, {
        type: GET_SERVICE_MEMBER.success,
        payload: sampleServiceMember,
      });

      expect(newState).toEqual({
        currentServiceMember: sampleServiceMember,
        hasSubmitError: false,
        hasSubmitSuccess: true,
      });
    });

    it('Should handle GET_SERVICE_MEMBER_FAILURE', () => {
      const initialState = { currentServiceMember: null };

      const newState = serviceMemberReducer(initialState, {
        type: GET_SERVICE_MEMBER.failure,
        error: 'No bueno.',
      });

      expect(newState).toEqual({
        currentServiceMember: null,
        hasSubmitError: true,
        hasSubmitSuccess: false,
        error: 'No bueno.',
      });
    });
  });
});
