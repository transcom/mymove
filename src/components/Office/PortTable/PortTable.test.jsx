import React from 'react';
import { render, screen } from '@testing-library/react';

import PortTable from './PortTable';

const poeLocation = {
  poeLocation: {
    portCode: 'PDX',
    portName: 'PORTLAND INTL',
    city: 'PORTLAND',
    state: 'OREGON',
    zip: '97220',
  },
  podLocation: null,
};

const podLocation = {
  poeLocation: null,
  podLocation: {
    portCode: 'SEA',
    portName: 'SEATTLE TACOMA INTL',
    city: 'SEATTLE',
    state: 'WASHINGTON',
    zip: '98158',
  },
};

const nullPortLocation = {
  poeLocation: null,
  podLocation: null,
};

describe('PortTable', () => {
  it('renders POE location if poeLocation is set', async () => {
    render(<PortTable {...poeLocation} />);
    expect(screen.getByText(/PDX - PORTLAND INTL/)).toBeInTheDocument();
    expect(screen.getByText(/Portland, Oregon 97220/)).toBeInTheDocument();
    expect(screen.queryByText(/SEA - SEATTLE TACOMA INTL/)).not.toBeInTheDocument();
    expect(screen.queryByText(/Seattle, Washington 98158/)).not.toBeInTheDocument();
  });

  it('renders POD location if podLocation is set', async () => {
    render(<PortTable {...podLocation} />);
    expect(screen.queryByText(/PDX - PORTLAND INTL/)).not.toBeInTheDocument();
    expect(screen.queryByText(/Portland, Oregon 97220/)).not.toBeInTheDocument();
    expect(screen.getByText(/SEA - SEATTLE TACOMA INTL/)).toBeInTheDocument();
    expect(screen.getByText(/Seattle, Washington 98158/)).toBeInTheDocument();
  });

  it('does not render port information if poeLocation and podLocation are null', async () => {
    render(<PortTable {...nullPortLocation} />);
    expect(screen.queryByText(/PDX - PORTLAND INTL/)).not.toBeInTheDocument();
    expect(screen.queryByText(/Portland, Oregon 97220/)).not.toBeInTheDocument();
    expect(screen.queryByText(/SEA - SEATTLE TACOMA INTL/)).not.toBeInTheDocument();
    expect(screen.queryByText(/Seattle, Washington 98158/)).not.toBeInTheDocument();
  });
});
