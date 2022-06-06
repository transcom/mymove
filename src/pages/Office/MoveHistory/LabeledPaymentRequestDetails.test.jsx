import React from 'react';
import { render, screen } from '@testing-library/react';

import LabeledPaymentRequestDetails from './LabeledPaymentRequestDetails';

describe('LabeledPaymentRequestDetails', () => {
  it('renders the labeled payment request details', async () => {
    const context = [{ name: 'Test service', price: '', shipment_id: '123' }];
    const labeledPaymentRequestDetails = {
      moveServices: 'Move management',
      shipmentServices: [
        { serviceItems: 'Test service', shipmentId: '123', shipmentType: 'HHG' },
        { serviceItems: 'Domestic uncrating', shipmentId: '456', shipmentType: 'HHG_INTO_NTS_DOMESTIC' },
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
    expect(screen.getByText('HHG shipment')).toBeInTheDocument();
    expect(screen.getByText(': Test service')).toBeInTheDocument();
    expect(screen.getByText('NTS shipment')).toBeInTheDocument();
    expect(screen.getByText(': Domestic uncrating')).toBeInTheDocument();
  });
});
