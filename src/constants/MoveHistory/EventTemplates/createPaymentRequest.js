import o from 'constants/MoveHistory/UIDisplay/Operations';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.INSERT,
  eventName: o.createPaymentRequest,
  tableName: t.payment_requests,
  detailsType: d.LABELED,
  getEventNameDisplay: ({ changedValues }) => `Submitted payment request ${changedValues?.payment_request_number}`,
  getDetailsLabeledDetails: ({ context }) => {
    let moveServices = '';
    let shipmentServices = '';
    context.forEach((serviceItem) => {
      if (serviceItem.name === 'Move management' || serviceItem.name === 'Counseling') {
        moveServices += `, ${serviceItem.name}`;
      } else {
        shipmentServices += `, ${serviceItem.name}`;
      }
    });
    return { move_services: moveServices.slice(2), shipment_services: shipmentServices.slice(2), shipment_type: 'HHG' };
  },
};
