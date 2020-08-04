import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';

import { MTOServiceItemShape, MTOShipmentShape, PaymentRequestShape, PaymentServiceItemShape } from 'types/index';
import { MatchShape, HistoryShape } from 'types/router';
import samplePDF from 'components/DocumentViewer/sample.pdf';
import styles from 'pages/TIO/PaymentRequestReview.module.scss';
import DocumentViewer from 'components/DocumentViewer/DocumentViewer';
import ReviewServiceItems from 'components/Office/ReviewServiceItems/ReviewServiceItems';
import {
  getPaymentRequest as getPaymentRequestAction,
  updatePaymentRequest as updatePaymentRequestAction,
} from 'shared/Entities/modules/paymentRequests';
import {
  getMTOShipments as getMTOShipmentsAction,
  selectMTOShipmentsByMTOId,
} from 'shared/Entities/modules/mtoShipments';
import {
  getMTOServiceItems as getMTOServiceItemsAction,
  selectMTOServiceItemsByMTOId,
} from 'shared/Entities/modules/mtoServiceItems';
import { patchPaymentServiceItemStatus as patchPaymentServiceItemStatusAction } from 'shared/Entities/modules/paymentServiceItems';
import { PAYMENT_REQUEST_STATUS } from 'shared/constants';

export class PaymentRequestReview extends Component {
  constructor(props) {
    super(props);

    this.state = {
      completeReviewError: undefined,
    };
  }

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
    const { patchPaymentServiceItemStatus, mtoServiceItems, paymentServiceItems } = this.props;
    const paymentServiceItemForRequest = paymentServiceItems.find((s) => s.id === paymentServiceItemID);
    patchPaymentServiceItemStatus(
      mtoServiceItems[0].moveTaskOrderID,
      paymentServiceItemID,
      values.status,
      paymentServiceItemForRequest.eTag,
      values.rejectionReason,
    );
  };

  handleCompleteReview = () => {
    const { updatePaymentRequest, paymentRequest, history } = this.props;
    const { completeReviewError } = this.state;
    // first reset error if there was one
    if (completeReviewError) this.setState({ completeReviewError: undefined });

    const newPaymentRequest = {
      paymentRequestID: paymentRequest.id,
      ifMatchETag: paymentRequest.eTag,
      status: PAYMENT_REQUEST_STATUS.REVIEWED,
    };

    updatePaymentRequest(newPaymentRequest)
      .then(() => {
        // TODO - show flash message?
        history.push(`/`); // Go home
      })
      .catch((e) => {
        const error = e.response?.response?.body;
        this.setState({
          completeReviewError: error,
        });
      });
  };

  handleClose = (moveOrderId) => {
    const { history } = this.props;
    history.push(`/moves/${moveOrderId}/payment-requests/`);
  };

  render() {
    // eslint-disable-next-line react/prop-types
    const { moveOrderId, mtoServiceItems, mtoShipments, paymentServiceItems } = this.props;
    const { completeReviewError } = this.state;

    const testFiles = [
      {
        filename: 'Test File.pdf',
        fileType: 'pdf',
        filePath: samplePDF,
      },
    ];

    const serviceItemCards = paymentServiceItems.map((item) => {
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
        <div className={styles.sidebar}>
          <ReviewServiceItems
            handleClose={() => this.handleClose(moveOrderId)}
            serviceItemCards={serviceItemCards}
            patchPaymentServiceItem={this.handleUpdatePaymentServiceItemStatus}
            onCompleteReview={this.handleCompleteReview}
            completeReviewError={completeReviewError}
          />
        </div>
      </div>
    );
  }
}

PaymentRequestReview.propTypes = {
  history: HistoryShape.isRequired,
  match: MatchShape.isRequired,
  getPaymentRequest: PropTypes.func.isRequired,
  getMTOServiceItems: PropTypes.func.isRequired,
  getMTOShipments: PropTypes.func.isRequired,
  paymentRequest: PaymentRequestShape,
  paymentServiceItems: PropTypes.arrayOf(PaymentServiceItemShape),
  patchPaymentServiceItemStatus: PropTypes.func.isRequired,
  updatePaymentRequest: PropTypes.func.isRequired,
  mtoServiceItems: PropTypes.arrayOf(MTOServiceItemShape),
  mtoShipments: PropTypes.arrayOf(MTOShipmentShape),
};

PaymentRequestReview.defaultProps = {
  paymentRequest: undefined,
  paymentServiceItems: [],
  mtoServiceItems: [],
  mtoShipments: [],
};

const mapStateToProps = (state, ownProps) => {
  const { moveOrderId, paymentRequestId } = ownProps.match.params;
  const paymentRequest = state.entities.paymentRequests && state.entities.paymentRequests[`${paymentRequestId}`];
  const paymentServiceItems = [];
  if (paymentRequest?.serviceItems) {
    Object.values(paymentRequest.serviceItems).forEach((id) => {
      if (state.entities.paymentServiceItems[`${id}`])
        paymentServiceItems.push(state.entities.paymentServiceItems[`${id}`]);
    });
  }

  return {
    paymentRequest,
    paymentServiceItems,
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
  updatePaymentRequest: updatePaymentRequestAction,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(PaymentRequestReview));
