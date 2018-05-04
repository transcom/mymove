import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import { loadServiceMember } from 'scenes/ServiceMembers/ducks';
import { showCurrentOrders } from './ducks';
import { no_op } from 'shared/utils';
import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import Uploader from 'shared/Uploader';

import './UploadOrders.css';

const formName = 'upload_orders';
// TODO: Replace no_op with form validation once we load existing uploads
const UploadOrdersWizardForm = reduxifyWizardForm(formName, no_op);

export class UploadOrders extends Component {
  componentDidUpdate(prevProps, prevState) {
    // If we don't have a service member yet, fetch one when loggedInUser loads.
    if (
      !prevProps.user.loggedInUser &&
      this.props.user.loggedInUser &&
      !this.props.currentServiceMember
    ) {
      const serviceMemberID = this.props.user.loggedInUser.service_member.id;
      this.props.loadServiceMember(serviceMemberID);
      this.props.showCurrentOrders(serviceMemberID);
    }
  }

  render() {
    const {
      pages,
      pageKey,
      hasSubmitSuccess,
      error,
      currentOrders,
    } = this.props;
    const initialValues = currentOrders ? currentOrders : null;
    return (
      <UploadOrdersWizardForm
        handleSubmit={no_op}
        className={formName}
        pageList={pages}
        pageKey={pageKey}
        hasSucceeded={hasSubmitSuccess}
        serverError={error}
        initialValues={initialValues}
      >
        <h1 className="sm-heading">Upload Photos or PDFs of Your Orders</h1>
        {currentOrders && (
          <Uploader
            ref={ref => (this.uploader = ref)}
            document={currentOrders.uploaded_orders}
          />
        )}
      </UploadOrdersWizardForm>
    );
  }
}

UploadOrders.propTypes = {
  hasSubmitSuccess: PropTypes.bool.isRequired,
  showCurrentOrders: PropTypes.func.isRequired,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ showCurrentOrders, loadServiceMember }, dispatch);
}
function mapStateToProps(state) {
  const props = {
    currentOrders: state.orders.currentOrders,
    user: state.loggedInUser,
    ...state.orders,
  };
  return props;
}
export default connect(mapStateToProps, mapDispatchToProps)(UploadOrders);
