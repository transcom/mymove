import React from 'react';
import { render, screen } from '@testing-library/react';

import LabeledPaymentRequestDetails from './LabeledPaymentRequestDetails';

describe('LabeledPaymentRequestDetails', () => {
  it('renders the labeled payment request details', async () => {
    const context = [
      { name: 'Test service', price: '', shipment_id: '123', shipment_id_abbr: 'acf7b' },
      { name: 'Domestic uncrating', price: '', shipment_id: '456', shipment_id_abbr: 'a1c2b' },
    ];
    const labeledPaymentRequestDetails = {
      moveServices: 'Move management',
      shipmentServices: [
        { serviceItems: 'Test service', shipmentId: '123', shipmentType: 'HHG', shipmentIdAbbr: 'ACF7B' },
        {
          serviceItems: 'Domestic uncrating',
          shipmentId: '456',
          shipmentType: 'HHG_INTO_NTS_DOMESTIC',
          shipmentIdAbbr: 'A1C2B',
        },
      ],
    };

    render(
      <LabeledPaymentRequestDetails
        context={context}
        getLabeledPaymentRequestDetails={() => labeledPaymentRequestDetails}
      />,
    );

    expect(screen.getByText('Move services')).toBeInTheDocument();
    expect(screen.getByText(': Move management')).toBeInTheDocument();
    expect(screen.getByText('HHG shipment #ACF7B')).toBeInTheDocument();
    expect(screen.getByText(': Test service')).toBeInTheDocument();
    expect(screen.getByText('NTS shipment #A1C2B')).toBeInTheDocument();
    expect(screen.getByText(': Domestic uncrating')).toBeInTheDocument();
  });
});
