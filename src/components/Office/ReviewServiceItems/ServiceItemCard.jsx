import React from 'react';
import PropTypes from 'prop-types';
import { Radio, Textarea, FormGroup, Fieldset, Label, Button, Form } from '@trussworks/react-uswds';
import { Formik } from 'formik';

import styles from './ServiceItemCard.module.scss';

import ShipmentContainer from 'components/Office/ShipmentContainer';
import { mtoShipmentTypeToFriendlyDisplay, toDollarString } from 'shared/formatters';
import { ShipmentOptionsOneOf } from 'types/shipment';
import { SERVICE_ITEM_STATUS } from 'shared/constants';

const ServiceItemCard = ({
  id,
  shipmentType,
  serviceItemName,
  amount,
  status,
  rejectionReason,
  patchPaymentServiceItem,
}) => {
  const { APPROVED, REJECTED } = SERVICE_ITEM_STATUS;

  return (
    <div data-testid="ServiceItemCard" className={styles.ServiceItemCard}>
      <Formik
        initialValues={{ status, rejectionReason }}
        onSubmit={(values) => {
          patchPaymentServiceItem(id, values);
        }}
      >
        {({ handleChange, submitForm, handleReset, values }) => {
          const handleApprovalChange = (event) => {
            handleChange(event);
            submitForm();
          };

          return (
            <Form className={styles.form} onSubmit={submitForm}>
              <ShipmentContainer shipmentType={shipmentType}>
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
                      id="approve"
                      checked={values.status === APPROVED}
                      value={APPROVED}
                      name="status"
                      label="Approve"
                      onChange={handleApprovalChange}
                    />
                  </div>
                  <div className={styles.statusOption}>
                    <Radio
                      id="reject"
                      checked={values.status === REJECTED}
                      value={REJECTED}
                      name="status"
                      label="Reject"
                      onChange={handleChange}
                    />

                    {values.status === REJECTED && (
                      <FormGroup>
                        <Label htmlFor="rejectReason">Reason for rejection</Label>
                        <Textarea
                          id="rejectReason"
                          name="rejectionReason"
                          onChange={handleChange}
                          value={values.rejectionReason}
                        />
                      </FormGroup>
                    )}
                  </div>

                  {(values.status === APPROVED || values.status === REJECTED) && (
                    <Button
                      type="button"
                      unstyled
                      data-testid="clearStatusButton"
                      className={styles.clearStatus}
                      onClick={() => {
                        handleReset();
                      }}
                    >
                      X Clear selection
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
  serviceItemName: PropTypes.string.isRequired,
  amount: PropTypes.number.isRequired,
  status: PropTypes.string,
  rejectionReason: PropTypes.string,
  patchPaymentServiceItem: PropTypes.func.isRequired,
};

ServiceItemCard.defaultProps = {
  shipmentType: null,
  status: undefined,
  rejectionReason: '',
};

export default ServiceItemCard;
