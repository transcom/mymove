import React, { useState } from 'react';
import { Alert, Button } from '@trussworks/react-uswds';
import { Link, useParams, withRouter } from 'react-router-dom';
import classnames from 'classnames';
import { queryCache, useMutation } from 'react-query';
import { func } from 'prop-types';
import { connect } from 'react-redux';

import styles from './MoveDetails.module.scss';

import { PRIME_SIMULATOR_MOVE } from 'constants/queryKeys';
import Shipment from 'components/PrimeUI/Shipment/Shipment';
import FlashGridContainer from 'containers/FlashGridContainer/FlashGridContainer';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import SectionWrapper from 'components/Customer/SectionWrapper';
import formStyles from 'styles/form.module.scss';
import descriptionListStyles from 'styles/descriptionList.module.scss';
import primeStyles from 'pages/PrimeUI/Prime.module.scss';
import { usePrimeSimulatorGetMove } from 'hooks/queries';
import { completeCounseling, deleteShipment } from 'services/primeApi';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';
import scrollToTop from 'shared/scrollToTop';

const MoveDetails = ({ setFlashMessage }) => {
  const { moveCodeOrID } = useParams();

  const [errorMessage, setErrorMessage] = useState();

  const { moveTaskOrder, isLoading, isError } = usePrimeSimulatorGetMove(moveCodeOrID);

  const [completeCounselingMutation] = useMutation(completeCounseling, {
    onSuccess: () => {
      setFlashMessage(
        `MSG_COMPLETE_COUNSELING${moveCodeOrID}`,
        'success',
        'Successfully completed counseling',
        '',
        true,
      );

      queryCache.setQueryData([PRIME_SIMULATOR_MOVE, moveCodeOrID], moveTaskOrder);
      queryCache.invalidateQueries([PRIME_SIMULATOR_MOVE, moveCodeOrID]).then(() => {});
    },
    onError: (error) => {
      const { response: { body } = {} } = error;

      if (body) {
        setErrorMessage({
          title: `Prime API: ${body.title} `,
          detail: `${body.detail}`,
        });
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

  const handleCompleteCounseling = () => {
    completeCounselingMutation({ moveTaskOrderID: moveTaskOrder.id, ifMatchETag: moveTaskOrder.eTag });
  };

  const [deleteShipmentMutation] = useMutation(deleteShipment, {
    onSuccess: () => {
      setFlashMessage(`MSG_DELETE_SHIPMENT${moveCodeOrID}`, 'success', 'Successfully deleted shipment', '', true);

      queryCache.setQueryData([PRIME_SIMULATOR_MOVE, moveCodeOrID], moveTaskOrder);
      queryCache.invalidateQueries([PRIME_SIMULATOR_MOVE, moveCodeOrID]).then(() => {});
    },
    onError: (error) => {
      const { response: { body } = {} } = error;

      if (body) {
        setErrorMessage({
          title: `Prime API: ${body.title} `,
          detail: `${body.detail}`,
        });
      } else {
        setErrorMessage({
          title: 'Unexpected error',
          detail: 'An unknown error has occurred, please check the state of the shipment for this move',
        });
      }
      scrollToTop();
    },
  });

  const handleDeleteShipment = (mtoShipmentID) => {
    deleteShipmentMutation({ mtoShipmentID });
  };

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
                    {!moveTaskOrder.primeCounselingCompletedAt && (
                      <Button onClick={handleCompleteCounseling}>Complete Counseling</Button>
                    )}
                  </div>
                  {errorMessage?.detail && (
                    <div className={primeStyles.errorContainer}>
                      <Alert slim type="error">
                        <span className={primeStyles.errorTitle}>{errorMessage.title}</span>
                        <span className={primeStyles.errorDetail}>{errorMessage.detail}</span>
                      </Alert>
                    </div>
                  )}
                  <div className={descriptionListStyles.row}>
                    <dt>Move Code:</dt>
                    <dd>{moveTaskOrder.moveCode}</dd>
                  </div>
                  <div className={descriptionListStyles.row}>
                    <dt>Move Id:</dt>
                    <dd>{moveTaskOrder.id}</dd>
                  </div>
                  {moveTaskOrder.primeCounselingCompletedAt && (
                    <div className={descriptionListStyles.row}>
                      <dt>Prime Counseling Completed At:</dt>
                      <dd>{moveTaskOrder.primeCounselingCompletedAt}</dd>
                    </div>
                  )}
                </dl>
              </SectionWrapper>
              <SectionWrapper className={formStyles.formSection}>
                <dl className={descriptionListStyles.descriptionList}>
                  <div className={styles.mainShipmentHeader}>
                    <h2>Shipments</h2>
                    <Link
                      to={`/simulator/moves/${moveTaskOrder.id}/shipments/new`}
                      className="usa-button usa-button-secondary"
                    >
                      Create Shipment
                    </Link>
                  </div>
                  {mtoShipments?.map((mtoShipment) => {
                    return (
                      <div key={mtoShipment.id}>
                        <Shipment shipment={mtoShipment} moveId={moveTaskOrder.id} onDelete={handleDeleteShipment} />
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

MoveDetails.propTypes = {
  setFlashMessage: func.isRequired,
};

const mapDispatchToProps = {
  setFlashMessage: setFlashMessageAction,
};

export default withRouter(connect(() => ({}), mapDispatchToProps)(MoveDetails));
