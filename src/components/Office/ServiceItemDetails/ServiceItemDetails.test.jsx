import React from 'react';
import { render, screen } from '@testing-library/react';

import ServiceItemDetails from './ServiceItemDetails';

const details = {
  description: 'some description',
  pickupPostalCode: '90210',
  SITPostalCode: '12345',
  sitEntryDate: '2020-10-10',
  reason: 'some reason',
  itemDimensions: { length: 1000, width: 2500, height: 3000 },
  crateDimensions: { length: 2000, width: 3500, height: 4000 },
  customerContacts: [
    { timeMilitary: '1200Z', firstAvailableDeliveryDate: '2020-09-15', dateOfContact: '2020-09-15' },
    { timeMilitary: '2300Z', firstAvailableDeliveryDate: '2020-09-21', dateOfContact: '2020-09-21' },
  ],
  estimatedWeight: 2500,
};

const serviceRequestDocs = [
  {
    uploads: [
      {
        filename: '/mto-service-item/ae1c6472-5e03-4f9c-bef5-55605dbeb31e/20230630161854-receipt.pdf',
      },
    ],
  },
];

const detailsRejectedServiceItem = { ...details, rejectionReason: 'some rejection reason' };

const nilDetails = {
  estimatedWeight: null,
};

describe('ServiceItemDetails Domestic Origin SIT', () => {
  it.each([['DOASIT'], ['DOPSIT']])('renders ZIP, reason, and service request documents', (code) => {
    render(<ServiceItemDetails id="1" code={code} details={details} serviceRequestDocs={serviceRequestDocs} />);

    expect(screen.getByText('ZIP:')).toBeInTheDocument();
    expect(screen.getByText('12345')).toBeInTheDocument();
    expect(screen.getByText('Reason:')).toBeInTheDocument();
    expect(screen.getByText('some reason')).toBeInTheDocument();
    expect(screen.getByText('Download service item documentation:')).toBeInTheDocument();
    const downloadLink = screen.getByText('receipt.pdf');
    expect(downloadLink).toBeInstanceOf(HTMLAnchorElement);
  });
});

describe('ServiceItemDetails for DOFSIT', () => {
  it('renders SIT entry date, ZIP, and reason', () => {
    render(<ServiceItemDetails id="1" code="DOFSIT" details={details} serviceRequestDocs={serviceRequestDocs} />);

    expect(screen.getByText('SIT entry date:')).toBeInTheDocument();
    expect(screen.getByText('10 Oct 2020')).toBeInTheDocument();
    expect(screen.getByText('ZIP:')).toBeInTheDocument();
    expect(screen.getByText('12345')).toBeInTheDocument();
    expect(screen.getByText('Reason:')).toBeInTheDocument();
    expect(screen.getByText('some reason')).toBeInTheDocument();
    expect(screen.getByText('Download service item documentation:')).toBeInTheDocument();
    const downloadLink = screen.getByText('receipt.pdf');
    expect(downloadLink).toBeInstanceOf(HTMLAnchorElement);
  });
});

describe('ServiceItemDetails Domestic Destination SIT', () => {
  it.each([['DDDSIT'], ['DDASIT'], ['DDFSIT']])(
    'renders first and second customer contact and available delivery date',
    (code) => {
      render(<ServiceItemDetails id="1" code={code} details={details} serviceRequestDocs={serviceRequestDocs} />);

      expect(screen.getByText('Customer contact attempt 1:')).toBeInTheDocument();
      expect(screen.getByText('15 Sep 2020, 1200Z')).toBeInTheDocument();
      expect(screen.getByText('15 Sep 2020')).toBeInTheDocument();
      expect(screen.getByText('Customer contact attempt 2:')).toBeInTheDocument();
      expect(screen.getByText('21 Sep 2020, 2300Z')).toBeInTheDocument();
      expect(screen.getByText('21 Sep 2020')).toBeInTheDocument();
      expect(screen.getByText('Download service item documentation:')).toBeInTheDocument();
      const downloadLink = screen.getByText('receipt.pdf');
      expect(downloadLink).toBeInstanceOf(HTMLAnchorElement);
    },
  );
});

