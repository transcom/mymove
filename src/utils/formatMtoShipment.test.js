import { formatMtoShipmentForAPI } from './formatMtoShipment';

import { SHIPMENT_OPTIONS } from 'shared/constants';

describe('formatMtoShipmentForAPI', () => {
  const mtoShipmentParams = {
    moveId: 'move123',
    customerRemarks: 'some mock remarks',
  };

  const pickupInfo = {
    requestedDate: 'Jan 7, 2026',
    address: {
      street_address_1: '123 main',
      city: 'legit human city',
      state: 'DC',
      postal_code: '20017',
    },
    agent: {
      firstName: 'mockFirstName',
      lastName: 'mockLastName',
      email: 'mockAgentEmail@example.com',
      phone: '2225551234',
    },
  };

  const deliveryInfo = {
    requestedDate: 'Jan 27, 2026',
    address: {
      street_address_1: '0011100010110101',
      city: 'R0B0T T0WN',
      state: 'CP',
      postal_code: '10101',
    },
    agent: {
      firstName: 'r0b0tBestFr1end',
      lastName: 'r0b0tBestFr1endLastName',
      email: 'r0b0t-fr1end@example.com',
      phone: '2225550101',
    },
  };

  it('can format an HHG shipment', () => {
    const params = {
      ...mtoShipmentParams,
      shipmentType: SHIPMENT_OPTIONS.HHG,
      pickup: { ...pickupInfo },
      delivery: { ...deliveryInfo },
    };

    const actual = formatMtoShipmentForAPI(params);

    expect(actual.shipmentType).toBe(SHIPMENT_OPTIONS.HHG);
    expect(actual.agents.length).toBe(2);
    expect(actual.requestedPickupDate).toBe('2026-01-07');
    expect(actual.agents[0].phone).toBe('222-555-1234');
    expect(actual.agents[0].agentType).toBe('RELEASING_AGENT');
    expect(actual.requestedDeliveryDate).toBe('2026-01-27');
    expect(actual.agents[1].phone).toBe('222-555-0101');
    expect(actual.agents[1].agentType).toBe('RECEIVING_AGENT');
    expect(actual.customerRemarks).toBe('some mock remarks');

    expect(actual.secondaryPickupAddress).toBeUndefined();
    expect(actual.secondaryDeliveryAddress).toBeUndefined();
  });

  it('can format an NTSr shipment', () => {
    const params = {
      ...mtoShipmentParams,
      shipmentType: SHIPMENT_OPTIONS.NTSR,
      pickup: { ...pickupInfo },
    };

    const actual = formatMtoShipmentForAPI(params);

    expect(actual.shipmentType).toBe(SHIPMENT_OPTIONS.NTSR);
    expect(actual.requestedPickupDate).toBe('2026-01-07');
    expect(actual.agents.length).toBe(1);
    expect(actual.agents[0].phone).toBe('222-555-1234');
    expect(actual.agents[0].agentType).toBe('RELEASING_AGENT');
    expect(actual.customerRemarks).toBe('some mock remarks');
  });

  it('can format an NTS shipment', () => {
    const params = {
      ...mtoShipmentParams,
      shipmentType: SHIPMENT_OPTIONS.NTS,
      delivery: { ...deliveryInfo },
    };
    const actual = formatMtoShipmentForAPI(params);
    expect(actual.shipmentType).toBe(SHIPMENT_OPTIONS.NTS);
    expect(actual.requestedDeliveryDate).toBe('2026-01-27');
    expect(actual.agents.length).toBe(1);
    expect(actual.agents[0].phone).toBe('222-555-0101');
    expect(actual.agents[0].agentType).toBe('RECEIVING_AGENT');
    expect(actual.customerRemarks).toBe('some mock remarks');
  });

  it('can format an HHG shipment with a secondary pickup/destination', () => {
    const params = {
      ...mtoShipmentParams,
      shipmentType: SHIPMENT_OPTIONS.HHG,
      pickup: { ...pickupInfo },
      secondaryPickup: {
        address: {
          street_address_1: '142 E Barrel Hoop Circle',
          street_address_2: '#4A',
          city: 'Corpus Christi',
          state: 'TX',
          postal_code: '78412',
        },
      },
      delivery: { ...deliveryInfo },
      secondaryDelivery: {
        address: {
          street_address_1: '441 SW Río de la Plata Drive',
          street_address_2: '',
          city: deliveryInfo.address.city,
          state: deliveryInfo.address.state,
          postal_code: deliveryInfo.address.postal_code,
        },
      },
    };

    const actual = formatMtoShipmentForAPI(params);

    expect(actual.secondaryPickupAddress).not.toBeUndefined();
    expect(actual.secondaryPickupAddress.street_address_1).toEqual('142 E Barrel Hoop Circle');

    expect(actual.secondaryDeliveryAddress).not.toBeUndefined();
    expect(actual.secondaryDeliveryAddress.street_address_1).toEqual('441 SW Río de la Plata Drive');
  });
});
