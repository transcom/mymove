import React from 'react';
import { render, screen } from '@testing-library/react';

import LabeledDetails from './LabeledDetails';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import {
  formatCustomerDate,
  formatEvaluationReportLocation,
  formatWeight,
  formatYesNoMoveHistoryValue,
  toDollarString,
} from 'utils/formatters';
import * as fieldDefault from 'constants/MoveHistory/Database/FieldMappings';

const FieldMappings = fieldDefault.default;

describe('LabeledDetails', () => {
  describe('for each changed value', () => {
    const historyRecord = {
      changedValues: {
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
      },
    };

    const testCases = [
      [FieldMappings.counselor_remarks, ': Test counselor remarks'],
      [FieldMappings.customer_remarks, ': Test customer remarks'],
      [FieldMappings.billable_weight_justification, ': Test TIO remarks'],
      [FieldMappings.billable_weight_cap, ': 400 lbs'],
      [FieldMappings.tac_type, ': HHG'],
      [FieldMappings.sac_type, ': NTS'],
      [FieldMappings.service_order_number, ': 1234'],
      [FieldMappings.authorized_weight, ': 500 lbs'],
      [FieldMappings.storage_in_transit, ': 5 days'],
      [FieldMappings.dependents_authorized, ': Yes'],
      [FieldMappings.pro_gear_weight, ': 100 lbs'],
      [FieldMappings.pro_gear_weight_spouse, ': 50 lbs'],
      [FieldMappings.required_medical_equipment_weight, ': 300 lbs'],
      [FieldMappings.orders_type, ': Permanent Change Of Station (PCS)'],
      [FieldMappings.orders_type_detail, ': Shipment of HHG Permitted'],
      [FieldMappings.origin_duty_location_name, ': Origin duty location name'],
      [FieldMappings.new_duty_location_name, ': New duty location name'],
      [FieldMappings.orders_number, ': 1111'],
      [FieldMappings.tac, ': 2222'],
      [FieldMappings.nts_tac, ': 3333'],
      [FieldMappings.sac, ': 4444'],
      [FieldMappings.nts_sac, ': 5555'],
      [FieldMappings.grade, ': E-1'],
      [FieldMappings.department_indicator, ': Air Force and Space Force'],
      [FieldMappings.actual_pickup_date, `: ${formatCustomerDate(historyRecord.changedValues.actual_pickup_date)}`],
      [FieldMappings.shipment_weight, ': 100 lbs'],
      [FieldMappings.destination_address_type, ': Home of selection (HOS)'],
      [
        FieldMappings.requested_delivery_date,
        `: ${formatCustomerDate(historyRecord.changedValues.requested_delivery_date)}`,
      ],
      [FieldMappings.sit_entry_date, `: ${formatCustomerDate(historyRecord.changedValues.sit_entry_date)}`],
      [FieldMappings.sit_departure_date, `: ${formatCustomerDate(historyRecord.changedValues.sit_departure_date)}`],
      [FieldMappings.sit_expected, `: ${formatYesNoMoveHistoryValue(historyRecord.changedValues.sit_expected)}`],
      [
        FieldMappings.sit_location,
        `: ${formatEvaluationReportLocation(historyRecord.changedValues.sit_location.toUpperCase())}`,
      ],
      [FieldMappings.sit_estimated_weight, `: ${formatWeight(historyRecord.changedValues.sit_estimated_weight)}`],
      [FieldMappings.sit_estimated_cost, `: ${toDollarString(historyRecord.changedValues.sit_estimated_cost / 100)}`],
      [FieldMappings.estimated_incentive, `: ${toDollarString(historyRecord.changedValues.estimated_incentive / 100)}`],
    ];

    it.each(testCases)('it renders %s%s', (displayName, value) => {
      render(<LabeledDetails historyRecord={historyRecord} />);
      const displayingElements = screen.getAllByText(displayName);
      const displayingElement = displayingElements.find((element) => element.parentElement.textContent.includes(value));
      const parent = displayingElement.parentElement;
      expect(parent.textContent).toContain(displayName);
      expect(parent.textContent).toContain(value);
    });
  });

  it('renders shipment_type as a header', async () => {
    const historyRecord = {
      changedValues: {
        billable_weight_cap: '200',
        billable_weight_justification: 'Test TIO Remarks',
        shipment_type: SHIPMENT_OPTIONS.NTSR,
        shipment_id_display: 'X9Y0Z',
      },
    };

    render(<LabeledDetails historyRecord={historyRecord} />);

    expect(screen.getByText('NTS-release shipment #X9Y0Z')).toBeInTheDocument();
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

  render(<LabeledDetails historyRecord={historyRecord} />);

  expect(screen.getByText('Counselor remarks')).toBeInTheDocument();
  expect(screen.getByText('—', { exact: false })).toBeInTheDocument();

  expect(await screen.queryByText('These remarks were deleted')).not.toBeInTheDocument();
});