describe('ServiceItemDetails Crating', () => {
  it('renders description and dimensions', () => {
    render(<ServiceItemDetails id="1" code="DCRT" details={details} serviceRequestDocs={serviceRequestDocs} />);

    expect(screen.getByText('some description')).toBeInTheDocument();
    expect(screen.getByText('Item size:')).toBeInTheDocument();
    expect(screen.getByText('1"x2.5"x3"')).toBeInTheDocument();
    expect(screen.getByText('Crate size:')).toBeInTheDocument();
    expect(screen.getByText('2"x3.5"x4"')).toBeInTheDocument();
    expect(screen.getByText('Download service item documentation:')).toBeInTheDocument();
    const downloadLink = screen.getByText('receipt.pdf');
    expect(downloadLink).toBeInstanceOf(HTMLAnchorElement);
  });
});

describe('ServiceItemDetails Domestic Shuttling', () => {
  it.each([['DOSHUT'], ['DDSHUT']])('renders formatted estimated weight and reason', (code) => {
    render(<ServiceItemDetails id="1" code={code} details={details} serviceRequestDocs={serviceRequestDocs} />);

    expect(screen.getByText('2,500 lbs')).toBeInTheDocument();
    expect(screen.getByText('estimated weight')).toBeInTheDocument();
    expect(screen.getByText('Reason:')).toBeInTheDocument();
    expect(screen.getByText('some reason')).toBeInTheDocument();
    expect(screen.getByText('Download service item documentation:')).toBeInTheDocument();
    const downloadLink = screen.getByText('receipt.pdf');
    expect(downloadLink).toBeInstanceOf(HTMLAnchorElement);
  });

  it.each([['DOSHUT'], ['DDSHUT']])('renders estimated weight nil values with an em dash', (code) => {
    render(<ServiceItemDetails id="1" code={code} details={nilDetails} />);

    expect(screen.getByText('— lbs')).toBeInTheDocument();
    expect(screen.getByText('estimated weight')).toBeInTheDocument();
  });
});

describe('ServiceItemDetails Crating Rejected', () => {
  it('renders the rejection reason field when it is populated with information', () => {
    render(
      <ServiceItemDetails
        id="1"
        code="DUCRT"
        details={detailsRejectedServiceItem}
        serviceRequestDocs={serviceRequestDocs}
      />,
    );

    expect(screen.getByText('Description:')).toBeInTheDocument();
    expect(screen.getByText('some description')).toBeInTheDocument();
    expect(screen.getByText('Item size:')).toBeInTheDocument();
    expect(screen.getByText('1"x2.5"x3"')).toBeInTheDocument();
    expect(screen.getByText('Crate size:')).toBeInTheDocument();
    expect(screen.getByText('2"x3.5"x4"')).toBeInTheDocument();
    expect(screen.getByText('some rejection reason')).toBeInTheDocument();
    expect(screen.getByText('Download service item documentation:')).toBeInTheDocument();
    const downloadLink = screen.getByText('receipt.pdf');
    expect(downloadLink).toBeInstanceOf(HTMLAnchorElement);
  });
});

describe('ServiceItemDetails Crating Rejected', () => {
  it('renders the rejection reason field when it is populated with information', () => {
    render(
      <ServiceItemDetails
        id="1"
        code="DCRT"
        details={detailsRejectedServiceItem}
        serviceRequestDocs={serviceRequestDocs}
      />,
    );

    expect(screen.getByText('Description:')).toBeInTheDocument();
    expect(screen.getByText('some description')).toBeInTheDocument();
    expect(screen.getByText('Item size:')).toBeInTheDocument();
    expect(screen.getByText('1"x2.5"x3"')).toBeInTheDocument();
    expect(screen.getByText('Crate size:')).toBeInTheDocument();
    expect(screen.getByText('2"x3.5"x4"')).toBeInTheDocument();
    expect(screen.getByText('some rejection reason')).toBeInTheDocument();
    expect(screen.getByText('Download service item documentation:')).toBeInTheDocument();
    const downloadLink = screen.getByText('receipt.pdf');
    expect(downloadLink).toBeInstanceOf(HTMLAnchorElement);
  });
});
