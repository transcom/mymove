import React from 'react';
import { render, screen } from '@testing-library/react';

import CustomerInfoList from './CustomerInfoList';

const info = {
  name: 'Smith, Kerry',
  dodId: '9999999999',
  phone: '+1 999-999-9999',
  email: 'ksmith@email.com',
  currentAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  backupContact: {
    name: 'Quinn Ocampo',
    email: 'quinnocampo@myemail.com',
    phone: '123-555-9898',
  },
};

describe('CustomerInfoList', () => {
  it('renders customer info', () => {
    render(<CustomerInfoList customerInfo={info} />);
    Object.keys(info)
      .filter((k) => k !== 'currentAddress' && k !== 'backupContact')
      .forEach((key) => {
        expect(screen.getByText(info[key])).toBeInTheDocument();
      });
  });

  it('renders formatted current address', () => {
    render(<CustomerInfoList customerInfo={info} />);
    expect(screen.getByText('812 S 129th St, San Antonio, TX 78234')).toBeInTheDocument();
  });

  it('renders formatted backup contact name', () => {
    render(<CustomerInfoList customerInfo={info} />);
    expect(screen.getByText('Quinn Ocampo')).toBeInTheDocument();
  });

  it('renders formatted backup contact email', () => {
    render(<CustomerInfoList customerInfo={info} />);
    expect(screen.getByText('quinnocampo@myemail.com')).toBeInTheDocument();
  });

  it('renders formatted backup contact phone', () => {
    render(<CustomerInfoList customerInfo={info} />);
    expect(screen.getByText('+1 123-555-9898')).toBeInTheDocument();
  });
});
