import React from 'react';
import { render, screen } from '@testing-library/react';

import ServiceItemDetails from './ServiceItemDetails';

const details = {
  description: 'some description',
  pickupPostalCode: '90210',
  SITPostalCode: '12345',
  reason: 'some reason',
  itemDimensions: { length: 1000, width: 2500, height: 3000 },
  crateDimensions: { length: 2000, width: 3500, height: 4000 },
  firstCustomerContact: { timeMilitary: '1200Z', firstAvailableDeliveryDate: '2020-09-15' },
  secondCustomerContact: { timeMilitary: '2300Z', firstAvailableDeliveryDate: '2020-09-21' },
  estimatedWeight: 2500,
};

const detailsRejectedServiceItem = { ...details, rejectionReason: 'some rejection reason' };

const nilDetails = {
  estimatedWeight: null,
};

describe('ServiceItemDetails Domestic Origin SIT', () => {
  it.each([['DOFSIT'], ['DOASIT'], ['DOPSIT']])('renders ZIP and reason', (code) => {
    render(<ServiceItemDetails id="1" code={code} details={details} />);

    expect(screen.getByText('ZIP:')).toBeInTheDocument();
    expect(screen.getByText('12345')).toBeInTheDocument();
    expect(screen.getByText('Reason:')).toBeInTheDocument();
    expect(screen.getByText('some reason')).toBeInTheDocument();
  });
});

describe('ServiceItemDetails Domestic Destination SIT', () => {
  it.each([['DDDSIT'], ['DDASIT'], ['DDFSIT']])(
    'renders first and second customer contact and available delivery date',
    (code) => {
      render(<ServiceItemDetails id="1" code={code} details={details} />);

      expect(screen.getByText('First Customer Contact:')).toBeInTheDocument();
      expect(screen.getByText('1200Z')).toBeInTheDocument();
      expect(screen.getByText('15 Sep 2020')).toBeInTheDocument();
      expect(screen.getByText('Second Customer Contact:')).toBeInTheDocument();
      expect(screen.getByText('2300Z')).toBeInTheDocument();
      expect(screen.getByText('21 Sep 2020')).toBeInTheDocument();
    },
  );
});

describe('ServiceItemDetails Crating', () => {
  it('renders description and dimensions', () => {
    render(<ServiceItemDetails id="1" code="DCRT" details={details} />);

    expect(screen.getByText('some description')).toBeInTheDocument();
    expect(screen.getByText('Item size:')).toBeInTheDocument();
    expect(screen.getByText('1"x2.5"x3"')).toBeInTheDocument();
    expect(screen.getByText('Crate size:')).toBeInTheDocument();
    expect(screen.getByText('2"x3.5"x4"')).toBeInTheDocument();
  });
});

describe('ServiceItemDetails Domestic Shuttling', () => {
  it.each([['DOSHUT'], ['DDSHUT']])('renders formatted estimated weight and reason', (code) => {
    render(<ServiceItemDetails id="1" code={code} details={details} />);

    expect(screen.getByText('2,500 lbs')).toBeInTheDocument();
    expect(screen.getByText('estimated weight')).toBeInTheDocument();
    expect(screen.getByText('Reason:')).toBeInTheDocument();
    expect(screen.getByText('some reason')).toBeInTheDocument();
  });

  it.each([['DOSHUT'], ['DDSHUT']])('renders estimated weight nil values with an em dash', (code) => {
    render(<ServiceItemDetails id="1" code={code} details={nilDetails} />);

    expect(screen.getByText('— lbs')).toBeInTheDocument();
    expect(screen.getByText('estimated weight')).toBeInTheDocument();
  });
});

describe('ServiceItemDetails Crating Rejected', () => {
  it('renders the rejection reason field when it is populated with information', () => {
    render(<ServiceItemDetails id="1" code="DCRT" details={detailsRejectedServiceItem} />);

    expect(screen.getByText('Description:')).toBeInTheDocument();
    expect(screen.getByText('some description')).toBeInTheDocument();
    expect(screen.getByText('Item size:')).toBeInTheDocument();
    expect(screen.getByText('1"x2.5"x3"')).toBeInTheDocument();
    expect(screen.getByText('Crate size:')).toBeInTheDocument();
    expect(screen.getByText('2"x3.5"x4"')).toBeInTheDocument();
    expect(screen.getByText('some rejection reason')).toBeInTheDocument();
  });
});
