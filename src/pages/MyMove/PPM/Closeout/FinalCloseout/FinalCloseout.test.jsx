import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { customerRoutes } from 'constants/routes';
import FinalCloseout from 'pages/MyMove/PPM/Closeout/FinalCloseout/FinalCloseout';
import { selectMTOShipmentById, selectServiceMemberAffiliation } from 'store/entities/selectors';
import { updateMTOShipment } from 'store/entities/actions';
import { selectMove } from 'shared/Entities/modules/moves';
import { MockProviders } from 'testUtils';
import { createPPMShipmentWithFinalIncentive } from 'utils/test/factories/ppmShipment';
import { getMTOShipmentsForMove, submitPPMShipmentSignedCertification } from 'services/internalApi';
import { ppmSubmissionCertificationText } from 'scenes/Legalese/legaleseText';
import { formatDateForSwagger } from 'shared/dates';
import affiliations from 'content/serviceMemberAgencies';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

const mockDispatch = jest.fn();
jest.mock('react-redux', () => ({
  ...jest.requireActual('react-redux'),
  useDispatch: jest.fn().mockImplementation(() => mockDispatch),
}));

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  getResponseError: jest.fn(),
  submitPPMShipmentSignedCertification: jest.fn(),
  getMTOShipmentsForMove: jest.fn(),
}));

jest.mock('store/entities/selectors', () => ({
  ...jest.requireActual('store/entities/selectors'),
  selectMTOShipmentById: jest.fn(),
  selectServiceMemberAffiliation: jest.fn(),
}));

jest.mock('store/entities/actions', () => ({
  ...jest.requireActual('store/entities/actions'),
  updateMTOShipment: jest.fn(),
}));

jest.mock('shared/Entities/modules/moves', () => ({
  ...jest.requireActual('store/entities/actions'),
  selectMove: jest.fn(),
}));

beforeEach(() => {
  jest.clearAllMocks();
});

describe('Final Closeout page', () => {
  const setUpMocks = (mtoShipment) => {
    const shipment = mtoShipment || createPPMShipmentWithFinalIncentive();
    const move = {
      closeout_office: {
        name: 'Altus AFB',
      },
    };
    selectMTOShipmentById.mockReturnValue(shipment);
    selectMove.mockReturnValue(move);
    selectServiceMemberAffiliation.mockReturnValue(affiliations.ARMY);
    getMTOShipmentsForMove.mockResolvedValueOnce(shipment);

    return { shipment, move };
  };

  it('loads the selected shipment from redux', async () => {
    const { shipment } = setUpMocks();

    const mockRoutingConfig = {
      path: customerRoutes.SHIPMENT_PPM_COMPLETE_PATH,
      params: { mtoShipmentId: shipment.id, moveId: shipment.moveTaskOrderId },
    };

    render(
      <MockProviders {...mockRoutingConfig}>
        <FinalCloseout />
      </MockProviders>,
    );

    await waitFor(() => {
      expect(selectMTOShipmentById).toHaveBeenCalledWith(expect.anything(), shipment.id);
    });
  });

  it('renders the page headings and closeout office name', async () => {
    const { shipment, move } = setUpMocks();

    const mockRoutingConfig = {
      path: customerRoutes.SHIPMENT_PPM_COMPLETE_PATH,
      params: { mtoShipmentId: shipment.id, moveId: shipment.moveTaskOrderId },
    };

    render(
      <MockProviders {...mockRoutingConfig}>
        <FinalCloseout />
      </MockProviders>,
    );

    await waitFor(() => {
      expect(screen.getByTestId('tag')).toHaveTextContent('PPM');
    });

    expect(screen.getByRole('heading', { level: 1, name: 'Complete PPM' })).toBeInTheDocument();
    expect(screen.getByRole('heading', { level: 2, name: /Your final estimated incentive: \$/ })).toBeInTheDocument();
    expect(screen.getByText(move.closeout_office.name, { exact: false }));
  });

  it('routes to the home page when the return to homepage link is clicked', async () => {
    const { shipment } = setUpMocks();

    const mockRoutingConfig = {
      path: customerRoutes.SHIPMENT_PPM_COMPLETE_PATH,
      params: { mtoShipmentId: shipment.id, moveId: shipment.moveTaskOrderId },
    };

    render(
      <MockProviders {...mockRoutingConfig}>
        <FinalCloseout />
      </MockProviders>,
    );

    await waitFor(async () => {
      await userEvent.click(screen.getByText('Return To Homepage'));
    });

    expect(mockNavigate).toHaveBeenCalledWith(customerRoutes.MOVE_HOME_PAGE);
  });

  it('submits the ppm signed certification', async () => {
    const { shipment } = setUpMocks();
    submitPPMShipmentSignedCertification.mockResolvedValueOnce(shipment.ppmShipment);

    const mockRoutingConfig = {
      path: customerRoutes.SHIPMENT_PPM_COMPLETE_PATH,
      params: { mtoShipmentId: shipment.id, moveId: shipment.moveTaskOrderId },
    };

    render(
      <MockProviders {...mockRoutingConfig}>
        <FinalCloseout />
      </MockProviders>,
    );
    await waitFor(() => {
      expect(screen.getByTestId('tag')).toHaveTextContent('PPM');
    });

    await userEvent.type(screen.getByRole('textbox', { name: 'Signature' }), 'Grace Griffin');
    await userEvent.click(screen.getByRole('button', { name: 'Submit PPM Documentation' }));

    await waitFor(() =>
      expect(submitPPMShipmentSignedCertification).toHaveBeenCalledWith(shipment.ppmShipment.id, {
        certification_text: ppmSubmissionCertificationText,
        signature: 'Grace Griffin',
        date: formatDateForSwagger(new Date()),
      }),
    );

    expect(updateMTOShipment).toHaveBeenCalledWith(shipment);
    expect(mockDispatch).toHaveBeenCalledTimes(2);

    expect(mockNavigate).toHaveBeenCalledWith(customerRoutes.MOVE_HOME_PAGE);
  });
});
