import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';

import { MTOServiceItemShape, MTOShipmentShape, PaymentRequestShape } from 'types/index';
import { MatchShape } from 'types/router';
import samplePDF from 'components/DocumentViewer/sample.pdf';
import styles from 'pages/TIO/PaymentRequestReview.module.scss';
import DocumentViewer from 'components/DocumentViewer/DocumentViewer';
import ReviewServiceItems from 'components/Office/ReviewServiceItems/ReviewServiceItems';
import { getPaymentRequest as getPaymentRequestAction } from 'shared/Entities/modules/paymentRequests';
import {
  getMTOShipments as getMTOShipmentsAction,
  selectMTOShipmentsByMTOId,
} from 'shared/Entities/modules/mtoShipments';
import {
  getMTOServiceItems as getMTOServiceItemsAction,
  selectMTOServiceItemsByMTOId,
} from 'shared/Entities/modules/mtoServiceItems';
import { patchPaymentServiceItemStatus as patchPaymentServiceItemStatusAction } from 'shared/Entities/modules/paymentServiceItems';

export class PaymentRequestReview extends Component {
  componentDidMount() {
    const { match, getPaymentRequest, getMTOServiceItems, getMTOShipments } = this.props;
    const { paymentRequestId } = match.params;
    getPaymentRequest(paymentRequestId).then(({ entities: { paymentRequests } }) => {
      const pr = paymentRequests[`${paymentRequestId}`];
      getMTOShipments(pr.moveTaskOrderID);
      getMTOServiceItems(pr.moveTaskOrderID);
    });
  }

  handleUpdatePaymentServiceItemStatus = (paymentServiceItemID, values) => {
    const { patchPaymentServiceItemStatus, mtoServiceItems, paymentRequest } = this.props;
    const paymentServiceItemForRequest = paymentRequest.serviceItems.find((s) => s.id === paymentServiceItemID);
    patchPaymentServiceItemStatus(
      mtoServiceItems[0].moveTaskOrderID,
      paymentServiceItemID,
      values.status,
      paymentServiceItemForRequest.eTag,
      values.rejectionReason,
    );
  };

  handleClose = (moveOrderId) => {
    // eslint-disable-next-line react/prop-types
    const { history } = this.props;

    // eslint-disable-next-line react/prop-types
    history.push(`/moves/${moveOrderId}/payment-requests/`);
  };

  render() {
    // eslint-disable-next-line react/prop-types
    const { moveOrderId, mtoServiceItems, mtoShipments, paymentRequest } = this.props;

    const testFiles = [
      {
        filename: 'Test File.pdf',
        fileType: 'pdf',
        filePath: samplePDF,
      },
    ];

    const serviceItemCards = paymentRequest?.serviceItems?.map((item) => {
      const mtoServiceItem = mtoServiceItems.find((s) => s.id === item.mtoServiceItemID);
      const itemShipment = mtoServiceItem && mtoShipments.find((s) => s.id === mtoServiceItem.mtoShipmentID);

      return {
        id: item.id,
        shipmentId: mtoServiceItem?.mtoShipmentID,
        shipmentType: itemShipment?.shipmentType,
        serviceItemName: mtoServiceItem?.reServiceName,
        amount: item.priceCents ? item.priceCents / 100 : 0,
        createdAt: item.createdAt,
        status: item.status,
        rejectionReason: item.rejectionReason,
      };
    });

    return (
      <div data-testid="PaymentRequestReview" className={styles.PaymentRequestReview}>
        <div className={styles.embed}>
          <DocumentViewer files={testFiles} />
        </div>
        {paymentRequest?.serviceItems && (
          <div className={styles.sidebar}>
            <ReviewServiceItems
              handleClose={() => this.handleClose(moveOrderId)}
              serviceItemCards={serviceItemCards}
              patchPaymentServiceItem={this.handleUpdatePaymentServiceItemStatus}
            />
          </div>
        )}
      </div>
    );
  }
}

PaymentRequestReview.propTypes = {
  match: MatchShape.isRequired,
  getPaymentRequest: PropTypes.func.isRequired,
  getMTOServiceItems: PropTypes.func.isRequired,
  getMTOShipments: PropTypes.func.isRequired,
  paymentRequest: PaymentRequestShape,
  patchPaymentServiceItemStatus: PropTypes.func.isRequired,
  mtoServiceItems: PropTypes.arrayOf(MTOServiceItemShape),
  mtoShipments: PropTypes.arrayOf(MTOShipmentShape),
};

PaymentRequestReview.defaultProps = {
  paymentRequest: null,
  mtoServiceItems: [],
  mtoShipments: [],
};

const mapStateToProps = (state, ownProps) => {
  const { moveOrderId, paymentRequestId } = ownProps.match.params;
  const paymentRequest = state.entities.paymentRequests && state.entities.paymentRequests[`${paymentRequestId}`];

  return {
    paymentRequest,
    mtoServiceItems: paymentRequest && selectMTOServiceItemsByMTOId(state, paymentRequest.moveTaskOrderID),
    mtoShipments: paymentRequest && selectMTOShipmentsByMTOId(state, paymentRequest.moveTaskOrderID),
    moveOrderId,
  };
};

const mapDispatchToProps = {
  getPaymentRequest: getPaymentRequestAction,
  getMTOServiceItems: getMTOServiceItemsAction,
  getMTOShipments: getMTOShipmentsAction,
  patchPaymentServiceItemStatus: patchPaymentServiceItemStatusAction,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(PaymentRequestReview));
