import moveHistoryOperations from 'constants/moveHistoryOperations';

const {
  default: getMoveHistoryEventTemplate,
  approveShipmentEvent,
  approveShipmentDiversionEvent,
  createOrdersEvent,
  createMTOShipmentAddressesEvent,
  createMTOShipmentAgentEvent,
  requestShipmentCancellationEvent,
  setFinancialReviewFlagEvent,
  updateMoveTaskOrderStatusEvent,
  updateMTOShipmentEvent,
  updateMTOShipmentAddressesEvent,
  updateMTOShipmentAgentEvent,
  updateServiceItemStatusEvent,
  acknowledgeExcessWeightRiskEvent,
  createStandardServiceItemEvent,
  createBasicServiceItemEvent,
  updateOrderEvent,
  requestShipmentReweighEvent,
  createPaymentRequestReweighUpdate,
  createPaymentRequestShipmentUpdate,
  undefinedEvent,
} = require('./moveHistoryEventTemplate');

const { detailsTypes } = require('constants/moveHistoryEventTemplate');

describe('moveHistoryEventTemplate', () => {
  describe('when given an Acknowledge excess weight risk history record', () => {
    const item = {
      action: 'UPDATE',
      eventName: 'acknowledgeExcessWeightRisk',
      tableName: 'moves',
    };
    it('correctly matches the Acknowledge excess weight risk event', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(acknowledgeExcessWeightRiskEvent);
      expect(result.getDetailsPlainText(item)).toEqual('Dismissed excess weight alert');
    });
  });

  describe('when given an Approved shipment history record', () => {
    const item = {
      action: 'UPDATE',
      changedValues: { status: 'APPROVED' },
      eventName: 'approveShipment',
      oldValues: { shipment_type: 'HHG' },
      tableName: 'mto_shipments',
    };
    it('correctly matches the Approved shipment event', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(approveShipmentEvent);
      expect(result.getDetailsPlainText(item)).toEqual('HHG shipment');
    });
  });

  describe('when given an Approved shipment diversion history record', () => {
    const item = {
      changedValues: { status: 'APPROVED' },
      eventName: 'approveShipmentDiversion',
      oldValues: { shipment_type: 'HHG' },
    };
    it('correctly matches the Approved shipment event', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(approveShipmentDiversionEvent);
      expect(result.getDetailsPlainText(item)).toEqual('HHG shipment');
    });
  });

  describe('when given a Submitted orders history record', () => {
    const item = {
      action: 'INSERT',
      changedValues: {
        status: 'DRAFT',
      },
      eventName: 'createOrders',
      tableName: 'orders',
    };
    it('correctly matches the Submitted orders event', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(createOrdersEvent);
    });
  });

  describe('when given a Create basic service item history record', () => {
    const item = {
      action: 'INSERT',
      context: [
        {
          name: 'Move management',
        },
      ],
      eventName: 'updateMoveTaskOrderStatus',
      tableName: 'mto_service_items',
    };
    it('correctly matches the Create basic service item event', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(createBasicServiceItemEvent);
      expect(result.getEventNameDisplay(result)).toEqual('Approved service item');
      expect(result.getDetailsPlainText(item)).toEqual('Move management');
    });
  });

  describe('when given a Create standard service item history record', () => {
    const item = {
      action: 'INSERT',
      context: [
        {
          shipment_type: 'HHG',
          name: 'Domestic linehaul',
        },
      ],
      eventName: 'approveShipment',
      tableName: 'mto_service_items',
    };
    it('correctly matches the Create standard service item event', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(createStandardServiceItemEvent);
      expect(result.getEventNameDisplay(result)).toEqual('Approved service item');
      expect(result.getDetailsPlainText(item)).toEqual('HHG shipment, Domestic linehaul');
    });
  });

  describe('when given a Request shipment cancellation history record', () => {
    const item = {
      action: 'UPDATE',
      changedValues: {
        status: 'DRAFT',
      },
      eventName: 'requestShipmentCancellation',
      oldValues: { shipment_type: 'PPM' },
      tableName: '',
    };
    it('correctly matches the Request shipment cancellation event', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(requestShipmentCancellationEvent);
      expect(result.getDetailsPlainText(item)).toEqual('Requested cancellation for PPM shipment');
    });
  });

  describe('when given a Request shipment reweigh history record', () => {
    const item = {
      action: 'INSERT',
      context: [{ shipment_type: 'HHG' }],
      eventName: 'requestShipmentReweigh',
      tableName: 'reweighs',
    };
    it('correctly matches the Request shipment reweigh event', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(requestShipmentReweighEvent);
      expect(result.getDetailsPlainText(item)).toEqual('HHG shipment, reweigh requested');
    });
  });

  describe('when given a Set financial review flag event for flagged move history record', () => {
    const item = {
      action: 'UPDATE',
      eventName: 'setFinancialReviewFlag',
      changedValues: { financial_review_flag: 'true' },
      tableName: 'moves',
    };
    it('correctly matches the Set financial review flag event', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(setFinancialReviewFlagEvent);
      expect(result.getDetailsPlainText(item)).toEqual('Move flagged for financial review');
    });
  });

  describe('when given a Set financial review flag event for unflagged move history record', () => {
    const item = {
      action: 'UPDATE',
      eventName: 'setFinancialReviewFlag',
      changedValues: { financial_review_flag: 'false' },
      tableName: 'moves',
    };
    it('correctly matches the Set financial review flag event', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(setFinancialReviewFlagEvent);
      expect(result.getDetailsPlainText(item)).toEqual('Move unflagged for financial review');
    });
  });

  describe('when given an Approved service item history record', () => {
    const item = {
      action: 'UPDATE',
      changedValues: { status: 'APPROVED' },
      context: [{ name: 'Domestic origin price', shipment_type: 'HHG_INTO_NTS_DOMESTIC' }],
      eventName: 'updateMTOServiceItemStatus',
      tableName: 'mto_service_items',
    };
    it('correctly matches the Approved service item event', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(updateServiceItemStatusEvent);
      expect(result.getDetailsPlainText(item)).toEqual('NTS shipment, Domestic origin price');
    });
  });

  describe('when given a Move approved history record', () => {
    const item = {
      action: 'UPDATE',
      changedValues: {
        available_to_prime_at: '2022-04-13T15:21:31.746028+00:00',
        status: 'APPROVED',
      },
      eventName: 'updateMoveTaskOrderStatus',
      tableName: 'moves',
    };
    it('correctly matches the Update move task order status event', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(updateMoveTaskOrderStatusEvent);
      expect(result.getEventNameDisplay(item)).toEqual('Approved move');
      expect(result.getDetailsPlainText(item)).toEqual('Created Move Task Order (MTO)');
    });
  });

  describe('when given a Move status update history record', () => {
    const item = {
      action: 'UPDATE',
      changedValues: { status: 'CANCELED' },
      eventName: 'updateMoveTaskOrderStatus',
      tableName: 'moves',
    };
    it('correctly matches the Update move task order status event', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(updateMoveTaskOrderStatusEvent);
      expect(result.getEventNameDisplay(item)).toEqual('Move status updated');
      expect(result.getDetailsPlainText(item)).toEqual('-');
    });
  });

  describe('when given an Order update history record', () => {
    const item = {
      action: 'UPDATE',
      eventName: 'updateOrder',
      tableName: 'orders',
      detailsType: detailsTypes.LABELED,
      changedValues: { old_duty_location_id: 'ID1', new_duty_location_id: 'ID2' },
      context: [{ old_duty_location_name: 'old name', new_duty_location_name: 'new name' }],
    };
    it('correctly matches the Update orders event', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(updateOrderEvent);
      // expect to have merged context and changedValues
      expect(result.getDetailsLabeledDetails({ context: item.context, changedValues: item.changedValues })).toEqual({
        old_duty_location_id: 'ID1',
        new_duty_location_id: 'ID2',
        old_duty_location_name: 'old name',
        new_duty_location_name: 'new name',
      });
    });
  });

  describe('when given an mto shipment update with mto shipment table history record', () => {
    const item = {
      action: 'UPDATE',
      eventName: moveHistoryOperations.updateMTOShipment,
      tableName: 'mto_shipments',
      detailsType: detailsTypes.LABELED,
      changedValues: {
        destination_address_type: 'HOME_OF_SELECTION',
        requested_delivery_date: '2020-04-14',
        requested_pickup_date: '2020-03-23',
      },
    };
    it('correctly matches the Update mto shipment event', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(updateMTOShipmentEvent);
    });
  });

  describe('when given an mto shipment update with address table history record', () => {
    const item = {
      action: 'UPDATE',
      eventName: moveHistoryOperations.updateMTOShipment,
      tableName: 'addresses',
      detailsType: detailsTypes.LABELED,
      changedValues: {
        city: 'Beverly Hills',
        postal_code: '90211',
        street_address_1: '12 Any Street',
        street_address_2: 'P.O. Box 1234',
      },
      oldValues: {
        city: 'Beverly Hills',
        postal_code: '90211',
        state: 'CA',
        street_address_1: '12 Any Street',
        street_address_2: 'P.O. Box 1234',
      },
      context: [{ addressType: 'pickupAddress' }],
    };

    it('correctly matches the Update mto shipment address event for pickup addresses', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(updateMTOShipmentAddressesEvent);
      // expect to have formatted the adddresses correctly
      expect(
        result.getDetailsLabeledDetails({
          changedValues: item.changedValues,
          oldValues: item.oldValues,
          context: item.context,
        }),
      ).toEqual({
        pickup_address: '12 Any Street, P.O. Box 1234, Beverly Hills, CA 90211',
        city: 'Beverly Hills',
        postal_code: '90211',
        street_address_1: '12 Any Street',
        street_address_2: 'P.O. Box 1234',
      });
    });

    it('correctly matches the Update mto shipment address event for destination addresses', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(updateMTOShipmentAddressesEvent);
      // expect to have formatted the adddresses correctly
      expect(
        result.getDetailsLabeledDetails({
          changedValues: item.changedValues,
          oldValues: item.oldValues,
          context: [{ addressType: 'destinationAddress' }],
        }),
      ).toEqual({
        destination_address: '12 Any Street, P.O. Box 1234, Beverly Hills, CA 90211',
        city: 'Beverly Hills',
        postal_code: '90211',
        street_address_1: '12 Any Street',
        street_address_2: 'P.O. Box 1234',
      });
    });
  });

  describe('when given an mto shipment insert with address table history record', () => {
    const item = {
      action: 'INSERT',
      eventName: '',
      tableName: 'addresses',
      detailsType: detailsTypes.LABELED,
      changedValues: {
        city: 'Beverly Hills',
        postal_code: '90211',
        street_address_1: '12 Any Street',
        street_address_2: 'P.O. Box 1234',
        state: 'CA',
      },
      context: [{ addressType: 'pickupAddress' }],
    };

    it('correctly matches the insert mto shipment address event for pickup addresses', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(createMTOShipmentAddressesEvent);
      // expect to have formatted the adddresses correctly
      expect(
        result.getDetailsLabeledDetails({
          changedValues: item.changedValues,
          context: item.context,
        }),
      ).toEqual({
        pickup_address: '12 Any Street, P.O. Box 1234, Beverly Hills, CA 90211',
        city: 'Beverly Hills',
        postal_code: '90211',
        street_address_1: '12 Any Street',
        street_address_2: 'P.O. Box 1234',
        state: 'CA',
      });
    });

    it('correctly matches the insert mto shipment address event for destination addresses', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(createMTOShipmentAddressesEvent);
      // expect to have formatted the adddresses correctly
      expect(
        result.getDetailsLabeledDetails({
          changedValues: item.changedValues,
          context: [{ addressType: 'destinationAddress' }],
        }),
      ).toEqual({
        destination_address: '12 Any Street, P.O. Box 1234, Beverly Hills, CA 90211',
        city: 'Beverly Hills',
        postal_code: '90211',
        street_address_1: '12 Any Street',
        street_address_2: 'P.O. Box 1234',
        state: 'CA',
      });
    });
  });

  describe('when given an mto shipment agents update with mto agents table history record', () => {
    const item = {
      action: 'UPDATE',
      eventName: moveHistoryOperations.updateMTOShipment,
      tableName: 'mto_agents',
      detailsType: detailsTypes.LABELED,
      changedValues: {
        email: 'grace@email.com',
        first_name: 'Grace',
        phone: '555-555-5555',
      },
      oldValues: {
        agent_type: 'RELEASING_AGENT',
        email: 'gracie@email.com',
        first_name: 'Gracie',
        last_name: 'Griffin',
        phone: '555-555-5551',
      },
    };

    it('correctly matches the Update mto shipment agent event for releasing agents', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(updateMTOShipmentAgentEvent);
      // expect to have formatted the agent correctly
      expect(
        result.getDetailsLabeledDetails({
          changedValues: item.changedValues,
          oldValues: item.oldValues,
        }),
      ).toEqual({
        releasing_agent: 'Grace Griffin, 555-555-5555, grace@email.com',
        email: 'grace@email.com',
        first_name: 'Grace',
        phone: '555-555-5555',
      });
    });

    it('correctly matches the Update mto shipment agent event for receiving agents', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(updateMTOShipmentAgentEvent);
      // expect to have formatted the agent correctly
      expect(
        result.getDetailsLabeledDetails({
          changedValues: item.changedValues,
          oldValues: { ...item.oldValues, agent_type: 'RECEIVING_AGENT' },
        }),
      ).toEqual({
        receiving_agent: 'Grace Griffin, 555-555-5555, grace@email.com',
        email: 'grace@email.com',
        first_name: 'Grace',
        phone: '555-555-5555',
      });
    });
  });

  describe('when given an mto shipment agents insert with mto agents table history record', () => {
    const item = {
      action: 'INSERT',
      eventName: moveHistoryOperations.updateMTOShipment,
      tableName: 'mto_agents',
      detailsType: detailsTypes.LABELED,
      changedValues: {
        email: 'grace@email.com',
        first_name: 'Grace',
        last_name: 'Griffin',
        phone: '555-555-5555',
        agent_type: 'RELEASING_AGENT',
      },
    };

    it('correctly matches the insert mto shipment agent event for releasing agents', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(createMTOShipmentAgentEvent);
      // expect to have formatted the agent correctly
      expect(
        result.getDetailsLabeledDetails({
          changedValues: item.changedValues,
        }),
      ).toEqual({
        releasing_agent: 'Grace Griffin, 555-555-5555, grace@email.com',
        email: 'grace@email.com',
        first_name: 'Grace',
        last_name: 'Griffin',
        phone: '555-555-5555',
        agent_type: 'RELEASING_AGENT',
      });
    });

    it('correctly matches the insert mto shipment agent event for receiving agents', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(createMTOShipmentAgentEvent);
      // expect to have formatted the agent correctly
      expect(
        result.getDetailsLabeledDetails({
          changedValues: { ...item.changedValues, agent_type: 'RECEIVING_AGENT' },
        }),
      ).toEqual({
        receiving_agent: 'Grace Griffin, 555-555-5555, grace@email.com',
        email: 'grace@email.com',
        first_name: 'Grace',
        last_name: 'Griffin',
        phone: '555-555-5555',
        agent_type: 'RECEIVING_AGENT',
      });
    });
  });

  describe('when given a payment request is created through reweigh', () => {
    const item = {
      action: 'INSERT',
      eventName: 'updateReweigh',
      tableName: 'payment_requests',
    };
    it('correctly matches the Request shipment reweigh event', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(createPaymentRequestReweighUpdate);
      expect(result.getStatusDetails(item)).toEqual('Pending');
    });
  });

  describe('when given a payment request is created through shipment update', () => {
    const item = {
      action: 'INSERT',
      eventName: 'updateMTOShipment',
      tableName: 'payment_requests',
    };
    it('correctly matches the Request shipment reweigh event', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(createPaymentRequestShipmentUpdate);
      expect(result.getStatusDetails(item)).toEqual('Pending');
    });
  });
});

describe('when given an unidentifiable move history record', () => {
  const item = {
    action: 'UPDATE',
    eventName: 'testEventName',
    tableName: 'imaginary_test_objects',
  };
  it('correctly matches the Undefined move history event', () => {
    const result = getMoveHistoryEventTemplate(item);
    expect(result).toEqual(undefinedEvent);
    expect(result.getEventNameDisplay(item)).toEqual('Updated item in imaginary test objects');
  });
});
