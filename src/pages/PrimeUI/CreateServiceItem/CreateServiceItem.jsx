import React from 'react';
import { useParams } from 'react-router-dom';
import classnames from 'classnames';

import { usePrimeSimulatorGetMove } from '../../../hooks/queries';
import LoadingPlaceholder from '../../../shared/LoadingPlaceholder';
import SomethingWentWrong from '../../../shared/SomethingWentWrong';
import styles from '../CreatePaymentRequest/CreatePaymentRequest.module.scss';
import CreateShipmentServiceItemForm from '../../../components/PrimeUI/CreateShipmentServiceItemForm/CreateShipmentServiceItemForm';

const CreateServiceItem = () => {
  const { moveCodeOrID, shipmentId } = useParams();
  // const history = useHistory();

  const { moveTaskOrder, isLoading, isError } = usePrimeSimulatorGetMove(moveCodeOrID);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const shipment = moveTaskOrder.mtoShipments.find((s) => s.id === shipmentId);

  return (
    <div className={classnames('grid-container-desktop-lg', 'usa-prose', styles.CreatePaymentRequest)}>
      <div className="grid-row">
        <div className="grid-col-12">
          <h1>Create Shipment Service Item</h1>
          <CreateShipmentServiceItemForm shipment={shipment} />
        </div>
      </div>
    </div>
  );
};

export default CreateServiceItem;
