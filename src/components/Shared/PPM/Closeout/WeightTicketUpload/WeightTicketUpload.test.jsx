import React, { createRef } from 'react';
import { render, screen } from '@testing-library/react';

import WeightTicketUpload from './WeightTicketUpload';

import { MockProviders } from 'testUtils';

const fullWeightTicketUploadMissingWeightTicket = {
  fieldName: 'document',
  missingWeightTicket: true,
  formikProps: { touched: {} },
  values: { document: [] },
  fileUploadRef: createRef(),
  onCreateUpload: jest.fn(),
  onUploadComplete: jest.fn(),
  onUploadDelete: jest.fn(),
};

const emptyWeightTicketUploadMissingWeightTicket = {
  ...fullWeightTicketUploadMissingWeightTicket,
  fieldName: 'emptyDocument',
};

const emptyWeightTicketUploadGunSafeWeightTicket = {
  ...fullWeightTicketUploadMissingWeightTicket,
  missingWeightTicket: false,
  fieldName: 'gunSafe',
};

describe('WeightTicketUpload', () => {
  it('populates form when the full weight ticket is missing', async () => {
    render(<WeightTicketUpload {...fullWeightTicketUploadMissingWeightTicket} />, { wrapper: MockProviders });
    expect(
      screen.getByText('If you do not upload legible certified weight tickets, your PPM incentive could be affected.'),
    ).toBeInTheDocument();
    expect(
      screen.getByText('Download the official government spreadsheet to calculate constructed weight.'),
    ).toBeInTheDocument();
  });

  it('populates form when empty weight ticket is missing', async () => {
    render(<WeightTicketUpload {...emptyWeightTicketUploadMissingWeightTicket} />, { wrapper: MockProviders });
    expect(
      screen.getByText('If you do not upload legible certified weight tickets, your PPM incentive could be affected.'),
    ).toBeInTheDocument();
    expect(
      screen.getByText(
        'Since you do not have a certified weight ticket, upload the registration or rental agreement for the vehicle used during the PPM',
      ),
    ).toBeInTheDocument();
  });

  it('displays correct label for gun safe upload', async () => {
    render(<WeightTicketUpload {...emptyWeightTicketUploadGunSafeWeightTicket} />, { wrapper: MockProviders });
    expect(screen.getByText("Upload your gun safe's weight tickets")).toBeInTheDocument();
  });
});
