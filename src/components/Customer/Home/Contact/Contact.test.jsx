/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen } from '@testing-library/react';

import Contact from './index';

describe('Contact component', () => {
  it('renders footer with given required props', () => {
    const header = 'Contact Info';
    const dutyLocationName = 'Headquarters';
    const officeType = 'Homebase';
    const telephone = '(777) 777-7777';
    const props = {
      header,
      dutyLocationName,
      officeType,
      telephone,
    };

    render(<Contact {...props} />);
    expect(screen.getByRole('heading', { level: 6, name: /Contact Info/i })).toBeInTheDocument();
    expect(screen.getByText(dutyLocationName)).toBeInTheDocument();
    expect(screen.getByText(officeType)).toBeInTheDocument();
    expect(screen.getByText(telephone)).toBeInTheDocument();
    expect(screen.getByRole('link')).toHaveAttribute(
      'href',
      'https://www.militaryonesource.mil/moving-housing/moving/planning-your-move/customer-service-contacts-for-military-pcs/',
    );
  });
});
