import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { Button, Fieldset, Form, FormGroup, Label, Radio, Textarea } from '@trussworks/react-uswds';
import { Formik } from 'formik';
import classnames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import { ShipmentPaymentSITBalanceShape } from '../../../types/serviceItems';

import styles from './ServiceItemCard.module.scss';

import ServiceItemCalculations from 'components/Office/ServiceItemCalculations/ServiceItemCalculations';
import { shipmentTypes, shipmentModificationTypes } from 'constants/shipments';
import ShipmentModificationTag from 'components/ShipmentModificationTag/ShipmentModificationTag';
import ShipmentContainer from 'components/Office/ShipmentContainer/ShipmentContainer';
import { formatDateFromIso } from 'shared/formatters';
import { toDollarString } from 'utils/formatters';
import { ShipmentOptionsOneOf } from 'types/shipment';
import { PAYMENT_SERVICE_ITEM_STATUS } from 'shared/constants';
import { PaymentServiceItemParam, MTOServiceItemShape } from 'types/order';
import { allowedServiceItemCalculations, SERVICE_ITEM_CODES } from 'constants/serviceItems';
import DaysInSITAllowance from 'components/Office/DaysInSITAllowance/DaysInSITAllowance';

const isAdditionalDaySIT = (mtoServiceItemCode) => {
  return mtoServiceItemCode === SERVICE_ITEM_CODES.DOASIT || mtoServiceItemCode === SERVICE_ITEM_CODES.DDASIT;
};

