import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { getPaymentRequestList, selectPaymentRequests } from 'shared/Entities/modules/paymentRequests';

class PaymentRequestIndex extends React.Component {
  componentDidMount() {
    this.props.getPaymentRequestList();
  }

  render() {
    console.log(this.props.paymentRequests);
    return <h1>Helloooooooo</h1>;
  }
}

const mapStateToProps = state => ({
  paymentRequests: selectPaymentRequests(state),
});

const mapDispatchToProps = dispatch => bindActionCreators({ getPaymentRequestList }, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(PaymentRequestIndex);
