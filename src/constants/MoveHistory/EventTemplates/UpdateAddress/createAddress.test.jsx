import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import createAddress from 'constants/MoveHistory/EventTemplates/UpdateAddress/createAddress';
import ADDRESS_TYPE from 'constants/MoveHistory/Database/AddressTypes';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

const LABEL = {
  backupMailingAddress: 'Backup mailing address',
  destinationAddress: 'Destination address',
  pickupAddress: 'Origin address',
  residentialAddress: 'Current mailing address',
  secondaryDestinationAddress: 'Secondary destination address',
  secondaryPickupAddress: 'Secondary origin address',
};

describe('when given an insert with address table history record', () => {
  const historyRecord = {
    action: a.INSERT,
    eventName: '',
    tableName: t.addresses,
  };
  const changedValues = {
    city: 'Beverly Hills',
    postal_code: '90211',
    street_address_1: '12 Any Street',
    street_address_2: 'P.O. Box 1234',
    state: 'CA',
  };
  const address = '12 Any Street, P.O. Box 1234, Beverly Hills, CA 90211';

  it('correctly matches the insert address event ', () => {
    const result = getTemplate(historyRecord);
    expect(result).toMatchObject(createAddress);
    expect(result.getEventNameDisplay()).toEqual('Updated address');
  });
  describe('when given a specific set of details', () => {
    const result = getTemplate(historyRecord);
    it.each(Object.keys(ADDRESS_TYPE).map((type) => [type, address]))(
      'for label %s it displays the proper details value %s',
      async (type, value) => {
        const newChangedValues = { ...changedValues, [type]: address };
        const context = [{ address_type: type }];
        render(result.getDetails({ ...historyRecord, changedValues: newChangedValues, context }));
        expect(screen.getByText(LABEL[type])).toBeInTheDocument();
        expect(screen.getByText(value, { exact: false })).toBeInTheDocument();
      },
    );
  });
});
