import React from 'react';
import { render, screen } from '@testing-library/react';

import CustomerInfoList from './CustomerInfoList';

const info = {
  name: 'Smith, Kerry',
  agency: 'COAST_GUARD',
  edipi: '9999999999',
  emplid: '7777777',
  phone: '999-999-9999',
  altPhone: '888-888-8888',
  email: 'ksmith@email.com',
  currentAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  backupAddress: {
    streetAddress1: '812½ S 129th St',
    streetAddress2: 'Apt B',
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
      .filter((k) => k !== 'currentAddress' && k !== 'backupAddress' && k !== 'backupContact' && k !== 'agency')
      .forEach((key) => {
        if (key === 'phone' || key === 'altPhone') {
          screen.getByText(`+1 ${info[key]}`);
        } else {
          expect(screen.getByText(info[key])).toBeInTheDocument();
        }
      });
  });

  it('renders formatted pickup address', () => {
    render(<CustomerInfoList customerInfo={info} />);
    expect(screen.getByText('812 S 129th St, San Antonio, TX 78234')).toBeInTheDocument();
  });

  it('renders formatted backup address', () => {
    render(<CustomerInfoList customerInfo={info} />);
    expect(screen.getByText('812½ S 129th St, Apt B, San Antonio, TX 78234')).toBeInTheDocument();
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
