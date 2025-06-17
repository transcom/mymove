/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import reactQuery from '@tanstack/react-query';
import routeData from 'react-router-dom';
import userEvent from '@testing-library/user-event';
import { render, screen } from '@testing-library/react';

import EvaluationReports from './EvaluationReports';

import { useEvaluationReportsQueries } from 'hooks/queries';
import { MockProviders, renderWithProviders } from 'testUtils';
import { qaeCSRRoutes } from 'constants/routes';
import { permissionTypes } from 'constants/permissions';
import { COUNSELING_EVALUATION_REPORTS, SHIPMENT_EVALUATION_REPORTS } from 'constants/queryKeys';

const mockRequestedMoveCode = 'LR4T8V';

jest.mock('hooks/queries', () => ({
  useEvaluationReportsQueries: jest.fn(),
}));

jest.mock('@tanstack/react-query', () => ({
  ...jest.requireActual('@tanstack/react-query'),
  useMutation: () => ({
    mutate: (_, handlers) => {
      handlers.onSuccess();
    },
  }),
  useQueryClient: () => ({
    refetchQueries: jest.fn(),
    invalidateQueries: jest.fn(),
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

      renderWithProviders(<EvaluationReports customerInfo={{}} grade="" />, mockRoutingOptions);

      const h2 = screen.getByRole('heading', { name: 'Loading, please wait...', level: 2 });
      expect(h2).toBeInTheDocument();
    });

    it('renders the Something Went Wrong component when the query errors', async () => {
      useEvaluationReportsQueries.mockReturnValue(errorReturnValue);

      renderWithProviders(<EvaluationReports customerInfo={{}} grade="" />, mockRoutingOptions);

      const errorMessage = screen.getByText(/Something went wrong./);
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

      renderWithProviders(<EvaluationReports customerInfo={{}} grade="" />, mockRoutingOptions);

      const h1 = screen.getByRole('heading', { name: 'Quality Assurance Reports', level: 1 });
      expect(h1).toBeInTheDocument();

      expect(screen.getByRole('heading', { name: 'Shipment QAE reports (0)', level: 2 })).toBeInTheDocument();
      expect(screen.getByRole('heading', { name: 'Counseling QAE reports (0)', level: 2 })).toBeInTheDocument();
    });

    it('create report button is disabled when move is locked', async () => {
      useEvaluationReportsQueries.mockReturnValue({
        isLoading: false,
        isError: false,
        shipmentEvaluationReports: [],
        counselingEvaluationReports: [],
        shipments: [],
      });
      const isMoveLocked = true;

      render(
        <MockProviders permissions={[permissionTypes.createEvaluationReport]}>
          <EvaluationReports customerInfo={{}} grade="" isMoveLocked={isMoveLocked} />,
        </MockProviders>,
      );

      const createReportBtn = screen.getByRole('button', { name: 'Create report' });
      expect(createReportBtn).toBeInTheDocument();
    });
  });

  describe('deleting report path', () => {
    it('refetches queries', async () => {
      const mockedQueryClient = {
        refetchQueries: jest.fn(),
        invalidateQueries: jest.fn(),
      };

      jest.spyOn(reactQuery, 'useQueryClient').mockReturnValue(mockedQueryClient);

      jest.spyOn(reactQuery, 'useMutation').mockImplementation((_, mutationHandlers) => ({
        mutate: (__, handlers) => {
          mutationHandlers.onSuccess();
          handlers.onSuccess();
        },
      }));

      useEvaluationReportsQueries.mockReturnValue({
        isLoading: false,
        isError: false,
        shipmentEvaluationReports: [],
        counselingEvaluationReports: [
          {
            createdAt: '2025-05-12T20:27:55.009Z',
            eTag: 'MjAyNS0wNS0xMlQyMDoyNzo1NS4wMDk4ODFa',
            id: 'ce13d5db-1263-474b-830a-00a0107370be',
            moveID: 'd70641a9-9225-4f5b-887b-dc4304dbaa31',
            moveReferenceID: '6589-0333',
            officeUser: {
              email: 'multi-role-20250512201621-8b4da589efbd@example.com',
              firstName: 'Alice',
              id: 'deb9a8f8-79de-4963-b08f-88be7585c428',
              lastName: 'Bob',
              phone: '333-333-3333',
            },
            type: 'COUNSELING',
            updatedAt: '2025-05-12T20:27:55.009Z',
          },
        ],
        shipments: [],
      });

      renderWithProviders(<EvaluationReports customerInfo={{}} grade="" />, mockRoutingOptions);

      await userEvent.click(screen.getByTestId('deleteReport'));
      await screen.findByText(/You cannot undo this action./);

      await userEvent.click(screen.getByRole('button', { name: /yes, delete/i }));

      expect(mockedQueryClient.refetchQueries).toBeCalledTimes(2);
      expect(mockedQueryClient.refetchQueries).toBeCalledWith([COUNSELING_EVALUATION_REPORTS], mockRequestedMoveCode);
      expect(mockedQueryClient.refetchQueries).toBeCalledWith([SHIPMENT_EVALUATION_REPORTS], mockRequestedMoveCode);
    });
  });

  describe('check the report status text', () => {
    const testCases = [
      { state: { showDeleteSuccess: true }, expectedText: /Your report has been deleted/ },
      { state: { showCanceledSuccess: true }, expectedText: /Your report has been canceled/ },
      { state: { showSaveDraftSuccess: true }, expectedText: /Your draft report has been saved/ },
      { state: { showSubmitSuccess: true }, expectedText: /Your report has been successfully submitted/ },
    ];
    testCases.forEach(({ state, expectedText }) => {
      it(`renders "${expectedText.source}" alert`, async () => {
        jest.spyOn(routeData, 'useLocation').mockReturnValue({
          pathname: `/moves/${mockRequestedMoveCode}/evaluation-reports`,
          state,
        });
        useEvaluationReportsQueries.mockReturnValue({
          isLoading: false,
          isError: false,
          shipmentEvaluationReports: [],
          counselingEvaluationReports: [],
          shipments: [],
        });
        renderWithProviders(<EvaluationReports customerInfo={{}} grade="" />, mockRoutingOptions);
        expect(screen.getByText(expectedText)).toBeInTheDocument();
      });
    });
  });
});
