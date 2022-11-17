import React from 'react';

import labeledStyles from './LabeledDetails.module.scss';

import { shipmentTypes } from 'constants/shipments';
import { PaymentRequestServicesShape } from 'constants/MoveHistory/UIDisplay/HistoryLogShape';

const LabeledPaymentRequestDetails = ({ services }) => {
  return (
    <>
      <div>
        <b>Move services</b>: {services.moveServices}
      </div>
      {services.shipmentServices?.map((shipmentService) => {
        const shipmentType = shipmentTypes[shipmentService.shipmentType];
        const shipmentID = shipmentService.shipmentIdAbbr;

        return (
          <div key={shipmentService.shipmentId}>
            <br />
            <span className={labeledStyles.shipmentType}>
              {shipmentType} shipment #{shipmentID}
            </span>
            <b>Shipment services</b>: {shipmentService.serviceItems}
          </div>
        );
      })}
    </>
  );
};

LabeledPaymentRequestDetails.propTypes = {
  services: PaymentRequestServicesShape,
};

LabeledPaymentRequestDetails.defaultProps = {
  services: { moveServices: null, shipmentServices: null },
};

export default LabeledPaymentRequestDetails;
