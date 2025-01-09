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
};

describe('PortTable', () => {
  it('renders port of embarkation and debarkation if one is set', async () => {
    render(<PortTable {...poeLocation} />);
    expect(screen.getByText('Port of Embarkation')).toBeInTheDocument();
    expect(screen.getByText('Port of Debarkation')).toBeInTheDocument();
  });
});
