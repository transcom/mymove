import React, { useState } from 'react';
import { Alert, Button } from '@trussworks/react-uswds';
import { Link, useParams } from 'react-router-dom';
import classnames from 'classnames';
import { useQueryClient, useMutation } from '@tanstack/react-query';
import { func } from 'prop-types';
import { connect } from 'react-redux';

import styles from './MoveDetails.module.scss';

import { PRIME_SIMULATOR_MOVE } from 'constants/queryKeys';
import Shipment from 'components/PrimeUI/Shipment/Shipment';
import FlashGridContainer from 'containers/FlashGridContainer/FlashGridContainer';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import Inaccessible from 'shared/Inaccessible';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import SectionWrapper from 'components/Customer/SectionWrapper';
import formStyles from 'styles/form.module.scss';
import descriptionListStyles from 'styles/descriptionList.module.scss';
import primeStyles from 'pages/PrimeUI/Prime.module.scss';
import { usePrimeSimulatorGetMove } from 'hooks/queries';
import { completeCounseling, deleteShipment, downloadMoveOrder } from 'services/primeApi';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';
import scrollToTop from 'shared/scrollToTop';
import { SIT_SERVICE_ITEMS_ALLOWED_UPDATE } from 'constants/serviceItems';
import { MoveOrderDocumentType } from 'shared/constants';
import { CHECK_SPECIAL_ORDERS_TYPES, SPECIAL_ORDERS_TYPES } from 'constants/orders';

