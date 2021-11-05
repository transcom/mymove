import React, { useState } from 'react';
import { useParams, useHistory } from 'react-router-dom';
import classnames from 'classnames';
import { useMutation } from 'react-query';
import { generatePath } from 'react-router';
import { Alert } from '@trussworks/react-uswds';

import { usePrimeSimulatorGetMove } from 'hooks/queries';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import CreateShipmentServiceItemForm from 'components/PrimeUI/CreateShipmentServiceItemForm/CreateShipmentServiceItemForm';
import { createServiceItem } from 'services/primeApi';
import { primeSimulatorRoutes } from 'constants/routes';
import scrollToTop from 'shared/scrollToTop';
import primeStyles from 'pages/PrimeUI/Prime.module.scss';

const CreateServiceItem = () => {
  const { moveCodeOrID, shipmentId } = useParams();
  const history = useHistory();

  const [errorMessage, setErrorMessage] = useState();

  const { moveTaskOrder, isLoading, isError } = usePrimeSimulatorGetMove(moveCodeOrID);

  const [createServiceItemMutation] = useMutation(createServiceItem, {
    onSuccess: () => {
      history.push(generatePath(primeSimulatorRoutes.VIEW_MOVE_PATH, { moveCodeOrID }));
    },
    onError: (error) => {
      const { response: { body } = {} } = error;

      if (body) {
        setErrorMessage({ title: body.title, detail: body.detail });
      } else {
        setErrorMessage({
          title: 'Unexpected error',
          detail:
            'An unknown error has occurred, please check the state of the shipment and service items data for this move',
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
          <h1>Create Shipment Service Item</h1>
          {errorMessage?.detail && (
            <div className={primeStyles.errorContainer}>
              <Alert slim type="error">
                <span className={primeStyles.errorTitle}>{errorMessage.title}</span>
                <span className={primeStyles.errorDetail}>{errorMessage.detail}</span>
              </Alert>
            </div>
          )}
          <CreateShipmentServiceItemForm shipment={shipment} createServiceItemMutation={createServiceItemMutation} />
        </div>
      </div>
    </div>
  );
};

export default CreateServiceItem;
