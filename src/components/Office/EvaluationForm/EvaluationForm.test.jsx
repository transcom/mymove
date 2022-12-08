import React from 'react';
import { cleanup, render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import EvaluationForm from './EvaluationForm';

import { MockProviders } from 'testUtils';

const mockMoveCode = 'LR4T8V';
const mockReportId = '58350bae-8e87-4e83-bd75-74027fb4333a';

const mockSaveEvaluationReport = jest.fn();
jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  saveEvaluationReport: (options) => mockSaveEvaluationReport(options),
}));

afterEach(() => {
  jest.resetAllMocks();
});

const evaluationReportCounseling = {
  type: 'COUNSELING',
};

const evaluationReportShipment = {
  type: 'SHIPMENT',
};

const mockEvaluationReport = {
  createdAt: '2022-09-07T15:17:37.484Z',
  eTag: 'MjAyMi0wOS0wN1QxODowNjozNy44NjQxNDJa',
  id: '6739d7fc-6067-4e84-996d-f4f70b8ec6fd',
  inspectionDate: '2022-09-08',
  inspectionType: 'DATA_REVIEW',
  location: 'ORIGIN',
  moveID: '551dd01f-90cf-44d6-addb-ff919433dd61',
  moveReferenceID: '4118-8295',
  officeUser: {
    email: 'qae_csr_role@office.mil',
    firstName: 'Leo',
    id: 'ef4f6d1f-4ac3-4159-a364-5403e7d958ff',
    lastName: 'Spaceman',
    phone: '415-555-1212',
  },
  remarks: 'test',
  shipmentID: '319e0751-1337-4ed9-b4b5-a15d4e6d272c',
  type: 'SHIPMENT',
  updatedAt: '2022-09-07T18:06:37.864Z',
  violationsObserved: false,
  evalStart: '10:30',
  evalEnd: '22:00',
};

const customerInfo = {
  agency: 'ARMY',
  backup_contact: { email: 'email@example.com', name: 'name', phone: '555-555-5555' },
  current_address: {
    city: 'Beverly Hills',
    country: 'US',
    eTag: 'MjAyMi0wOC0xNVQxNjoxMToyNi4zMzIwOTFa',
    id: '28f11990-7ced-4d01-87ad-b16f2c86ea83',
    postalCode: '90210',
    state: 'CA',
    streetAddress1: '123 Any Street',
    streetAddress2: 'P.O. Box 12345',
    streetAddress3: 'c/o Some Person',
  },
  dodID: '5052247544',
  eTag: 'MjAyMi0wOC0xNVQxNjoxMToyNi4zNTkzNFo=',
  email: 'leo_spaceman_sm@example.com',
  first_name: 'Leo',
  id: 'ea557b1f-2660-4d6b-89a0-fb1b5efd2113',
  last_name: 'Spacemen',
  phone: '555-555-5555',
  userID: 'f4bbfcdf-ef66-4ce7-92f8-4c1bf507d596',
};

const grade = 'E_4';
const moveCode = 'FAKEIT';
const destinationDutyLocationPostalCode = '90210';

const mockPush = jest.fn();
jest.mock('react-router', () => ({
  ...jest.requireActual('react-router'),
  useHistory: () => ({
    push: mockPush,
  }),
  useParams: () => ({ moveCode: 'LR4T8V', reportId: '58350bae-8e87-4e83-bd75-74027fb4333a' }),
}));

afterEach(() => {
  jest.resetAllMocks();
  cleanup();
});

const renderForm = (props) => {
  const defaultProps = {
    evaluationReport: mockEvaluationReport,
    customerInfo,
    grade,
    moveCode,
    destinationDutyLocationPostalCode,
  };

  return render(
    <MockProviders initialEntries={[`/moves/${mockMoveCode}/evaluation-reports/${mockReportId}`]}>
      <EvaluationForm {...defaultProps} {...props} />
    </MockProviders>,
  );
};

