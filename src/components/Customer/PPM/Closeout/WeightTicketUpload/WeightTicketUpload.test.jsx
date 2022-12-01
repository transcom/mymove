import React, { createRef } from 'react';
import { render, screen } from '@testing-library/react';

import WeightTicketUpload from './WeightTicketUpload';

import { MockProviders } from 'testUtils';

const weightTicketUploadMissingWeightTicket = {
  fieldName: 'document',
  missingWeightTicket: true,
  formikProps: { touched: {} },
  values: { document: [] },
  fileUploadRef: createRef(),
  onCreateUpload: jest.fn(),
  onUploadComplete: jest.fn(),
  onUploadDelete: jest.fn(),
};
describe('WeightTicketUpload', () => {
  it('populates form when weight ticket is missing', async () => {
    render(<WeightTicketUpload {...weightTicketUploadMissingWeightTicket} />, { wrapper: MockProviders });
    expect(
      screen.getByText('Download the official government spreadsheet to calculate constructed weight.'),
    ).toBeInTheDocument();
  });
});
