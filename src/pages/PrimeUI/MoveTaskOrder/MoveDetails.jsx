import React from 'react';
import { useParams, Link } from 'react-router-dom';
import classnames from 'classnames';

import styles from './MoveDetails.module.scss';

import Shipment from 'components/PrimeUI/Shipment/Shipment';
import FlashGridContainer from 'containers/FlashGridContainer/FlashGridContainer';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import SectionWrapper from 'components/Customer/SectionWrapper';
import formStyles from 'styles/form.module.scss';
import descriptionListStyles from 'styles/descriptionList.module.scss';
import { usePrimeSimulatorGetMove } from 'hooks/queries';

const MoveDetails = () => {
  const { moveCodeOrID } = useParams();

  const { moveTaskOrder, isLoading, isError } = usePrimeSimulatorGetMove(moveCodeOrID);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const { mtoShipments, paymentRequests } = moveTaskOrder;

  return (
    <div>
      <div className={classnames('grid-container-desktop-lg', 'usa-prose', styles.MoveDetails)}>
        <div className="grid-row">
          <div className="grid-col-12">
            <FlashGridContainer className={styles.flashContainer} data-testid="move-details-flash-grid-container">
              <SectionWrapper className={formStyles.formSection}>
                <dl className={descriptionListStyles.descriptionList}>
                  <div className={styles.moveHeader}>
                    <h2>Move</h2>
                    <Link to="payment-requests/new" className="usa-button usa-button-secondary">
                      Create Payment Request
                    </Link>
                  </div>
                  <div className={descriptionListStyles.row}>
                    <dt>Move Code:</dt>
                    <dd>{moveTaskOrder.moveCode}</dd>
                  </div>
                  <div className={descriptionListStyles.row}>
                    <dt>Move Id:</dt>
                    <dd>{moveTaskOrder.id}</dd>
                  </div>
                </dl>
              </SectionWrapper>
              <SectionWrapper className={formStyles.formSection}>
                <dl className={descriptionListStyles.descriptionList}>
                  <h2>Shipments</h2>
                  <Link
                    to={`/simulator/moves/${moveTaskOrder.id}/shipments/new`}
                    className="usa-button usa-button-secondary"
                  >
                    Create Shipment
                  </Link>
                  {mtoShipments?.map((mtoShipment) => {
                    return (
                      <div key={mtoShipment.id}>
                        <Shipment shipment={mtoShipment} moveId={moveTaskOrder.id} />
                      </div>
                    );
                  })}
                </dl>
              </SectionWrapper>
              {paymentRequests?.length > 0 && (
                <SectionWrapper className={formStyles.formSection}>
                  <dl className={descriptionListStyles.descriptionList}>
                    <h2>Payment Requests</h2>
                    {paymentRequests?.map((paymentRequest) => {
                      return (
                        <div className={styles.paymentRequestRow} key={paymentRequest.id}>
                          <div data-testid="paymentRequestNumber">{paymentRequest.paymentRequestNumber}</div>
                          <Link
                            to={`payment-requests/${paymentRequest.id}/upload`}
                            className="usa-button usa-button-secondary"
                          >
                            Upload Document
                          </Link>
                        </div>
                      );
                    })}
                  </dl>
                </SectionWrapper>
              )}
            </FlashGridContainer>
          </div>
        </div>
      </div>
    </div>
  );
};

export default MoveDetails;
