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

class PaymentRequestReview extends Component {
  componentDidMount() {
    const { match, getPaymentRequest, getMTOServiceItems, getMTOShipments } = this.props;
    const { paymentRequestId } = match.params;
    getPaymentRequest(paymentRequestId).then(({ entities: { paymentRequests } }) => {
      const pr = paymentRequests[`${paymentRequestId}`];
      getMTOShipments(pr.moveTaskOrderID);
      getMTOServiceItems(pr.moveTaskOrderID);
    });
  }

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

    const serviceItemCards = mtoServiceItems.map((item) => {
      const itemShipment = mtoShipments.find((s) => s.id === item.mtoShipmentID);
      const itemPaymentServiceItem = paymentRequest?.serviceItems?.find((s) => s.mtoServiceItemID === item.id);

      return {
        id: item.id,
        shipmentId: item.mtoShipmentID,
        shipmentType: itemShipment?.shipmentType,
        serviceItemName: item.reServiceName,
        amount: itemPaymentServiceItem?.priceCents ? itemPaymentServiceItem.priceCents / 100 : 0,
        createdAt: item.createdAt,
        status: item.status,
      };
    });

    return (
      <div data-testid="PaymentRequestReview" className={styles.PaymentRequestReview}>
        <div className={styles.embed}>
          <DocumentViewer files={testFiles} />
        </div>
        <div className={styles.sidebar}>
          <ReviewServiceItems handleClose={() => this.handleClose(moveOrderId)} serviceItemCards={serviceItemCards} />
        </div>
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
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(PaymentRequestReview));
