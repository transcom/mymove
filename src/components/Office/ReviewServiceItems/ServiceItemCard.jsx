import React from 'react';
import PropTypes from 'prop-types';
import { Radio, Textarea, FormGroup, Fieldset, Label, Button } from '@trussworks/react-uswds';

import styles from './ServiceItemCard.module.scss';

import ShipmentContainer from 'components/Office/ShipmentContainer';
import { mtoShipmentTypeToFriendlyDisplay, toDollarString } from 'shared/formatters';
import { ShipmentOptionsOneOf } from 'types/shipment';
import { SERVICE_ITEM_STATUS } from 'shared/constants';

const ServiceItemCard = ({ id, shipmentType, serviceItemName, amount, onChange, value, clearValues }) => {
  const { status, rejectionReason } = value;
  const { APPROVED, REJECTED } = SERVICE_ITEM_STATUS;

  return (
    <div data-testid="ServiceItemCard" className={styles.ServiceItemCard}>
      <ShipmentContainer className={styles.shipmentContainerCard} shipmentType={shipmentType}>
        <h6 className={styles.cardHeader}>{mtoShipmentTypeToFriendlyDisplay(shipmentType) || 'BASIC SERVICE ITEMS'}</h6>
        <dl>
          <dt>Service item</dt>
          <dd data-cy="serviceItemName">{serviceItemName}</dd>

          <dt>Amount</dt>
          <dd data-cy="serviceItemAmount">{toDollarString(amount)}</dd>
        </dl>
        <Fieldset>
          <div className={styles.statusOption}>
            <Radio
              id="approve"
              checked={status === APPROVED}
              value={APPROVED}
              name={`${id}.status`}
              label="Approve"
              onChange={onChange}
            />
          </div>
          <div className={styles.statusOption}>
            <Radio
              id="reject"
              checked={status === REJECTED}
              value={REJECTED}
              name={`${id}.status`}
              label="Reject"
              onChange={onChange}
            />

            {status === REJECTED && (
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
          </div>

          {(status === APPROVED || status === REJECTED) && (
            <Button
              type="button"
              unstyled
              data-testid="clearStatusButton"
              className={styles.clearStatus}
              onClick={() => {
                clearValues(id);
              }}
            >
              X Clear selection
            </Button>
          )}
        </Fieldset>
      </ShipmentContainer>
    </div>
  );
};

ServiceItemCard.propTypes = {
  id: PropTypes.string.isRequired,
  shipmentType: ShipmentOptionsOneOf,
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
  shipmentType: null,
  onChange: () => {},
  value: {
    status: undefined,
    rejectionReason: '',
  },
  clearValues: () => {},
};

export default ServiceItemCard;
