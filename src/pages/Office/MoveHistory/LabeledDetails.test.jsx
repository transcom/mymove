import React from 'react';
import { render, screen } from '@testing-library/react';

import LabeledDetails from './LabeledDetails';

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
        dependents_authorized: 'true',
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
        sac: '3333',
        nts_tac: '4444',
        nts_sac: '5555',
        department_indicator: 'AIR_FORCE',
        grade: 'E_1',
        shipment_type: '(Shipment Type) HHG',
      },
    };
    it.each([
      ['Customer remarks', ': Test customer remarks'],
      ['Counselor remarks', ': Test counselor remarks'],
      ['Billable weight remarks', ': Test TIO remarks'],
      ['Billable weight', ': 400 lbs'],
      ['TAC type', ': HHG'],
      ['SAC type', ': NTS'],
      ['Service order number', ': 1234'],
      ['Authorized weight', ': 500 lbs'],
      ['Storage in transit (SIT)', ': 5 days'],
      ['Dependents', ': true'],
      ['Pro-gear', ': 100 lbs'],
      ['Spouse pro-gear', ': 50 lbs'],
      ['RME', ': 300 lbs'],
      ['Orders type', ': Permanent Change Of Station (PCS)'],
      ['Orders type detail', ': Shipment of HHG Permitted'],
      ['Origin duty location name', ': Origin duty location name'],
      ['New duty location name', ': New duty location name'],
      ['Orders number', ': 1111'],
      ['HHG TAC', ': 2222'],
      ['NTS TAC', ': 3333'],
      ['HHG SAC', ': 4444'],
      ['NTS SAC', ': 5555'],
      ['Dept. indicator', ': Air Force'],
      ['Shipment type', ': (Shipment Type) HHG'],
    ])('it renders %s%s', (displayName, value) => {
      render(<LabeledDetails historyRecord={historyRecord} />);

      expect(screen.getByText(displayName)).toBeInTheDocument();

      expect(screen.getByText(value)).toBeInTheDocument();
    });
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
