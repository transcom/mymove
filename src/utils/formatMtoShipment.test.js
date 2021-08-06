import { formatMtoShipmentForAPI, formatMtoShipmentForDisplay } from './formatMtoShipment';

import { MTOAgentType, SHIPMENT_OPTIONS } from 'shared/constants';

describe('formatMtoShipmentForDisplay', () => {
  const emptyAgentShape = {
    firstName: '',
    lastName: '',
    email: '',
    phone: '',
  };

  const emptyAddressShape = {
    street_address_1: '',
    street_address_2: '',
    city: '',
    state: '',
    postal_code: '',
  };

  const mtoShipment = {
    requestedPickupDate: '2026-01-20',
    pickupAddress: {
      street_address_1: '123 main',
      city: 'legit human city',
      state: 'DC',
      postal_code: '20017',
    },
    requestedDeliveryDate: '2026-01-27',
    moveTaskOrderID: 'move123',
  };

  const customerRemarks = 'some mock remarks';
  const counselorRemarks = 'all looks good';

  const releasingAgent = {
    firstName: 'mockFirstName',
    lastName: 'mockLastName',
    email: 'mockAgentEmail@example.com',
    phone: '222-555-1234',
  };

  const receivingAgent = {
    firstName: 'r0b0tBestFr1end',
    lastName: 'r0b0tBestFr1endLastName',
    email: 'r0b0t-fr1end@example.com',
    phone: '222-555-0101',
  };

  const destinationAddress = {
    street_address_1: '0011100010110101',
    city: 'R0B0T T0WN',
    state: 'CP',
    postal_code: '10101',
  };

  const secondaryPickupAddress = {
    street_address_1: '142 E Barrel Hoop Circle',
    street_address_2: '#4A',
    city: 'Corpus Christi',
    state: 'TX',
    postal_code: '78412',
  };

  const secondaryDeliveryAddress = {
    street_address_1: '441 SW Río de la Plata Drive',
    street_address_2: '',
    city: destinationAddress.city,
    state: destinationAddress.state,
    postal_code: destinationAddress.postal_code,
  };

  const checkAddressesAreEqual = (address1, address2) => {
    expect(address1.street_address_1 === address2.street_address_1);
    expect(address1.street_address_2 === address2.street_address_2);
    expect(address1.street_address_3 === address2.street_address_3);
    expect(address1.city === address2.city);
    expect(address1.state === address2.state);
    expect(address1.postal_code === address2.postal_code);
  };

  const checkAgentsAreEqual = (agent1, agent2) => {
    expect(agent1.firstName === agent2.firstName);
    expect(agent1.lastName === agent2.lastName);
    expect(agent1.email === agent2.email);
    expect(agent1.phone === agent2.phone);
    expect(agent1.agentType === agent2.agentType);
  };

  it.each([[SHIPMENT_OPTIONS.HHG], [SHIPMENT_OPTIONS.NTSR], [SHIPMENT_OPTIONS.NTS]])(
    'can format a shipment (type: %s)',
    (shipmentType) => {
      const params = {
        ...mtoShipment,
        shipmentType,
      };

      const displayValues = formatMtoShipmentForDisplay(params);

      expect(displayValues.shipmentType).toBe(shipmentType);
      expect(displayValues.moveTaskOrderID).toBe(mtoShipment.moveTaskOrderID);
      expect(displayValues.customerRemarks).toBe('');
      expect(displayValues.counselorRemarks).toBe('');

      expect(displayValues.pickup.requestedDate.toDateString()).toBe('Tue Jan 20 2026');

      const expectedPickupAddress = { ...emptyAddressShape, ...mtoShipment.pickupAddress };
      checkAddressesAreEqual(displayValues.pickup.address, expectedPickupAddress);
      checkAgentsAreEqual(displayValues.pickup.agent, emptyAgentShape);

      expect(displayValues.delivery.requestedDate.toDateString()).toBe('Tue Jan 27 2026');
      checkAddressesAreEqual(displayValues.delivery.address, emptyAddressShape);
      checkAgentsAreEqual(displayValues.delivery.agent, emptyAgentShape);
      expect(displayValues.hasDeliveryAddress).toBe('no');

      checkAddressesAreEqual(displayValues.secondaryPickup.address, emptyAddressShape);
      expect(displayValues.hasSecondaryPickup).toBe('no');

      checkAddressesAreEqual(displayValues.secondaryDelivery.address, emptyAddressShape);
      expect(displayValues.hasSecondaryDelivery).toBe('no');

      expect(displayValues.agents).toBeUndefined();
    },
  );

  it('can format a shipment with remarks', () => {
    const params = {
      ...mtoShipment,
      shipmentType: SHIPMENT_OPTIONS.HHG,
      customerRemarks,
      counselorRemarks,
    };

    const displayValues = formatMtoShipmentForDisplay(params);

    expect(displayValues.customerRemarks).toBe(customerRemarks);
    expect(displayValues.counselorRemarks).toBe(counselorRemarks);
  });

  it('can format a shipment with agents', () => {
    const params = {
      ...mtoShipment,
      shipmentType: SHIPMENT_OPTIONS.HHG,
      agents: [
        { ...releasingAgent, agentType: MTOAgentType.RELEASING },
        { ...receivingAgent, agentType: MTOAgentType.RECEIVING },
      ],
    };

    const displayValues = formatMtoShipmentForDisplay(params);

    checkAgentsAreEqual(displayValues.pickup.agent, releasingAgent);
    checkAgentsAreEqual(displayValues.delivery.agent, receivingAgent);
  });

  it('can format a shipment with a destination, secondary pickup, and secondary destination', () => {
    const params = {
      ...mtoShipment,
      shipmentType: SHIPMENT_OPTIONS.HHG,
      destinationAddress,
      secondaryPickupAddress,
      secondaryDeliveryAddress,
    };

    const displayValues = formatMtoShipmentForDisplay(params);

    const expectedDeliveryAddress = { ...emptyAddressShape, ...destinationAddress };
    checkAddressesAreEqual(displayValues.delivery.address, expectedDeliveryAddress);
    expect(displayValues.hasDeliveryAddress).toBe('yes');

    const expectedSecondaryPickupAddress = { ...emptyAddressShape, ...secondaryPickupAddress };
    checkAddressesAreEqual(displayValues.secondaryPickup.address, expectedSecondaryPickupAddress);
    expect(displayValues.hasSecondaryPickup).toBe('yes');

    const expectedSecondaryDeliveryAddress = { ...emptyAddressShape, ...secondaryDeliveryAddress };
    checkAddressesAreEqual(displayValues.secondaryDelivery.address, expectedSecondaryDeliveryAddress);
    expect(displayValues.hasSecondaryDelivery).toBe('yes');
  });
});

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
      phone: '222-555-1234',
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
      phone: '222-555-0101',
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
