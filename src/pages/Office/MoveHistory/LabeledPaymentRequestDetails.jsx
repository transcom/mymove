import React from 'react';
import PropTypes from 'prop-types';

import labeledStyles from './LabeledDetails.module.scss';

import { shipmentTypes } from 'constants/shipments';
import { HistoryLogContextShape } from 'constants/MoveHistory/UIDisplay/HistoryLogShape';

const LabeledPaymentRequestDetails = ({ context, getLabeledPaymentRequestDetails }) => {
  let valuesToDisplay = context;

  if (getLabeledPaymentRequestDetails) {
    valuesToDisplay = getLabeledPaymentRequestDetails(context);
  }

  return (
    <>
      <div>
        <b>Move services</b>: {valuesToDisplay.moveServices}
      </div>
      {valuesToDisplay.shipmentServices?.map((shipmentService) => {
        const shipmentType = shipmentTypes[shipmentService.shipmentType];

        return (
          <div key={shipmentService.shipmentId}>
            <br />
            <span className={labeledStyles.shipmentType}>{shipmentType} shipment</span>
            <b>Shipment services</b>: {shipmentService.serviceItems}
          </div>
        );
      })}
    </>
  );
};

LabeledPaymentRequestDetails.propTypes = {
  context: HistoryLogContextShape,
  getLabeledPaymentRequestDetails: PropTypes.func,
};

LabeledPaymentRequestDetails.defaultProps = {
  context: {},
  getLabeledPaymentRequestDetails: null,
};

export default LabeledPaymentRequestDetails;