/** This component represents a Payment Request Service Item */
const ServiceItemCard = ({
  id,
  mtoShipmentType,
  mtoShipmentDepartureDate,
  mtoShipmentPickupAddress,
  mtoShipmentDestinationAddress,
  mtoShipmentModificationType,
  mtoServiceItemCode,
  mtoServiceItemName,
  amount,
  status,
  rejectionReason,
  patchPaymentServiceItem,
  requestComplete,
  paymentServiceItemParams,
  additionalServiceItemData,
  shipmentSITBalance,
}) => {
  const [calculationsVisible, setCalulationsVisible] = useState(false);
  const [canEditRejection, setCanEditRejection] = useState(!rejectionReason);

  const { APPROVED, DENIED } = PAYMENT_SERVICE_ITEM_STATUS;

  const toggleCalculations =
    mtoServiceItemCode &&
    allowedServiceItemCalculations.includes(mtoServiceItemCode) &&
    paymentServiceItemParams.length > 0 ? (
      <>
        <Button
          className={styles.toggleCalculations}
          type="button"
          data-testid="toggleCalculations"
          aria-expanded={calculationsVisible}
          unstyled
          onClick={() => {
            setCalulationsVisible((isVisible) => {
              return !isVisible;
            });
          }}
        >
          {calculationsVisible ? 'Hide calculations' : 'Show calculations'}
        </Button>
        {calculationsVisible && (
          <div className={styles.calculationsContainer}>
            <ServiceItemCalculations
              totalAmountRequested={amount * 100}
              serviceItemParams={paymentServiceItemParams}
              additionalServiceItemData={additionalServiceItemData}
              itemCode={mtoServiceItemCode}
              shipmentType={mtoShipmentType}
              tableSize="small"
            />
          </div>
        )}
      </>
    ) : null;

  if (requestComplete) {
    return (
      <div data-testid="ServiceItemCard" id={`card-${id}`} className={styles.ServiceItemCard}>
        <ShipmentContainer className={styles.shipmentContainerCard} shipmentType={mtoShipmentType}>
          <div className={styles.cardHeader}>
            <h3>
              {shipmentTypes[`${mtoShipmentType}`] || 'BASIC SERVICE ITEMS'}
              {mtoShipmentModificationType && (
                <ShipmentModificationTag shipmentModificationType={mtoShipmentModificationType} />
              )}
            </h3>
            {(mtoShipmentDepartureDate || mtoShipmentPickupAddress || mtoShipmentPickupAddress) && (
              <small className={styles.addressBlock}>
                {mtoShipmentDepartureDate && (
                  <div>
                    <span>Departed</span> {formatDateFromIso(mtoShipmentDepartureDate, 'DD MMM YYYY')}
                  </div>
                )}
                {mtoShipmentPickupAddress && (
                  <div>
                    <span>From</span> {mtoShipmentPickupAddress}
                  </div>
                )}
                {mtoShipmentPickupAddress && (
                  <div>
                    <span>To</span> {mtoShipmentDestinationAddress}
                  </div>
                )}
              </small>
            )}
          </div>
          <hr className="divider" />
          <dl>
            <dt>Service item</dt>
            <dd data-testid="serviceItemName">{mtoServiceItemName}</dd>

            <dt>Amount</dt>
            <dd data-testid="serviceItemAmount">{toDollarString(amount)}</dd>
          </dl>
          {toggleCalculations}
          <div data-testid="completeSummary" className={styles.completeContainer}>
            {status === APPROVED ? (
              <div data-testid="statusHeading" className={classnames(styles.statusHeading, styles.statusApproved)}>
                <FontAwesomeIcon icon="check" />
                Accepted
              </div>
            ) : (
              <>
                <div data-testid="statusHeading" className={classnames(styles.statusHeading, styles.statusRejected)}>
                  <FontAwesomeIcon icon="times" aria-hidden />
                  Rejected
                </div>
                {rejectionReason && (
                  <p data-testid="rejectionReason" className={styles.rejectionReason}>
                    {rejectionReason}
                  </p>
                )}
              </>
            )}
          </div>
        </ShipmentContainer>
      </div>
    );
  }

  return (
    <div data-testid="ServiceItemCard" id={`card-${id}`} className={styles.ServiceItemCard}>
      <Formik
        initialValues={{ status, rejectionReason }}
        onSubmit={(values) => {
          return patchPaymentServiceItem(id, values);
        }}
        enableReinitialize
      >
        {({ initialValues, handleReset, handleChange, submitForm, values, setValues }) => {
          const handleApprovalChange = (event) => {
            handleChange(event);
            submitForm().then(() => {
              setCanEditRejection(true);
            });
          };

          const handleRejectChange = (event) => {
            handleChange(event);
            submitForm().then(() => {
              setCanEditRejection(false);
            });
          };

          const handleRejectCancel = (event) => {
            if (initialValues.rejectionReason) {
              setCanEditRejection(false);
            }

            handleReset(event);
          };

          const handleFormReset = () => {
            setValues({
              status: 'REQUESTED',
              rejectionReason: '',
            });
            submitForm().then(() => {
              setCanEditRejection(true);
            });
          };

          return (
            <Form className={styles.form} onSubmit={submitForm}>
              <ShipmentContainer className={styles.shipmentContainerCard} shipmentType={mtoShipmentType}>
                <div className={styles.cardHeader}>
                  <h3>
                    {shipmentTypes[`${mtoShipmentType}`] || 'BASIC SERVICE ITEMS'}
                    {mtoShipmentModificationType && (
                      <ShipmentModificationTag shipmentModificationType={mtoShipmentModificationType} />
                    )}
                  </h3>
                  {(mtoShipmentDepartureDate || mtoShipmentPickupAddress || mtoShipmentPickupAddress) && (
                    <small className={styles.addressBlock}>
                      {mtoShipmentDepartureDate && (
                        <div>
                          <span>Departed</span> {formatDateFromIso(mtoShipmentDepartureDate, 'DD MMM YYYY')}
                        </div>
                      )}
                      {mtoShipmentPickupAddress && (
                        <div>
                          <span>From</span> {mtoShipmentPickupAddress}
                        </div>
                      )}
                      {mtoShipmentPickupAddress && (
                        <div>
                          <span>To</span> {mtoShipmentDestinationAddress}
                        </div>
                      )}
                    </small>
                  )}
                </div>
                <hr className={styles.divider} />
                <dl>
                  <dt>Service item</dt>
                  <dd data-testid="serviceItemName">{mtoServiceItemName}</dd>
                  {isAdditionalDaySIT(mtoServiceItemCode) && (
                    <>
                      <dt className={styles.daysInSIT}>SIT days invoiced</dt>
                      <dd>
                        <DaysInSITAllowance
                          className={styles.daysInSITDetails}
                          shipmentPaymentSITBalance={shipmentSITBalance}
                        />
                      </dd>
                    </>
                  )}
                  <dt>Amount</dt>
                  <dd data-testid="serviceItemAmount">{toDollarString(amount)}</dd>
                </dl>
                {toggleCalculations}
                <Fieldset>
                  <div className={classnames(styles.statusOption, { [styles.selected]: values.status === APPROVED })}>
                    <Radio
                      id={`approve-${id}`}
                      checked={values.status === APPROVED}
                      value={APPROVED}
                      name="status"
                      label="Approve"
                      onChange={handleApprovalChange}
                      data-testid="approveRadio"
                    />
                  </div>
                  <div className={classnames(styles.statusOption, { [styles.selected]: values.status === DENIED })}>
                    <Radio
                      id={`reject-${id}`}
                      checked={values.status === DENIED}
                      value={DENIED}
                      name="status"
                      label="Reject"
                      onChange={handleChange}
                      data-testid="rejectRadio"
                    />

                    {values.status === DENIED && (
                      <FormGroup>
                        <Label htmlFor={`rejectReason-${id}`}>Reason for rejection</Label>
                        {!canEditRejection && (
                          <>
                            <p data-testid="rejectionReasonReadOnly">{values.rejectionReason}</p>
                            <Button
                              type="button"
                              unstyled
                              data-testid="editReasonButton"
                              className={styles.clearStatus}
                              onClick={() => setCanEditRejection(true)}
                              aria-label="Edit reason button"
                            >
                              <span className="icon">
                                <FontAwesomeIcon icon="pen" title="Edit reason" alt="" />
                              </span>
                              <span aria-hidden="true">Edit reason</span>
                            </Button>
                          </>
                        )}

                        {!requestComplete && canEditRejection && (
                          <>
                            <Textarea
                              id={`rejectReason-${id}`}
                              name="rejectionReason"
                              onChange={handleChange}
                              value={values.rejectionReason}
                            />
                            <div className={styles.rejectionButtonGroup}>
                              <Button
                                id="rejectionSaveButton"
                                type="button"
                                data-testid="rejectionSaveButton"
                                onClick={handleRejectChange}
                                disabled={!values.rejectionReason}
                                aria-label="Rejection save button"
                              >
                                Save
                              </Button>
                              <Button
                                data-testid="cancelRejectionButton"
                                secondary
                                onClick={handleRejectCancel}
                                type="button"
                                aria-label="Cancel rejection button"
                              >
                                Cancel
                              </Button>
                            </div>
                          </>
                        )}
                      </FormGroup>
                    )}
                  </div>

                  {(values.status === APPROVED || values.status === DENIED) && (
                    <Button
                      type="button"
                      unstyled
                      data-testid="clearStatusButton"
                      className={styles.clearStatus}
                      onClick={handleFormReset}
                      aria-label="Clear status"
                    >
                      <span className="icon">
                        <FontAwesomeIcon icon="times" title="Clear status" alt=" " />
                      </span>
                      <span aria-hidden="true">Clear selection</span>
                    </Button>
                  )}
                </Fieldset>
              </ShipmentContainer>
            </Form>
          );
        }}
      </Formik>
    </div>
  );
};

