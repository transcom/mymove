import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { Radio, Textarea, FormGroup, Fieldset, Label, Button, Form } from '@trussworks/react-uswds';
import { Formik } from 'formik';
import classnames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './ServiceItemCard.module.scss';

import ShipmentContainer from 'components/Office/ShipmentContainer';
import { toDollarString } from 'shared/formatters';
import { ShipmentOptionsOneOf } from 'types/shipment';
import { PAYMENT_SERVICE_ITEM_STATUS } from 'shared/constants';
import { mtoShipmentTypes } from 'constants/shipments';
import ServiceItemCalculations from 'components/Office/ServiceItemCalculations/ServiceItemCalculations';
import { PaymentServiceItemParam } from 'types/order';
import { allowedServiceItemCalculations } from 'constants/serviceItems';

/** This component represents a Payment Request Service Item */
const ServiceItemCard = ({
  id,
  mtoShipmentType,
  mtoServiceItemCode,
  mtoServiceItemName,
  amount,
  status,
  rejectionReason,
  patchPaymentServiceItem,
  requestComplete,
  paymentServiceItemParams,
}) => {
  const [calculationsVisible, setCalulationsVisible] = useState(false);

  const { APPROVED, DENIED } = PAYMENT_SERVICE_ITEM_STATUS;

  const toggleCalculations =
    allowedServiceItemCalculations.includes(mtoServiceItemCode) && paymentServiceItemParams.length > 0 ? (
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
              itemCode={mtoServiceItemCode}
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
          <h6 className={styles.cardHeader}>{mtoShipmentTypes[`${mtoShipmentType}`] || 'BASIC SERVICE ITEMS'}</h6>
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
              <ShipmentContainer className={styles.shipmentContainerCard} shipmentType={mtoShipmentType}>
                <h6 className={styles.cardHeader}>{mtoShipmentTypes[`${mtoShipmentType}`] || 'BASIC SERVICE ITEMS'}</h6>
                <dl>
                  <dt>Service item</dt>
                  <dd data-testid="serviceItemName">{mtoServiceItemName}</dd>

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
  mtoServiceItemCode: PropTypes.string.isRequired,
  mtoShipmentType: ShipmentOptionsOneOf,
  mtoServiceItemName: PropTypes.string,
  amount: PropTypes.number.isRequired,
  status: PropTypes.string,
  rejectionReason: PropTypes.string,
  patchPaymentServiceItem: PropTypes.func.isRequired,
  requestComplete: PropTypes.bool,
  paymentServiceItemParams: PropTypes.arrayOf(PaymentServiceItemParam),
};

ServiceItemCard.defaultProps = {
  mtoShipmentType: null,
  mtoServiceItemName: null,
  status: undefined,
  rejectionReason: '',
  requestComplete: false,
  paymentServiceItemParams: [],
};

export default ServiceItemCard;
