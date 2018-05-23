import PropTypes from 'prop-types';
import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { get } from 'lodash';

import { loadServiceMember } from 'scenes/ServiceMembers/ducks';
import { showCurrentOrders, deleteUpload, addUploads } from './ducks';
import Uploader from 'shared/Uploader';
import UploadsTable from 'shared/Uploader/UploadsTable';
import WizardPage from 'shared/WizardPage';

import './UploadOrders.css';

export class UploadOrders extends Component {
  constructor(props) {
    super(props);

    this.state = {
      newUploads: [],
      showAmendedOrders: false,
    };

    this.onChange = this.onChange.bind(this);
    this.deleteFile = this.deleteFile.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.setShowAmendedOrders = this.setShowAmendedOrders.bind(this);
  }

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

  handleSubmit() {
    this.props.addUploads(this.state.newUploads);
  }

  setShowAmendedOrders(show) {
    this.setState({ showAmendedOrders: show });
  }

  onChange(files) {
    this.setState({
      newUploads: files,
    });
  }

  deleteFile(e, uploadId) {
    e.preventDefault();
    this.props.deleteUpload(uploadId);
  }

  render() {
    const {
      pages,
      pageKey,
      hasSubmitSuccess,
      error,
      currentOrders,
      uploads,
    } = this.props;
    const isValid = Boolean(uploads.length || this.state.newUploads.length);
    const isDirty = Boolean(this.state.newUploads.length);
    return (
      <WizardPage
        handleSubmit={this.handleSubmit}
        pageList={pages}
        pageKey={pageKey}
        pageIsValid={isValid}
        pageIsDirty={isDirty}
        hasSucceeded={hasSubmitSuccess}
        error={error}
      >
        <div>
          <h1 className="sm-heading">Upload Your Orders</h1>
          <p>
            In order to schedule your move, we need to have a complete copy of
            your orders.
          </p>
          <p>
            You can upload a PDF, or you can take a picture of each page and
            upload the images.
          </p>
        </div>
        {Boolean(uploads.length) && (
          <Fragment>
            <br />
            <UploadsTable uploads={uploads} onDelete={this.deleteFile} />
          </Fragment>
        )}
        {currentOrders && (
          <div className="uploader-box">
            <Uploader
              document={currentOrders.uploaded_orders}
              onChange={this.onChange}
            />
            <div className="hint">(Each page must be clear and legible)</div>
          </div>
        )}

        {/* TODO: Uncomment when we support upload of amended orders */}
        {/* <div className="amended-orders">
          <p>
            Do you have amended orders? If so, you need to upload those as well.
          </p>
          <YesNoBoolean
            value={showAmendedOrders}
            onChange={this.setShowAmendedOrders}
          />
          {this.state.showAmendedOrders && (
            <div className="uploader-box">
              <h4>Upload amended orders</h4>
              <Uploader document={{}} onChange={no_op} />
              <div className="hint">(Each page must be clear and legible)</div>
            </div>
          )}
        </div> */}
      </WizardPage>
    );
  }
}

UploadOrders.propTypes = {
  hasSubmitSuccess: PropTypes.bool.isRequired,
  showCurrentOrders: PropTypes.func.isRequired,
  deleteUpload: PropTypes.func.isRequired,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    { showCurrentOrders, loadServiceMember, deleteUpload, addUploads },
    dispatch,
  );
}
function mapStateToProps(state) {
  const props = {
    currentOrders: state.orders.currentOrders,
    uploads: get(state, 'orders.currentOrders.uploaded_orders.uploads', []),
    user: state.loggedInUser,
    ...state.orders,
  };
  return props;
}
export default connect(mapStateToProps, mapDispatchToProps)(UploadOrders);
