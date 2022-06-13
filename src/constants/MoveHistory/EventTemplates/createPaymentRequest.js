import o from 'constants/MoveHistory/UIDisplay/Operations';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.INSERT,
  eventName: o.createPaymentRequest,
  tableName: t.payment_requests,
  detailsType: d.LABELED_PAYMENT_REQUEST,
  getEventNameDisplay: ({ changedValues }) => `Submitted payment request ${changedValues?.payment_request_number}`,
  getLabeledPaymentRequestDetails: (context) => {
    let moveServices = '';
    const shipmentServices = {};
    context.forEach((serviceItem) => {
      if (serviceItem.name === 'Move management' || serviceItem.name === 'Counseling') {
        moveServices += `, ${serviceItem.name}`;
      } else {
        const shipmentId = serviceItem.shipment_id;
        if (shipmentServices[shipmentId]) {
          const { serviceItems } = shipmentServices[shipmentId];
          shipmentServices[shipmentId].serviceItems = `${serviceItems}, ${serviceItem.name}`;
        } else {
          shipmentServices[shipmentId] = {
            serviceItems: serviceItem.name,
            shipmentType: serviceItem.shipment_type,
            shipmentId,
          };
        }
      }
    });
    return { moveServices: moveServices.slice(2), shipmentServices: Object.values(shipmentServices) };
  },
};
