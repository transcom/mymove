import React from 'react';
import { render, screen } from '@testing-library/react';

import LabeledPaymentRequestDetails from './LabeledPaymentRequestDetails';

describe('LabeledPaymentRequestDetails', () => {
  it('renders the labeled payment request details', async () => {
    const labeledPaymentRequestDetails = {
      moveServices: 'Move management',
      shipmentServices: [
        { serviceItems: 'Test service', shipmentId: '123', shipmentType: 'HHG', shipmentIdAbbr: 'ACF7B' },
        {
          serviceItems: 'Domestic uncrating',
          shipmentId: '456',
          shipmentType: 'HHG_INTO_NTS',
          shipmentIdAbbr: 'A1C2B',
        },
      ],
    };

    render(<LabeledPaymentRequestDetails services={labeledPaymentRequestDetails} />);

    expect(screen.getByText('Move services')).toBeInTheDocument();
    expect(screen.getByText(': Move management')).toBeInTheDocument();
    expect(screen.getByText('HHG shipment #ACF7B')).toBeInTheDocument();
    expect(screen.getByText(': Test service')).toBeInTheDocument();
    expect(screen.getByText('NTS shipment #A1C2B')).toBeInTheDocument();
    expect(screen.getByText(': Domestic uncrating')).toBeInTheDocument();
  });
});
