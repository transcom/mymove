import { formatMtoShipmentForAPI } from './formatMtoShipment';

import { SHIPMENT_OPTIONS } from 'shared/constants';

describe('formatMtoShipmentForAPI', () => {
  it('can format an HHG shipment', () => {
    const params = {
      moveId: 'move123',
      shipmentType: SHIPMENT_OPTIONS.HHG,
      customerRemarks: 'some mock remarks',
      pickup: {
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
      },
      delivery: {
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
      },
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
  });

  it('can format an NTSr shipment', () => {
    const params = {
      moveId: 'move123',
      shipmentType: SHIPMENT_OPTIONS.NTSR,
      customerRemarks: 'some mock remarks',
      pickup: {
        requestedDate: 'Jan 27, 2026',
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
      },
    };
    const actual = formatMtoShipmentForAPI(params);
    expect(actual.shipmentType).toBe(SHIPMENT_OPTIONS.NTSR);
    expect(actual.requestedPickupDate).toBe('2026-01-27');
    expect(actual.agents.length).toBe(1);
    expect(actual.agents[0].phone).toBe('222-555-1234');
    expect(actual.agents[0].agentType).toBe('RELEASING_AGENT');
    expect(actual.customerRemarks).toBe('some mock remarks');
  });

  it('can format an NTS shipment', () => {
    const params = {
      moveId: 'move123',
      shipmentType: SHIPMENT_OPTIONS.NTS,
      customerRemarks: 'some mock remarks',
      delivery: {
        requestedDate: 'Jan 27, 2026',
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
      },
    };
    const actual = formatMtoShipmentForAPI(params);
    expect(actual.shipmentType).toBe(SHIPMENT_OPTIONS.NTS);
    expect(actual.requestedDeliveryDate).toBe('2026-01-27');
    expect(actual.agents.length).toBe(1);
    expect(actual.agents[0].phone).toBe('222-555-1234');
    expect(actual.agents[0].agentType).toBe('RECEIVING_AGENT');
    expect(actual.customerRemarks).toBe('some mock remarks');
  });
});
