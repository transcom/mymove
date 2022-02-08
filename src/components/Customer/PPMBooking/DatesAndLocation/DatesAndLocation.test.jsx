import React from 'react';
import { render, screen } from '@testing-library/react';

import DatesAndLocation from './DatesAndLocation';

const defaultProps = {
  onSubmit: jest.fn(),
  onBack: jest.fn(),
  serviceMember: {
    id: '123',
    residentialAddress: {
      postalCode: '90210',
    },
  },
  destinationDutyStation: {
    address: {
      postalCode: '94611',
    },
  },
  postalCodeValidator: jest.fn(),
};

describe('DatesAndLocation component', () => {
  it('renders blank form on load', async () => {
    render(<DatesAndLocation {...defaultProps} />);
    expect(await screen.getByRole('heading', { level: 2, name: 'Origin' })).toBeInTheDocument();
    expect(screen.getAllByLabelText('ZIP')[0]).toBeInstanceOf(HTMLInputElement);
    expect(screen.getAllByLabelText('Yes')[0]).toBeInstanceOf(HTMLInputElement);
    expect(screen.getAllByLabelText('No')[0]).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByRole('heading', { level: 2, name: 'Destination' })).toBeInTheDocument();
    expect(screen.getAllByLabelText('ZIP')[1]).toBeInstanceOf(HTMLInputElement);
    expect(screen.getAllByLabelText('Yes')[1]).toBeInstanceOf(HTMLInputElement);
    expect(screen.getAllByLabelText('No')[1]).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByRole('heading', { level: 2, name: 'Storage' })).toBeInTheDocument();
    expect(screen.getAllByLabelText('Yes')[2]).toBeInstanceOf(HTMLInputElement);
    expect(screen.getAllByLabelText('No')[2]).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByRole('heading', { level: 2, name: 'Departure date' })).toBeInTheDocument();
    expect(screen.getByLabelText('When do you plan to start moving your PPM?')).toBeInstanceOf(HTMLInputElement);
  });
});
