import React from 'react';
import PropTypes from 'prop-types';
import { Radio, Textarea, FormGroup, Fieldset, Label, Button } from '@trussworks/react-uswds';

import styles from './ServiceItemCard.module.scss';

import ShipmentContainer from 'components/Office/ShipmentContainer';
import { mtoShipmentTypeToFriendlyDisplay, toDollarString } from 'shared/formatters';

const ServiceItemCard = ({ id, shipmentType, serviceItemName, amount, onChange, value, clearValues }) => {
  const { status, rejectionReason } = value;

  return (
    <div data-testid="ServiceItemCard" className={styles.ServiceItemCard}>
      <ShipmentContainer shipmentType={shipmentType}>
        <>
          <h6 data-cy="shipmentTypeHeader" className={styles.cardHeader}>
            {mtoShipmentTypeToFriendlyDisplay(shipmentType)?.toUpperCase() || 'BASIC SERVICE ITEMS'}
          </h6>
          <div className="usa-label">Service item</div>
          <div data-cy="serviceItemName" className={styles.textValue}>
            {serviceItemName}
          </div>
          <div className="usa-label">Amount</div>
          <div data-cy="serviceItemAmount" className={styles.textValue}>
            {toDollarString(amount)}
          </div>
          <Fieldset className={styles.statusOption}>
            <Radio
              id="approve"
              checked={status === 'APPROVED'}
              value="APPROVED"
              name={`${id}.status`}
              label="Approve"
              onChange={onChange}
            />
          </Fieldset>
          <Fieldset className={styles.statusOption}>
            <Radio
              id="reject"
              checked={status === 'REJECTED'}
              value="REJECTED"
              name={`${id}.status`}
              label="Reject"
              onChange={onChange}
            />

            {status === 'REJECTED' && (
              <FormGroup>
                <Label htmlFor="rejectReason">Reason for rejection</Label>
                <Textarea
                  id="rejectReason"
                  name={`${id}.rejectionReason`}
                  onChange={onChange}
                  value={rejectionReason}
                />
              </FormGroup>
            )}
          </Fieldset>
          {(status === 'APPROVED' || status === 'REJECTED') && (
            <Button
              type="button"
              unstyled
              className={styles.clearStatus}
              onClick={() => {
                clearValues(id);
              }}
            >
              X Clear selection
            </Button>
          )}
        </>
      </ShipmentContainer>
    </div>
  );
};

ServiceItemCard.propTypes = {
  id: PropTypes.string.isRequired,
  shipmentType: PropTypes.string,
  serviceItemName: PropTypes.string.isRequired,
  amount: PropTypes.number.isRequired,
  onChange: PropTypes.func,
  value: PropTypes.shape({
    status: PropTypes.string,
    rejectionReason: PropTypes.string,
  }),
  clearValues: PropTypes.func,
};

ServiceItemCard.defaultProps = {
  shipmentType: '',
  onChange: () => {},
  value: {
    status: undefined,
    rejectionReason: '',
  },
  clearValues: () => {},
};

export default ServiceItemCard;
