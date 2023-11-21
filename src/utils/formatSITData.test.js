import { formatSITData } from './formatSITData';

describe('formatSITData', () => {
  const historyRecord = {
    context: [
      {
        name: 'Domestic destination SIT delivery',
        shipment_type: 'HHG',
        shipment_id_abbr: 'e4285',
        sit_destination_address_final: `{"id":"14a265d6-95b4-4842-a2ed-e020ba7da3fb","street_address_1":"676 Destination Sit Req","street_address_2":null,"city":"Florence","state":"MT","postal_code":"59805","created_at":"2023-11-21T02:56:56.832038","updated_at":"2023-11-21T02:56:56.832038","street_address_3":null,"country":null}`,
        sit_destination_address_initial: `{"id":"ff666bfe-1a2c-45e0-b38a-18c138958f16","street_address_1":"4 Delivery address init","street_address_2":null,"city":"Great Falls","state":"MT","postal_code":"59402","created_at":"2023-11-21T02:56:08.299416","updated_at":"2023-11-21T02:56:08.299416","street_address_3":null,"country":null}`,
        contractor_remarks: 'Need to store in Florence',
      },
    ],
  };

  const address = {
    context: [
      {
        sit_destination_address_final: `{"id":"14a265d6-95b4-4842-a2ed-e020ba7da3fb","street_address_1":"676 Destination Sit Req","street_address_2":null,"city":"Florence","state":"MT","postal_code":"59805","created_at":"2023-11-21T02:56:56.832038","updated_at":"2023-11-21T02:56:56.832038","street_address_3":null,"country":null}`,
      },
    ],
  };

  it('formats JSON into a readable address', () => {
    expect(formatSITData(address)).toEqual({
      sit_destination_address_final: '676 Destination Sit Req, Florence, MT 59805',
    });
  });

  it('formats a historyRecord object properly', () => {
    expect(formatSITData(historyRecord)).toEqual({
      sit_destination_address_final: '676 Destination Sit Req, Florence, MT 59805',
      sit_destination_address_initial: '4 Delivery address init, Great Falls, MT 59402',
      contractor_remarks: 'Need to store in Florence',
    });
  });
});
