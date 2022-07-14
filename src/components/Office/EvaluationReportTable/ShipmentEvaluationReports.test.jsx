import React from 'react';
import { render, screen } from '@testing-library/react';

import ShipmentEvaluationReports from './ShipmentEvaluationReports';

describe('ShipmentEvaluationReports', () => {
  it('renders with no shipments', () => {
    render(<ShipmentEvaluationReports shipments={[]} reports={[]} />);
    expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Shipment QAE reports (0)');
  });
});
