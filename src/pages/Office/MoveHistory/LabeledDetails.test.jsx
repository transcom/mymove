import React from 'react';
import { render, screen } from '@testing-library/react';

import LabeledDetails from './LabeledDetails';

describe('LabeledDetails', () => {
  describe('for each changed value', () => {
    const changedValues = {
      customer_remarks: 'Test customer remarks',
      counselor_remarks: 'Test counselor remarks',
      billable_weight_cap: '400',
      tac_type: 'HHG',
      sac_type: 'NTS',
      service_order_number: '1234',
      authorized_weight: '500',
      storage_in_transit: '5',
      dependents_authorized: 'true',
      pro_gear_weight: '100',
      pro_gear_weight_spouse: '50',
      required_medical_equipment_weight: '300',
      organizational_clothing_and_individual_equipment: 'false',
    };
    it.each([
      ['Customer remarks', ': Test customer remarks'],
      ['Counselor remarks', ': Test counselor remarks'],
      ['Billable weight cap', ': 400 lbs'],
      ['TAC type', ': HHG'],
      ['SAC type', ': NTS'],
      ['Service order number', ': 1234'],
      ['Authorized weight', ': 500 lbs'],
      ['Storage in transit (SIT)', ': 5 days'],
      ['Dependents', ': true'],
      ['Pro-gear', ': 100 lbs'],
      ['Spouse pro-gear', ': 50 lbs'],
      ['RME', ': 300 lbs'],
      ['OCIE', ': false'],
    ])('it renders %s%s', (displayName, value) => {
      render(<LabeledDetails changedValues={changedValues} />);

      expect(screen.getByText(displayName)).toBeInTheDocument();

      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });

  it('does not render any text for changed values that are blank', async () => {
    const changedValues = {
      billable_weight_cap: '200',
      customer_remarks: 'Test customer remarks',
      counselor_remarks: '',
    };

    render(<LabeledDetails changedValues={changedValues} />);

    expect(screen.getByText('Billable weight cap')).toBeInTheDocument();

    expect(screen.getByText(200, { exact: false })).toBeInTheDocument();

    expect(screen.getByText('Customer remarks')).toBeInTheDocument();

    expect(screen.getByText('Test customer remarks', { exact: false })).toBeInTheDocument();

    expect(await screen.queryByText('Counselor remarks')).not.toBeInTheDocument();
  });
});
