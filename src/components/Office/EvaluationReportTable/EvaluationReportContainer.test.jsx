import { render, screen } from '@testing-library/react';
import React from 'react';

import EvaluationReportContainer from './EvaluationReportContainer';

describe('Evaluation Report Container', () => {
  it('renders the sample text', async () => {
    render(<EvaluationReportContainer />);

    const evaluationReportContainer = await screen.findByTestId('EvaluationReportContainer');

    expect(evaluationReportContainer).toBeInTheDocument();
  });
});
