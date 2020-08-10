import * as moveSetupHelpers from './moveSetupFlow';

describe('moveSetupFlow utils', () => {
  describe('getFullAgentName', () => {
    const agent = {
      firstName: 'Bob',
      lastName: 'Bobson',
    };
    const fullName = moveSetupHelpers.getFullAgentName(agent);
    it('should be concatenated as expected', () => {
      expect(fullName).toEqual('Bob Bobson');
    });
  });
  describe('getFullSMName', () => {
    const serviceMember = {
      first_name: 'Bob',
      middle_name: 'Belilah',
      last_name: 'Bobson',
      suffix: 'Jr.',
    };
    const fullName = moveSetupHelpers.getFullSMName(serviceMember);
    it('should be concatenated as expected', () => {
      expect(fullName).toEqual('Bob Belilah Bobson Jr.');
    });
  });
});