ServiceItemCard.propTypes = {
  id: PropTypes.string.isRequired,
  mtoServiceItemCode: PropTypes.string,
  mtoShipmentType: ShipmentOptionsOneOf,
  mtoShipmentDepartureDate: PropTypes.string,
  mtoShipmentDestinationAddress: PropTypes.node,
  mtoShipmentPickupAddress: PropTypes.node,
  mtoShipmentModificationType: PropTypes.oneOf(Object.values(shipmentModificationTypes)),
  mtoServiceItemName: PropTypes.string,
  amount: PropTypes.number.isRequired,
  status: PropTypes.string,
  rejectionReason: PropTypes.string,
  patchPaymentServiceItem: PropTypes.func.isRequired,
  requestComplete: PropTypes.bool,
  paymentServiceItemParams: PropTypes.arrayOf(PaymentServiceItemParam),
  additionalServiceItemData: MTOServiceItemShape,
  shipmentSITBalance: ShipmentPaymentSITBalanceShape,
};

ServiceItemCard.defaultProps = {
  mtoServiceItemCode: null,
  mtoShipmentType: null,
  mtoShipmentDepartureDate: '',
  mtoShipmentDestinationAddress: '',
  mtoShipmentPickupAddress: '',
  mtoShipmentModificationType: undefined,
  mtoServiceItemName: null,
  status: undefined,
  rejectionReason: '',
  requestComplete: false,
  paymentServiceItemParams: [],
  additionalServiceItemData: {},
  shipmentSITBalance: undefined,
};

export default ServiceItemCard;
