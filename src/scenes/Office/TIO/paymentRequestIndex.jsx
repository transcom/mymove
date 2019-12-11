import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { getPaymentRequestList } from 'shared/Entities/modules/paymentRequests';

class PaymentRequestIndex extends React.Component {
  componentDidMount() {
    this.props.getPaymentRequestList();
  }

  render() {
    return <h1>Helloooooooo</h1>;
  }
}

const mapStateToProps = (state, ownProps) => ({});

const mapDispatchToProps = dispatch => bindActionCreators({ getPaymentRequestList }, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(PaymentRequestIndex);
