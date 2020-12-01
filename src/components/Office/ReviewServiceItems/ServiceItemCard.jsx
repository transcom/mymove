import React from 'react';
import PropTypes from 'prop-types';
import { Radio, Textarea, FormGroup, Fieldset, Label, Button, Form } from '@trussworks/react-uswds';
import { Formik } from 'formik';
import classnames from 'classnames';

import styles from './ServiceItemCard.module.scss';

import ShipmentContainer from 'components/Office/ShipmentContainer';
import { mtoShipmentTypeToFriendlyDisplay, toDollarString } from 'shared/formatters';
import { ShipmentOptionsOneOf } from 'types/shipment';
import { PAYMENT_SERVICE_ITEM_STATUS } from 'shared/constants';
import { ReactComponent as XHeavyIcon } from 'shared/icon/x-heavy.svg';
import { ReactComponent as XLightIcon } from 'shared/icon/x-light.svg';
import { ReactComponent as FormCheckmarkIcon } from 'shared/icon/form-checkmark.svg';

/** This component represents a Payment Request Service Item */
const ServiceItemCard = ({
  id,
  shipmentType,
  serviceItemName,
  amount,
  status,
  rejectionReason,
  patchPaymentServiceItem,
  requestComplete,
}) => {
  const { APPROVED, DENIED } = PAYMENT_SERVICE_ITEM_STATUS;

  if (requestComplete) {
    return (
      <div data-testid="ServiceItemCard" id={`card-${id}`} className={styles.ServiceItemCard}>
        <ShipmentContainer className={styles.shipmentContainerCard} shipmentType={shipmentType}>
          <h6 className={styles.cardHeader}>
            {mtoShipmentTypeToFriendlyDisplay(shipmentType) || 'BASIC SERVICE ITEMS'}
          </h6>
          <dl>
            <dt>Service item</dt>
            <dd data-testid="serviceItemName">{serviceItemName}</dd>

            <dt>Amount</dt>
            <dd data-testid="serviceItemAmount">{toDollarString(amount)}</dd>
          </dl>

          <div data-testid="completeSummary" className={styles.completeContainer}>
            {status === APPROVED ? (
              <div data-testid="statusHeading" className={classnames(styles.statusHeading, styles.statusApproved)}>
                <FormCheckmarkIcon />
                Accepted
              </div>
            ) : (
              <>
                <div data-testid="statusHeading" className={classnames(styles.statusHeading, styles.statusRejected)}>
                  <XHeavyIcon />
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
          patchPaymentServiceItem(id, values);
        }}
      >
        {({ handleChange, submitForm, values, setValues }) => {
          const handleApprovalChange = (event) => {
            handleChange(event);
            submitForm();
          };

          const handleFormReset = () => {
            setValues({
              status: 'REQUESTED',
              rejectionReason: undefined,
            });
            submitForm();
          };

          return (
            <Form className={styles.form} onSubmit={submitForm}>
              <ShipmentContainer className={styles.shipmentContainerCard} shipmentType={shipmentType}>
                <h6 className={styles.cardHeader}>
                  {mtoShipmentTypeToFriendlyDisplay(shipmentType) || 'BASIC SERVICE ITEMS'}
                </h6>
                <dl>
                  <dt>Service item</dt>
                  <dd data-testid="serviceItemName">{serviceItemName}</dd>

                  <dt>Amount</dt>
                  <dd data-testid="serviceItemAmount">{toDollarString(amount)}</dd>
                </dl>
                <Fieldset>
                  <div className={styles.statusOption}>
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
                  <div className={styles.statusOption}>
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
                        <Label htmlFor="rejectReason">Reason for rejection</Label>
                        <Textarea
                          id={`rejectReason-${id}`}
                          name="rejectionReason"
                          onChange={handleChange}
                          value={values.rejectionReason}
                        />
                        {!requestComplete && (
                          <div className={styles.rejectionButtonGroup}>
                            <Button type="button" data-testid="rejectionSaveButton" onClick={submitForm}>
                              Save
                            </Button>
                            <Button
                              data-testid="cancelRejectionButton"
                              secondary
                              onClick={handleFormReset}
                              type="button"
                            >
                              Cancel
                            </Button>
                          </div>
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
                    >
                      <span className="icon">
                        <XLightIcon />
                      </span>
                      Clear selection
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
  shipmentType: ShipmentOptionsOneOf,
  serviceItemName: PropTypes.string,
  amount: PropTypes.number.isRequired,
  status: PropTypes.string,
  rejectionReason: PropTypes.string,
  patchPaymentServiceItem: PropTypes.func.isRequired,
  requestComplete: PropTypes.bool,
};

ServiceItemCard.defaultProps = {
  shipmentType: null,
  serviceItemName: null,
  status: undefined,
  rejectionReason: '',
  requestComplete: false,
};

export default ServiceItemCard;
