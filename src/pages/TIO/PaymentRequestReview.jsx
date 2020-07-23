import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';

import { MatchShape } from 'types/router';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import samplePDF from 'components/DocumentViewer/sample.pdf';
import styles from 'pages/TIO/PaymentRequestReview.module.scss';
import DocumentViewer from 'components/DocumentViewer/DocumentViewer';
import ReviewServiceItems from 'components/Office/ReviewServiceItems/ReviewServiceItems';
import { getPaymentRequest as getPaymentRequestAction } from 'shared/Entities/modules/paymentRequests';

class PaymentRequestReview extends Component {
  componentDidMount() {
    const { match, getPaymentRequest } = this.props;
    const { paymentRequestId } = match.params;
    getPaymentRequest(paymentRequestId);
  }

  handleClose = (moveOrderId) => {
    // eslint-disable-next-line react/prop-types
    const { history } = this.props;

    // eslint-disable-next-line react/prop-types
    history.push(`/moves/${moveOrderId}/payment-requests/`);
  };

  render() {
    // eslint-disable-next-line react/prop-types
    const { serviceItemCards, moveOrderId } = this.props;
    const testFiles = [
      {
        filename: 'Test File.pdf',
        fileType: 'pdf',
        filePath: samplePDF,
      },
    ];

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
};

const mapStateToProps = (state, ownProps) => {
  const { moveOrderId } = ownProps.match.params;
  // TODO need to select from redux store and construct data based on ServiceItemCardsShape
  // test data for now
  const serviceItemCardsTestData = [
    {
      id: '1',
      shipmentType: SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC,
      serviceItemName: 'Domestic linehaul',
      amount: 1234.56,
      status: 'APPROVED',
      createdAt: Date(),
    },
    {
      id: '2',
      shipmentType: SHIPMENT_OPTIONS.NTS,
      serviceItemName: 'Domestic linehaul',
      amount: 1234.56,
      status: 'SUBMITTED',
      createdAt: Date(),
    },
    {
      id: '3',
      shipmentType: null, // to indicate basic service item
      serviceItemName: 'Domestic linehaul',
      amount: 1234.56,
      status: 'REJECTED',
      createdAt: Date(),
    },
  ];

  return {
    serviceItemCards: serviceItemCardsTestData,
    moveOrderId,
  };
};

const mapDispatchToProps = {
  getPaymentRequest: getPaymentRequestAction,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(PaymentRequestReview));
