import React from 'react';
import { render, screen } from '@testing-library/react';

import ShipmentRemarks from './ShipmentRemarks';

describe('ShipmentRemarks', () => {
  const remarks = 'Please treat gently';
  const title = 'Customer remarks';
  it('renders remarks', async () => {
    render(<ShipmentRemarks title={title} remarks={remarks} />);
    const renderedRemark = await screen.getByText(remarks);
    expect(renderedRemark).toBeInTheDocument();
  });
  it('renders title', async () => {
    render(<ShipmentRemarks title={title} remarks={remarks} />);
    const renderedTitle = await screen.getByText(title);
    expect(renderedTitle).toBeInTheDocument();
  });
});
