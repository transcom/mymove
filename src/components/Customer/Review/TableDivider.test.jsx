import React from 'react';
import { render, screen } from '@testing-library/react';

import { TableDivider } from 'components/Customer/Review/TableDivider';

describe('TableDivider', () => {
  it('verify table divider display', async () => {
    render(<TableDivider className="Test" />);

    expect(await screen.findByTestId('tableDivider')).toBeInTheDocument();
  });
});
