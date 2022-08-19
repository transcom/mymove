import React from 'react';
import { render, screen } from '@testing-library/react';

import EvaluationReportShipmentInfo from './EvaluationReportShipmentInfo';

import { ORDERS_BRANCH_OPTIONS } from 'constants/orders';

describe('EvaluationReportShipmentInfo', () => {
  it('renders the correct content', async () => {
    const mockCustomerInfo = {
      last_name: 'Smith',
      first_name: 'John',
      phone: '+441234567890',
      email: 'abc@123.com',
      agency: ORDERS_BRANCH_OPTIONS.NAVY,
    };

    const mockOrders = {
      grade: 'E_1',
    };

    const mockReport = {
      officeUser: {
        firstName: 'John',
        lastName: 'Smith',
        phone: '+441234567890',
        email: 'abc@123.com',
      },
    };

    const mockShipment = {
      actualPickupDate: '2020-03-16',
      approvedDate: '2022-08-16T00:00:00.000Z',
      billableWeightCap: 4000,
      billableWeightJustification: 'heavy',
      createdAt: '2022-08-16T00:00:22.316Z',
      customerRemarks: 'Please treat gently',
      destinationAddress: {
        city: 'Fairfield',
        country: 'US',
        eTag: 'MjAyMi0wOC0xNlQwMDowMDoyMi4zMTQ0MDha',
        id: '1cf638df-1c1e-4c03-916a-bd20f7ec58ce',
        postalCode: '94535',
        state: 'CA',
        streetAddress1: '987 Any Avenue',
        streetAddress2: 'P.O. Box 9876',
        streetAddress3: 'c/o Some Person',
      },
      eTag: 'MjAyMi0wOC0xNlQwMDowMDoyMi4zMTY2N1o=',
      id: 'c37ccf04-637c-4afc-9ef6-dee1555e16ef',
      moveTaskOrderID: '35eb1c36-8916-46f4-a72a-32267c9cb070',
      pickupAddress: {
        city: 'Beverly Hills',
        country: 'US',
        eTag: 'MjAyMi0wOC0xNlQwMDowMDoyMi4zMTIzOTZa',
        id: 'c0cf15bb-96ee-443a-982e-0e9ef2b9a80d',
        postalCode: '90210',
        state: 'CA',
        streetAddress1: '123 Any Street',
        streetAddress2: 'P.O. Box 12345',
        streetAddress3: 'c/o Some Person',
      },
      primeActualWeight: 2000,
      primeEstimatedWeight: 1400,
      requestedDeliveryDate: '2020-03-15',
      requestedPickupDate: '2020-03-15',
      scheduledPickupDate: '2020-03-16',
      shipmentType: 'HHG',
      status: 'APPROVED',
      updatedAt: '2022-08-16T00:00:22.316Z',
    };

    render(
      <EvaluationReportShipmentInfo
        customerInfo={mockCustomerInfo}
        orders={mockOrders}
        report={mockReport}
        shipment={mockShipment}
      />,
    );

    // Section headings
    expect(screen.getByRole('heading', { name: 'Shipment information', level: 2 })).toBeInTheDocument();

    // Shipment info is included, just check it is on page, details handled by EvaluationReportShipmentDisplay component
    expect(screen.getByText('HHG')).toBeInTheDocument();

    expect(screen.getByTestId('qaeAndCustomerInfo')).toBeInTheDocument();

    // Customer info
    expect(screen.getByText('Customer information')).toBeInTheDocument();

    // QAE info
    expect(screen.getByText('QAE')).toBeInTheDocument();
  });
});
