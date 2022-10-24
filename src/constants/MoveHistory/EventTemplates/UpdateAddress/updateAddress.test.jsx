import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import updateAddress from 'constants/MoveHistory/EventTemplates/UpdateAddress/updateAddress';
import ADDRESS_TYPE from 'constants/MoveHistory/Database/AddressTypes';
import { shipmentTypes } from 'constants/shipments';
import { formatMoveHistoryFullAddress } from 'utils/formatters';

const LABEL = {
  backupMailingAddress: 'Backup mailing address',
  destinationAddress: 'Destination address',
  pickupAddress: 'Origin address',
  residentialAddress: 'Current mailing address',
  secondaryDestinationAddress: 'Secondary destination address',
  secondaryPickupAddress: 'Secondary origin address',
};

describe('when given a Update basic service item address history record', () => {
  const historyRecord = {
    action: a.UPDATE,
    eventName: '',
    tableName: t.addresses,
    oldValues: {
      city: 'Beverly Hills',
      postal_code: '90211',
      street_address_1: '12 Any Street',
      street_address_2: 'P.O. Box 1234',
      state: 'CA',
    },
  };
  const changedValues = {
    city: 'San Diego',
    postal_code: '92134',
    street_address_1: '123 Test Street',
    street_address_2: 'Apt 345',
    state: 'GA',
  };

  const template = getTemplate(historyRecord);
  it('correctly matches the update address event', () => {
    expect(template).toMatchObject(updateAddress);
    expect(template.getEventNameDisplay()).toEqual('Updated address');
  });
  describe('when given a specific address type', () => {
    const result = getTemplate(historyRecord);
    const newAddress = formatMoveHistoryFullAddress(changedValues);
    // test each address type available in ADDRESS_TYPES
    it.each(Object.keys(ADDRESS_TYPE).map((type) => [type, newAddress]))(
      'for label %s it displays the proper details value %s',
      async (type, value) => {
        // add the address type and new valuei in changeValues
        const newChangedValues = { ...changedValues, [type]: newAddress };
        // set the address type in context
        const context = [{ address_type: type }];
        // add shipment labels where needed
        if (!LABEL[type].includes('mailing')) {
          newChangedValues.shipment_type = 'HHG';
        }
        render(result.getDetails({ ...historyRecord, changedValues: newChangedValues, context }));
        expect(screen.getByText(LABEL[type])).toBeInTheDocument();
        expect(screen.getByText(value, { exact: false })).toBeInTheDocument();
        if (changedValues.shipment_type) {
          expect(screen.getByText('HHG Shipment', { exact: false })).toBeInTheDocument();
        }
      },
    );
  });
  describe('when given one change to address, the updated address is correct', () => {
    const result = getTemplate(historyRecord);
    // test changes in city postal code, state and street individually based off changedValues above
    it.each(Object.keys(changedValues).map((type) => [type, changedValues[type]]))(
      'for label %s it displays the proper details value %s',
      async (label, value) => {
        const newChangedValues = { [label]: changedValues[label] };
        const context = [{ address_type: 'pickupAddress' }];
        const newHistoryRecord = { ...historyRecord, changedValues: newChangedValues, context };
        const expectedAddress = formatMoveHistoryFullAddress(newChangedValues);
        render(result.getDetails(newHistoryRecord));
        expect(screen.getByText(value, { exact: false })).toBeInTheDocument();
        expect(screen.getByText(expectedAddress, { exact: false })).toBeInTheDocument();
      },
    );
  });
  describe('when given a shipment type, the correct label renders', () => {
    const result = getTemplate(historyRecord);
    // test changes in city postal code, state and street individually based off changedValues above
    it.each(Object.keys(shipmentTypes).map((type) => [type, shipmentTypes[type]]))(
      'for label %s it displays the proper details value %s',
      async (label, value) => {
        const newChangedValues = { ...changedValues, shipment_type: label };
        const context = [{ address_type: 'pickupAddress' }];
        const newHistoryRecord = { ...historyRecord, changedValues: newChangedValues, context };
        render(result.getDetails(newHistoryRecord));
        expect(screen.getByText(`${value} shipment`, { exact: false })).toBeInTheDocument();
      },
    );
  });
});
