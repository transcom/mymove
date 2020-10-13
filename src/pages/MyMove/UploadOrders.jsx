/* eslint-disable */
import PropTypes from 'prop-types';
import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { get } from 'lodash';

import { loadServiceMember } from 'scenes/ServiceMembers/ducks';
import {
  fetchLatestOrders,
  selectActiveOrLatestOrders,
  selectUploadsForActiveOrders,
} from 'shared/Entities/modules/orders';

import { createUpload, deleteUpload, selectDocument } from 'shared/Entities/modules/documents';
import OrdersUploader from 'components/OrdersUploader';
import UploadsTable from 'shared/Uploader/UploadsTable';
import WizardPage from 'shared/WizardPage';
import { documentSizeLimitMsg } from 'shared/constants';

import './UploadOrders.css';
import { no_op } from 'shared/utils';

const uploaderLabelIdle = 'Drag & drop or <span class="filepond--label-action">click to upload orders</span>';

export class UploadOrders extends Component {
  constructor(props) {
    super(props);

    this.state = {
      newUploads: [],
      showAmendedOrders: false,
    };

    this.onChange = this.onChange.bind(this);
    this.deleteFile = this.deleteFile.bind(this);
    this.setShowAmendedOrders = this.setShowAmendedOrders.bind(this);
  }

  componentDidMount() {
    const { serviceMemberId } = this.props;
    this.props.fetchLatestOrders(serviceMemberId);
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
    if (this.props.currentOrders) {
      this.props.deleteUpload(uploadId);
    }
  }

  render() {
    const { pages, pageKey, error, currentOrders, uploads, document, additionalParams } = this.props;
    const isValid = Boolean(uploads.length || this.state.newUploads.length);
    const isDirty = Boolean(this.state.newUploads.length);
    return (
      <WizardPage
        additionalParams={additionalParams}
        dirty={isDirty}
        error={error}
        handleSubmit={no_op}
        pageIsValid={isValid}
        pageKey={pageKey}
        pageList={pages}
      >
        <div>
          <h1>Upload your orders</h1>
          <p>In order to schedule your move, we need to have a complete copy of your orders.</p>
          <p>You can upload a PDF, or you can take a picture of each page and upload the images.</p>
          <p>{documentSizeLimitMsg}</p>
        </div>
        {Boolean(uploads.length) && (
          <Fragment>
            <br />
            <UploadsTable uploads={uploads} onDelete={this.deleteFile} />
          </Fragment>
        )}
        {currentOrders && (
          <div className="uploader-box">
            <OrdersUploader
              createUpload={this.props.createUpload}
              deleteUpload={this.props.deleteUpload}
              document={document}
              onChange={this.onChange}
              options={{ labelIdle: uploaderLabelIdle }}
            />
            <div className="hint">(Each page must be clear and legible.)</div>
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
  deleteUpload: PropTypes.func.isRequired,
};

function mapStateToProps(state) {
  const serviceMemberId = get(state, 'serviceMember.currentServiceMember.id');
  const currentOrders = selectActiveOrLatestOrders(state);

  const props = {
    serviceMemberId: serviceMemberId,
    currentOrders,
    uploads: selectUploadsForActiveOrders(state),
    document: selectDocument(state, currentOrders.uploaded_orders),
  };
  return props;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ fetchLatestOrders, createUpload, deleteUpload, loadServiceMember }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(UploadOrders);