const MoveDetails = ({ setFlashMessage }) => {
  const { moveCodeOrID } = useParams();

  const [errorMessage, setErrorMessage] = useState();

  const [documentTypeKey, setDocumentTypeKey] = useState(MoveOrderDocumentType.ALL);

  const { moveTaskOrder, isLoading, isError, errors } = usePrimeSimulatorGetMove(moveCodeOrID);

  const queryClient = useQueryClient();
  const { mutate: completeCounselingMutation } = useMutation(completeCounseling, {
    onSuccess: () => {
      setFlashMessage(
        `MSG_COMPLETE_COUNSELING${moveCodeOrID}`,
        'success',
        'Successfully completed counseling',
        '',
        true,
      );

      queryClient.setQueryData([PRIME_SIMULATOR_MOVE, moveCodeOrID], moveTaskOrder);
      queryClient.invalidateQueries([PRIME_SIMULATOR_MOVE, moveCodeOrID]).then(() => {});
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

  const { mutate: downloadMoveOrderMutation } = useMutation(downloadMoveOrder, {
    onSuccess: (response) => {
      // dynamically update DOM to trigger browser to display SAVE AS download file modal
      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;
      const disposition = response.headers['content-disposition'];
      const filenameRegex = /filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/;
      let filename = 'moverOrder.pdf';
      const matches = filenameRegex.exec(disposition);
      if (matches != null && matches[1]) {
        filename = matches[1].replace(/['"]/g, '');
      }
      link.setAttribute('download', filename);

      // Append to html link element page
      document.body.appendChild(link);

      // Start download
      link.click();

      // Clean up and remove the link
      link.parentNode.removeChild(link);

      // erase error messages from previous if exists
      setErrorMessage(null);
    },
    onError: (error) => {
      const { response: { body } = {} } = error;
      if (body) {
        setErrorMessage({
          title: `Prime API: ${body.title} `,
          detail: `${body.detail}`,
        });
      } else {
        // Error message is coming in as byte array(PDF).
        // Need to convert byte array into text/json.
        (async () => {
          let title = 'Unexpected Error: ';
          if (error.response.status === 422) {
            title = 'Unprocessable Entity Error: ';
          }
          const text = await error.response.data.text();
          setErrorMessage({
            title,
            detail: JSON.parse(text).detail,
          });
        })();
      }
    },
  });

  const handleCompleteCounseling = () => {
    completeCounselingMutation({ moveTaskOrderID: moveTaskOrder.id, ifMatchETag: moveTaskOrder.eTag });
  };

  const { mutate: deleteShipmentMutation } = useMutation(deleteShipment, {
    onSuccess: () => {
      setFlashMessage(`MSG_DELETE_SHIPMENT${moveCodeOrID}`, 'success', 'Successfully deleted shipment', '', true);

      queryClient.setQueryData([PRIME_SIMULATOR_MOVE, moveCodeOrID], moveTaskOrder);
      queryClient.invalidateQueries([PRIME_SIMULATOR_MOVE, moveCodeOrID]).then(() => {});
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
  if (isError) {
    return errors?.[0]?.response?.body?.message ? <Inaccessible /> : <SomethingWentWrong />;
  }

  const { mtoShipments, paymentRequests, mtoServiceItems } = moveTaskOrder;

  const handleDownloadOrders = () => {
    downloadMoveOrderMutation({ locator: moveTaskOrder.moveCode, type: documentTypeKey });
  };

  const handleDocumentTypeChange = (e) => {
    setDocumentTypeKey(e.target.value);
  };

  return (
    <div>
      <div className={classnames('grid-container-desktop-lg', 'usa-prose', styles.MoveDetails)}>
        <div className="grid-row">
          <div className="grid-col-12">
            <FlashGridContainer className={styles.flashContainer} data-testid="move-details-flash-grid-container">
              {CHECK_SPECIAL_ORDERS_TYPES(moveTaskOrder?.order?.ordersType) ? (
                <div className={styles.specialMovesLabel}>
                  {SPECIAL_ORDERS_TYPES[`${moveTaskOrder?.order?.ordersType}`]}
                </div>
              ) : null}
              <SectionWrapper className={formStyles.formSection}>
                <dl className={descriptionListStyles.descriptionList}>
                  <div className={styles.moveHeader}>
                    <h2>Move</h2>
                    <Link to="../payment-requests/new" relative="path" className="usa-button usa-button-secondary">
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
                  <div className={descriptionListStyles.row}>
                    <dt>Gun Safe:</dt>
                    <dd>{moveTaskOrder.order.entitlement.gunSafe ? 'yes' : 'no'}</dd>
                  </div>
                  <div className={descriptionListStyles.row}>
                    <Button onClick={handleDownloadOrders}>Download Move Orders</Button>
                    <select
                      onChange={handleDocumentTypeChange}
                      className="usa-select"
                      name="moveOrderDocumentType"
                      id="moveOrderDocumentType"
                      title="moveOrderDocumentType"
                    >
                      <option value={MoveOrderDocumentType.ALL}>ALL</option>
                      <option value={MoveOrderDocumentType.ORDERS}>ORDERS</option>
                      <option value={MoveOrderDocumentType.AMENDMENTS}>AMENDMENTS</option>
                    </select>
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
                    <Link to="../shipments/new" relative="path" className="usa-button usa-button-secondary">
                      Create Shipment
                    </Link>
                  </div>
                  {mtoShipments?.map((mtoShipment) => {
                    return (
                      <div key={mtoShipment.id}>
                        <Shipment
                          shipment={mtoShipment}
                          moveId={moveTaskOrder.id}
                          onDelete={handleDeleteShipment}
                          mtoServiceItems={mtoServiceItems}
                        />
                        <div className={styles.serviceItemHeader}>
                          {moveTaskOrder.mtoServiceItems?.length > 0 && <h2>Service Items</h2>}
                        </div>
                        {moveTaskOrder.mtoServiceItems?.map((serviceItem) => {
                          if (serviceItem.mtoShipmentID === mtoShipment.id) {
                            return (
                              <div className={styles.paymentRequestRows} key={serviceItem.id}>
                                <h3 className={styles.serviceItemHeading}>
                                  {serviceItem.reServiceCode} - {serviceItem.reServiceName}
                                </h3>
                                <div className={styles.uploadBtn}>
                                  {SIT_SERVICE_ITEMS_ALLOWED_UPDATE.includes(serviceItem.reServiceCode) ? (
                                    <Link
                                      className={classnames(styles.editButton, 'usa-button usa-button--outline')}
                                      to={`../mto-service-items/${serviceItem.id}/update`}
                                      relative="path"
                                    >
                                      Edit
                                    </Link>
                                  ) : null}
                                  <Link
                                    to={`../mto-service-items/${serviceItem.id}/upload`}
                                    relative="path"
                                    className="usa-button usa-button-secondary"
                                  >
                                    Upload Document for {serviceItem.reServiceName}
                                  </Link>
                                </div>
                              </div>
                            );
                          }
                          return <div />;
                        })}
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
                            to={`../payment-requests/${paymentRequest.id}/upload`}
                            relative="path"
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

export default connect(() => ({}), mapDispatchToProps)(MoveDetails);
