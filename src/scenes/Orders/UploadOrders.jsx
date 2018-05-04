import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import { loadOrders } from './ducks';
import { no_op } from 'shared/utils';
import WizardPage from 'shared/WizardPage';
import Uploader from 'shared/Uploader';

import './UploadOrders.css';

export class UploadOrders extends Component {
  componentDidMount() {
    this.props.loadOrders(this.props.match.params.serviceMemberId);
  }

  render() {
    const {
      pages,
      pageKey,
      hasSubmitSuccess,
      error,
      currentOrders,
    } = this.props;
    return (
      <WizardPage
        handleSubmit={no_op}
        isAsync={true}
        pageList={pages}
        pageKey={pageKey}
        pageIsValid={true}
        hasSucceeded={hasSubmitSuccess}
        error={error}
      >
        <h1 className="sm-heading">Upload Photos or PDFs of Your Orders</h1>
        {currentOrders && (
          <Uploader
            ref={ref => (this.uploader = ref)}
            document={currentOrders.uploaded_orders}
          />
        )}
      </WizardPage>
    );
  }
}

UploadOrders.propTypes = {
  hasSubmitSuccess: PropTypes.bool.isRequired,
  loadOrders: PropTypes.func.isRequired,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ loadOrders }, dispatch);
}
function mapStateToProps(state) {
  console.log(state.orders.currentOrders);
  const props = {
    currentOrders: state.orders.currentOrders,
    ...state.orders,
  };
  return props;
}
export default connect(mapStateToProps, mapDispatchToProps)(UploadOrders);
