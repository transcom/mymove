import { SearchDutyLocations, ShowAddress } from './api';

jest.mock('shared/Swagger/api.js', () => ({
  ...jest.requireActual('shared/Swagger/api.js'),
  getClient: async () => {
    return {
      apis: {
        addresses: {
          showAddress: ({ addressId }) => {
            if (addressId === 'broken') {
              return { ok: false };
            }

            return { ok: true, body: `address ${addressId}` };
          },
        },
        duty_locations: {
          searchDutyLocations: ({ search }) => {
            if (search === 'broken') {
              return { ok: false };
            }

            return { ok: true, body: `queried ${search}` };
          },
        },
      },
    };
  },
}));

describe('scenes ServiceMembers api', () => {
  describe('SearchDutyLocations', () => {
    it('retrieves a response from the server', async () => {
      const response = await SearchDutyLocations('ok');
      expect(response).toEqual('queried ok');
    });

    it('throws an error when appropriate', async () => {
      await expect(async () => {
        await SearchDutyLocations('broken');
      }).rejects.toThrow('failed to query duty locations due to server error');
    });
  });

  describe('ShowAddress', () => {
    it('retrieves a response from the server', async () => {
      const response = await ShowAddress('ok');
      expect(response).toEqual('address ok');
    });

    it('throws an error when appropriate', async () => {
      await expect(async () => {
        await ShowAddress('broken');
      }).rejects.toThrow('failed to query address for duty location');
    });
  });
});
