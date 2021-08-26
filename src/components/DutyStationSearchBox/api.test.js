import { SearchDutyStations, ShowAddress } from 'scenes/ServiceMembers/api';

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
        duty_stations: {
          searchDutyStations: ({ search }) => {
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
  describe('SearchDutyStations', () => {
    it('retrieves a response from the server', async () => {
      const response = await SearchDutyStations('ok');
      expect(response).toEqual('queried ok');
    });

    it('throws an error when appropriate', async () => {
      await expect(async () => {
        await SearchDutyStations('broken');
      }).rejects.toThrow('failed to query duty stations due to server error');
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
      }).rejects.toThrow('failed to query address for duty station');
    });
  });
});
