/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { screen } from '@testing-library/react';

import EvaluationReports from './EvaluationReports';

import { renderWithRouter } from 'testUtils';
import { useEvaluationReportsQueries } from 'hooks/queries';
import { qaeCSRRoutes } from 'constants/routes';

const mockRequestedMoveCode = 'LR4T8V';

jest.mock('hooks/queries', () => ({
  useEvaluationReportsQueries: jest.fn(),
}));

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useLocation: () => ({
    pathname: `/moves/${mockRequestedMoveCode}/evaluation-reports`,
    state: { showDeleteSuccess: true },
  }),
}));
const mockRoutingOptions = {
  path: qaeCSRRoutes.BASE_EVALUATION_REPORTS_PATH,
  params: { moveCode: mockRequestedMoveCode },
};

const loadingReturnValue = {
  isLoading: true,
  isError: false,
  isSuccess: false,
};

const errorReturnValue = {
  isLoading: false,
  isError: true,
  isSuccess: false,
};

describe('EvaluationReports', () => {
  describe('check loading and error component states', () => {
    it('renders the Loading Placeholder when the query is still loading', async () => {
      useEvaluationReportsQueries.mockReturnValue(loadingReturnValue);

      renderWithRouter(<EvaluationReports customerInfo={{}} grade="" />, mockRoutingOptions);

      const h2 = await screen.getByRole('heading', { name: 'Loading, please wait...', level: 2 });
      expect(h2).toBeInTheDocument();
    });

    it('renders the Something Went Wrong component when the query errors', async () => {
      useEvaluationReportsQueries.mockReturnValue(errorReturnValue);

      renderWithRouter(<EvaluationReports customerInfo={{}} grade="" />, mockRoutingOptions);

      const errorMessage = await screen.getByText(/Something went wrong./);
      expect(errorMessage).toBeInTheDocument();
    });
  });

  describe('check the component renders the correct content', () => {
    it('renders the component', async () => {
      useEvaluationReportsQueries.mockReturnValue({
        isLoading: false,
        isError: false,
        shipmentEvaluationReports: [],
        counselingEvaluationReports: [],
        shipments: [],
      });

      renderWithRouter(<EvaluationReports customerInfo={{}} grade="" />, mockRoutingOptions);

      const h1 = await screen.getByRole('heading', { name: 'Quality assurance reports', level: 1 });
      expect(h1).toBeInTheDocument();

      expect(await screen.getByRole('heading', { name: 'Counseling QAE reports (0)', level: 2 })).toBeInTheDocument();
      expect(await screen.getByRole('heading', { name: 'Shipment QAE reports (0)', level: 2 })).toBeInTheDocument();
    });
  });

  describe('check the delete report confirmation', () => {
    it('renders the "report has been deleted" alert', async () => {
      useEvaluationReportsQueries.mockReturnValue({
        isLoading: false,
        isError: false,
        shipmentEvaluationReports: [],
        counselingEvaluationReports: [],
        shipments: [],
        showDeleteSuccess: true,
      });

      renderWithRouter(<EvaluationReports customerInfo={{}} grade="" />, mockRoutingOptions);

      const alert = await screen.getByText(/Your report has been deleted/);
      expect(alert).toBeInTheDocument();
    });
  });
});
