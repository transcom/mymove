import React from 'react';
import { withRouter } from 'react-router-dom-old';
import { Formik, Form, Field } from 'formik';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { selectPaymentRequest, getPaymentRequest, updatePaymentRequest } from 'shared/Entities/modules/paymentRequests';

class PaymentRequestShow extends React.Component {
  componentDidMount() {
    this.props.getPaymentRequest(this.props.id);
  }

  updatePaymentRequest = (paymentRequest = {}) => {
    paymentRequest.status = 'REVIEWED';
    paymentRequest.paymentRequestID = this.props.id;
    paymentRequest.ifMatchETag = this.props.paymentRequest.eTag;
    this.props.updatePaymentRequest(paymentRequest);
  };

  render() {
    const {
      id,
      paymentRequest: { isFinal, rejectionReason, serviceItemIDs, status },
    } = this.props;
    return (
      <>
        <div>
          <h1>Payment Request Id {id}</h1>
          <ul>
            <li>isFinal: {`${isFinal}`}</li>
            <li>rejectionReason: {rejectionReason}</li>
            <li>serviceItemIds: {serviceItemIDs}</li>
            <li>status: {status}</li>
          </ul>
          <button className="usa-button usa-button--secondary" onClick={this.updatePaymentRequest}>
            Approve
          </button>

          <Formik
            initialValues={{ rejectionReason: '' }}
            onSubmit={(values, { setSubmitting }) => {
              this.updatePaymentRequest({ rejectionReason: values.rejectionReason });
              setSubmitting(false);
            }}
          >
            {({ isSubmitting }) => (
              <Form>
                <Field type="text" name="rejectionReason" />
                <button className="usa-button usa-button--secondary" type="submit" disabled={isSubmitting}>
                  Reject
                </button>
              </Form>
            )}
          </Formik>
        </div>
      </>
    );
  }
}
const mapStateToProps = (state, props) => {
  const id = props.match.params.id;
  return {
    id,
    paymentRequest: selectPaymentRequest(state, id),
  };
};

const mapDispatchToProps = (dispatch) => bindActionCreators({ getPaymentRequest, updatePaymentRequest }, dispatch);

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(PaymentRequestShow));
