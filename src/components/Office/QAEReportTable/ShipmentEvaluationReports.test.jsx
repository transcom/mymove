import React from 'react';
import { render, screen } from '@testing-library/react';

import ShipmentEvaluationReports from './ShipmentEvaluationReports';

const customerInfo = {
  agency: 'ARMY',
  backup_contact: { email: 'email@example.com', name: 'name', phone: '555-555-5555' },
  current_address: {
    city: 'Beverly Hills',
    country: 'US',
    eTag: 'MjAyMi0wOC0xNVQxNjoxMToyNi4zMzIwOTFa',
    id: '28f11990-7ced-4d01-87ad-b16f2c86ea83',
    postalCode: '90210',
    state: 'CA',
    streetAddress1: '123 Any Street',
    streetAddress2: 'P.O. Box 12345',
    streetAddress3: 'c/o Some Person',
  },
  dodID: '5052247544',
  eTag: 'MjAyMi0wOC0xNVQxNjoxMToyNi4zNTkzNFo=',
  email: 'leo_spaceman_sm@example.com',
  first_name: 'Leo',
  id: 'ea557b1f-2660-4d6b-89a0-fb1b5efd2113',
  last_name: 'Spacemen',
  phone: '555-555-5555',
  userID: 'f4bbfcdf-ef66-4ce7-92f8-4c1bf507d596',
};

describe('ShipmentEvaluationReports', () => {
  it('renders with no shipments', () => {
    render(
      <ShipmentEvaluationReports
        shipments={[]}
        reports={[]}
        moveCode="Test123"
        customerInfo={customerInfo}
        grade="E_4"
      />,
    );
    expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Shipment QAE reports (0)');
  });
});