describe('EvaluationForm', () => {
  it('renders the form components', async () => {
    renderForm();

    // Headers
    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Evaluation form');
      const h3s = screen.getAllByRole('heading', { level: 3 });
      expect(h3s).toHaveLength(3);

      expect(screen.getByText('Evaluation information')).toBeInTheDocument();
      expect(screen.getByText('Violations')).toBeInTheDocument();
      expect(screen.getByText('QAE remarks')).toBeInTheDocument();

      // // Form components
      expect(screen.getByTestId('evaluationReportForm')).toBeInTheDocument();

      expect(screen.getByText('Date of inspection')).toBeInTheDocument();
      expect(screen.getByText('Evaluation type')).toBeInTheDocument();
      expect(screen.getByText('Evaluation location')).toBeInTheDocument();
      expect(screen.getByText('Violations observed')).toBeInTheDocument();
      expect(screen.getByText('Evaluation remarks')).toBeInTheDocument();
      expect(screen.getByText('Time evaluation started')).toBeInTheDocument();
      expect(screen.getByText('Time evaluation ended')).toBeInTheDocument();

      // Conditionally shown fields should not be displayed initially
      expect(screen.queryByText('Time departed for evaluation')).not.toBeInTheDocument();
      expect(screen.queryByText('Observed pickup date')).not.toBeInTheDocument();
      expect(screen.queryByText('Observed delivery date')).not.toBeInTheDocument();

      // Form buttons
      expect(screen.getByText('Cancel')).toBeInTheDocument();
      expect(screen.getByText('Save draft')).toBeInTheDocument();
      expect(screen.getByText('Review and submit')).toBeInTheDocument();
    });
  });

  it('renders conditionally displayed shipment form components correctly', async () => {
    renderForm({ evaluationReport: evaluationReportShipment });

    // Initially no conditional fields shown
    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Evaluation form');
      expect(screen.getAllByTestId('textarea')).toHaveLength(1);

      expect(screen.queryByText('Time departed for evaluation')).not.toBeInTheDocument();
      expect(screen.queryByText('Observed pickup date')).not.toBeInTheDocument();
      expect(screen.queryByText('Observed delivery date')).not.toBeInTheDocument();
    });

    // Select Physical Evaluation type and origin location, ensure time depart, evaluation start and end are in the doc
    await waitFor(() => {
      userEvent.click(screen.getByText('Physical'));
      userEvent.click(screen.getByText('Origin'));
      expect(screen.getByText('Time departed for evaluation')).toBeInTheDocument();
      expect(screen.getByText('Time evaluation started')).toBeInTheDocument();
      expect(screen.getByText('Time evaluation ended')).toBeInTheDocument();
      expect(screen.queryByText('Observed delivery date')).not.toBeInTheDocument();
      expect(screen.queryByText('Observed pickup date')).toBeInTheDocument();
      expect(screen.getAllByTestId('textarea')).toHaveLength(1);
    });

    await waitFor(() => {
      userEvent.click(screen.getByText('Destination'));
      expect(screen.getByText('Observed delivery date')).toBeInTheDocument();
      expect(screen.queryByText('Observed pickup date')).not.toBeInTheDocument();
      expect(screen.getAllByTestId('textarea')).toHaveLength(1);
    });

    await waitFor(() => {
      userEvent.click(screen.getByText('Other'));
      expect(screen.queryByText('Observed delivery date')).not.toBeInTheDocument();
      expect(screen.queryByText('Observed pickup date')).not.toBeInTheDocument();
      expect(screen.getAllByTestId('textarea')).toHaveLength(2);
    });

    // If not 'Physical' eval type, no conditional time fields should be shown
    await waitFor(() => {
      userEvent.click(screen.getByText('Virtual'));
      expect(screen.queryByText('Time departed for evaluation')).not.toBeInTheDocument();
      expect(screen.queryByText('Observed delivery date')).not.toBeInTheDocument();
      expect(screen.queryByText('Observed pickup date')).not.toBeInTheDocument();
    });
  });

  it('displays the delete confirmation on cancel', async () => {
    renderForm({ evaluationReport: evaluationReportCounseling });

    expect(await screen.getByRole('heading', { level: 2 })).toHaveTextContent('Evaluation form');

    // Buttons
    await userEvent.click(await screen.getByRole('button', { name: 'Cancel' }));

    expect(
      await screen.findByRole('heading', { level: 3, name: 'Are you sure you want to cancel this report?' }),
    ).toBeInTheDocument();
  });

  it('updates the submit button when there are violations', async () => {
    renderForm();

    expect(await screen.findByRole('button', { name: 'Review and submit' })).toBeInTheDocument();
    expect(
      screen.queryByText('You will select the specific PWS paragraphs violated on the next screen.'),
    ).not.toBeInTheDocument();

    await waitFor(() => {
      userEvent.click(screen.getByTestId('yesViolationsRadioOption'));

      expect(screen.getByRole('button', { name: 'Next: select violations' })).toBeInTheDocument();
      expect(
        screen.getByText('You will select the specific PWS paragraphs violated on the next screen.'),
      ).toBeInTheDocument();
    });

    await waitFor(() => {
      userEvent.click(screen.getByRole('button', { name: 'Next: select violations' }));
    });
    expect(mockSaveEvaluationReport).toHaveBeenCalledTimes(1);
    expect(mockSaveEvaluationReport).toHaveBeenCalledWith({
      body: {
        inspectionDate: mockEvaluationReport.inspectionDate,
        inspectionType: mockEvaluationReport.inspectionType,
        location: mockEvaluationReport.location,
        locationDescription: undefined,
        observedShipmentDeliveryDate: undefined,
        observedShipmentPhysicalPickupDate: undefined,
        remarks: mockEvaluationReport.remarks,
        violationsObserved: true,
        evalStart: mockEvaluationReport.evalStart,
        evalEnd: mockEvaluationReport.evalEnd,
        timeDepart: undefined,
      },
      ifMatchETag: mockEvaluationReport.eTag,
      reportID: mockReportId,
    });
  });

  it('can save a draft and reroute back to the eval reports', async () => {
    renderForm();

    // Click save draft button
    await userEvent.click(await screen.findByRole('button', { name: 'Save draft' }));

    // Verify that report was saved and page rerouted
    await waitFor(() => {
      expect(mockPush).toHaveBeenCalledTimes(1);
      expect(mockSaveEvaluationReport).toHaveBeenCalledTimes(1);
      expect(mockSaveEvaluationReport).toHaveBeenCalledWith({
        body: {
          inspectionDate: '2022-09-08',
          inspectionType: 'DATA_REVIEW',
          location: 'ORIGIN',
          locationDescription: undefined,
          observedShipmentDeliveryDate: undefined,
          observedShipmentPhysicalPickupDate: undefined,
          remarks: 'test',
          violationsObserved: false,
          evalStart: mockEvaluationReport.evalStart,
          evalEnd: mockEvaluationReport.evalEnd,
          timeDepart: undefined,
        },
        ifMatchETag: 'MjAyMi0wOS0wN1QxODowNjozNy44NjQxNDJa',
        reportID: mockReportId,
      });
      expect(mockPush).toHaveBeenCalledWith(`/moves/${mockMoveCode}/evaluation-reports`, {
        showSaveDraftSuccess: true,
      });
    });
  });
});
