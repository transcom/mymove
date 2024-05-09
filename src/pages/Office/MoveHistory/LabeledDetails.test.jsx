import React from 'react';
import { render, screen } from '@testing-library/react';

import LabeledDetails, { retrieveTextToDisplay } from './LabeledDetails';

import { SHIPMENT_OPTIONS } from 'shared/constants';

describe('LabeledDetails', () => {
  describe('for each changed value', () => {
    const changedValues = {
      customer_remarks: 'Test customer remarks',
      counselor_remarks: 'Test counselor remarks',
      billable_weight_justification: 'Test TIO remarks',
      billable_weight_cap: '400',
      tac_type: 'HHG',
      sac_type: 'NTS',
      service_order_number: '1234',
      authorized_weight: '500',
      storage_in_transit: '5',
      dependents_authorized: true,
      pro_gear_weight: '100',
      pro_gear_weight_spouse: '50',
      required_medical_equipment_weight: '300',
      organizational_clothing_and_individual_equipment: 'false',
      orders_type: 'PERMANENT_CHANGE_OF_STATION',
      orders_type_detail: 'HHG_PERMITTED',
      origin_duty_location_name: 'Origin duty location name',
      new_duty_location_name: 'New duty location name',
      orders_number: '1111',
      tac: '2222',
      sac: '4444',
      nts_tac: '3333',
      nts_sac: '5555',
      department_indicator: 'AIR_AND_SPACE_FORCE',
      grade: 'E_1',
      actual_pickup_date: '2022-01-01',
      prime_actual_weight: '100 lbs',
      destination_address_type: 'HOME_OF_SELECTION',
      affiliation: 'COAST_GUARD',
      requested_delivery_date: '2023-02-05',
      sit_entry_date: '2023-04-05',
      sit_departure_date: '2023-03-05',
      sit_expected: true,
      sit_location: 'Destination',
      sit_estimated_weight: '500',
      sit_estimated_cost: '120000',
      estimated_incentive: '850',
      shipment_weight: '100',
    };

    const testCases = Object.entries(changedValues).map(([fieldName, value]) => {
      const { displayName, displayValue } = retrieveTextToDisplay(fieldName, value);
      return [fieldName, value, displayName, displayValue];
    });

    it.each(testCases)("it renders [%s %s] as '%s: %s'", (fieldName, value, displayName, displayValue) => {
      const { baseElement } = render(<LabeledDetails historyRecord={{ changedValues: { [fieldName]: value } }} />);
      expect(baseElement).toHaveTextContent(`${displayName}: ${displayValue}`, { normalizeWhitespace: true });
    });
  });

  it('renders shipment_type as a header', async () => {
    const historyRecord = {
      changedValues: {
        billable_weight_cap: '200',
        billable_weight_justification: 'Test TIO Remarks',
        shipment_type: SHIPMENT_OPTIONS.NTSR,
        shipment_locator: 'RQ38D4-01',
        shipment_id_display: 'X9Y0Z',
      },
    };

    render(<LabeledDetails historyRecord={historyRecord} />);

    expect(screen.getByText('NTS-release shipment #RQ38D4-01')).toBeInTheDocument();
  });

  it('does not render any text for changed values that are blank', async () => {
    const historyRecord = {
      changedValues: {
        billable_weight_cap: '200',
        customer_remarks: 'Test customer remarks',
        counselor_remarks: '',
      },
    };

    render(<LabeledDetails historyRecord={historyRecord} />);

    expect(screen.getByText('Billable weight')).toBeInTheDocument();

    expect(screen.getByText(200, { exact: false })).toBeInTheDocument();

    expect(screen.getByText('Customer remarks')).toBeInTheDocument();

    expect(screen.getByText('Test customer remarks', { exact: false })).toBeInTheDocument();

    expect(await screen.queryByText('Counselor remarks')).not.toBeInTheDocument();
  });
});

it('does render text for changed values that are blank when they exist in the old values (deleted values)', async () => {
  const historyRecord = {
    changedValues: {
      billable_weight_cap: '200',
      counselor_remarks: '',
    },
    oldValues: {
      counselor_remarks: 'These remarks were deleted',
    },
  };

  const { baseElement } = render(<LabeledDetails historyRecord={historyRecord} />);

  expect(baseElement).toHaveTextContent('Counselor remarks');
  expect(baseElement).toHaveTextContent('—');
});
