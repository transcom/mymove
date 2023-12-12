import React, { useState } from 'react';
import { useParams, useNavigate, generatePath } from 'react-router-dom';
import { Alert } from '@trussworks/react-uswds';
import classnames from 'classnames';
import { useMutation } from '@tanstack/react-query';
import { connect } from 'react-redux';
import { func } from 'prop-types';

import { createSITExtensionRequest } from 'services/primeApi';
import scrollToTop from 'shared/scrollToTop';
import CreateSITExtensionRequestForm from 'components/PrimeUI/CreateSITExtensionRequestForm/CreateSITExtensionRequestForm';
import { primeSimulatorRoutes } from 'constants/routes';
import primeStyles from 'pages/PrimeUI/Prime.module.scss';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { usePrimeSimulatorGetMove } from 'hooks/queries';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';
import Shipment from 'components/PrimeUI/Shipment/Shipment';

const CreateSITExtensionRequest = ({ setFlashMessage }) => {
  const { moveCodeOrID, shipmentId } = useParams();
  const navigate = useNavigate();

  const [errorMessage, setErrorMessage] = useState();

  const { moveTaskOrder, isLoading, isError } = usePrimeSimulatorGetMove(moveCodeOrID);

  const { mutate: createSITExtensionRequestMutation } = useMutation(createSITExtensionRequest, {
    onSuccess: () => {
      setFlashMessage(
        `MSG_CREATE_SIT_EXTENSION_REQUEST_SUCCESS${moveCodeOrID}`,
        'success',
        'Successfully created SIT extension request',
        '',
        true,
      );

      navigate(generatePath(primeSimulatorRoutes.VIEW_MOVE_PATH, { moveCodeOrID }));
    },
    onError: (error) => {
      const { response: { body } = {} } = error;

      if (body) {
        let additionalDetails = '';
        if (body.invalidFields) {
          Object.keys(body.invalidFields).forEach((key) => {
            additionalDetails += `:\n${key} - ${body.invalidFields[key]}`;
          });
        }

        setErrorMessage({
          title: `Prime API: ${body.title} `,
          detail: `${body.detail}${additionalDetails}`,
        });
      } else {
        setErrorMessage({
          title: 'Unexpected error',
          detail: 'An unknown error has occurred, please check the state of the shipment data for this move',
        });
      }
      scrollToTop();
    },
  });

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const shipment = moveTaskOrder.mtoShipments.find((s) => s.id === shipmentId);

  return (
    <div className={classnames('grid-container-desktop-lg', 'usa-prose', primeStyles.primeContainer)}>
      <div className="grid-row">
        <div className="grid-col-12">
          <h1>Create SIT Extension Request</h1>
          {errorMessage?.detail && (
            <div className={primeStyles.errorContainer}>
              <Alert headingLevel="h4" slim type="error">
                <span className={primeStyles.errorTitle}>{errorMessage.title}</span>
                <span className={primeStyles.errorDetail}>{errorMessage.detail}</span>
              </Alert>
            </div>
          )}
          <Shipment shipment={shipment} />
          <CreateSITExtensionRequestForm shipment={shipment} submission={createSITExtensionRequestMutation} />
        </div>
      </div>
    </div>
  );
};

CreateSITExtensionRequest.propTypes = {
  setFlashMessage: func.isRequired,
};

const mapDispatchToProps = {
  setFlashMessage: setFlashMessageAction,
};

export default connect(() => ({}), mapDispatchToProps)(CreateSITExtensionRequest);
